package main

import (
	"os"
	"regexp"
	"strings"
)

func processFile(path string, info os.FileInfo) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	replacedContent := replaceURLFor(string(content))

	err = os.WriteFile(path, []byte(replacedContent), info.Mode().Perm())
	if err != nil {
		return err
	}
	return nil
}

var re = regexp.MustCompile(`{{\s*url_for\s*\(\s*'([^']+)'(?:,\s*path='([^']*)')?\s*\)\s*}}`)

func replaceURLFor(input string) string {
	var replace string
	for {
		replace = transformURLFor(input)
		if replace == "" {
			return input
		}
		input = replace
	}
}

func transformURLFor(input string) string {
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return ""
	}

	if len(matches) == 3 && matches[2] != "" {
		path := matches[1] + "/" + strings.Trim(matches[2], "/")
		return strings.ReplaceAll(input, matches[0], path)
	}

	return strings.ReplaceAll(input, matches[0], matches[1])
}
