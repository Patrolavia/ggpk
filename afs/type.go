package afs

import (
	"encoding/binary"
	"log"
	"os"

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
	err = binary.Read(ggpk, binary.LittleEndian, data)
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
