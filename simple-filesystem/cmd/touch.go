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
			err := fmt.Errorf("filename cannot have more than 55 characters including the extension. Use a shorter name please")
			return err
		}

		disk, err := os.OpenFile(DiskFile, os.O_RDWR, 0o644)
		if err != nil {
			return fmt.Errorf("error opening disk")
		}
		defer disk.Close()

		superBlockBuf := make([]byte, BlockSize)
		_, err = disk.ReadAt(superBlockBuf, 0)
		if err != nil {
			return fmt.Errorf("error reading disk")
		}

		freeInodes := binary.LittleEndian.Uint32(superBlockBuf[16:20])

		if freeInodes == 0 {
			return fmt.Errorf("the disk is full, simple-fs only supports a total 13 files")
		}

		fmt.Println("free inodes count", freeInodes)

		inodeBlockBuf := make([]byte, 1024)
		_, err = disk.ReadAt(inodeBlockBuf, 1024)
		if err != nil {
			return fmt.Errorf("error reading disk")
		}

		bitMapBuf := make([]byte, BlockSize)
		_, err = disk.ReadAt(bitMapBuf, 2048)
		if err != nil {
			return fmt.Errorf("error reading disk")
		}

		fileNameBuf := make([]byte, 56)
		copy(fileNameBuf, nameBytes)

		binary.LittleEndian.PutUint32(superBlockBuf[16:20], uint32(freeInodes-1))

		var dataBlockIndex int
		for i := 0; i < 13; i++ {
			entry := inodeBlockBuf[i*64]

			if entry == 0 {
				fmt.Printf("writing at offest: %v", i)
				buf := make([]byte, 64)
				copy(buf, fileNameBuf)

				binary.LittleEndian.PutUint32(buf[56:60], 0)
				binary.LittleEndian.PutUint32(buf[60:64], uint32(i+1))

				copy(inodeBlockBuf[i*64:(i+1)*64], buf)
				dataBlockIndex = i
				break
			}
		}

		bitMap := binary.LittleEndian.Uint16(bitMapBuf[:2])
		mask := uint16(1 << (dataBlockIndex + 3))
		newBitmap := bitMap | mask
		binary.LittleEndian.PutUint16(bitMapBuf[:2], newBitmap)

		_, err = disk.WriteAt(inodeBlockBuf, 1024)
		if err != nil {
			return fmt.Errorf("some error creating file, please try again")
		}

		_, err = disk.WriteAt(superBlockBuf, 0)
		if err != nil {
			return fmt.Errorf("some error creating file, please try again")
		}

		_, err = disk.WriteAt(bitMapBuf, 2048)
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
