package db

import "strings"

func IsSupportedEngines(engine string) bool {
	switch strings.ToLower(engine) {
	case "sqlite", "postgresql", "mysql":
		return true
	}
	return false
}
