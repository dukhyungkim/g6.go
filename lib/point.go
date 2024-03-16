package lib

import (
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/util"
	"log"
	"time"
)

func InsertPoint(dbConn *db.Database, request util.Request, mbId string, point int, content, relTable, relId, relAction string, expire int) {
	cfg := request.State.Config

	if cfg.CfUsePoint == 0 {
		return
	}
	if point == 0 {
		return
	}
	if mbId == "" {
		return
	}

	if relTable != "" || relAction != "" {
		clauses := model.Point{
			MbID:        mbId,
			PoRelTable:  relTable,
			PoRelID:     relId,
			PoRelAction: relAction,
		}
		var count int64
		err := dbConn.Where(&clauses).Count(&count).Error
		if err != nil {
			log.Println(err)
			return
		}
		if count == 1 {
			return
		}
	}

	now := time.Now()
	poExpireDate := time.Date(9999, 12, 31, 0, 0, 0, 0, time.Local)
	if cfg.CfPointTerm > 0 {
		expireDays := cfg.CfPointTerm
		if expire > 0 {
			expireDays = expire
		}
		afterDatetime := time.Duration(expireDays-1) * 24 * time.Hour
		poExpireDate = now.Add(afterDatetime)
	}

	mbPoint := getPointSum(request, mbId)
	poExpired := 0
	if point < 0 {
		poExpired = 1
		poExpireDate = now
	}
	poMbPoint := mbPoint + point

	newPoint := model.Point{
		MbID:         mbId,
		PoDatetime:   now,
		PoContent:    content,
		PoPoint:      point,
		PoUsePoint:   0,
		PoExpired:    poExpired,
		PoExpireDate: poExpireDate,
		PoMbPoint:    poMbPoint,
		PoRelTable:   relTable,
		PoRelID:      relId,
		PoRelAction:  relAction,
	}
	dbConn.Create(&newPoint)

	dbConn.Where("mb_id = ?", mbId).Update("mb_point", poMbPoint)
}

func getPointSum(request util.Request, mbId string) int {
	// TODO
	return 0
}
