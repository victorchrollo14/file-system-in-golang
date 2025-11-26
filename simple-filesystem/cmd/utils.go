package cmd

import (
	"fmt"
	"os"
)

// ReadBlock - blocks range from 0 - 15
func ReadBlock(disk *os.File, index int64) ([]byte, error) {
	if index < 0 || index > 15 {
		return nil, fmt.Errorf("only blocks 0 to 15 allowed")
	}

	offset := (index * int64(BlockSize))
	buf := make([]byte, BlockSize)

	_, err := disk.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func WriteBlock(disk *os.File, index int64, block []byte) (int, error) {
	if index < 0 || index > 15 {
		return 0, fmt.Errorf("only blocks 0 to 15 allowed")
	}

	offset := index * int64(BlockSize)
	n, err := disk.WriteAt(block, offset)
	if err != nil {
		return 0, err
	}

	return n, nil
}
