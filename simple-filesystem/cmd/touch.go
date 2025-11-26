/*
Copyright © 2025 NAME HERE <victor20030214@gmail.com>
*/

package cmd

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// touchCmd represents the touch command
var touchCmd = &cobra.Command{
	Use:   "touch <filename>",
	Short: "Creates a new empty file in the filesystem",
	Long:  `Creates a new empty file in the filesystem.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameBytes := []byte(args[0])
		if len(nameBytes) > 55 {
			return fmt.Errorf("filename cannot have more than 55 characters including the extension. Use a shorter name please")
		}

		if len(nameBytes) == 0 {
			return fmt.Errorf("enter a valid filename")
		}

		disk, err := os.OpenFile(DiskFile, os.O_RDWR, 0o644)
		if err != nil {
			return fmt.Errorf("error opening disk")
		}
		defer disk.Close()

		superBlock, err := ReadBlock(disk, SuperIndex)
		if err != nil {
			return fmt.Errorf("error reading disk")
		}

		freeInodes := binary.LittleEndian.Uint32(superBlock[20:24])

		if freeInodes == 0 {
			return fmt.Errorf("the disk is full, simple-fs only supports a total 13 files")
		}
		fmt.Println("free inodes count", freeInodes)

		inodeBlock, err := ReadBlock(disk, InodeIndex)
		if err != nil {
			return fmt.Errorf("error reading disk")
		}

		bitMapBlock, err := ReadBlock(disk, BitmapIndex)
		if err != nil {
			return fmt.Errorf("error reading disk")
		}

		fileNameBuf := make([]byte, 56)
		copy(fileNameBuf, nameBytes)

		dataBlockIndex := -1

		for i := range make([]int, 12) {
			entry := inodeBlock[i*64]

			fmt.Printf("%08b\n", entry)

			if entry == 0 {
				fmt.Printf("writing at position: %v", i)
				buf := make([]byte, 64)
				copy(buf, fileNameBuf)

				binary.LittleEndian.PutUint32(buf[56:60], 0)
				binary.LittleEndian.PutUint32(buf[60:64], uint32(i+1))

				copy(inodeBlock[i*64:(i+1)*64], buf)
				dataBlockIndex = i
				break
			}
		}

		if dataBlockIndex == -1 {
			return fmt.Errorf("no free inodes found")
		}

		bitMap := binary.LittleEndian.Uint16(bitMapBlock[:2])
		mask := uint16(1 << (dataBlockIndex + 3))
		newBitmap := bitMap | mask
		binary.LittleEndian.PutUint16(bitMapBlock[:2], newBitmap)

		// decrementing freeInodes count
		binary.LittleEndian.PutUint32(superBlock[20:24], uint32(freeInodes-1))

		_, err = WriteBlock(disk, InodeIndex, inodeBlock)
		if err != nil {
			return fmt.Errorf("some error creating file, please try again")
		}

		_, err = WriteBlock(disk, BitmapIndex, bitMapBlock)
		if err != nil {
			return fmt.Errorf("some error creating file, please try again")
		}

		_, err = WriteBlock(disk, SuperIndex, superBlock)
		if err != nil {
			return fmt.Errorf("some error creating file, please try again")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(touchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.
	// touchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// touchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
