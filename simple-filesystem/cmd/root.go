/*
Copyright © 2025 NAME HERE <victor20030214@gmail.com>

*/

// Package cmd
package cmd

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	DiskFile                  = "disk/virtual_disk.img"
	MagicNumber        uint32 = 0x5346594D // "MYFS"
	BlockSize          uint32 = 1024
	TotalBlocks        uint32 = 16
	TotalInodes        uint32 = 13 // we could be storing 16 files, but since we decided to 1 block per file, we can have only 13inodes, so (3*64) bytes would be wasted
	FreeInodes         uint32 = 13
	FreeBlocks         uint32 = 13
	InodeStart         uint32 = 1
	BitmapStart        uint32 = 2
	DataBlockStart     uint32 = 3
	InitialBitmapBlock uint16 = uint16(0b0000000000000111) // writing the first byte in reverse, since we are using LittleEndian which writes the least significant bit first
)

func initDisk(disk *os.File) (ok bool, err error) {
	zeroBlocks := make([]byte, 16*1024)
	_, err = disk.WriteAt(zeroBlocks, 0)
	if err != nil {
		fmt.Errorf("error initializing the disk blocks %w", err)
		return false, nil
	}

	fmt.Println("Initialize virtual disk")

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, MagicNumber)
	_, err = disk.WriteAt(buf, 0)
	if err != nil {
		return false, err
	}

	// block size
	binary.LittleEndian.PutUint32(buf, BlockSize)
	_, err = disk.WriteAt(buf, 4)
	if err != nil {
		return false, err
	}

	binary.LittleEndian.PutUint32(buf, TotalBlocks)
	_, err = disk.WriteAt(buf, 8)
	if err != nil {
		return false, err
	}

	binary.LittleEndian.PutUint32(buf, TotalInodes)
	if _, err = disk.WriteAt(buf, 12); err != nil {
		return false, err
	}

	binary.LittleEndian.PutUint32(buf, FreeBlocks)
	if _, err = disk.WriteAt(buf, 16); err != nil {
		return false, err
	}

	binary.LittleEndian.PutUint32(buf, FreeInodes)
	if _, err = disk.WriteAt(buf, 20); err != nil {
		return false, err
	}

	fmt.Println("Superblock initialized and written to disk")

	twoBytesBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(twoBytesBuf, InitialBitmapBlock)

	// offset is the start of 3rd block
	if _, err = disk.WriteAt(twoBytesBuf, 2*1024); err != nil {
		return false, err
	}

	fmt.Println("Bitmap Block Initialized and written to disk")

	return true, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simple-fs",
	Short: "A simple filesytem with 16kb disk",
	Long:  "simple-fs is a simple filesystem with 16kb disk, it has 16 blocks where each block is 1kb i.e 1024 bytes.\n\nIt only supports files and exposes the following api's - mkfs, refmt, stat, rm, touch, ls, read, write",
}

// when we run mkfs this should create a new disk wih 16kb under the disk folder
// it should include the superblock data and all other blocks data ( initial representation representation )
// if disk already exists, throw a warning that disk is already created and can be used
// if user needs to reset the disk they can use a seperate command called refmt
var mkfsCmd = &cobra.Command{
	Use:   "mkfs",
	Short: "Initialize a 16kb virtual disk",
	Long: strings.TrimSpace(`
		Initializes a 16KB virtual disk with a superblock, inode table, bitmap, and data blocks.

		Filesystem Layout
		──────────────────
		Superblock (1KB)
			- Magic number: identifies your FS type (0x1400) (4 bytes)
			- Block size: 1024 bytes (4 bytes)
			- Total blocks: 16 (4 bytes)
			- Inode count: 13 (4 bytes)
			- Free inode count: dynamically updated (4 bytes)
			- Data block count: remaining blocks after metadata (13)
			- Free data block count: dynamically updated (4 bytes)
			- Start block indexes for inode table, bitmap, data area

		Inode Table (1KB)
			- 64 bytes per file (16 files)
			- filename: max 55 chars + null terminator (56 bytes)
			- filesize: 4 bytes
			- datablock: 4 bytes (location of data block)

		Free Block Bitmap (1KB)
			- Bit 0: superblock (always 1)
			- Bit 1: inode table (always 1)
			- Bit 2: bitmap (always 1)
			- Bit 3–15: data blocks (0 = free, 1 = used)

		Data Blocks (3–15)
  		- 1KB each, used to store file data

		`),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(DiskFile)
		if err == nil {
			err := fmt.Errorf("disk already exists. To reformat the disk use `simple-fs refmt`")
			fmt.Println(err)

			return
		}

		if !os.IsNotExist(err) {
			fmt.Printf("Error checking  disk %v", err)
			return
		}

		disk, err := os.Create(DiskFile)
		if err != nil {
			fmt.Printf("Error creating the disk %v", err)
			return
		}
		defer disk.Close()

		_, err = initDisk(disk)
		if err != nil {
			fmt.Printf("Error writing superblock %v", err)
			return
		}

		fmt.Println("Created a new virtual disk successfully")
	},
}

var refmtCmd = &cobra.Command{
	Use:   "refmt",
	Short: "reformats the disk",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(DiskFile)
		if err != nil {
			fmt.Println("Virtual Disk doesn't exist, use `simple-fs mkfs` to initialize a disk")
		}

		err = os.Remove(DiskFile)
		if err != nil {
			fmt.Println("Error reformatting the disk")
			return
		}

		disk, err := os.Create(DiskFile)
		if err != nil {
			fmt.Printf("Error creating the disk %v", err)
			return
		}
		defer disk.Close()

		fmt.Println("Erase all the data on virual disk")

		_, err = initDisk(disk)
		if err != nil {
			fmt.Printf("Error writing superblock %v", err)
			return
		}

		fmt.Println("Created a new virtual disk successfully")
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
	rootCmd.AddCommand(refmtCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simple-filesystem.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
