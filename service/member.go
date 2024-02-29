package service

import (
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
)

type MemberService struct {
	dbConn *db.Database
}

func NewMemberService(dbConn *db.Database) *MemberService {
	return &MemberService{dbConn: dbConn}
}

func (m *MemberService) CreateById(id string) (*model.Member, error) {
	member := model.Member{}
	if err := m.dbConn.Where("mb_id = ?", id).Scan(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}
