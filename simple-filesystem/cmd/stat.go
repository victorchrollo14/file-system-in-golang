/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statCmd represents the stat command
var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Gives a brief description of the file",
	Long:  `Shows information about the file such as fileName, Size, Data block, Inode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stat called")
	},
}

func init() {
	rootCmd.AddCommand(statCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
