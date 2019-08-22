package ops

func ReadMe() string  {
	str := `
<MySQL故障信息采集工具>
【快速开始]
1. 采集数据
./mysqltool collect snap -H127.0.0.1 -P3306 -uroot -pxxx -r 10,6 -f "/tmp"
2. 分析数据
./mysqltool report

文档说明：
在MySQL数据库出问题以后，可以利用该工具快速搜集当前数据库的一些性能信息，方便后续处理
目前该工具支持三个子命令：show, collect ,report
show命令：
    show命令是用来查询当前数据库的一些状态信息，即时从当前操作系统中查询出最新数据
    可以利用./mysqltool show -h  查看show支持哪些操作
    比如可以利用./mysqltool show instances查看当前操作系统上运行了哪些数据库实例
        [root@myhost1 mysqltools]# ./mysqltool show instances
        当前一共:1个实例
        打印第:1个实例信息
        {
        "User": {
            "Uid": "1000",
            "Gid": "1000",
            "Username": "mysql",
            "Name": "",
            "HomeDir": "/home/mysql"
        },
        "PID": 4320,
        "MysqldPath": "/opt/mysql/bin/mysqld",
        "NetStat": {
            "Stat": "LISTEN",
            "ProjName": "mysqld",
            "Port": 3306,
            "SocketFile": "/data/mysql/my3306/mysql.sock"
        },
        "Mycnf": "/etc/my.cnf",
        "ShortVersion": "5.7.25"
        }

collect命令：
    collect命令是在故障时刻搜集当时数据库相关性能信息并存放到快照文件中，以用来保存分析数据供后续进行分析
    可以利用./mysqltool collect -h  查看collect支持哪些操作，目前只支持搜集快照操作即：
    ./mysqltool collect snap 在执行该命令的时候需要指定：IP、端口号、用户名、用户密码、数据采集频率、快照存放路径 相关参数，也可以不指定采用默认值。
    用户和密码信息可以输入一次，在后续的使用中不必输入，下次使用会自动获取之前的用户和密码信息
    示例：
        1. 获取本机MySQL默认端口下的性能指标信息并存放到/tmp/3306（注意这个路径必须存在)中
        [root@myhost1 mysqltools]# ./mysqltool collect snap -uroot -proot -r 2,5 -f "/tmp/3306"
        2019/08/01 01:35:50 savedata.go:49: 开始采集数据,数据存放路径为: /tmp/3306
        2019/08/01 01:35:50 savedata.go:62: 指标信息存放在sqlite文件中: /tmp/3306/sqlite3_mysql_myhost1_3306_2019-08-01-01:35:50.db
        2019/08/01 01:35:50 savedata.go:78: 一共采集5轮数据
        2019/08/01 01:35:50 savedata.go:81: 开始第:1轮的数据采集...
        2019/08/01 01:35:52 savedata.go:81: 开始第:2轮的数据采集...
        2019/08/01 01:35:54 savedata.go:81: 开始第:3轮的数据采集...
        2019/08/01 01:35:56 savedata.go:81: 开始第:4轮的数据采集...
        2019/08/01 01:35:58 savedata.go:81: 开始第:5轮的数据采集...
        2019/08/01 01:36:00 savedata.go:91: 保存diag日志文件路径为:/tmp/3306/sqlite3_mysql_myhost1_3306_diag_2019-08-01.log
        2019/08/01 01:36:00 savedata.go:94: slow_query_log is OFF
        2. 因为上面已经输入了用户名和密码，那么下次不必输入即可进行搜集数据
        [root@myhost1 mysqltools]# ./mysqltool collect snap -r 2,1 -f "/tmp/3306"
        2019/08/01 01:38:19 savedata.go:49: 开始采集数据,数据存放路径为: /tmp/3306
        2019/08/01 01:38:19 savedata.go:62: 指标信息存放在sqlite文件中: /tmp/3306/sqlite3_mysql_myhost1_3306_2019-08-01-01:38:19.db
        2019/08/01 01:38:19 savedata.go:78: 一共采集1轮数据
        2019/08/01 01:38:19 savedata.go:81: 开始第:1轮的数据采集...
        2019/08/01 01:38:21 savedata.go:91: 保存diag日志文件路径为:/tmp/3306/sqlite3_mysql_myhost1_3306_diag_2019-08-01.log
        2019/08/01 01:38:21 savedata.go:94: slow_query_log is OFF
        [root@myhost1 mysqltools]# 

report命令：
    report命令用来分析由collect生成的snap快照信息，从中得到一些指标的分析结果
    ./mysqltool report  直接运行该命令会从collect搜集的快照信息中进行解析快照文件并生成汇总信息，也可以自己进入快照文件中进行查询。
    ./mysqltool report --printSQL  会打印常用的一些查询快照文件的SQL语句，快照文件为sqlite3文本数据库文件，
                                可以通过sqlite3 -header -column xxx.db 进入并执行相应SQL语句
监控指标采集说明：
工具会将当前操作系统和数据库中相关性能数据采集到
sqlite> .table
global_status      innodb_status      sql_explain        vm_memstat       
global_variables   innodb_trx         vm_cpuinfo         vm_procstat_top10
innodb_lock_waits  processlist        vm_cpustat         vm_swapstat      
innodb_locks       slave_status       vm_diskstat  
上述这些表中，按照参数配置的采集频率进行采集（操作系统相关指标为内部配置,固定每5秒采集一次)
sql_expain表并不是每次都进行采集，而是先判断是否应做了采集，如果已经在快照中存在则不进行采集。


`
	return str
}