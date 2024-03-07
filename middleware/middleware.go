package middleware

import (
	"context"
	"github.com/dukhyungkim/gonuboard/config"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/lib"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/service"
	"github.com/dukhyungkim/gonuboard/util"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var mbIdChecker = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func MainMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !shouldRunMiddleware(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		request := r.Context().Value(KeyRequest).(util.Request)

		engine := os.Getenv("DB_ENGINE")
		dbConn, err := db.NewDB(engine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/install") {
			if config.NotExistENV {
				renderAlertTemplate(w, request, ".env 파일이 없습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				return
			}

			ok, err := dbConn.HasTable(model.Prefix + "config")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !ok {
				renderAlertTemplate(w, request, "DB 또는 테이블이 존재하지 않습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				return
			}
		} else {
			next.ServeHTTP(w, r)
			return
		}

		cfg := model.Config{}
		if dbConn.Take(&cfg).Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		request.State.Config = cfg
		request.State.Title = cfg.CfTitle

		request.State.Editor = cfg.CfEditor
		request.State.UseEditor = false
		if cfg.CfEditor != "" {
			request.State.UseEditor = true
		}

		request.State.CookieDomain = os.Getenv("COOKIE_DOMAIN")

		var member *model.Member
		isAutoLogin := false
		ssMbKey := ""
		sessionMbId := request.Session["ss_mb_id"]
		cookieMbId := request.Cookies["ck_mb_id"]
		clientIP := lib.GetClientIp(r)

		memberService := service.NewMemberService(dbConn)
		if sessionMbId != "" {
			member, err = memberService.CreateById(sessionMbId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if member.IsInterceptOrLeave() {
				request.Session = make(map[string]string)
			}
		} else if cookieMbId != "" {
			mbId := mbIdChecker.ReplaceAllString(cookieMbId, "")
			member, err = memberService.CreateById(mbId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !lib.IsSuperAdmin(request, mbId) &&
				member.IsEmailCertify(cfg.CfUseEmailCertify == 1) &&
				member.IsInterceptOrLeave() {
				ssMbKey = lib.SessionMemberKey(r, member)
				if request.Cookies["ck_auto"] == ssMbKey {
					request.Session["ss_mb_id"] = cookieMbId
					isAutoLogin = true
				}
			}
		}

		if member != nil {
			nowDate := time.Now().Format(time.DateOnly)
			if member.MbTodayLogin.Format(time.DateOnly) != nowDate {
				insertPoint(dbConn, request, member, nowDate+" 첫로그인", "@login", nowDate)
				member.MbTodayLogin = time.Now()
				member.MbLoginIP = clientIP
				dbConn.Model(member).Select("mb_today_login", "mb_login_ip").Updates(member)
			}
		}

		request.State.LoginMember = member
		request.State.IsSuperAdmin = lib.IsSuperAdmin(request, member.MbID)

		if !lib.IsPossibleIP(request, clientIP) {
			lib.SendHTML(w, "<meta charset=utf-8>접근이 허용되지 않은 IP 입니다.")
			return
		}
		if lib.IsInterceptIP(request, clientIP) {
			lib.SendHTML(w, "<meta charset=utf-8>접근이 차단된 IP 입니다.")
			return
		}

		const secondsOfDay = 60 * 60 * 24
		cookieDomain := request.State.CookieDomain

		if isAutoLogin == true && request.Session["ss_mb_id"] != "" {
			http.SetCookie(w, &http.Cookie{
				Name:   "ck_mb_id",
				Value:  cookieMbId,
				Domain: cookieDomain,
				MaxAge: secondsOfDay * 30,
			})

			http.SetCookie(w, &http.Cookie{
				Name:   "ck_auto",
				Value:  ssMbKey,
				Domain: cookieDomain,
				MaxAge: secondsOfDay * 30,
			})
		}

		ckVisitIP, err := r.Cookie("ck_visit_ip")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ckVisitIP.Value != clientIP {
			http.SetCookie(w, &http.Cookie{
				Name:   "ck_visit_ip",
				Value:  clientIP,
				Domain: cookieDomain,
				MaxAge: secondsOfDay,
			})
			lib.RecordVisit(r)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func insertPoint(dbConn *db.Database, request util.Request, member *model.Member, content string, relTable string, relAction string) {
	cfg := request.State.Config

	if cfg.CfUsePoint == 0 {
		return
	}

	if cfg.CfLoginPoint == 0 {
		return
	}

	if member == nil {
		return
	}

	if relTable != "" || relAction != "" {
		// TODO
		dbConn.Where("mb_id")
	}

}

func shouldRunMiddleware(path string) bool {
	switch path {
	case "/generate_token":
		return false
	}
	return true
}

func renderAlertTemplate(w http.ResponseWriter, request util.Request, message string, statusCode int, url string) {
	tpl, err := util.AlertTemplate(request, message, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(tpl)
}

type CtxKey string

const (
	KeyRequest CtxKey = "request"
)

func RequestMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		request := util.NewRequest(r)
		ctx := context.WithValue(r.Context(), KeyRequest, request)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func UrlForMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		setUrlMapForInstall(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func setUrlMapForInstall(r *http.Request) {
	request := r.Context().Value(KeyRequest).(util.Request)
	installURL, _ := url.JoinPath(request.BaseURL, "install")
	util.UrlMap.Store("install_license", installURL+"/license")
	util.UrlMap.Store("install_form", installURL+"/form")
	util.UrlMap.Store("install", installURL)
	util.UrlMap.Store("index", "/")
}
