package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/util"
	"github.com/jellydator/ttlcache/v3"
	"github.com/mileusna/useragent"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func CreateDynamicWriteTable(dbConn *db.Database, tableName string) error {
	writeTable := model.NewWriteTable(tableName)
	err := dbConn.Table(writeTable.TableName()).AutoMigrate(writeTable)
	if err != nil {
		return err
	}

	const numReplyIndex = "idx_wr_num_reply"
	newNumReplyIndex := fmt.Sprintf("%s_%s", numReplyIndex, tableName)
	err = dbConn.Table(writeTable.TableName()).Migrator().RenameIndex(writeTable, numReplyIndex, newNumReplyIndex)
	if err != nil {
		return err
	}

	const commentIndex = "idx_wr_is_comment"
	newCommentIndex := fmt.Sprintf("%s_%s", commentIndex, tableName)
	err = dbConn.Table(writeTable.TableName()).Migrator().RenameIndex(writeTable, commentIndex, newCommentIndex)
	if err != nil {
		return err
	}

	return nil
}

func GetClientIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	if strings.Contains(ip, ":") {
		ip = net.ParseIP(ip).To4().String()
	}

	return ip
}

func IsSuperAdmin(request util.Request, mbId string) bool {
	cfg := request.State.Config
	cfAdmin := strings.TrimSpace(strings.ToLower(cfg.CfAdmin))

	if cfAdmin == "" {
		return false
	}

	if mbId == "" {
		mbId = request.Session["ss_mb_id"]
	}

	lowerMbId := strings.TrimSpace(strings.ToLower(mbId))
	if mbId != "" && lowerMbId == cfAdmin {
		return true
	}
	return false
}

func SessionMemberKey(r *http.Request, member *model.Member) string {
	ssMbKeyInput := member.MbDatetime.Format(time.DateTime) + GetClientIp(r) + r.Header.Get("User-Agent")

	hash := md5.Sum([]byte(ssMbKeyInput))
	byteSlice := hash[:]
	return hex.EncodeToString(byteSlice)
}

func IsPossibleIP(request util.Request, clientIP string) bool {
	ipList := request.State.Config.CfPossibleIP
	return checkIPList(request, clientIP, ipList, true)
}

func IsInterceptIP(request util.Request, clientIP string) bool {
	ipList := request.State.Config.CfInterceptIP
	return checkIPList(request, clientIP, ipList, false)
}

func checkIPList(request util.Request, clientIP string, ipList string, allow bool) bool {
	if request.State.IsSuperAdmin {
		return allow
	}

	ipList = strings.TrimSpace(ipList)
	if ipList == "" {
		return allow
	}

	ipPatterns := strings.Split(ipList, "\n")
	for _, pattern := range ipPatterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			return false
		}
		pattern = strings.ReplaceAll(pattern, ".", `\.`)
		pattern = strings.ReplaceAll(pattern, "+", `[0-9\.]+`)
		pattern = fmt.Sprintf("^%s$", pattern)
		isMatch, err := regexp.MatchString(pattern, clientIP)
		if err != nil {
			log.Println(err)
			return false
		}
		if isMatch {
			return true
		}
	}
	return false
}

func RecordVisit(r *http.Request) error {
	dbConn := db.GetInstance()
	viIP := GetClientIp(r)

	var count int64
	today := time.Now().Format(time.DateOnly)
	err := dbConn.Model(&model.Visit{}).Where("vi_date = ? and vi_ip = ?", today, viIP).Count(&count).Error
	if err != nil {
		log.Println(err)
		return err
	}

	if count != 0 {
		return nil
	}

	referer := r.Referer()
	userAgent := r.UserAgent()
	ua := useragent.Parse(userAgent)
	browser := ua.Name
	os := ua.OS
	device := ua.Device
	now := time.Now()
	viDate := now.Format(time.DateOnly)
	viTime := now.Format(time.TimeOnly)
	visit := model.Visit{
		ViIP:      viIP,
		ViDate:    viDate,
		ViTime:    viTime,
		ViReferer: referer,
		ViAgent:   userAgent,
		ViBrowser: browser,
		ViOs:      os,
		ViDevice:  device,
	}
	err = dbConn.Create(&visit).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var visitCountToday int64
	err = dbConn.Model(&visit).Where("vi_date = ?", today).Count(&visitCountToday).Error
	if err != nil {
		log.Println(err)
		return err
	}

	visitSum := model.VisitSum{
		VsDate:  viDate,
		VsCount: visitCountToday,
	}
	err = dbConn.Save(&visitSum).Error
	if err != nil {
		log.Println(err)
		return err
	}

	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	var yesterdayVisitCount int64
	err = dbConn.Model(&visit).Where("vi_date = ?", yesterday).Count(&visitCountToday).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var maxVisitCount int64
	err = dbConn.Model(&visitSum).Select("max(vs_count)").Scan(&maxVisitCount).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var totalVisitCount int64
	err = dbConn.Model(&visitSum).Select("sum(vs_count)").Scan(&totalVisitCount).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var cfg model.Config
	dbConn.First(&cfg)
	cfg.CfVisit = fmt.Sprintf("오늘:%d,어제:%d,최대:%d,전체:%d", visitCountToday, yesterdayVisitCount, maxVisitCount, totalVisitCount)
	err = dbConn.Model(&cfg).Where("cf_id = ?", cfg.CfID).Update("cf_visit", cfg.CfVisit).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const defaultLoginMinute = 10

func GetCurrentLoginCount(request util.Request) [2]int {
	cfg := request.State.Config

	loginMinute := defaultLoginMinute
	if cfg.CfLoginMinutes == 0 {
		loginMinute = cfg.CfLoginMinutes
	}
	baseDate := time.Now().Add(time.Duration(loginMinute) * time.Minute)

	result := struct {
		Login  int `gorm:"login"`
		Member int `gorm:"member"`
	}{}

	dbConn := db.GetInstance()
	dbConn.Model(&model.Login{}).
		Where("mb_id = ? AND lo_ip NOT '' AND lo_datetime > ?", cfg.CfAdmin, baseDate).
		Select("COUNT(mb_id) AS login, CASE mb_id WHEN NOT '' THEN 1 ELSE 0 END member").
		Scan(&result)

	return [2]int{result.Login, result.Member}
}

var menusCache = ttlcache.New[string, []*model.Menu](
	ttlcache.WithTTL[string, []*model.Menu](60*time.Second),
	ttlcache.WithCapacity[string, []*model.Menu](1),
)

func GetMenus() []*model.Menu {
	item := menusCache.Get("menus")
	if item != nil {
		return item.Value()
	}

	//dbConn := db.GetInstance()
	var menus []*model.Menu

	// TODO 부모 메뉴 조회
	// TODO 자식 메뉴 조회

	menusCache.Set("menus", menus, ttlcache.DefaultTTL)
	return menus
}

var lfuCache = ttlcache.New[string, *model.Poll](
	ttlcache.WithCapacity[string, *model.Poll](128))

func GetRecentPoll() *model.Poll {
	if pollItem := lfuCache.Get("poll"); pollItem != nil {
		return pollItem.Value()
	}

	dbConn := db.GetInstance()

	var poll model.Poll
	err := dbConn.Model(&model.Poll{}).Where("po_use == 1").Order("po_id desc").First(&poll).Error
	if err != nil {
		return nil
	}

	lfuCache.Set("poll", &poll, ttlcache.DefaultTTL)
	return &poll
}
