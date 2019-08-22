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
  "github.com/spf13/cobra"
  "mysqltools/ops"
  "os"
  "path/filepath"
)





// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   filepath.Base(os.Args[0]),
  Short: "MySQL数据搜集工具",
  Long:  "搜集数据库中processlist、status、variables、show engine innodb status、innodb_trx、innodb_locks等相关视图信息" +
    "\n以及搜集语句期间执行的SQL语句的执行计划信息、slow.log和error.log信息，\n慢日志和错误日志为大小不会超过20MB的最新数据\n\n" +
    "Author:wanlgei@SoftwareGroup.DataCenter",
  Version:  "1.0.0",
  Run: func(cmd *cobra.Command, args []string) {
    readme,_ := cmd.Flags().GetBool("readme")
    if readme {
      fmt.Println(ops.ReadMe())
    }

  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  rootCmd.Flags().BoolP("readme","",false,"查看工具的帮助")
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {

}



