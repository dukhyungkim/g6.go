package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func parseFlag() *string {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s -path <target_folder>\n", os.Args[0])
		flag.PrintDefaults()
	}

	targetPath := flag.String("path", "", "Path to the target folder")
	flag.Parse()

	if *targetPath == "" {
		fmt.Println("Error: -path is required.")
		flag.Usage()
		os.Exit(1)
	}
	return targetPath
}

func main() {
	targetPath := parseFlag()

	err := filepath.Walk(*targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		err = processFile(path, info)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
