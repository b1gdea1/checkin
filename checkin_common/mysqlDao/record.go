package mysqlDao

import "log"

type Record struct {
	ID     int
	SID    string
	Answer string
	Date   string
}

func UpdateRecord(record Record) bool {
	stmt, _ := MysqlDb.Prepare("insert into records (sid, answer, date) VALUES (?,?,?)")
	result, err := stmt.Exec(record.SID, record.Answer, record.Date)
	if err != nil {
		return false
	}
	if n, err := result.RowsAffected(); err == nil {
		log.Printf("inserted row:%d", n)
	}
	if n, err := result.LastInsertId(); err == nil {
		log.Printf("last insert row:%d", n)
	}
	return true
}
