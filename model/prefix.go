package model

import "os"

var Prefix string

const WriteTablePrefix = "write_"

func init() {
	Prefix = os.Getenv("DB_TABLE_PREFIX")
	if Prefix == "" {
		Prefix = "g6_"
	}
}
