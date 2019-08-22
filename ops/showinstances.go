package ops

import (
	"encoding/json"
	"fmt"
	"mysqltools/instance"
	"log"
)

func ShowInstances()  {
	instances,err := instance.GetMySQLInstances()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("当前一共:%d个实例\n",len(instances))
	for i,eachInstances := range instances {
		fmt.Printf("打印第:%d个实例信息\n",i+1)
		bs,_ := json.MarshalIndent(eachInstances,"","  ")
		fmt.Println(string(bs))
	}
}
