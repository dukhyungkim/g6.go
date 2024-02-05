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
	"strings"
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

const envPath = ".env"

func installDatabase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := copyFile("example.env", envPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		engine := r.FormValue("db_engine")
		for _, setKey := range []func() error{
			setKeyToEnv(envPath, "DB_ENGINE", engine),
			setKeyToEnv(envPath, "DB_HOST", r.FormValue("db_host")),
			setKeyToEnv(envPath, "DB_PORT", r.FormValue("db_port")),
			setKeyToEnv(envPath, "DB_USER", r.FormValue("db_user")),
			setKeyToEnv(envPath, "DB_PASSWORD", r.FormValue("db_password")),
			setKeyToEnv(envPath, "DB_NAME", r.FormValue("db_name")),
			setKeyToEnv(envPath, "DB_TABLE_PREFIX", r.FormValue("db_table_prefix")),
			setKeyToEnv(envPath, "SESSION_SECRET_KEY", r.FormValue("session_secret_key")),
			setKeyToEnv(envPath, "COOKIE_DOMAIN", r.FormValue("cookie_domain")),
		} {
			if err = setKey(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		if !db.IsSupportedEngines(engine) {
			http.Error(w, "지원 가능한 데이터베이스 엔진을 선택해주세요.", http.StatusBadRequest)
			return
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

func setKeyToEnv(filePath, key string, value any) func() error {
	return func() error {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		fileContent := string(content)

		if strings.Contains(fileContent, key) {
			lines := strings.Split(fileContent, "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, key) {
					lines[i] = fmt.Sprintf("%s=%v", key, value)
					break
				}
			}
			fileContent = strings.Join(lines, "\n")
		} else {
			fileContent += fmt.Sprintf("\n%s=%v", key, value)
		}

		err = os.WriteFile(envPath, []byte(fileContent), os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
}
