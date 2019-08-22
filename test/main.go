package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strconv"
)

var (
	db *gorm.DB
	err error
)
func init() {
	user := "root"
	password := "root"
	host := "10.0.0.10"
	port := 3306
	mysql_conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8&loc=Asia%%2FShanghai&parseTime=true",user,password,host,port)
	db,err = gorm.Open("mysql",mysql_conn_str)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(300)
}

//测试MySQL的并发情况
//创建测试表
type TestTab struct {
	Idd int    `gorm:"column:Idd",type:int"`
	Name string
}

func (TestTab)TableName() string  {
	return "test.test"
}

//插入数据
func InsertData(db *gorm.DB)  {
	for i:=0;i<1000000;i++{
		db.Create(&TestTab{i,strconv.Itoa(i)})
	}
}
//进行查询
func QueryData(db *gorm.DB)  {
	var ts []TestTab
	var cnt int
	for {
		db.Find(&ts).Count(&cnt)
		fmt.Println(cnt)
	}

}

func main() {
	if ! db.HasTable(&TestTab{}) {
		db.CreateTable(&TestTab{})
	}/*else {
		db.DropTable(&TestTab{})
		db.Create(&TestTab{})
	}*/
	//向数据库中插入数据
	for i:=0;i<20;i++ {
		go InsertData(db)
	}
	for i:=0;i<40;i++ {
		go QueryData(db)
	}
	select {}
}