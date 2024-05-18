package version

import (
	"log"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	Version        = "Go누보드6.0.0"
	RuntimeVersion = runtime.Version()
	RouterVersion  = ""
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build info")
		return
	}

	for _, dep := range bi.Deps {
		if !strings.Contains(dep.Path, "gin-gonic") {
			continue
		}
		RouterVersion = dep.Version
		break
	}
}
