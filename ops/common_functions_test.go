package ops

import "testing"

func TestReportResult(t *testing.T) {
	db := "/tmp/sqlite3_mysql_myhost1_3306_2019-07-31-05:13:25.db"
	//db := "/Users/wency/Downloads/sqlite3_mysql_t-test-db01_3306_2019-07-28-20:47:20.db"
	ReportResult(db)
}
func TestPutParmToFile(t *testing.T) {
	p := &Parm{"user","pwd","path"}
	if err := PutParmToFile(p,ParmFile);err != nil {
		t.Error(err)
	}
	if err := LoadParmFromFile(ParmFile,p);err != nil {
		t.Error(err)
	}else {
		t.Log(p)
	}
}
