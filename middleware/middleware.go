package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dukhyungkim/gonuboard/config"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/lib"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/service"
	"github.com/dukhyungkim/gonuboard/util"
	"github.com/gin-gonic/gin"
	"github.com/nikolalohinski/gonja/v2/exec"
	"gorm.io/gorm"
)

var mbIdChecker = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldRunMiddleware(c.Request.URL.Path) {
			c.Next()
			return
		}

		request := c.MustGet(KeyRequest).(util.Request)

		engine := config.Global.DbEngine
		dbConn, err := db.NewDB(engine)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
			return
		}

		if !strings.HasPrefix(c.Request.URL.Path, "/install") {
			if !config.ExistENV {
				fmt.Println(".env 파일이 없습니다. 설치를 진행해 주세요.")
				util.RenderAlertTemplate(c, request, ".env 파일이 없습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				c.Abort()
				return
			}

			ok, err := dbConn.HasTable(config.Global.DbTablePrefix + "config")
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
				return
			}
			if !ok {
				util.RenderAlertTemplate(c, request, "DB 또는 테이블이 존재하지 않습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				c.Abort()
				return
			}
		} else {
			c.Set(KeyTemplateCtx, newDefaultTemplateCtx(request))
			c.Next()
			return
		}

		cfg := model.Config{}
		if dbConn.Take(&cfg).Error != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
			return
		}
		request.State.Config = cfg
		request.State.Title = cfg.CfTitle

		request.State.Editor = cfg.CfEditor
		request.State.UseEditor = false
		if cfg.CfEditor != "" {
			request.State.UseEditor = true
		}

		request.State.CookieDomain = config.Global.CookieDomain

		var member *model.Member
		isAutoLogin := false
		ssMbKey := ""
		sessionMbId := request.Session["ss_mb_id"]
		cookieMbId := request.Cookies["ck_mb_id"]
		clientIP := c.ClientIP()

		memberService := service.NewMemberService(dbConn)
		if sessionMbId != "" {
			member, err = memberService.CreateById(sessionMbId)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
				return
			}
			if member.IsInterceptOrLeave() {
				request.Session = make(map[string]string)
			}
		} else if cookieMbId != "" {
			mbId := mbIdChecker.ReplaceAllString(cookieMbId, "")
			member, err = memberService.CreateById(mbId)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
				return
			}
			if !lib.IsSuperAdmin(request, mbId) &&
				member.IsEmailCertify(cfg.CfUseEmailCertify == 1) &&
				member.IsInterceptOrLeave() {
				ssMbKey = lib.SessionMemberKey(c.Request, member)
				if request.Cookies["ck_auto"] == ssMbKey {
					request.Session["ss_mb_id"] = cookieMbId
					isAutoLogin = true
				}
			}
		}

		if member != nil {
			nowDate := time.Now().Format(time.DateOnly)
			if member.MbTodayLogin.Format(time.DateOnly) != nowDate {
				lib.InsertPoint(dbConn, request, member.MbID, member.MbPoint, nowDate+" 첫로그인", "@login", member.MbID, nowDate, 0)
				member.MbTodayLogin = time.Now()
				member.MbLoginIP = clientIP
				dbConn.Model(member).Select("mb_today_login", "mb_login_ip").Updates(member)
			}
		}

		request.State.LoginMember = member
		if member != nil {
			request.State.IsSuperAdmin = lib.IsSuperAdmin(request, member.MbID)
		}

		if !lib.IsPossibleIP(request, clientIP) {
			lib.SendHTML(c.Writer, "<meta charset=utf-8>접근이 허용되지 않은 IP 입니다.")
			c.Abort()
			return
		}
		if lib.IsInterceptIP(request, clientIP) {
			lib.SendHTML(c.Writer, "<meta charset=utf-8>접근이 차단된 IP 입니다.")
			c.Abort()
			return
		}

		const secondsOfDay = 60 * 60 * 24
		cookieDomain := request.State.CookieDomain

		if isAutoLogin == true && request.Session["ss_mb_id"] != "" {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:   "ck_mb_id",
				Value:  cookieMbId,
				Domain: cookieDomain,
				MaxAge: secondsOfDay * 30,
			})

			http.SetCookie(c.Writer, &http.Cookie{
				Name:   "ck_auto",
				Value:  ssMbKey,
				Domain: cookieDomain,
				MaxAge: secondsOfDay * 30,
			})
		}

		ckVisitIP, err := c.Request.Cookie("ck_visit_ip")
		if err != nil {
			ckVisitIP = &http.Cookie{}
		}
		if ckVisitIP.Value != clientIP {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:   "ck_visit_ip",
				Value:  clientIP,
				Domain: cookieDomain,
				MaxAge: secondsOfDay,
			})
			err = lib.RecordVisit(c.Request)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
				return
			}
		}

		if !request.State.IsSuperAdmin && !strings.HasPrefix(c.Request.URL.Path, "/admin") {
			mbId := ""
			if member != nil {
				mbId = member.MbID
			}

			var currentLogin model.Login
			err = dbConn.Where("lo_ip = ?", clientIP).Take(&currentLogin).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newLogin := model.Login{
					LoID:       25,
					LoIP:       clientIP,
					MbID:       mbId,
					LoDatetime: time.Now(),
					LoLocation: c.Request.URL.Path,
					LoURL:      c.Request.URL.Path,
				}
				dbConn.Create(&newLogin)
			} else if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
				return
			}

			currentLogin.MbID = mbId
			currentLogin.LoDatetime = time.Now()
			currentLogin.LoLocation = c.Request.URL.Path
			currentLogin.LoLocation = c.Request.URL.Path
			dbConn.Save(&currentLogin)
		}

		timeDelta := time.Unix(int64(cfg.CfLoginMinutes)*60, 0)
		dbConn.Where("lo_datetime < ?", timeDelta).Delete(&model.Login{})

		c.Next()
	}
}

func shouldRunMiddleware(path string) bool {
	for _, hasPrefix := range []bool{
		strings.HasPrefix(path, "/generate_token"),
		strings.HasPrefix(path, "/device/change"),
		strings.HasPrefix(path, "/static"),
		strings.HasPrefix(path, "/theme_static"),
		endsWith(path, []string{".css", ".js", ".jpg", ".png", ".gif", ".webp"}),
	} {
		if hasPrefix {
			return false
		}
	}
	return true
}

func endsWith(path string, ends []string) bool {
	for _, end := range ends {
		return strings.HasSuffix(path, end)
	}
	return false
}

const (
	KeyRequest     = "request"
	KeyTemplateCtx = "templateCtx"
)

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := util.NewRequest(c.Request)
		c.Set(KeyRequest, request)
		c.Next()
	}
}

func newDefaultTemplateCtx(request util.Request) *exec.Context {
	return exec.NewContext(map[string]interface{}{
		"current_login_count": lib.GetCurrentLoginCount(request),
		"menus":               lib.GetMenus(),
		"poll":                lib.GetRecentPoll(),
		// TODO
		//"populars": get_populars(),
		//"render_latest_posts": render_latest_posts,
		//"render_visit_statistics": render_visit_statistics,
	})
}

func UrlForMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.MustGet(KeyRequest).(util.Request)

		installURL, _ := url.JoinPath(request.BaseURL, "install")
		util.UrlMap.Store("install_license", installURL+"/license")
		util.UrlMap.Store("install_form", installURL+"/form")
		util.UrlMap.Store("install", installURL)
		util.UrlMap.Store("index", "/")
		log.Println(request.BaseURL)

		c.Next()
	}
}
