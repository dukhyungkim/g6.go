package model

import "os"

var Prefix = ""

const WriteTablePrefix = "write_"

func init() {
	Prefix = os.Getenv("DB_TABLE_PREFIX")
}
