/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// touchCmd represents the touch command
var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "Creates a new empty file in the filesystem",
	Long:  `Creates a new empty file in the filesystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("touch called")
	},
}

func init() {
	rootCmd.AddCommand(touchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// touchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// touchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
