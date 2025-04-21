package chat

import (
	"database/sql"
	"log"
	"meetingagent/dao"
)

func SaveContext(role string, content string) {
	InsertMessage(dao.SqliteDB, role, content)
}
func InsertMessage(db *sql.DB, role string, content string) {
	_, err := db.Exec("INSERT INTO messages (role, content) VALUES (?, ?)", role, content)
	if err != nil {
		log.Println("插入消息失败", err)
	}
}
