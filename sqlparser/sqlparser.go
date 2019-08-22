package sqlparser

import (
	ast "github.com/ruiaylin/sqlparser/ast"
	"github.com/ruiaylin/sqlparser/parser"
	"fmt"
)

//format sql
func GetTables()  {
	sql := "SELECT /*+ TIDB_SMJ(employees) */ emp_no, first_name, last_name " +
		"FROM employees USE INDEX (last_name) " +
		"where last_name='Aamodt' and gender='F' and birth_date > '1960-01-01' and a=10 and b between 10 and 20"
	sqlParser := parser.New()
	stmtnodes,_ := sqlParser.Parse(sql,"","")
	for _,nodes := range stmtnodes {
		nodes.Accept(&Visitor{})
	}
}

type Visitor struct {}

func (v *Visitor)Enter(in ast.Node) (out ast.Node, skipChildren bool)  {
	fmt.Printf("%T\n",in)
	switch in.(type) {
	case *ast.TableName:
		fmt.Println("tablename:",in.(*ast.TableName).Name.String())
	case *ast.ValueExpr:
		fmt.Println(in.(*ast.ValueExpr).GetString())
		in.(*ast.ValueExpr).SetString("?")
		fmt.Println(in.(*ast.ValueExpr).GetString())
	case *ast.BinaryOperationExpr:
		fmt.Println("xxx",in.(*ast.BinaryOperationExpr).Op.String())
	}

	return in, false
}

func (v *Visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}