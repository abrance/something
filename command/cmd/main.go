package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var CfgPath string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sc",
		Short: "Stone Rhino Controller",
		Long: `石犀-数据流动治理平台SR-DFGP：集石犀引擎SR-Engine、石犀商城SR-Store、
石犀总控SR-Controller一体，解决数据流动带来的可视、管控、安全、优化等挑战，
并深度挖掘数据流动以实现高效、可控、防泄密等业务价值，帮助组织简化网格结构，
降低采购、管理、更新等全流程成本，提高资源利用率和业务效能。`,
		Run: func(cmd *cobra.Command, args []string) {
			// 处理配置文件逻辑
			if CfgPath != "" {
				fmt.Printf("Loading config from: %s\n", CfgPath)
				// 加载和解析配置文件的逻辑
			} else {
				fmt.Println("No config file provided, using default settings.")
			}
			fmt.Println("Welcome to sc!")
		},
	}
	rootCmd.PersistentFlags().StringVar(&CfgPath, "config", "", "Config file path")

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
