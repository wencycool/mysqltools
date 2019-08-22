/*
Author:wanglei
直接将当前数据库中的信息show出来
*/
package cmd

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"mysqltools/ops"
	"fmt"
)

// showCmd represents the show command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "分析采集的快照,打印MySQL相关报告,目前只支持打印查询命令",
	Long:	"该命令目前提供查询历史记录的SQL命令，帮助分析故障时间段采集的指标信息",
	Run: func(cmd *cobra.Command, args []string) {
		printSQL,_ := cmd.Flags().GetBool("printSQL")
		if printSQL {
			str := ops.ReportCmdToString()
			fmt.Println(str)
		}else {
			//直接执行SQL，打印结果
			filename,_ := cmd.Flags().GetString("filename")
			var p = &ops.Parm{}
			if filename == "" {
				if err := ops.LoadParmFromFile(ops.ParmFile,p);err == nil {
					filename = p.SnapPath
				}
			}
			log.Println("快照文件为:",filename)
			ops.ReportResult(filename)
		}

	},
}

/*
var reportSystemCmd = &cobra.Command{
	Use:                        "system",
	Aliases:                    nil,
	SuggestFor:                 nil,
	Short:                      "操作系统相关运行信息",
	Long:                       "该命令分析指定的sqlite数据库，从中读取操作系统部分信息，并打印生成报告，默认打印到当前窗口中(目前只打印查询命令)",
	Example:                    fmt.Sprintf("%s report system",filepath.Base(os.Args[0])),
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
	PreRun: 					nil,
	PreRunE:                    nil,
	Run: func(cmd *cobra.Command, args []string) {
		ops.ReportCmd("sqlite3.db")
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
*/
func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().BoolP("printSQL","p",false,"打印分析SQL语句")
	reportCmd.Flags().StringP("filename","f","","需要分析的sqlite文件路径")
	//reportCmd.AddCommand(reportSystemCmd)
	//reportCmd.AddCommand(reportMySQLCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
