/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"mysqltools/instance"
	"mysqltools/ops"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"sync"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "搜集指定数据库的快照信息和日志信息",
	Long: `搜集processlist,innodb_trx,innodb_locks,innodb_lock_waits,
global_status,variables,show engine innodb status等系统表信息，
搜集最近执行的error.log日志信息以及slow.log日志信息
使用方式例如:
collect snap --host localhost --port 3306 --user root --password pwd -repeat 10,6
`,
}

var collectSnapCmd = &cobra.Command{
	Use:                        "snap",
	Aliases:                    nil,
	SuggestFor:                 nil,
	Short:                      "搜集当前数据库信息",
	Long:                       "",
	Example:                    fmt.Sprintf("%s collect snap -h <127.0.0.1> -P <3306> -u <root> -p <pwd> -repeat 10,6",filepath.Base(os.Args[0])),
	ValidArgs:                  nil,
	Args:                       nil,
	ArgAliases:                 nil,
	BashCompletionFunction:     "",
	Deprecated:                 "",
	Hidden:                     false,
	Annotations:                nil,
	Version:                    "",
	PersistentPreRun:           nil,
	PersistentPreRunE:          nil,
	PreRunE:                    nil,
	Run: func(cmd *cobra.Command, args []string) {
		//保存数据信息
		//flagNames := []string{"host","port","user","password","repeat","path"}
		var (
			host string
			port int
			user string
			password string
			repeat string
			numsec int = 10 //默认每10秒采集1次
			count int = 6  //默认一共采集6次
			path string
			allInstanceFlag bool
			saveDataParm *ops.SaveDataParm
		)
		saveDataParm = &ops.SaveDataParm{}
		host,_ = cmd.Flags().GetString("host")
		port,_ = cmd.Flags().GetInt("port")
		user,_ = cmd.Flags().GetString("user")
		password,_ = cmd.Flags().GetString("password")
		repeat,_ = cmd.Flags().GetString("repeat")
		path,_ = cmd.Flags().GetString("path")
		allInstanceFlag,_ = cmd.Flags().GetBool("all")
		//处理repeat
		if fields := strings.Split(repeat,",");len(fields) == 2 {
			if result1,err := strconv.Atoi(fields[0]);err == nil {
				numsec = result1
			}
			if result2,err := strconv.Atoi(fields[1]);err == nil {
				count = result2
			}
		}
		if ! allInstanceFlag {
			mysql_conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8&loc=Asia%%2FShanghai&parseTime=true",user,password,host,port)
			db_mysql,err := gorm.Open("mysql",mysql_conn_str)
			if err != nil {
				//log.Println(err)
				//log.Println("尝试从参数文件中查找:",ops.ParmFile)
				//尝试从参数文件中获取用户名和密码进行重试
				var p = &ops.Parm{}
				if err1 := ops.LoadParmFromFile(ops.ParmFile,p);err1 != nil {
					log.Fatal(err)
				}else {
					user = p.User
					password = p.Password
					mysql_conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8&loc=Asia%%2FShanghai&parseTime=true",user,password,host,port)
					db_mysql,err = gorm.Open("mysql",mysql_conn_str)
					if err != nil {
						log.Fatal(err)
					}
				}

			}
			if err := ops.SaveDatas(db_mysql,numsec,count,path,saveDataParm);err != nil {
				log.Println(err)
			}else {
				ops.PutParmToFile(&ops.Parm{user,password,saveDataParm.Sqlite3_dbfile},ops.ParmFile)
			}
		}else {
			//当指定了all，将会获取所有实例
			if myinstances,err := instance.GetMySQLInstances();err != nil {
				log.Fatal(err)
			}else {
				var wg *sync.WaitGroup = &sync.WaitGroup{}
				f := func(user,password,host string,port ,numsec,count int ,path string,wg *sync.WaitGroup) {
					log.Printf("[开始采集端口号为:%d的MySQL实例]\n",port)
					defer wg.Add(-1)
					mysql_conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8&loc=Asia%%2FShanghai&parseTime=true",user,password,host,port)
					if db_mysql,err := gorm.Open("mysql",mysql_conn_str);err != nil {
						log.Printf("处理实例：%d的时候出错:%s\n",port,err.Error())
					}else {
						if err := ops.SaveDatas(db_mysql,numsec,count,path,saveDataParm);err != nil {
							log.Printf("处理实例：%d的时候出错:%s\n",port,err.Error())
						}
					}
				}
				for _,eachI := range myinstances {
					port := eachI.NetStat.Port
					host := "127.0.0.1"
					wg.Add(1)
					go f(user,password,host,port,numsec,count,path,wg)
				wg.Wait()
				}
			}

		}

	},
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
	TraverseChildren:           false,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.AddCommand(collectSnapCmd)
	collectSnapCmd.Flags().StringP("host","H","127.0.0.1","IP地址")
	collectSnapCmd.Flags().IntP("port","P",3306,"端口号")
	collectSnapCmd.Flags().StringP("user","u","","用户名")
	collectSnapCmd.Flags().StringP("password","p","","用户密码")
	collectSnapCmd.Flags().StringP("repeat","r","10,6","执行时间间隔，执行次数,中间用逗号分开，两者之间不能用空白符")
	collectSnapCmd.Flags().StringP("path","f",".","信息存放路径,默认为当前路径")
	collectSnapCmd.Flags().BoolP("all","A",false,"是否在所有实例上执行，指定此参数，则忽略host和port")
	log.SetFlags(log.Lshortfile|log.LstdFlags)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
