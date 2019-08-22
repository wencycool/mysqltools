package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"mysqltools/measure/vm"
	"time"
)

var db *gorm.DB
func init() {
	var err error
	db,err = gorm.Open("sqlite3","test.db?cache=shared")
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(10)
}
func main() {
	//测试将数据写入到sqlite3中
	//测试将cpu使用率,内存使用情况写入到sqlite3中
	go vm.SaveCPUInfo(db)
	go vm.SaveDiskUsageStat(db,time.Second)
	go vm.SaveCPUStat(db,time.Second)
	go vm.SaveVirtualMemoryStat(db,time.Second)
	go vm.SaveSwapMemoryStat(db,time.Second)
	go vm.SaveProcStat(db,time.Second)
	go func() {
		var i = 1
		for {
			fmt.Println(i)
			i++
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(time.Second* time.Duration(60))
}

func xx()  {
	var i = 0
	for {
		i++
		i--
	}
}