package install

import (
	"errors"
	"fmt"
	"github.com/dukhyungkim/gonuboard/config"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/lib"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/plugin"
	"github.com/dukhyungkim/gonuboard/util"
	"github.com/dukhyungkim/gonuboard/version"
	"github.com/go-chi/chi/v5"
	"github.com/jellydator/ttlcache/v3"
	"github.com/nikolalohinski/gonja/v2/exec"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var formCache = ttlcache.New[string, *installForm](
	ttlcache.WithTTL[string, *installForm](60*time.Second),
	ttlcache.WithCapacity[string, *installForm](1),
)

func DefaultRouter(r chi.Router) {
	r.Get("/", indexHandler())
	r.Post("/", installDatabase())
	r.Get("/license", licenseHandler())
	r.Post("/form", formHandler())
	r.Get("/process", installProcess())
}

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const templatePath = "install/templates/main.html"
		data := exec.NewContext(map[string]any{
			"python_version":  version.RuntimeVersion,
			"fastapi_version": version.RouterVersion,
		})

		util.RenderTemplate(w, templatePath, data)
	}
}

func licenseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		license, err := readLicense()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		const templatePath = "install/templates/license.html"
		data := exec.NewContext(map[string]any{
			"license": license,
		})

		util.RenderTemplate(w, templatePath, data)
	}
}

func readLicense() (string, error) {
	license, err := os.ReadFile("LICENSE")
	if err != nil {
		return "", err
	}
	return string(license), nil
}

func formHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const templatePath = "install/templates/form.html"
		data := exec.NewContext(map[string]any{})

		util.RenderTemplate(w, templatePath, data)
	}
}

func installDatabase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, err := parseInstallForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionSecretKey, err := util.TokenURLSafe(50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = copyFile("example.env", util.EnvPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, setKey := range []func() error{
			util.SetKeyToEnv(util.EnvPath, "DB_ENGINE", form.DBEngine),
			util.SetKeyToEnv(util.EnvPath, "DB_HOST", form.DBHost),
			util.SetKeyToEnv(util.EnvPath, "DB_PORT", form.DBPort),
			util.SetKeyToEnv(util.EnvPath, "DB_USER", form.DBUser),
			util.SetKeyToEnv(util.EnvPath, "DB_PASSWORD", form.DBPassword),
			util.SetKeyToEnv(util.EnvPath, "DB_NAME", form.DBName),
			util.SetKeyToEnv(util.EnvPath, "DB_TABLE_PREFIX", form.DBTablePrefix),
			util.SetKeyToEnv(util.EnvPath, "SESSION_SECRET_KEY", sessionSecretKey),
			util.SetKeyToEnv(util.EnvPath, "COOKIE_DOMAIN", ""),
		} {
			if err = setKey(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		isResponsive, err := strconv.ParseBool(os.Getenv("IS_RESPONSIVE"))
		if err != nil {
			isResponsive = false
		}
		config.IsResponsive = isResponsive

		model.Prefix = form.DBTablePrefix
		// TODO use db handler
		_, err = db.NewDB(form.DBEngine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pluginList, err := plugin.ReadPluginState()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = plugin.WritePluginState(pluginList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formCache.Set("form", form, ttlcache.DefaultTTL)

		const templatePath = "install/templates/result.html"
		util.RenderTemplate(w, templatePath, nil)
	}
}

type installForm struct {
	DBEngine      string
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	DBTablePrefix string

	AdminId       string
	AdminPassword string
	AdminName     string
	AdminEmail    string

	Reinstall bool
}

func parseInstallForm(r *http.Request) (*installForm, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	engine := r.FormValue("db_engine")
	if !db.IsSupportedEngines(engine) {
		return nil, errors.New("지원 가능한 데이터베이스 엔진을 선택해주세요.")
	}

	dbPort, err := strconv.Atoi(r.FormValue("db_port"))
	if err != nil {
		return nil, err
	}

	return &installForm{
		DBEngine:      engine,
		DBHost:        r.FormValue("db_host"),
		DBPort:        dbPort,
		DBUser:        r.FormValue("db_user"),
		DBPassword:    r.FormValue("db_password"),
		DBName:        r.FormValue("db_name"),
		DBTablePrefix: r.FormValue("db_table_prefix"),
		AdminId:       r.FormValue("admin_id"),
		AdminPassword: r.FormValue("admin_password"),
		AdminName:     r.FormValue("admin_name"),
		AdminEmail:    r.FormValue("admin_email"),
		Reinstall:     false,
	}, nil
}

func copyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return nil
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func installProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSSEHeader(w)

		form := formCache.Get("form").Value()

		dbConn, err := db.NewDB(form.DBEngine)
		if err != nil {
			sendSSE(w, failedInstallMessage(err))
			return
		}
		sendSSE(w, "데이터베이스 연결 완료")

		if form.Reinstall {
			// TODO reinstall
			sendSSE(w, "기존 데이터베이스 테이블 삭제 완료")
		}

		err = dbConn.MigrateTables()
		if err != nil {
			sendSSE(w, failedInstallMessage(err))
			return
		}
		sendSSE(w, "데이터베이스 테이블 생성 완료")

		err = setupConfig(dbConn, form.AdminId, form.AdminEmail)
		if err != nil {
			sendSSE(w, failedInstallMessage(err))
			return
		}
		err = setupAdminMember(dbConn, form.AdminId, form.AdminName, form.AdminPassword, form.AdminEmail)
		if err != nil {
			sendSSE(w, failedInstallMessage(err))
			return
		}
		err = setupContent(dbConn)
		if err != nil {
			sendSSE(w, failedInstallMessage(err))
			return
		}
		// TODO setup default config
		sendSSE(w, "기본설정 정보 입력 완료")

		// TODO creat boards
		sendSSE(w, "게시판 테이블 생성 완료")

		// TODO create data path
		sendSSE(w, "데이터 경로 생성 완료")

		sendSSE(w, fmt.Sprintf("[success] 축하합니다. %s 설치가 완료되었습니다.", version.Version))
	}
}

func setupConfig(dbConn *db.Database, adminId, adminEmail string) error {
	var exists bool
	err := dbConn.Model(&model.Config{}).
		Select("count(*) > 0").
		Where("cf_id = 1").
		Find(&exists).
		Error
	if err != nil {
		return err
	}

	if !exists {
		adminConfig := defaultConfig
		adminConfig.CfAdmin = adminId
		adminConfig.CfAdminEmail = adminEmail

		err = dbConn.Create(&adminConfig).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func setupAdminMember(dbConn *db.Database, adminId string, adminName string, adminPassword string, adminEmail string) error {
	var adminMember model.Member
	err := dbConn.
		Where("mb_id = ?", adminId).
		First(&adminMember).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAdminMember := defaultMember
			newAdminMember.MbID = adminId
			newAdminMember.MbPassword = lib.CreateHash(adminPassword)
			newAdminMember.MbName = adminName
			newAdminMember.MbNick = adminName
			newAdminMember.MbEmail = adminEmail
			return dbConn.Create(&newAdminMember).Error
		}
		return err
	}

	adminMember.MbPassword = lib.CreateHash(adminPassword)
	adminMember.MbName = adminName
	adminMember.MbEmail = adminEmail
	return dbConn.Save(&adminMember).Error
}

func setupContent(dbConn *db.Database) error {
	return dbConn.Clauses(clause.OnConflict{DoNothing: true}).Create(defaultContents).Error
}

func failedInstallMessage(err error) string {
	return fmt.Sprintf("[error] 설치가 실패했습니다. %v", err)
}

func setSSEHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

func sendSSE(w http.ResponseWriter, message string) {
	_, _ = fmt.Fprintf(w, "data: %s\n\n", message)
	w.(http.Flusher).Flush()
}
