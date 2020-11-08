# jia4kao3xi4tong3
/*项目流程简述
1、随机生成考生姓名，将名字写入管道chNames；

2、从chNames遍历取出考生姓名，
3、每个考生都开一个协程运行考试函数，而考试函数中的chLanes管道容量，限制了实际运行的考试函数
的协程数量;考试函数中产生考试随机成绩、建立name-score的map、考生姓名切片examers、将成绩不合格者写进违纪管道chFouls；
4、同时开协程从违纪管道chFouls中读取违纪者姓名并通报；

5、所有考试者考完后，【实际生产中，不用等所有考生考完，可以全流程的所有环节均并发运行】
将name-score的map写入MySQL数据库；
*/

/*体会：
1、在考试函数内加一个额外管道，利用管道容量来控制考试函数的协程的并发数量
2、产生随机数的函数被多个协程使用，故产生随机数的函数内需加互斥锁
3、上述两者有异曲同工之妙：互斥锁相当于容量为1的控制产生随机数的函数的协程并发数量的管道！！！ :)
*/
