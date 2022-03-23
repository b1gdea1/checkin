package mysqlDao

import (
	"fmt"
	"log"
)

type Admin struct {
	ID         uint
	EmployeeID string
	Password   string
	Telephone  string
}

func GetAdministratorID(EID string, password string) Admin {
	var id Admin
	rows, _ := MysqlDb.Query("select id from admins where employee_id=? and password=?", EID, password)
	if rows.Next() {
		err := rows.Scan(&(id.ID))
		if err != nil {
			log.Println(err)
			return Admin{}
		}
		id.EmployeeID = EID
	}
	return id
}
func AdminExist(EID string) bool {
	rows, err := MysqlDb.Query("select count(*) from admins where employee_id=?", EID)
	if err != nil {
		return false
	}
	var n int
	if rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			return false
		}
	}
	return n != 0
}
func InsertAdminSimple(username string, password string) bool {
	stmt, _ := MysqlDb.Prepare(`INSERT INTO admins (employee_id, password) VALUES (?, ?)`)
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
