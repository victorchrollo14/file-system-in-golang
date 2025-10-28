/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/

// Package cmd
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simple-fs",
	Short: "A simple filesytem with 16kb disk",
	Long:  "simple-fs is a simple filesystem with 16kb disk, it has 16 blocks where each block is 1kb i.e 1024 bytes.\n\nIt only supports files and exposes the following api's - mkfs, rm, touch, ls, read, write",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("hello cobra")
	// },
}

// when we run mkfs this should create a new disk wih 16kb under the disk folder
// it should include the superblock data and all other blocks data ( initial representation representation )
// if disk already exists, throw a warning that disk is already created and can be used
// if user needs to reset the disk they can use a seperate command called refmt
var mkfsCmd = &cobra.Command{
	Use:   "mkfs",
	Short: "Initialize a 16kb virtual disk",
	Long: `Initializes a 16kb virtual disk, with a superblock, inode table, data bitmap, data blocks. 
	
	Superblock ( 1kb ) 
		magic number → identifies your FS type. ( 14 )
		block size → 1024 B.
		total blocks → 16.	
		inode count → 32.
		data block count → remaining blocks after superblock, inode table, bitmap.
		free inode count → how many unused inodes.
		free data block count → how many unused data blocks.
		start block indexes for inode table, bitmap, data area.

  Inode table ( 1KB) ( 64 bytes per file - 16 files)
	  filename - 55 characters max, string termination marker ( \0 ) ( 56 bytes )
	  filesize - 4 bytes ( size of the file )
		datablock - 4 bytes ( location of the data block where the file's data is stored )

	Free block bitmap ( 1KB )
	  Bit 0 - superblock ( always 1 )
	  Bit 1 - inode table ( always 1 )
		Bit 2 - free block bitmap block ( always 1 )
	  Bit 3-15 - ( use 0 for not used, 1 for used )

	Data block ( 3 - 15 ) ( 1KB each )
		contains the data of the files
	`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("disk/virtual_disk.img")
		if err == nil {
			err := fmt.Errorf("disk already exists. To reformat the disk use `simple-fs refmt`")
			fmt.Println(err)

			return
		}

		if !os.IsNotExist(err) {
			fmt.Printf("Error checking  disk %v", err)
			return
		}

		disk, err := os.Create("disk/virtual_disk.img")
		if err != nil {
			fmt.Printf("Error creating the disk %v", err)
			return
		}
		defer disk.Close()

		fmt.Println("Created a new virutal disk")

		zeroBlocks := make([]byte, 16*1024)
		fmt.Printf("zero blocks %v", zeroBlocks)

		os.Remove("disk/virtual_disk.img")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(mkfsCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simple-filesystem.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
