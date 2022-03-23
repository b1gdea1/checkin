package mysqlDao

import "log"

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
