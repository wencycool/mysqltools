package measure

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"os"
	"time"
)

//将mysql中表信息全部放入到sqlite3中

//传入需要创建表的结构体指针
func CreateTables(tables []interface{},db *gorm.DB)  {
	for _,eachTable := range tables {
		if !db.HasTable(eachTable) {
			db.CreateTable(eachTable)
		}
	}
}

func Md5Convert(str string) (md5code string)  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//InsertDatas,当前版本不做SQL合并处理,innodbTrxTGT,ProcesslistTGT表存在sqlid,t为插入数据的时间
func InsertData(targetTables []interface{},db_src,db_tgt *gorm.DB,t time.Time)  {
	t_now := t
	for _,targetTable := range targetTables {
		tx := db_tgt.Begin()
		switch targetTable.(type) {
		case ProcesslistTGT:
			var src_records []Processlist
			db_src.Find(&src_records)
			for _, eachRecords := range src_records {
				//针对每一条SQL都进行MD5处理
				var sqlid string
				if eachRecords.COMMAND == "Query" {
					sqlid = Md5Convert(eachRecords.INFO)
				}else {
					sqlid = ""
				}
				tx.Create(&ProcesslistTGT{eachRecords, sqlid,t_now})
			}
		case InnodbTrxTGT:
			var src_records []InnodbTrx
			db_src.Find(&src_records)
			for _, eachRecords := range src_records {
				var sqlid string
				if len(eachRecords.Trx_query) > 0 {
					sqlid = Md5Convert(eachRecords.Trx_query)
				}else {
					sqlid = ""
				}
				tx.Create(&InnodbTrxTGT{eachRecords, sqlid,t_now})
			}
		case InnodbLocksTGT:
			var src_records []InnodbLocks
			db_src.Find(&src_records)
			for _, eachRecords := range src_records {
				tx.Create(&InnodbLocksTGT{eachRecords, t_now})
			}
		case InnodbLockWaitsTGT:
			var src_records []InnodbLockWaits
			db_src.Find(&src_records)
			for _, eachRecords := range src_records {
				tx.Create(&InnodbLockWaitsTGT{eachRecords, t_now})
			}
		case TablesTGT:
			//检查innodb_stats_on_metadata参数是否设置为OFF,并且该表只存储一次记录
			var src_var_records []GlobalVariables
			db_src.Raw("show variables like 'innodb_stats_on_metadata'").Scan(&src_var_records)
			var cnt int
			db_tgt.Table(TablesTGT{}.TableName()).Count(&cnt)
			if len(src_var_records) != 1 || src_var_records[0].Value != "OFF" ||cnt > 0 {continue}
			var src_records []Tables
			db_src.Find(&src_records)
			for _,eachRecords := range src_records {
				tx.Create(&TablesTGT{eachRecords, t_now})
			}
		case IndexesTGT:
			//检查innodb_stats_on_metadata参数是否设置为OFF，并且该表只存储一次记录
			var src_var_records []GlobalVariables
			db_src.Raw("show variables like 'innodb_stats_on_metadata'").Scan(&src_var_records)
			var cnt int
			db_tgt.Table(IndexesTGT{}.TableName()).Count(&cnt)
			if len(src_var_records) != 1 || src_var_records[0].Value != "OFF" ||cnt > 0 {continue}
			var src_records []Indexes
			db_src.Find(&src_records)
			for _,eachRecords := range src_records {
				tx.Create(&IndexesTGT{eachRecords, t_now})
			}
		case InnodbStatusTGT:
			var src_records []InnodbStatus
			db_src.Raw("show engine innodb status").Scan(&src_records)
			for _, eachRecords := range src_records {
				tx.Create(&InnodbStatusTGT{eachRecords, t_now})
			}
		case GlobalStatusTGT:
			var src_records []GlobalStatus
			db_src.Raw("show global status").Scan(&src_records)
			for _, eachRecords := range src_records {
				tx.Create(&GlobalStatusTGT{eachRecords, t_now})
			}
		case GlobalVariablesTGT:
			//只搜集一次
			var cnt int
			db_tgt.Table(GlobalVariablesTGT{}.TableName()).Count(&cnt)
			if cnt > 0 {continue}
			var src_records []GlobalVariables
			db_src.Raw("show variables").Scan(&src_records)
			for _, eachRecords := range src_records {
				tx.Create(&GlobalVariablesTGT{eachRecords, t_now})
			}
		case SlaveStatusTGT:
			var src_records []SlaveStatus
			db_src.Raw("show slave status").Scan(&src_records)
			for _, eachRecords := range src_records {
				tx.Create(&SlaveStatusTGT{eachRecords, t_now})
			}
		}
		tx.Commit()
	}
}

//获取执行计划信息,针对当前正在执行的SQL获取explain信息并存放到目标表中
func InsertExplainData(db_src,db_tgt *gorm.DB)  {
	//为防止重复查询processlist，针对已经插入到sqlite3数据库中的最后一批数据查找执行计划信息
	if !db_tgt.HasTable(&SQLExplainTGT{}) {
		db_tgt.CreateTable(&SQLExplainTGT{})
	}
	//对当前目标数据库中的processlist表
	var processlistTGTs []ProcesslistTGT
	sql := fmt.Sprintf("select * from %s where insert_time=(select max(insert_time) from %s) and COMMAND='Query'",ProcesslistTGT{}.TableName(),ProcesslistTGT{}.TableName())
	db_tgt.Raw(sql).Scan(&processlistTGTs)
	for _,eachRow := range processlistTGTs {
		sqlid := eachRow.SqlId
		sql_stmt := eachRow.INFO
		insert_time := eachRow.InsertTime
		if db:= eachRow.DB;len(db) > 0 {
			//db_src.Raw(fmt.Sprintf("use %s", eachRow.DB))
			db_src.Exec(fmt.Sprintf("use %s", eachRow.DB))  //设置当前连接所用数据库
		}
		if len(sqlid) == 0 {
			continue
		}
		var sqlexplains []SQLExplain

		db_src.Raw(fmt.Sprintf("explain %s",eachRow.INFO)).Scan(&sqlexplains)
		for _,eachRow := range sqlexplains {
			//先检查sqlid是否已经存在
			var cnt int = 0
			db_tgt.Model(&SQLExplainTGT{}).Where("sql_id=?",sqlid).Count(&cnt)
			if cnt > 0 {
				continue
			}
			db_tgt.Create(&SQLExplainTGT{eachRow,sqlid,sql_stmt,insert_time})
		}

	}
}


func SaveSlowlog(db *gorm.DB,logname string,logsize int64) error  {
	if slowlogpath,err := findSlowlogPath(db);err != nil {
		return err
	}else {
		if err := saveLog(slowlogpath,logname,logsize);err != nil {
			return err
		}
	}
	return nil
}

func SaveDiag(db *gorm.DB,logname string,logsize int64) error  {
	if diaglogpath,err := findDiaglogPath(db);err != nil {
		return err
	}else {
		if err := saveLog(diaglogpath,logname,logsize);err != nil {
			return err
		}
	}
	return nil
}




func findSlowlogPath(db *gorm.DB) (string,error) {
	var (
		Variable_name 	string
		Value 			string
	)
	db.Raw("show variables like 'slow_query_log'").Row().Scan(&Variable_name,&Value)
	if Value == "ON" {
		if err := db.Raw("show variables like 'slow_query_log_file'").Row().Scan(&Variable_name,&Value);err != nil {
			return "",err
		}else {
			return Value,nil
		}

	}else {
		return "",errors.New("slow_query_log is OFF")
	}

}


func findDiaglogPath(db *gorm.DB) (string,error) {
	var (
		Variable_name 	string
		Value 			string
	)
	if err := db.Raw("show variables like 'log_error'").Row().Scan(&Variable_name,&Value);err != nil {
		return "",err
	}else {
		return Value,nil
	}

}


func saveLog(src_logpath,target_logpath string,size int64) error {
	var (
		src_file *os.File
		target_file *os.File
		err error
		seekpos int64 = 0
	)
	if src_file,err = os.Open(src_logpath);err != nil {
		return err
	}
	stat,_ := src_file.Stat()
	filesize := stat.Size()
	if size < filesize {
		seekpos = filesize - size
	}
	src_file.Seek(seekpos,0)
	reader := bufio.NewReader(src_file)
	target_file,err = os.OpenFile(target_logpath,os.O_CREATE|os.O_TRUNC|os.O_RDWR,os.ModePerm)
	if _,err = reader.WriteTo(target_file);err != nil {
		return err
	}
	target_file.Close()
	return nil
}
