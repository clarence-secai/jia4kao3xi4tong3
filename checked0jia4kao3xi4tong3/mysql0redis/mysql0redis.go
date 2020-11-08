package mysql0redis

import (
	"database/sql"
	"fmt"
	"go8jia4kao3xi4tong3/checked0jia4kao3xi4tong3/utils"
)


/*将全员考试成绩单写入MySQL数据库*/
func WriteScore2Mysql(scoreMap map[string]int) {

	db, err := sql.Open("mysql", "root:413188ok@tcp(localhost:3306)/driving_exam")
	utils.HandlerError(err, "mysql0redis包52行")
	defer db.Close()
	for name, score := range scoreMap { //todo:此处可以拼接sql语句进行批量插入
		_, err := db.Exec("insert into score(name,score) values(?,?);", name, score)
		utils.HandlerError(err, "mysql0redis包57行")
		fmt.Println("插入成功！")
	}
	fmt.Println("成绩录入完毕！")
}

func QueryFromMysql(name string) int {
	fmt.Println("QueryScoreFromMysql...")
	db, err := sql.Open("mysql", "root:413188ok@tcp(localhost:3306)/driving_exam")
	utils.HandlerError(err,"mysql0redis包44行")
	defer db.Close()

	sql := "select * from score where name=?"
	row := db.QueryRow(sql,name)
	var score int
	row.Scan(&score)

	return score
}



