package measure

import (
	"time"
)


//采集当前正在执行的processlist信息,information_schema.processlist
//注意gorm会自动将首字母转为大写其余小写当成字段名称，因此需要进行转换
//information_schema.processlist
type Processlist struct {
	ID 		int 			`gorm:"column:ID;		type:bigint(21)"`
	USER 	string 			`gorm:"column:USER; 	type:varchar(32)"`
	HOST 	string 			`gorm:"column:HOST;		type:varchar(64)"`
	DB 		string 			`gorm:"column:DB;		type:varchar(64)"`
	COMMAND string			`gorm:"column:COMMAND;	type:varchar(100)"`
	TIME 	int 			`gorm:"column:TIME;		type:int(7)"`
	STATE 	string 			`gorm:"column:STATE;	type:varchar(64)"`
	INFO 	string 			`gorm:"column:INFO;		type:longtext"`
}

func (Processlist) TableName() string {
	return "information_schema.processlist"
}


//目标表定义
type ProcesslistTGT struct {
	Processlist
	SqlId				string		`gorm:"column:sql_id;			type:varchar(100)"`		//记录SQL的MD5值
	InsertTime 		time.Time		`grom:"column:insert_time;		type:datetime;		index"`
}
func (ProcesslistTGT) TableName() string {
	return "processlist"
}




//information_schema.INNODB_TRX
type InnodbTrx struct {
	Trx_id 								string 		`gorm:"column:trx_id;						type:varchar(18)"`
	Trx_state 							string 		`gorm:"column:trx_state;					type:varchar(13)"`
	Trx_started 						time.Time 	`gorm:"column:trx_started;					type:datetime"`
	Trx_requested_lock_id 				string 		`gorm:"column:trx_requested_lock_id;		type:varchar(81)"`
	Trx_wait_started 					time.Time 	`gorm:"column:trx_wait_started;				type:datetime"`
	Trx_weight 							int 		`gorm:"column:trx_weight;					type:int"`
	Trx_mysql_thread_id 				int 		`gorm:"column:trx_mysql_thread_id;			type:int"`
	Trx_query 							string 		`gorm:"column:trx_query;					type:varchar(8192)"`	//执行的SQL语句
	Trx_operation_state 				string 		`gorm:"column:trx_operation_state;			type:varchar(64)"`
	Trx_tables_in_use 					int 		`gorm:"column:trx_tables_in_use;			type:int"`
	Trx_tables_locked 					int 		`gorm:"column:trx_tables_locked;			type:int"`
	Trx_lock_structs 					int 		`gorm:"column:trx_lock_structs;				type:int"`
	Trx_lock_memory_bytes 				int 		`gorm:"column:trx_lock_memory_bytes;		type:int"`
	Trx_rows_locked 					int 		`gorm:"column:Trx_rows_locked;				type:int"`
	Trx_rows_modified 					int 		`gorm:"column:trx_rows_modified;			type:int"`
	Trx_concurrency_tickets 			int 		`gorm:"column:trx_concurrency_tickets;		type:int"`
	Trx_isolation_level 				string 		`gorm:"column:trx_isolation_level;			type:varchar(16)"`
	Trx_unique_checks 					int 		`gorm:"column:trx_unique_checks;			type:int"`
	Trx_foreign_key_checks 				int 		`gorm:"column:trx_foreign_key_checks;		type:int"`
	Trx_last_foreign_key_error 			string 		`gorm:"column:trx_last_foreign_key_error;	type:varchar(256)"`
	Trx_adaptive_hash_latched 			string 		`gorm:"column:trx_adaptive_hash_latched;	type:int"`
	Trx_adaptive_hash_timeout 			int 		`gorm:"column:trx_adaptive_hash_timeout;	type:int"`
	Trx_is_read_only 					int 		`gorm:"column:trx_is_read_only;				type:int"`
	Trx_autocommit_non_locking 			int 		`gorm:"column:trx_autocommit_non_locking;	type:int"`

}

func (InnodbTrx) TableName() string {
	return "information_schema.innodb_trx"
}

type InnodbTrxTGT struct {
	InnodbTrx
	SqlId			string			`gorm:"column:sql_id;			type:varchar(100)"`		//记录SQL的MD5值
	InsertTime 		time.Time		`gorm:"column:insert_time;		type:datetime;		index"`
}
func (InnodbTrxTGT) TableName() string {
	return "innodb_trx"
}

//information_schema.INNODB_LOCKS
type InnodbLocks struct {
	Lock_id 		string 		`gorm:"column:lock_id;			type:varchar(81)"`
	Lock_trx_id 	string 		`gorm:"column:lock_trx_id;		type:varchar(18)"`
	Lock_mode 		string 		`gorm:"column:lock_mode;		type:varchar(32)"`
	Lock_type 		string		`gorm:"column:lock_type;		type:varchar(32)"`
	Lock_table 		string 		`gorm:"column:lock_table;		type:varchar(1024)"`
	Lock_index 		string 		`gorm:"column:lock_index;		type:varchar(1024)"`
	Lock_space 		int 		`gorm:"column:lock_space;		type:int"`
	Lock_page 		int 		`gorm:"column:lock_page;		type:int"`
	lock_rec 		int 		`gorm:"column:lock_rec;			type:int"`
	Lock_data 		string 		`gorm:"column:lock_data;		type:varchar(8192)"`
}

func (InnodbLocks) TableName() string {
	return "information_schema.innodb_locks"
}

type InnodbLocksTGT struct {
	InnodbLocks
	InsertTime 		time.Time		`gorm:"column:insert_time;		type:datetime;		index"`
}

func (InnodbLocksTGT) TableName() string{
	return "innodb_locks"

}


//information_schema.INNODB_LOCK_WAITS;
type InnodbLockWaits struct {
	Requesting_trx_id 	string 		`gorm:"column:requesting_trx_id;	type:varchar(18)"`
	Requested_lock_id 	string 		`gorm:"column:requested_lock_id;	type:varchar(81)"`
	Blocking_trx_id 	string 		`gorm:"column:blocking_trx_id;		type:varchar(18)"`
	Blocking_lock_id 	string	 	`gorm:"column:blocking_lock_id;		type:varchar(81)"`
}

func (InnodbLockWaits) TableName() string {
	return "information_schema.innodb_lock_waits"
}

type InnodbLockWaitsTGT struct {
	InnodbLockWaits
	InsertTime 		time.Time		`gorm:"column:insert_time;		type:datetime;		index"`
}

func (InnodbLockWaitsTGT) TableName() string {
	return "innodb_lock_waits"
}



//show global status
type GlobalStatus struct {
	Variable_name 	string 		`gorm:"column:Variable_name;	type:varchar(64)"`
	Value 			string 		`gorm:"column:Value;			type:varchar(1024)"`
	//InsertTime time.Time `type:"datetime"`
}

func (GlobalStatus)TableName() string  {
	return "information_schema.global_status"
}

type GlobalStatusTGT struct {
	GlobalStatus
	InsertTime 		time.Time 		`gorm:"column:insert_time;		type:datetime;		index"`
}
//设置目标表名
func (GlobalStatusTGT) TableName()string {
	return "global_status"
}



//show variables
type GlobalVariables struct {
	Variable_name 	string 		`gorm:"column:Variable_name;	type:varchar(64)"`
	Value 			string 		`gorm:"column:Value;			type:varchar(1024)"`
	//InsertTime time.Time `type:"datetime"`
}

func (GlobalVariables)TableName()string  {
	return "information_schema.global_variables"
}

type GlobalVariablesTGT struct {
	GlobalVariables
	InsertTime 		time.Time		`gorm:"column:insert_time;		type:datetime;		index"`
}

func (GlobalVariablesTGT) TableName() string {
	return "global_variables"
}

//show engine innodb status

type InnodbStatus struct {
	Type   		string		`gorm:"column:Type;		type:varchar(100)"`
	Name   		string		`gorm:"column:Name;		type:varchar(100)"`
	Status 		string		`gorm:"column:Status;	type:varchar(8192)"`
}

func (InnodbStatus) TableName() string {
	return "innodb_status"
}

type InnodbStatusTGT struct {
	InnodbStatus
	InsertTime 		time.Time	`gorm:"column:insert_time;		type:datetime;		index"`
}

func (InnodbStatusTGT) TableName() string {
	return "innodb_status"
}


//show slave status

type SlaveStatus struct {
	Slave_IO_State 					string			`gorm:"column:Slave_IO_State;					type:varchar(100)"`
	Master_Host						string			`gorm:"column:Master_Host;						type:varchar(100)"`
	Master_User						string			`gorm:"column:Master_User;						type:varchar(100)"`
	Master_Port						int				`gorm:"column:Master_Port;						type:int"`
	Connect_Retry					int				`gorm:"column:Connect_Retry;					type:int"`
	Master_Log_File					string			`gorm:"column:Master_Log_File;					type:varchar(100)"`
	Read_Master_Log_Pos				int				`gorm:"column:Read_Master_Log_Pos;				type:int"`
	Relay_Log_File					string			`gorm:"column:Relay_Log_File;					type:varchar(100)"`
	Relay_Log_Pos					int				`gorm:"column:Relay_Log_Pos;					type:int"`
	Relay_Master_Log_File			string			`gorm:"column:Relay_Master_Log_File;			type:varchar(100)"`
	Slave_IO_Running				string			`gorm:"column:Slave_IO_Running;					type:varchar(100)"`
	Slave_SQL_Running				string			`gorm:"column:Slave_SQL_Running;				type:varchar(100)"`
	Replicate_Do_DB					string			`gorm:"column:Replicate_Do_DB;					type:varchar(100)"`
	Replicate_Ignore_DB				string			`gorm:"column:Replicate_Ignore_DB;				type:varchar(1000)"`
	Replicate_Do_Table				string			`gorm:"column:Replicate_Do_Table;				type:varchar(100)"`
	Replicate_Ignore_Table			string			`gorm:"column:Replicate_Ignore_Table;			type:varchar(1000)"`
	Replicate_Wild_Do_Table			string			`gorm:"column:Replicate_Wild_Do_Table;			type:varchar(1000)"`
	Replicate_Wild_Ignore_Table		string			`gorm:"column:Replicate_Wild_Ignore_Table;		type:varchar(1000)"`
	Last_Errno						int				`gorm:"column:Last_Errno;						type:int"`
	Last_Error						string			`gorm:"column:Last_Error;						type:varchar(1024)"`
	Skip_Counter					int				`gorm:"column:Skip_Counter;						type:int"`
	Exec_Master_Log_Pos				int				`gorm:"column:Exec_Master_Log_Pos;				type:int"`
	Relay_Log_Space					int				`gorm:"column:Relay_Log_Space;					type:int"`
	Until_Condition					string			`gorm:"column:Until_Condition;					type:varchar(100)"`
	Until_Log_File					string			`gorm:"column:Until_Log_File;					type:varchar(100)"`
	Until_Log_Pos					int				`gorm:"column:Until_Log_Pos;					type:int"`
	Master_SSL_Allowed				string			`gorm:"column:Master_SSL_Allowed;				type:varchar(100)"`
	Master_SSL_CA_File				string			`gorm:"column:Master_SSL_CA_File;				type:varchar(255)"`
	Master_SSL_CA_Path				string			`gorm:"column:Master_SSL_CA_Path;				type:varchar(255)"`
	Master_SSL_Cert					string			`gorm:"column:Master_SSL_Cert;					type:varchar(255)"`
	Master_SSL_Cipher				string			`gorm:"column:Master_SSL_Cipher;				type:varchar(255)"`
	Master_SSL_Key					string			`gorm:"column:Master_SSL_Key;					type:varchar(255)"`
	Seconds_Behind_Master			int				`gorm:"column:Seconds_Behind_Master;			type:int"`
	Master_SSL_Verify_Server_Cert	string			`gorm:"column:Master_SSL_Verify_Server_Cert;	type:varchar(255)"`
	Last_IO_Errno					int				`gorm:"column:Last_IO_Errno;					type:int"`
	Last_IO_Error					string			`gorm:"column:Last_IO_Error;					type:varchar(8192)"`
	Last_SQL_Errno					int				`gorm:"column:Last_SQL_Errno;					type:int"`
	Last_SQL_Error					string			`gorm:"column:Last_SQL_Error;					type:varchar(8192)"`
	Replicate_Ignore_Server_Ids		string			`gorm:"column:Replicate_Ignore_Server_Ids;		type:varchar(100)"`
	Master_Server_Id				string			`gorm:"column:Master_Server_Id;					type:varchar(100)"`
	Master_UUID						string			`gorm:"column:Master_UUID;						type:varchar(100)"`
	Master_Info_File				string			`gorm:"column:Master_Info_File;					type:varchar(255)"`
	SQL_Delay						int				`gorm:"column:SQL_Delay;						type:int"`
	SQL_Remaining_Delay				int				`gorm:"column:SQL_Remaining_Delay;				type:int"`
	Slave_SQL_Running_State			string			`gorm:"column:Slave_SQL_Running_State;			type:varchar(25)"`
	Master_Retry_Count				int				`gorm:"column:Master_Retry_Count;				type:int"`
	Master_Bind						string			`gorm:"column:Master_Bind;						type:varchar(100)"`
	Last_IO_Error_Timestamp			string			`gorm:"column:Last_IO_Error_Timestamp;			type:varchar(100)"`
	Last_SQL_Error_Timestamp		string			`gorm:"column:Last_SQL_Error_Timestamp;			type:varchar(100)"`
	Master_SSL_Crl					string			`gorm:"column:Master_SSL_Crl;					type:varchar(255)"`
	Master_SSL_Crlpath				string			`gorm:"column:Master_SSL_Crlpath;				type:varchar(255)"`
	Retrieved_Gtid_Set				string			`gorm:"column:Retrieved_Gtid_Set;				type:varchar(100)"`
	Executed_Gtid_Set				string			`gorm:"column:Executed_Gtid_Set;				type:varchar(100)"`
	Auto_Position					int				`gorm:"column:Auto_Position;					type:int"`
	Replicate_Rewrite_DB			string			`gorm:"column:Replicate_Rewrite_DB;				type:varchar(100)"`
	Channel_Name					string			`gorm:"column:Channel_Name;						type:varchar(255)"`
	Master_TLS_Version				string			`gorm:"column:Master_TLS_Version;				type:varchar(100)"`


}

func (SlaveStatus) TableName() string {
	return "slave_status"
}

type SlaveStatusTGT struct {
	SlaveStatus
	InsertTime 		time.Time		`gorm:"column:insert_time;		type:datetime;		index"`
}

func (SlaveStatusTGT) TableName() string {
	return "slave_status"
}

//存放explain信息
type SQLExplain struct {
	Id 					int				`gorm:"column:id;					type:int"`
	SelectType 			string			`gorm:"column:select_type;			type:varchar(30)"`
	Table 				string			`gorm:"column:table;				type:varchar(255)"`
	Partitions 			string			`gorm:"column:partitions;			type:varchar(255)"`
	Type 				string			`gorm:"column:type;					type:varchar(30)"`
	Possible_keys		string			`gorm:"column:possible_keys;		type:varchar(255)"`
	Key 				string			`gorm:"column:key;					type:varchar(255)"`
	KeyLen 				int				`gorm:"column:key_len;				type:int"`
	Ref 				string			`gorm:"column:ref;					type:varchar(255)"`
	Rows 				string			`gorm:"column:rows;					type:int"`
	filtered            float32			`gorm:"column:filtered;				type:decimal(10,2)"`
	Extra               string			`gorm:"column:Extra;				type:varchar(1024)"`
}

func (SQLExplain) TableName()string {
	return "sql_explain"
}

type SQLExplainTGT struct {
	SQLExplain
	SqlId				string			`gorm:"column:sql_id;				type:varchar(100);	primary_key"`		//记录SQL的MD5值
	SqlStmt  			string			`gorm:"column:sql_stmt;				type:varchar(8192)"`
	InsertTime			time.Time		`gorm:"column:insert_time;			type:datetime;		index"`
}

func (SQLExplainTGT)TableName()string  {
	return "sql_explain"
}

//当参数innodb_stats_on_metadata为OFF时候
//执行show index from tablename或者执行show table stauts、查询information_schema.tables\statistics并不会触发统计信息搜集
//如果该参数设置为ON，则不进行搜集表和索引的相关元数据信息

type Tables struct {
	Table_catalog 		string  		`gorm:"column:TABLE_CATALOG;		type:varchar(512);"`
	Table_schema 		string  		`gorm:"column:TABLE_SCHEMA;			type:varchar(64);"`
	Table_name 			string  		`gorm:"column:TABLE_NAME;			type:varchar(64);"`
	Table_type 			string  		`gorm:"column:TABLE_TYPE;			type:varchar(64);"`
	Engine 				string  		`gorm:"column:ENGINE;				type:varchar(64);"`
	Version				int				`gorm:"column:VERSION;				type:int;		"`
	Row_format			string			`gorm:"column:ROW_FORMAT;			type:varchar(10);"`
	Table_rows			int				`gorm:"column:TABLE_ROWS;			type:int;		"`
	Avg_row_length		int				`gorm:"column:AVG_ROW_LENGTH;		type:int;		"`
	Data_length			int				`gorm:"column:DATA_LENGTH;			type:int;		"`
	Max_data_length		int				`gorm:"column:MAX_DATA_LENGTH;		type:int;		"`
	Index_length		int				`gorm:"column:INDEX_LENGTH;			type:int;		"`
	Data_free			int				`gorm:"column:DATA_FREE;			type:int;		"`
	Auto_increment		int				`gorm:"column:AUTO_INCREMENT;		type:int;		"`
	Create_time			time.Time		`gorm:"column:CREATE_TIME;			type:datetime;	"`
	Update_time			time.Time		`gorm:"column:UPDATE_TIME;			type:datetime;	"`
	Check_time			time.Time		`gorm:"column:CHECK_TIME;			type:datetime;	"`
	Table_collation		string			`gorm:"column:TABLE_COLLATION;		type:varchar(32);"`
	Checksum			int 			`gorm:"column:CHECKSUM;				type:int;		"`
	Create_options		string			`gorm:"column:CREATE_OPTIONS;		type:varchar(255)"`
	Table_comment		string			`gorm:"column:TABLE_COMMENT;		type:varchar(2048)"`

}

func (Tables) TableName()string  {
	return "information_schema.tables"
}

type TablesTGT struct {
	Tables
	InsertTime			time.Time		`gorm:"column:insert_time;			type:datetime;		index"`
}

func (TablesTGT) TableName()string  {
	return "tables"
}

type Indexes struct {
	Table_catalog 		string  		`gorm:"column:TABLE_CATALOG;		type:varchar(512);"`
	Table_schema 		string  		`gorm:"column:TABLE_SCHEMA;			type:varchar(64);"`
	Table_name 			string  		`gorm:"column:TABLE_NAME;			type:varchar(64);"`
	NON_UNIQUE			int				`gorm:"column:NON_UNIQUE;			type:int		"`
	INDEX_SCHEMA		string			`gorm:"column:INDEX_SCHEMA;			type:varchar(64);"`
	INDEX_NAME			string			`gorm:"column:INDEX_NAME;			type:varchar(64);"`
	SEQ_IN_INDEX		int				`gorm:"column:SEQ_IN_INDEX;			type:int;		"`
	COLUMN_NAME			string			`gorm:"column:COLUMN_NAME;			type:varchar(64);"`
	COLLATION			string			`gorm:"column:COLLATION;			type:varchar(10);"`
	CARDINALITY			int				`gorm:"column:CARDINALITY;			type:int		;"`
	SUB_PART			int				`gorm:"column:SUB_PART;				type:int		;"`
	PACKED				string			`gorm:"column:PACKED;				type:varchar(10);"`
	NULLABLE			string			`gorm:"column:NULLABLE;				type:varchar(3) ;"`
	INDEX_TYPE			string			`gorm:"column:INDEX_TYPE;			type:varchar(16);"`
	COMMENT				string			`gorm:"column:COMMENT;				type:varchar(16);"`
	INDEX_COMMENT		string			`gorm:"column:INDEX_COMMENT;		type:varchar(1024);"`
}

func (Indexes)TableName()string	 {
	return "information_schema.statistics"
}

type IndexesTGT struct {
	Indexes
	InsertTime			time.Time		`gorm:"column:insert_time;			type:datetime;		index"`
}

func (IndexesTGT) TableName()string {
	return "indexes"
}
