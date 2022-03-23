package mysqlDao

type Question struct {
	ID         string
	QuestionID int32
	Type       int32
	Content    string
	Answers    string
	Depends    string
}

// func UpdateQuestions(questions []*Question) (ret bool) {
// 	tx, _ := MysqlDb.Begin()
// 	defer tx.Rollback()
// 	tx.Exec("delete from questions")
// 	for _, question := range questions {
// 		stmt, _ := tx.Prepare("insert into questions (question_id, type, content, answers, depends) VALUES (?,?,?,?,?)")
// 		result, err := stmt.Exec(question.QuestionID, question.Type, question.Content, question.Answers, question.Depends)
// 		if err != nil {
// 			return false
// 		}
// 		if LastInsertId, err := result.LastInsertId(); nil == err {
// 			fmt.Println("LastInsertId:", LastInsertId)
// 		}
// 		if RowsAffected, err := result.RowsAffected(); nil == err {
// 			fmt.Println("RowsAffected:", RowsAffected)
// 		}
// 	}
// 	tx.Commit()
// 	return true
// }

func GetQuestions() (questions []*Question, err error) {
	rows, _ := MysqlDb.Query("select * from questions")
	for rows.Next() {
		q := new(Question)
		err = rows.Scan(&q.ID, &q.QuestionID, &q.Type, &q.Content, &q.Answers, &q.Depends)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return
}
