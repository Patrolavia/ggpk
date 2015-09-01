package generate

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/Patrolavia/ggpk/afs"
	"github.com/Patrolavia/ggpk/record"
)

func w(f *os.File, data interface{}) error {
	return binary.Write(f, binary.LittleEndian, data)
}

type GGPKFile struct {
	Header record.RecordHeader
	Record record.FileRecord
	Parent *record.DirectoryEntry
	Orig   *afs.File
}

func NewGGPKFile(file *afs.File, parent *record.DirectoryEntry) (ret GGPKFile) {
	ret.Record = record.FileRecord{
		NameLength: uint32(len(file.Name) + 1),
		Digest:     file.Digest,
		Name:       file.Name,
	}
	ret.Header = record.RecordHeader{
		Length: uint32(ret.Record.ByteLength()) + uint32(file.Size),
		Tag:    "FILE",
	}
	ret.Header.Length += uint32(ret.Header.ByteLength())
	ret.Parent = parent
	ret.Orig = file
	parent.Timestamp = file.Timestamp
	return

}

func (file GGPKFile) Save(f *os.File) {
	path := file.Orig.Path
	if err := file.Header.Save(f); err != nil {
		log.Fatalf("While writing header of %s: %s", path, err)
	}

	if err := file.Record.Save(f); err != nil {
		log.Fatalf("While writing info of %s: %s", path, err)
	}

	data, err := file.Orig.Content()
	if err != nil {
		log.Fatalf("While reading content of %s: %s", path, err)
	}

	if err := w(f, data); err != nil {
		log.Fatalf("While writing content of %s: %s", path, err)
	}
}

func (file GGPKFile) Size() uint32 {
	return file.Header.Length - uint32(file.Header.ByteLength() + file.Record.ByteLength())
}

type GGPKDirectory struct {
	Header record.RecordHeader
	Record record.DirectoryRecord
	Parent *record.DirectoryEntry
}

func NewGGPKDirectory(dir *afs.Directory, parent *record.DirectoryEntry) (ret GGPKDirectory) {
	ret.Record = record.DirectoryRecord{
		NameLength: uint32(len(dir.Name) + 1),
		ChildCount: uint32(len(dir.Subfolders) + len(dir.Files)),
		Digest:     dir.Digest(),
		Name:       dir.Name,
		Entries:    make([]record.DirectoryEntry, len(dir.Subfolders)+len(dir.Files)),
	}
	ret.Header = record.RecordHeader{
		Length: 0,
		Tag:    "PDIR",
	}
	ret.Header.Length = uint32(ret.Header.ByteLength() + ret.Record.ByteLength())
	ret.Parent = parent
	if parent != nil {
		parent.Timestamp = dir.Timestamp
	}

	return
}

func (dir GGPKDirectory) Save(f *os.File) {
	if err := dir.Header.Save(f); err != nil {
		log.Fatalf("Failed to save directory header of %s: %s", dir.Record.Name, err)
	}
	if err := dir.Record.Save(f); err != nil {
		log.Fatalf("Failed to save directory record of %s: %s", dir.Record.Name, err)
	}
}
