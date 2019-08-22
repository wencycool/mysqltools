package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	var (
		db *gorm.DB
		err error
	)
	db,err = gorm.Open("mysql","root:root@tcp(10.0.0.10:3306)/mysql")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = db.Exec("use test")
	tx := db.Begin()
	var (
		id int
		name string
	)
	row := tx.Raw("select * from test0 limit 10").Row()
	row.Scan(&id,&name)
	fmt.Println(id,name)

	var ts []T
	tx.Raw("select id,name from test0 limit 10").Scan(&ts)
	fmt.Println(ts)
	tx.Commit()
}
type T struct {
	id int
	name string
}