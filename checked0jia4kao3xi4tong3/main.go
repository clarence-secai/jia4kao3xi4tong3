
package main

import (
	"fmt"
	"go8jia4kao3xi4tong3/checked0jia4kao3xi4tong3/mysql0redis"
	"go8jia4kao3xi4tong3/checked0jia4kao3xi4tong3/utils"
	"sync"
	"time"
)

/*项目流程简述
随机生成考生姓名，将名字写入管道chNames；

从chNames遍历取出考生姓名，
每个考生都开一个协程运行考试函数，而考试函数中的chLanes管道容量，限制了实际运行的考试函数
的协程数量;考试函数中产生考试随机成绩、建立name-score的map、考生姓名切片examers、将成绩不合格者写进违纪管道chFouls；
同时开协程从违纪管道chFouls中读取违纪者姓名并通报；

所有考试者考完后，【实际生产中，不用等所有考生考完，可以全流程的所有环节均并发运行】
将name-score的map写入MySQL数据库；
同时开协程读取MySql数据库内的name-score成绩
*/

/*体会：
1、在考试函数内加一个额外管道，利用管道容量来控制考试函数的协程的并发数量
2、产生随机数的函数被多个协程使用，故产生随机数的函数内需加互斥锁
3、上述两者有异曲同工之妙：互斥锁相当于容量为1的控制产生随机数的函数的协程并发数量的管道！！！ :)
*/

var (
	chNames = make(chan string, 100)
	examers = make([]string, 0)
	//考试函数协程的并发数量控制，只有5条考试车道，同时正在考试的只能有5个
	chLanes = make(chan int, 5)
	//违纪者
	chFouls = make(chan string, 100) //并非一定得设置为100，大于1的都行
	//考试成绩
	scoreMap = make(map[string]int)

	wg sync.WaitGroup
)

func main() {
	for i := 0; i < 20; i++ {
		chNames <- utils.GetRandomName() //todo:此处也可开协程，毕竟该操作里的GetRandomInt有互斥锁
	}
	close(chNames)

	//巡考
	go Patrol() //todo:因不知道何时关闭chFouls，故采用的for select

	//考生并发考试
	for name := range chNames {
		wg.Add(1)
		go func(name string) {
			TakeExam(name)
			wg.Done()
		}(name)
	}

	wg.Wait()
	fmt.Println("考试完毕！")//todo 实际场景中也可以全流程各个环节流水线式并发，而不必等考试结束才运行下面的流程
	for k,v := range scoreMap{
		fmt.Println(k,"=",v)
	}

	//录入成绩
	mysql0redis.WriteScore2Mysql(scoreMap)

	//发放成绩
	for _, name := range examers {
		go func(name string) { //todo:数据库自带锁，本身就是可以被多协程同时读取数据！！！
			mysql0redis.QueryFromMysql(name)
		}(name)
	}
	fmt.Println("END")
}

//todo 巡考
func Patrol() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		//fmt.Println("正在巡考...")
		select {
		case name := <-chFouls:  //违纪人员会被写进管道chFouls，此处是在取出并通报
			fmt.Println(name, "考试违纪!!!!! ")
		default:
			//fmt.Println("考场秩序良好")
		}
		<-ticker.C //每隔一秒巡考一次
	}
}

//todo 考试  每个考生开一个考试协程
func TakeExam(name string) {
	chLanes <- 123 //todo:chLanes控制同时考试的人数最多为5条跑道的，即管道控制了瞬间同时运行的并发协程数量
	fmt.Println(name, "正在考试...")

	//记录一下参与考试的考生姓名
	examers = append(examers, name) //todo:由于TakeExam被开20个协程，会导致这里被竞争操作吗？

	//生成考试成绩
	score := utils.GetRandomInt(0, 100)
	scoreMap[name] = score  //构建考生--成绩单的map//todo:会协程竞争操作scoreMap!!!

	//违纪考生写进违纪管道chFouls
	if score < 10 {
		score = 0
		chFouls <- name  //将违纪考生丢入违纪管道
		//fmt.Println(name, "考试违纪！！！", score)
	}

	//考试持续时间5秒左右
	<-time.After(400 * time.Millisecond)
	//todo:违纪的也必须考完400毫秒，它的123没法中途被从chLanes管道中取出

	<-chLanes //todo:后文对该处函数开启goroutine，这里取出的123不一定是38行放入的一一对应的123吧？
}
