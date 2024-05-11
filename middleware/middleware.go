package middleware

import (
	"context"
	"errors"
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
	"github.com/labstack/echo/v4"
	"github.com/nikolalohinski/gonja/v2/exec"
	"gorm.io/gorm"
)

var mbIdChecker = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func MainMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !shouldRunMiddleware(c.Request().URL.Path) {
			return next(c)
		}

		request := c.Get(KeyRequest).(util.Request)

		engine := config.Global.DbEngine
		dbConn, err := db.NewDB(engine)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if !strings.HasPrefix(c.Request().URL.Path, "/install") {
			if !config.ExistENV {
				util.RenderAlertTemplate(c.Response(), request, ".env 파일이 없습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				return nil
			}

			ok, err := dbConn.HasTable(config.Global.DbTablePrefix + "config")
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if !ok {
				util.RenderAlertTemplate(c.Response(), request, "DB 또는 테이블이 존재하지 않습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				return nil
			}
		} else {
			setTemplateCtx(c.Request())
			return next(c)
		}

		cfg := model.Config{}
		if dbConn.Take(&cfg).Error != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		clientIP := lib.GetClientIp(c.Request())

		memberService := service.NewMemberService(dbConn)
		if sessionMbId != "" {
			member, err = memberService.CreateById(sessionMbId)
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if member.IsInterceptOrLeave() {
				request.Session = make(map[string]string)
			}
		} else if cookieMbId != "" {
			mbId := mbIdChecker.ReplaceAllString(cookieMbId, "")
			member, err = memberService.CreateById(mbId)
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if !lib.IsSuperAdmin(request, mbId) &&
				member.IsEmailCertify(cfg.CfUseEmailCertify == 1) &&
				member.IsInterceptOrLeave() {
				ssMbKey = lib.SessionMemberKey(c.Request(), member)
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
			lib.SendHTML(c.Response(), "<meta charset=utf-8>접근이 허용되지 않은 IP 입니다.")
			return nil
		}
		if lib.IsInterceptIP(request, clientIP) {
			lib.SendHTML(c.Response(), "<meta charset=utf-8>접근이 차단된 IP 입니다.")
			return nil
		}

		const secondsOfDay = 60 * 60 * 24
		cookieDomain := request.State.CookieDomain

		if isAutoLogin == true && request.Session["ss_mb_id"] != "" {
			http.SetCookie(c.Response(), &http.Cookie{
				Name:   "ck_mb_id",
				Value:  cookieMbId,
				Domain: cookieDomain,
				MaxAge: secondsOfDay * 30,
			})

			http.SetCookie(c.Response(), &http.Cookie{
				Name:   "ck_auto",
				Value:  ssMbKey,
				Domain: cookieDomain,
				MaxAge: secondsOfDay * 30,
			})
		}

		ckVisitIP, err := c.Request().Cookie("ck_visit_ip")
		if err != nil {
			ckVisitIP = &http.Cookie{}
		}
		if ckVisitIP.Value != clientIP {
			http.SetCookie(c.Response(), &http.Cookie{
				Name:   "ck_visit_ip",
				Value:  clientIP,
				Domain: cookieDomain,
				MaxAge: secondsOfDay,
			})
			err = lib.RecordVisit(c.Request())
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		if !request.State.IsSuperAdmin && !strings.HasPrefix(c.Request().URL.Path, "/admin") {
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
					LoLocation: c.Request().URL.Path,
					LoURL:      c.Request().URL.Path,
				}
				dbConn.Create(&newLogin)
			} else if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			currentLogin.MbID = mbId
			currentLogin.LoDatetime = time.Now()
			currentLogin.LoLocation = c.Request().URL.Path
			currentLogin.LoLocation = c.Request().URL.Path
			dbConn.Save(&currentLogin)
		}

		timeDelta := time.Unix(int64(cfg.CfLoginMinutes)*60, 0)
		dbConn.Where("lo_datetime < ?", timeDelta).Delete(&model.Login{})

		return next(c)
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

func RequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := util.NewRequest(c.Request())
		ctx := context.WithValue(c.Request().Context(), KeyRequest, request)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func setTemplateCtx(r *http.Request) {
	request := r.Context().Value(KeyRequest).(util.Request)
	ctx := context.WithValue(r.Context(), KeyTemplateCtx, newDefaultTemplateCtx(request))
	r.WithContext(ctx)
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

func UrlForMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		setUrlMapForInstall(c.Request())
		return next(c)
	}
}

func setUrlMapForInstall(r *http.Request) {
	request := r.Context().Value(KeyRequest).(util.Request)
	installURL, _ := url.JoinPath(request.BaseURL, "install")
	util.UrlMap.Store("install_license", installURL+"/license")
	util.UrlMap.Store("install_form", installURL+"/form")
	util.UrlMap.Store("install", installURL)
	util.UrlMap.Store("index", "/")
}
