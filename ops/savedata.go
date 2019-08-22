package ops

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"mysqltools/measure/vm"
	"os"
	"path/filepath"
	"time"
	"errors"
	"mysqltools/measure"
	_ "github.com/mattn/go-sqlite3"
)
type SaveDataParm struct {
	Sqlite3_dbfile string //数据文件所在地
}
func SaveDatas(db_mysql *gorm.DB,numsec,count int,path string, p *SaveDataParm) error {
	db_mysql.LogMode(false)
	//获取mysql数据库中当前正在使用的端口号信息
	var varName,varValue string
	if err := db_mysql.Raw("SHOW VARIABLES WHERE Variable_name = 'port'").Row().Scan(&varName,&varValue);err != nil {
		return err
	}else if len(varValue) == 0{
		varValue = "3306_not_find"
	}

	//检查path路径是否存在
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	if finfo,err := os.Stat(path);err != nil {
		return err
	}else if !finfo.IsDir() {
		return errors.New(fmt.Sprintf("Path:%s is not Path",path))
	}
	/*
	path = filepath.Join(path,varValue)
	if finfo,err := os.Stat(path);err != nil {
		//路径不存在，进行创建
		if err := os.Mkdir(path,os.ModePerm);err != nil {
			return err
		}
	}else if !finfo.IsDir(){
		//不是路径
		return errors.New(fmt.Sprintf("path:%s 不是路径！"))
	}
	 */
	abs_path,_ := filepath.Abs(path)
	log.Println("开始采集数据,数据存放路径为:",abs_path)
	//创建表
	hostname,_ := os.Hostname()
	timeTemplate1 := "2006-01-02-15:04:05"
	timeTemplate2 := "2006-01-02"
	targetdbname := filepath.Join(path,fmt.Sprintf("sqlite3_mysql_%s_%s_%s.db",hostname, varValue,time.Now().Format(timeTemplate1)))
	diaglogname  := filepath.Join(path,fmt.Sprintf("sqlite3_mysql_%s_%s_diag_%s.log",hostname,varValue,time.Now().Format(timeTemplate2)))
	slowlogname  := filepath.Join(path,fmt.Sprintf("sqlite3_mysql_%s_%s_slow_%s.log",hostname,varValue,time.Now().Format(timeTemplate2)))
	db_sqlite3,err := gorm.Open("sqlite3",fmt.Sprintf("%s?cache=shared&_busy_timeout=300000",targetdbname)) //basy_timeout设置为30秒
	if err != nil {
		return err
	}
	p.Sqlite3_dbfile = targetdbname
	log.Println("指标信息存放在sqlite文件中:",targetdbname)
	//测试将cpu使用率,内存使用情况写入到sqlite3中
	//设置捕获时间间隔，最小1秒
	t_delta := time.Duration(5)
	go vm.SaveCPUInfo(db_sqlite3)
	go vm.SaveCPUStat(db_sqlite3,time.Second * t_delta)
	go vm.SaveDiskUsageStat(db_sqlite3,time.Second * t_delta)
	go vm.SaveVirtualMemoryStat(db_sqlite3,time.Second * t_delta)
	go vm.SaveSwapMemoryStat(db_sqlite3,time.Second * t_delta)
	go vm.SaveProcStat(db_sqlite3,time.Second * t_delta)
	var tablist = []interface{}{measure.ProcesslistTGT{},measure.InnodbTrxTGT{},measure.InnodbLocksTGT{},
		measure.InnodbLockWaitsTGT{},measure.InnodbStatusTGT{},measure.GlobalVariablesTGT{},
		measure.GlobalStatusTGT{},measure.SlaveStatusTGT{},measure.TablesTGT{},measure.IndexesTGT{},
	}
	//创建监控表
	measure.CreateTables(tablist,db_sqlite3)
	log.Printf("一共采集%d轮数据",count)
	for i:=0;i<count;i++ {
		//插入数据
		log.Printf("开始第:%d轮的数据采集...\n",i+1)
		measure.InsertData(tablist,db_mysql,db_sqlite3,time.Now())
		measure.InsertExplainData(db_mysql,db_sqlite3)
		time.Sleep(time.Duration(numsec) * time.Second)
	}

	//保存diag日志文件和slowlog文件，最大20MB
	if err := measure.SaveDiag(db_mysql,diaglogname, 20<<20);err != nil {
		log.Println(err)
	}else{
		log.Printf("保存diag日志文件路径为:%s\n",diaglogname)
	}
	if err := measure.SaveSlowlog(db_mysql,slowlogname,20<<20);err != nil {
		log.Println(err)
	}else {
		log.Printf("保存慢日志路径为:%s\n",slowlogname)
	}
	return nil
}
