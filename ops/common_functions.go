package ops

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

const ParmFile  = "/tmp/.sqlite3_mysql_parm.txt"
const AESKey    = "1qaz@WSX1qaz@WSX"  //长度必须为16, 24或者32

//定义加密存储参数文件的结构体
type Parm struct {
	User		string //用户名
	Password	string //密码
	SnapPath 	string //存放快照路径
}


//将结构体参数加密存储到文件中,默认参数存放的地址为:/tmp/.sqlite3_mysql_parm.txt
func PutParmToFile(parm *Parm,filename string) error  {
	f,err := os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_TRUNC,os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(filename,os.ModePerm)
	defer f.Close()
	parm.encrypt(AESKey)
	var bs = new(bytes.Buffer)
	if err := gob.NewEncoder(bs).Encode(&parm);err != nil {
		return err
	}
	f.Write(bs.Bytes())
	return nil

}
//将结构体参数从文件加载到程序中
func LoadParmFromFile(filename string,p *Parm) (err error)  {
	f,err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var bs = new(bytes.Buffer)
	cnt,err := bs.ReadFrom(f)
	if err != nil ||cnt <=0 {
		return errors.New("No Data")
	}
	err = gob.NewDecoder(bs).Decode(p)
	p.decrypt(AESKey)
	if err == io.EOF {
		return nil
	}else if err != nil {
		return err
	}
	return nil
}
func (p *Parm)encrypt(key string)  {
	p.User = aesEncrypt(p.User,key)
	p.Password = aesEncrypt(p.Password,key)
	//p.SnapPath = aesEncrypt(p.SnapPath,key)
	return
}
func (p *Parm)decrypt(key string)  {
	p.User = aesDecrypt(p.User,key)
	p.Password = aesDecrypt(p.Password,key)
	//p.SnapPath = aesDecrypt(p.SnapPath,key)
	return
}

func aesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func aesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = pKCS7UnPadding(orig)
	return string(orig)
}
//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
//去码
func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


//定义打印SQL语句的结构体
type ReportSQL struct {
	Name 	string  //SQL指标解释
	SQL 	string  //SQL语句
}

func ReportCmd()  (reportSQLs []ReportSQL){
	reportSQLs = append(reportSQLs, ReportSQL{"查询CPU使用情况",
	`select round(totalused*100,2) as TotalUsed,round(user*100,2) as User,round(system*100,2) as System,
		round(iowait*100,2) as IOWait,round(idle*100,2) as IDLE,
		round(load1) as Load1,round(load5) as Load5,round(load15) as Load15,datetime(insert_time,'localtime') as Time
		from vm_cpustat order by insert_time asc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查询内存使用情况",
	`select total/1024/1024 as Total_MB,available/1024/1024 as Aval_MB,used/1024/1024 as Used_MB,
    round(used_percent,2) as UsedPercent,
    free/1024/1024 as free_MB,active/1024/1024 as active_MB,inactive/1024/1024 as inactive_MB ,
    datetime(insert_time,'localtime') as Time
from vm_memstat order by insert_time asc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查询交换空间使用情况",
	`select total/1024/1024 as Total_MB,used/1024/1024 as Used_MB,free/1024/1024 as Free_MB,round(used_percent,2) as UsedPercent,
datetime(insert_time,'localtime') as Time
from vm_swapstat  order by insert_time asc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看Top进程的使用情况",
	`select pid,background,round(cpu_percent,2) as CPUPercent,cmd_line,create_time,user,
    read_bytes__s/1024/1024 as read_MBs,write_bytes__s/1024/1024 as write_MBs,
    bytes_sent__s/1024/1024 as NetSend_MBs,bytes_recv__s/1024/1024 as NetRecv_MBs,
    datetime(insert_time,'localtime') as Time
from vm_procstat_top10 order by pid,insert_time asc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看数据库整体状态",
	`select conExec.value,conTotal.value,lockWaits.value,datetime(conExec.insert_time,'localtime') as Time from
(select value ,insert_time  from global_status where variable_name ='Threads_running') as conExec
left join 
(select value ,insert_time  from global_status where variable_name = 'Threads_connected') as conTotal
on conExec.insert_time = conTotal.insert_time
left join 
(select value,insert_time from global_status where variable_name = 'Innodb_row_lock_current_waits') as lockWaits
on conExec.insert_time = lockWaits.insert_time order by conExec.insert_time asc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看快照中执行次数最多的一次的SQL执行情况",
	`select * from processlist where insert_time = (select max(insert_time) from processlist where insert_time = (
    select insert_time from (
        select insert_time,count(*) cnt from processlist where COMMAND='Query' group by insert_time order by cnt limit 1
    ) a 
)) and command='Query' order by TIME desc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看快照中正在执行的而且执行执行时间最长的TOP10语句",
	`select * from processlist where COMMAND='Query' and insert_time=(select max(insert_time) from processlist)
order by time desc limit 10`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看快照中processlist状态信息",
	`select state,count(*) as cnt,datetime(insert_time,'localtime') as Time 
from processlist group by state,insert_time order by insert_time`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看主从复制状态",
	`select Master_Host,Slave_IO_Running,Slave_SQL_Running ,datetime(insert_time,'localtime') as Time
from slave_status order by Master_Host,insert_time asc`,
	})
	reportSQLs = append(reportSQLs, ReportSQL{"查看处于锁等待状态的SQL语句",
	`select c.trx_id as holding_trx_id,c.trx_requested_lock_id as holding_trx_requested_lock_id,
    c.trx_state as holding_trx_state,
    c.trx_started as holding_trx_started,
    c.trx_wait_started as holding_trx_wait_started,
    c.trx_query as holding_trx_query,
    c.sql_id as holding_sql_id,
    b.trx_id as request_trx_id,
    b.trx_mysql_thread_id as request_trx_mysql_thread_id,
    b.trx_state as request_trx_state,
    b.trx_started as request_trx_started,
    b.trx_wait_started as request_trx_wait_started,
    b.trx_query as request_trx_query,
    b.sql_id as request_sql_id,
    strftime('%%s',substr(a.insert_time,1,19)) - strftime('%%s',substr(b.trx_wait_started,1,19)) as wait_time_s,
    datetime(a.insert_time,'localtime') as Time
from innodb_lock_waits a
left join innodb_trx b 
    on a.requesting_trx_id=b.trx_id and a.insert_time = b.insert_time
left join innodb_trx c 
    on a.blocking_trx_id=c.trx_id and a.insert_time = c.insert_time
order by a.insert_time asc,wait_time_s desc`,
	})
	return reportSQLs
}

func ReportCmdToString()  string {
	var reportSQLs = ReportCmd()
	str := `连接数据库:
说明:连接到数据库中 .mod line 代表采用行模式展现，.mod column 代表采用列模式展现
sqlite3 -header -column sqlite3_mysql_xxx.db`
	for _,eachR := range reportSQLs {
		str = str + "\n--" + eachR.Name + "\n" + eachR.SQL + ";"
	}
	return str
}
//这里直接采用命令行打印，不通过调用sqlite接口操作
func ReportResult(db string) {
	if finfo,err := os.Stat(db);err != nil {
		fmt.Printf("snap文件目录:[%s]不存在\n",db)
		return
	}else if finfo.IsDir() {
		fmt.Printf("snap:[%s]不是一个文件\n",db)
		return
	}else if ! strings.HasSuffix(finfo.Name(),"db") {
		fmt.Printf("文件:[%s]不是snap文件\n",db)
		return
	}
	var reportSQLs = ReportCmd()
	for _,eachR := range reportSQLs {
		sql := eachR.SQL
		cmd := exec.Command("sqlite3","-header","-column",db,sql)
		rb,_ := cmd.CombinedOutput()
		fmt.Println(eachR.Name)
		fmt.Println(string(rb))
	}

}