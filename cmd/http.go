/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Launching EcoBay App",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")
		if err := gateway.Server(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
