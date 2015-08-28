package main

import (
	"log"
	"os"

	"github.com/Patrolavia/ggpk/afs"
	"github.com/Patrolavia/ggpk/record"
)

type F struct {
	h      record.RecordHeader
	f      record.FileRecord
	parent *record.DirectoryEntry
	orig   *afs.File
}

func fileRecord(file *afs.File, parent *record.DirectoryEntry) (ret F) {
	ret.f = record.FileRecord{
		NameLength: uint32(len(file.Name) + 1),
		Digest:     file.Digest,
		Name:       file.Name,
	}
	ret.h = record.RecordHeader{
		Length: uint32(ret.f.ByteLength()) + uint32(file.Size),
		Tag:    "FILE",
	}
	ret.h.Length += uint32(ret.h.ByteLength())
	ret.parent = parent
	ret.orig = file
	return

}

func (file F) save(f *os.File, o *os.File) {
	path := file.orig.Path
	if err := file.h.Save(f); err != nil {
		log.Fatalf("While writing header of %s: %s", path, err)
	}

	if err := file.f.Save(f); err != nil {
		log.Fatalf("While writing info of %s: %s", path, err)
	}

	data, err := file.orig.Content(o)
	if err != nil {
		log.Fatalf("While reading content of %s: %s", path, err)
	}

	if err := W(f, data); err != nil {
		log.Fatalf("While writing content of %s: %s", path, err)
	}
}

type D struct {
	h      record.RecordHeader
	d      record.DirectoryRecord
	parent *record.DirectoryEntry
}

func dirRecord(dir *afs.Directory, parent *record.DirectoryEntry) (ret D) {
	ret.d = record.DirectoryRecord{
		NameLength: uint32(len(dir.Name) + 1),
		ChildCount: uint32(len(dir.Subfolders) + len(dir.Files)),
		Digest:     dir.Digest,
		Name:       dir.Name,
		Entries:    make([]record.DirectoryEntry, len(dir.Subfolders)+len(dir.Files)),
	}
	ret.h = record.RecordHeader{
		Length: 0,
		Tag:    "PDIR",
	}
	ret.h.Length = uint32(ret.h.ByteLength() + ret.d.ByteLength())
	ret.parent = parent

	return
}

func (dir D) save(f *os.File) {
	if err := dir.h.Save(f); err != nil {
		log.Fatalf("Failed to save directory header of %s: %s", dir.d.Name, err)
	}
	if err := dir.d.Save(f); err != nil {
		log.Fatalf("Failed to save directory record of %s: %s", dir.d.Name, err)
	}
}
