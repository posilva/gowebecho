package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BindPFlag utility to bind flags using viper
func BindPFlag(name string, cmd *cobra.Command) {
	viper.BindPFlag(name, cmd.Flags().Lookup(name))
}
