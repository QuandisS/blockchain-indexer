package cmd

import (
	"blockchain-indexer/internal/indexer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rpc     string
	start   uint64
	out     string
	limit   uint64
	workers int
	runCmd  = &cobra.Command{
		Use:   "run",
		Short: "Start an indexer",
		Run: func(cmd *cobra.Command, args []string) {
			indexer.Start(
				viper.GetString("rpc"),
				viper.GetUint64("start"),
				viper.GetString("out"),
				viper.GetUint64("limit"),
				viper.GetInt("workers"),
			)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVar(&rpc, "rpc", "", "RPC endpoint")
	runCmd.PersistentFlags().Uint64Var(&start, "start", 0, "Start block")
	runCmd.PersistentFlags().StringVar(&out, "out", "", "Output file")
	runCmd.PersistentFlags().Uint64Var(&limit, "limit", 0, "Block limit for indexing")
	runCmd.PersistentFlags().IntVar(&workers, "workers", 0, "Number of workers")
	runCmd.MarkPersistentFlagRequired("rpc")

	viper.BindPFlag("rpc", runCmd.PersistentFlags().Lookup("rpc"))
	viper.BindPFlag("start", runCmd.PersistentFlags().Lookup("start"))
	viper.BindPFlag("out", runCmd.PersistentFlags().Lookup("out"))
	viper.BindPFlag("limit", runCmd.PersistentFlags().Lookup("limit"))
	viper.BindPFlag("workers", runCmd.PersistentFlags().Lookup("workers"))

	viper.SetDefault("out", "blocks.log")
	viper.SetDefault("limit", 5)
	viper.SetDefault("start", 1)
	viper.SetDefault("workers", 5)
}
