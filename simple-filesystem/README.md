
# simple-fs

A minimal 16KB virtual filesystem implemented in Go.
It simulates a disk with fixed-size blocks, superblock, an inode table, bitmap and data blocks for learning how filesystems work internally.

Info dump (bunch of stuff I understood)

- An actual hard disk is like a file filled with only 0 and 1, in Hdds 0's and 1's are represented with magnetic polarity and in ssds we use electrical charge to represent 0's and 1's.
- essentially all your data is stored in binary, the filesystem is an abstractions that kind of tells the user your data is stored in directories, files and gives you a bunch of utilities to create, delete, update, read data stored on disk.
- Data on disk is usually stored in block (sectors) usually (512/4096)bytes, when we read data from a disk, the os actually loads the whole block into memory and same goes with writes, if you update a few bytes, os writes those bytes to the block in memory then updates the whole block on disk;

---

## limitations

- supports only files, we will have only 13 files max since we have 13 data blocks and I've decided to give 1 block per file, so each file can only be 1024 bytes long
- file names can only 55 characters.
- we are not doing all the file permission, timestamps and stuff

---

## 🧩 Filesystem Layout

```
16KB Virtual Disk
├── Superblock (1KB)
│   ├── Magic number (4B)
│   ├── Block size (4B)
│   ├── Total blocks (4B)
│   ├── Inode count (4B)
│   ├── Free inode count (4B)
│   ├── Data block count (4B)
│   ├── Free data block count (4B)
│   └── Start block indexes (inode, bitmap, data) (4B each)
│
├── Inode Table (1KB)
│   ├── 64 bytes per file (16 total)
│   ├── filename (55 chars + '\o')
│   ├── filesize (4B)
│   └── data block pointer (4B)
│
├── Free Block Bitmap (1KB) (using only 2 bytes - 1022 bytes wasted) 
│   ├── Bit 0: superblock (1)
│   ├── Bit 1: inode table (1)
│   ├── Bit 2: bitmap (1)
│   └── Bits 3–15: data blocks (0 = free, 1 = used)
│
└── Data Blocks (3–15)
    └── Each 1KB, stores file contents
```

---

## ⚙️ Commands

| Command           | Description                                                       |
| ----------------- | ----------------------------------------------------------------- |
| `simple-fs mkfs`  | Initializes a new 16KB virtual disk under `disk/virtual_disk.img` |
| `simple-fs refmt` | (Planned) Reformat or reset the disk                              |
| `simple-fs ls`    | (Planned) List files in the filesystem                            |
| `simple-fs touch` | (Planned) Create an empty file                                    |
| `simple-fs read`  | (Planned) Read file contents                                      |
| `simple-fs write` | (Planned) Write data to a file                                    |
| `simple-fs rm`    | (Planned) Delete a file                                           |
| `simple-fs stat`  | (Planned) Show file metadata                                      |

---

## 🚀 Usage

```bash
# Build the binary
go build -o simple-fs

# Initialize the virtual disk
./simple-fs mkfs

# Inspect binary data (e.g. bitmap at block 3)
xxd -s $((3*1024)) -l 16 disk/virtual_disk.img
```

---

## 🧠 Notes

- Uses **little-endian** encoding for all multibyte fields.
- The first three blocks (superblock, inode table, bitmap) are always marked as used.
- Each bit in the bitmap corresponds to one 1KB block.
- Designed for educational purposes — not a production filesystem.

---

## 📁 Project Structure

```
simple-fs/
├── cmd/
│   └── mkfs.go      # mkfs command & disk initialization logic
├── disk/
│   └── virtual_disk.img  # created after running mkfs
├── go.mod
└── README.md
```

---
