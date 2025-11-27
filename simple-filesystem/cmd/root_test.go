package cmd

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDiskFile = "test_disk.img"
)

func TestInitDisk(t *testing.T) {
	disk, err := os.Create(testDiskFile)
	if err != nil {
		fmt.Printf("Error creating the test disk %v", err)
		return
	}
	defer func() {
		_ = disk.Close()
		_ = os.Remove(testDiskFile)
	}()

	if _, err = initDisk(disk); err != nil {
		fmt.Printf("Error initializing the test disk %v", err)
		return
	}

	readDisk, err := os.OpenFile(testDiskFile, os.O_RDONLY, 0o644)
	if err != nil {
		fmt.Printf("Error reading the initialized disk %v", err)
		return
	}

	superBlockBuf := make([]byte, BlockSize)
	if _, err = readDisk.ReadAt(superBlockBuf, 0); err != nil {
		fmt.Printf("Error reading the initialized disk %v", err)
		return
	}

	assert.Equal(t, string(superBlockBuf[:4]), "MYFS")
	assert.Equal(t, binary.LittleEndian.Uint32(superBlockBuf[4:8]), BlockSize)
	assert.Equal(t, binary.LittleEndian.Uint32(superBlockBuf[8:12]), TotalBlocks)
	assert.Equal(t, binary.LittleEndian.Uint32(superBlockBuf[12:16]), TotalInodes)
	assert.Equal(t, binary.LittleEndian.Uint32(superBlockBuf[16:20]), FreeBlocks)
	assert.Equal(t, binary.LittleEndian.Uint32(superBlockBuf[20:24]), FreeInodes)

	bitMapBuf := make([]byte, BlockSize)
	if _, err = readDisk.ReadAt(bitMapBuf, BitmapIndex*int64(BlockSize)); err != nil {
		fmt.Printf("Error reading the initialized disk %v", err)
		return
	}

	assert.Equal(t, binary.LittleEndian.Uint16(bitMapBuf[:2]), InitialBitmapBlock)
}
