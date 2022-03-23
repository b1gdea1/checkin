package mysqlDao

import (
	"fmt"
	"log"
)

type Student struct {
	ID        uint
	Username  string
	Password  string
	Telephone string
}

func GetStudentID(username string, password string) Student {
	var stu Student
	row := MysqlDb.QueryRow("select id from students where username=? and password =?", username, password)
	err := row.Scan(&(stu.ID))
	if err != nil {
		log.Println(err)
		return Student{ID: 0}
	}
	return stu
}

func InsertStudentFull(username string, password string, telephone string) bool {
	stmt, _ := MysqlDb.Prepare(`INSERT INTO students (username, password,telephone) VALUES (?, ?, ?)`)
	ret, err := stmt.Exec(username, password, telephone)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return false
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		fmt.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:", RowsAffected)
	}
	return true
}

func InsertStudentSimple(username string, password string) bool {
	stmt, _ := MysqlDb.Prepare(`INSERT INTO students (username, password) VALUES (?, ?)`)
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return false
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		fmt.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:", RowsAffected)
	}
	return true
	
}
