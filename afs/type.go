package afs

import (
	"crypto/sha256"
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Patrolavia/ggpk/record"
)

// File represents virtual file
type File struct {
	Path      string
	Name      string
	Timestamp uint32
	Digest    []byte
	Size      uint64
	Offset    uint64
	OrigFile  *os.File
}

// FromFileRecord creates File from ggpk record
func FromFileRecord(h record.RecordHeader, f record.FileRecord, t uint32) *File {
	return &File{
		Path:      "",
		Name:      f.Name,
		Timestamp: t,
		Digest:    f.Digest,
		Size:      uint64(h.Length) - uint64(h.ByteLength()+f.ByteLength()),
		Offset:    h.Offset + uint64(f.ByteLength()),
		OrigFile:  f.OrigFile,
	}
}

// FromFile create afs file from physic file
func FromFile(f *os.File) (ret *File, err error) {
	info, err := f.Stat()
	if err != nil {
		return
	}

	if _, err = f.Seek(0, 0); err != nil {
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	sum := sha256.Sum256(data)
	digest := make([]byte, len(sum))
	for k, v := range sum {
		digest[k] = v
	}

	ret = &File{
		Path:      "",
		Name:      info.Name(),
		Timestamp: uint32(info.ModTime().Unix()),
		Digest:    digest,
		Size:      uint64(info.Size()),
		Offset:    0,
		OrigFile:  f,
	}
	return
}

// Dump will dump some info for debug
func (f *File) Dump() {
	log.Print(f.Path)
}

// Content reads file content from ggpk file
func (f *File) Content() (data []byte, err error) {
	if _, err = f.OrigFile.Seek(int64(f.Offset), 0); err != nil {
		return
	}

	data = make([]byte, f.Size)
	err = binary.Read(f.OrigFile, binary.LittleEndian, data)
	return
}

// Directory represents virtual directory
type Directory struct {
	Path       string
	Name       string
	Timestamp  uint32
	Digest     []byte
	Subfolders []*Directory
	Files      []*File
	Offset     uint64
}

// Root creates empty root record
func Root() *Directory {
	return &Directory{
		Path:       "",
		Name:       "",
		Timestamp:  uint32(time.Now().Unix()),
		Digest:     make([]byte, 0),
		Subfolders: make([]*Directory, 0),
		Files:      make([]*File, 0),
		Offset:     0,
	}
}

// FromDirectoryRecord creates Directory from ggpk record
func FromDirectoryRecord(h record.RecordHeader, d record.DirectoryRecord, t uint32) *Directory {
	return &Directory{
		Path:       "",
		Name:       d.Name,
		Timestamp:  t,
		Digest:     d.Digest,
		Subfolders: make([]*Directory, 0),
		Files:      make([]*File, 0),
		Offset:     h.Offset,
	}
}

// Dump will dump some info for debug
func (d *Directory) Dump() {
	log.Print(d.Path)
	for _, f := range d.Files {
		f.Dump()
	}
	for _, f := range d.Subfolders {
		f.Dump()
	}
}

// ByName can sort files by filename
type ByName []*File

func (b ByName) Len() int           { return len(b) }
func (b ByName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByName) Less(i, j int) bool { return b[i].Name < b[j].Name }

// ByPath can sort directories by their name
type ByPath []*Directory

func (b ByPath) Len() int           { return len(b) }
func (b ByPath) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPath) Less(i, j int) bool { return b[i].Name < b[j].Name }
