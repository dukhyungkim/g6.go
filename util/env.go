package util

import (
	"fmt"
	"os"
	"strings"
)

const EnvPath = ".env"

func SetKeyToEnv(filePath, key string, value any) func() error {
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
					switch value.(type) {
					case int, int32, int64:
						lines[i] = fmt.Sprintf(`%s=%v`, key, value)
					case string:
						lines[i] = fmt.Sprintf(`%s="%v"`, key, value)
					default:
						lines[i] = fmt.Sprintf(`%s="%v"`, key, value)
					}
					break
				}
			}
			fileContent = strings.Join(lines, "\n")
		} else {
			fileContent += fmt.Sprintf("\n%s=%v", key, value)
		}

		err = os.WriteFile(EnvPath, []byte(fileContent), os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
}
