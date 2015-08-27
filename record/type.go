// package record represents ggpk structure
package record

import (
	"encoding/binary"
	"io"
	"os"
)

// RecordHeader represents record header
type RecordHeader struct {
	Length uint32 // bytes of this record
	Tag    string // ascii string
}

// Header reads header from stream
func Header(r io.Reader) (ret RecordHeader, err error) {
	var l uint32
	if err = binary.Read(r, binary.LittleEndian, &l); err != nil {
		return
	}

	t := make([]byte, 4)
	if err = binary.Read(r, binary.LittleEndian, t); err != nil {
		return
	}

	ret = RecordHeader{l, string(t)}
	return
}

// GGGRecord is root record
type GGGRecord struct {
	Header    RecordHeader
	NodeCount uint32   // how many
	Offsets   []uint64 // file position of child node
}

// Children reads child node header from file
func (g GGGRecord) Children(f *os.File) (ret []RecordHeader, err error) {
	for i := uint32(0); i < g.NodeCount; i++ {
		if _, err = f.Seek(int64(g.Offsets[i]), 0); err != nil {
			return
		}
		node, err := Header(f)
		if err != nil {
			return ret, err
		}
		ret = append(ret, node)
	}
	return
}

// GGG reads GGGRecord from stream
func GGG(r io.Reader) (ret GGGRecord, err error) {
	if ret.Header, err = Header(r); err != nil {
		return
	}

	var c uint32
	if err = binary.Read(r, binary.LittleEndian, &c); err != nil {
		return
	}
	pos := make([]uint64, c)
	if err = binary.Read(r, binary.LittleEndian, pos); err != nil {
		return
	}
	ret.NodeCount = c
	ret.Offsets = pos
	return
}

// DirectoryEntry is an item in directory
type DirectoryEntry struct {
	Timestamp uint32
	Offset    uint64
}

func readDirectoryEntry(r io.Reader) (ret DirectoryEntry, err error) {
	err = binary.Read(r, binary.LittleEndian, &ret)
	return
}

// FileRecord is a file
type FileRecord struct {
	NameLength uint32
	Digest     [32]byte
	Name       []uint16 // file name in utf16le
}

// File reads FileRecord from stream
func File(r io.Reader) (ret FileRecord, err error) {
	var l uint32
	if err = binary.Read(r, binary.LittleEndian, &l); err != nil {
		return
	}

	var d [32]byte
	if err = binary.Read(r, binary.LittleEndian, &d); err != nil {
		return
	}

	name := make([]uint16, l)
	if err = binary.Read(r, binary.LittleEndian, name); err != nil {
		return
	}

	ret = FileRecord{l, d, name}
	return
}

// DirectoryRecord is a directory
type DirectoryRecord struct {
	NameLength uint32
	ChildCount uint32
	Digest     [32]byte
	Name       []uint16
	Children   []DirectoryEntry
}

// Directory reads DirectoryRecord from stream
func Directory(r io.Reader) (ret DirectoryRecord, err error) {
	var l uint32
	if err = binary.Read(r, binary.LittleEndian, &l); err != nil {
		return
	}

	var c uint32
	if err = binary.Read(r, binary.LittleEndian, &c); err != nil {
		return
	}

	var d [32]byte
	if err = binary.Read(r, binary.LittleEndian, &d); err != nil {
		return
	}

	n := make([]uint16, l)
	if err = binary.Read(r, binary.LittleEndian, n); err != nil {
		return
	}

	child := make([]DirectoryEntry, c)
	for i := uint32(0); i < c; i++ {
		de, err := readDirectoryEntry(r)
		if err != nil {
			return ret, err
		}
		child[i] = de
	}

	ret = DirectoryRecord{l, c, d, n, child}
	return
}

// FreeRecord is free space
type FreeRecord uint64

// Free reads FreeRecord from stream
func Free(r io.Reader) (ret FreeRecord, err error) {
	err = binary.Read(r, binary.LittleEndian, &ret)
	return
}
