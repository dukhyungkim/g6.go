package install

import (
	"fmt"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/util"
	"github.com/dukhyungkim/gonuboard/version"
	"github.com/go-chi/chi/v5"
	"github.com/nikolalohinski/gonja/v2/exec"
	"io"
	"net/http"
	"os"
	"strconv"
)

func DefaultRouter(r chi.Router) {
	r.Get("/", indexHandler())
	r.Post("/", installDatabase())
	r.Get("/license", licenseHandler())
	r.Post("/form", formHandler())
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
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		engine := r.FormValue("db_engine")
		if !db.IsSupportedEngines(engine) {
			http.Error(w, "지원 가능한 데이터베이스 엔진을 선택해주세요.", http.StatusBadRequest)
			return
		}

		dbPort, err := strconv.Atoi(r.FormValue("db_port"))
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
			util.SetKeyToEnv(util.EnvPath, "DB_ENGINE", engine),
			util.SetKeyToEnv(util.EnvPath, "DB_HOST", r.FormValue("db_host")),
			util.SetKeyToEnv(util.EnvPath, "DB_PORT", dbPort),
			util.SetKeyToEnv(util.EnvPath, "DB_USER", r.FormValue("db_user")),
			util.SetKeyToEnv(util.EnvPath, "DB_PASSWORD", r.FormValue("db_password")),
			util.SetKeyToEnv(util.EnvPath, "DB_NAME", r.FormValue("db_name")),
			util.SetKeyToEnv(util.EnvPath, "DB_TABLE_PREFIX", r.FormValue("db_table_prefix")),
			util.SetKeyToEnv(util.EnvPath, "SESSION_SECRET_KEY", sessionSecretKey),
			util.SetKeyToEnv(util.EnvPath, "COOKIE_DOMAIN", ""),
		} {
			if err = setKey(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		//const templatePath = "install/templates/result.html"
		//util.RenderTemplate(w, templatePath, nil)
	}
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
