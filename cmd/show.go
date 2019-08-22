/*
Author:wanglei
直接将当前数据库中的信息show出来
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mysqltools/ops"
	"os"
	"path/filepath"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "查看当前MySQL数据库相关信息",
}

var showInstancesCmd = &cobra.Command{
	Use:                        "instances",
	Aliases:                    []string{"I"},
	SuggestFor:                 nil,
	Short:                      "查看当前机器一共多少活动实例,并打印实例配置信息",
	Long:                       "",
	Example:                    fmt.Sprintf("%s show instances",filepath.Base(os.Args[0])),
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
		ops.ShowInstances()
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
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(showInstancesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
