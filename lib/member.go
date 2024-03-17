package lib

import (
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"log"
)

func getMember(mbId string) *model.Member {
	var member model.Member
	err := db.GetInstance().Where("mb_id = ?", mbId).Take(&member).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	return &member
}
