package generate

import (
	"log"

	"github.com/Patrolavia/ggpk/afs"
	"github.com/Patrolavia/ggpk/record"
)

func generate(root *afs.Directory, parent *record.DirectoryEntry) (dirs []GGPKDirectory, files []GGPKFile) {
	dirs = append(dirs, NewGGPKDirectory(root, parent))
	idx := 0
	me := &dirs[len(dirs)-1]
	for _, f := range root.Files {
		files = append(files, NewGGPKFile(f, &me.Record.Entries[idx]))
		idx++
	}

	for _, dir := range root.Subfolders {
		d, f := generate(dir, &me.Record.Entries[idx])
		idx++
		dirs = append(dirs, d...)
		files = append(files, f...)
	}
	if idx != len(me.Record.Entries) {
		log.Fatalf("%d, %d", idx, len(me.Record.Entries))
	}

	return
}

func FromAFS(root *afs.Directory, offset uint64) (dirs []GGPKDirectory, files []GGPKFile) {
	dirs, files = generate(root, nil)
	curOffset := uint64(dirs[0].Header.Length) + offset
	for idx := 1; idx < len(dirs); idx++ {
		dirs[idx].Parent.Offset = curOffset
		curOffset += uint64(dirs[idx].Header.Length)
	}
	for idx := 0; idx < len(files); idx++ {
		files[idx].Parent.Offset = curOffset
		curOffset += uint64(files[idx].Header.Length)
	}
	return
}
