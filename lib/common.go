package lib

import (
	"fmt"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"net/http"
)

func CreateDynamicWriteTable(dbConn *db.Database, tableName string) error {
	writeTable := model.NewWriteTable(tableName)
	err := dbConn.Table(writeTable.TableName()).AutoMigrate(writeTable)
	if err != nil {
		return err
	}

	const numReplyIndex = "idx_wr_num_reply"
	newNumReplyIndex := fmt.Sprintf("%s_%s", numReplyIndex, tableName)
	err = dbConn.Table(writeTable.TableName()).Migrator().RenameIndex(writeTable, numReplyIndex, newNumReplyIndex)
	if err != nil {
		return err
	}

	const commentIndex = "idx_wr_is_comment"
	newCommentIndex := fmt.Sprintf("%s_%s", commentIndex, tableName)
	err = dbConn.Table(writeTable.TableName()).Migrator().RenameIndex(writeTable, commentIndex, newCommentIndex)
	if err != nil {
		return err
	}

	return nil
}

func GetClientIp(r *http.Request) string {
	if clientIp := r.Header.Get("X-FORWARDED-FOR"); clientIp != "" {
		return clientIp
	}
	return r.RemoteAddr
}
