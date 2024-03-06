package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/util"
	"net/http"
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
	if clientIp := r.Header.Get("X-FORWARDED-FOR"); clientIp != "" {
		return clientIp
	}
	return r.RemoteAddr
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
	// TODO
	//cfPossibleIP := request.State.Config.CfPossibleIP
	return true
}

func IsInterceptIP(request util.Request, clientIP string) bool {
	// TODO
	return false
}
