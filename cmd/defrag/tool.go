package main

import (
	"log"

	"github.com/Patrolavia/ggpk/afs"
	"github.com/Patrolavia/ggpk/record"
)

func generate(root *afs.Directory, parent *record.DirectoryEntry) (dirs []D, files []F) {
	dirs = append(dirs, dirRecord(root, parent))
	idx := 0
	me := &dirs[len(dirs)-1]
	for _, f := range root.Files {
		files = append(files, fileRecord(f, &me.d.Entries[idx]))
		idx++
	}

	for _, dir := range root.Subfolders {
		d, f := generate(dir, &me.d.Entries[idx])
		idx++
		dirs = append(dirs, d...)
		files = append(files, f...)
	}
	if idx != len(me.d.Entries) {
		log.Fatalf("%d, %d", idx, len(me.d.Entries))
	}
	return
}

func compute(root *afs.Directory, offset uint64) (dirs []D, files []F) {
	dirs, files = generate(root, nil)
	curOffset := uint64(dirs[0].h.Length) + offset
	for idx := 1; idx < len(dirs); idx++ {
		dirs[idx].parent.Offset = curOffset
		curOffset += uint64(dirs[idx].h.Length)
	}
	for idx := 0; idx < len(files); idx++ {
		files[idx].parent.Offset = curOffset
		curOffset += uint64(files[idx].h.Length)
	}
	return
}
