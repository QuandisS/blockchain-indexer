/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var (
	rpc    string
	start  int
	out    string
	limit  int
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start an indexer",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("run called")
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVar(&rpc, "rpc", "", "RPC endpoint")
	runCmd.PersistentFlags().IntVar(&start, "start", 0, "Start block")
	runCmd.PersistentFlags().StringVar(&out, "out", "", "Output file")
	runCmd.PersistentFlags().IntVar(&limit, "limit", 0, "Block limit for indexing")
	runCmd.MarkPersistentFlagRequired("rpc")

	viper.BindPFlag("rpc", runCmd.PersistentFlags().Lookup("rpc"))
	viper.BindPFlag("start", runCmd.PersistentFlags().Lookup("start"))
	viper.BindPFlag("out", runCmd.PersistentFlags().Lookup("out"))
	viper.BindPFlag("limit", runCmd.PersistentFlags().Lookup("limit"))

	viper.SetDefault("out", "blocks.log")
	viper.SetDefault("limit", 5)
	viper.SetDefault("start", 1)
}
