package afs

import (
	"errors"
	"log"
	"os"
	"sort"

	"github.com/Patrolavia/ggpk/record"
)

// FromGGPK builds afs structure from ggpk file
func FromGGPK(f *os.File) (root *Directory, err error) {
	// test if we can seek, also ensure we are at very beginning of file
	if _, err = f.Seek(0, 0); err != nil {
		return
	}

	// read GGPK sign
	rootNode, err := record.GGG(f)
	if err != nil {
		return
	}
	if rootNode.Header.Tag != "GGPK" {
		return root, errors.New("This file is not GGPK file")
	}

	// find root directory
	nodes, err := rootNode.Children(f)
	if err != nil {
		log.Fatalf("Cannot read root nodes from ggpk: %s", err)
	}
	var rootDirNode record.RecordHeader
	for _, n := range nodes {
		if n.Tag == "PDIR" {
			rootDirNode = n
			break
		}
	}

	if rootDirNode.Tag != "PDIR" {
		return root, errors.New("Cannot find root directory from ggpk")
	}

	// create afs root
	rootdir, err := record.ReadDir(f, rootDirNode)
	if err != nil {
		return
	}
	if rootdir.Name != "" {
		return root, errors.New("root dir name is not empty")
	}
	root = FromDirectoryRecord(rootDirNode, rootdir, 0)
	root.Path = "/"

	for _, e := range rootdir.Entries {
		if err = doEntry(f, e, root); err != nil {
			return
		}
	}

	sort.Sort(ByName(root.Files))
	sort.Sort(ByPath(root.Subfolders))
	return
}

func doEntry(f *os.File, e record.DirectoryEntry, cur *Directory) error {
	if _, err := f.Seek(int64(e.Offset), 0); err != nil {
		return err
	}
	h, err := record.Header(f)
	if err != nil {
		return err
	}
	h.Offset = e.Offset + uint64(h.ByteLength())
	return doHeader(f, h, cur, e.Timestamp)
}

func doDir(f *os.File, h record.RecordHeader, cur *Directory, t uint32) error {
	dir, err := record.ReadDir(f, h)
	if err != nil {
		return err
	}

	me := FromDirectoryRecord(h, dir, t)
	me.Path = cur.Path + me.Name + "/"
	cur.Subfolders = append(cur.Subfolders, me)

	for _, e := range dir.Entries {
		if err := doEntry(f, e, me); err != nil {
			return err
		}
	}

	sort.Sort(ByName(me.Files))
	sort.Sort(ByPath(me.Subfolders))
	return nil
}

func doFile(f *os.File, h record.RecordHeader, cur *Directory, t uint32) error {
	file, err := record.ReadFile(f, h)
	if err != nil {
		return err
	}

	me := FromFileRecord(h, file, t)
	me.Path = cur.Path + me.Name
	cur.Files = append(cur.Files, me)
	return nil
}

func doHeader(f *os.File, h record.RecordHeader, cur *Directory, t uint32) error {
	switch h.Tag {
	case "PDIR":
		return doDir(f, h, cur, t)
	case "FILE":
		return doFile(f, h, cur, t)
	case "FREE":
	default:
		log.Print("This record is unknown type")
	}
	return nil
}
