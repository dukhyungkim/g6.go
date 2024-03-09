package lib

import (
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/util"
)

func InsertPoint(dbConn *db.Database, request util.Request, member *model.Member, content string, relTable string, relAction string) {
	cfg := request.State.Config

	if cfg.CfUsePoint == 0 {
		return
	}

	if cfg.CfLoginPoint == 0 {
		return
	}

	if member == nil {
		return
	}

	if relTable != "" || relAction != "" {
		// TODO
		dbConn.Where("mb_id")
	}
}
