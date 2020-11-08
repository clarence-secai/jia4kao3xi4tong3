package utils

import (
"fmt"
"math/rand"
"os"
"time"
_ "github.com/go-sql-driver/mysql"
"sync"
)

//处理错误
func HandlerError(err error, when string) {
	if err != nil {
		fmt.Println(when, err)
		os.Exit(1)
	}
}

//考试成绩
type ExamScore struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Score int    `db:"score"`
}

var (
	//姓氏
	familyNames = []string{"赵", "钱", "孙", "李", "周", "吴", "郑", "王",
		"冯", "陈", "楚", "卫", "蒋", "沈", "韩", "杨", "张", "欧阳", "东门",
		"西门", "上官", "诸葛", "司徒", "司空", "夏侯"}
	//辈分（宗的永其光...）
	middleNamesMap = map[string][]string{}
	//名字
	lastNames = []string{"春", "夏", "秋", "冬", "风", "霜", "雨", "雪", "木",
		"禾", "米", "竹", "山", "石", "田", "土", "福", "禄", "寿", "喜", "文",
		"武", "才", "华"}
)

//初始化姓氏和对应的辈分
func init() {
	for _, x := range familyNames {
		if x != "欧阳" {
			middleNamesMap[x] = []string{"德", "惟", "守", "世", "令", "子", "伯", "师", "希", "与", "孟", "由", "宜", "顺", "元", "允", "宗", "仲", "士", "不", "善", "汝", "崇", "必", "良", "友", "季", "同"}
		} else {
			middleNamesMap[x] = []string{"宗", "的", "永", "其", "光"}
		}
	}
}
//获得随机姓名
func GetRandomName() (name string) {
	familyName := familyNames[GetRandomInt(0, len(familyNames)-1)]
	middleName := middleNamesMap[familyName][GetRandomInt(0,
		len(middleNamesMap[familyName])-1)]
	lastName := lastNames[GetRandomInt(0, len(lastNames)-1)]
	return familyName + middleName + lastName
}

var (
	//随机数互斥锁（确保GetRandomInt不能被并发访问）
	randomMutex sync.Mutex
)
//获取[start,end]范围内的随机数
func GetRandomInt(start, end int) int {
	randomMutex.Lock() //todo:因多个协程都需要使用该随机数生成函数，故加互斥锁很重要
	<-time.After(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := start + r.Intn(end-start+1)
	randomMutex.Unlock()
	return n
}




