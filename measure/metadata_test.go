package measure

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

var db_mysql *gorm.DB
var db_sqlite3 *gorm.DB

func init()  {
	var err error
	db_mysql,err = gorm.Open("mysql","root:root@tcp(192.168.40.200:3306)/mysql?charset=utf8&loc=Asia%2FShanghai&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	db_mysql.DB().SetMaxOpenConns(10)
	db_sqlite3,err = gorm.Open("sqlite3","test.db")
	if err != nil {
		log.Fatal(err)
	}
	db_sqlite3.DB().SetMaxOpenConns(10)
}

func TestGetSlaveStatus(t *testing.T) {
	//测试show engine innodb
	var innodb_status []InnodbStatus
	db_mysql.Raw("show engine innodb status").Scan(&innodb_status)
	fmt.Println(innodb_status[0].Status)
}

func TestCreateTables(t *testing.T) {
	var tablist = []interface{}{&ProcesslistTGT{},&InnodbTrxTGT{},&InnodbLocksTGT{},&InnodbLockWaitsTGT{},&InnodbStatusTGT{},
		&GlobalVariablesTGT{},&GlobalStatusTGT{},&SlaveStatusTGT{},
	}
	CreateTables(tablist,db_sqlite3)
}

func TestInsertData(t *testing.T) {
	InsertData([]interface{}{Processlist{},},db_mysql,db_sqlite3)

}
