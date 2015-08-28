// package record represents ggpk structure
package record

import (
	"encoding/binary"
	"io"
	"os"
	"unicode/utf16"
)

func w(f *os.File, data interface{}, err error) (e error) {
	e = err
	if e == nil {
		e = binary.Write(f, binary.LittleEndian, data)
	}
	return
}

// RecordHeader represents record header
type RecordHeader struct {
	Length uint32 // bytes of this record
	Tag    string // ascii string
	Offset uint64 // file offset of data
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

	ret = RecordHeader{l, string(t), 0}
	return
}

// Save saves header to ggpk
func (h RecordHeader) Save(f *os.File) (err error) {
	err = w(f, h.Length, err)
	data := []byte(h.Tag)
	err = w(f, data, err)
	return
}

// ByteLength returns how many bytes occupied in ggpk file
func (h RecordHeader) ByteLength() int {
	return 8
}

// GGGRecord is root record
type GGGRecord struct {
	Header    RecordHeader
	NodeCount uint32   // how many
	Offsets   []uint64 // file position of child node
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

// Save record to file
func (g GGGRecord) Save(f *os.File) (err error) {
	err = g.Header.Save(f)
	err = w(f, g.NodeCount, err)
	err = w(f, g.Offsets, err)
	return
}

// ByteLength returns how many bytes occupied in ggpk file
func (g GGGRecord) ByteLength() int {
	return g.Header.ByteLength() + 4 + int(g.NodeCount)*8
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
		node.Offset = g.Offsets[i] + uint64(node.ByteLength())
		ret = append(ret, node)
	}
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

// Save directory entry to file
func (d DirectoryEntry) Save(f *os.File) (err error) {
	err = w(f, d.Timestamp, err)
	err = w(f, d.Offset, err)
	return
}

// ByteLength returns how many bytes occupied in ggpk file
func (d DirectoryEntry) ByteLength() int {
	return 4 + 8
}

// FileRecord is a file
type FileRecord struct {
	NameLength uint32
	Digest     []byte
	Name       string // file name in utf16le, null ended
}

// File reads FileRecord from stream
func File(r io.Reader) (ret FileRecord, err error) {
	var l uint32
	if err = binary.Read(r, binary.LittleEndian, &l); err != nil {
		return
	}

	d := make([]byte, 32)
	if err = binary.Read(r, binary.LittleEndian, &d); err != nil {
		return
	}

	name := make([]uint16, l)
	if err = binary.Read(r, binary.LittleEndian, name); err != nil {
		return
	}
	utf8Name := utf16.Decode(name)

	ret = FileRecord{l, d, string(utf8Name[:len(utf8Name)-1])}
	return
}

// Save file record to ggpk file
func (r FileRecord) Save(f *os.File) (err error) {
	name := utf16.Encode([]rune(r.Name))
	name = append(name, 0)
	err = w(f, r.NameLength, err)
	err = w(f, r.Digest, err)
	err = w(f, name, err)
	return
}

// ByteLength returns how many bytes occupied in ggpk file
func (f FileRecord) ByteLength() int {
	return 4 + 32 + int(f.NameLength*2)
}

// DirectoryRecord is a directory
type DirectoryRecord struct {
	NameLength uint32
	ChildCount uint32
	Digest     []byte
	Name       string
	Entries    []DirectoryEntry
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

	d := make([]byte, 32)
	if err = binary.Read(r, binary.LittleEndian, &d); err != nil {
		return
	}

	n := make([]uint16, l)
	if err = binary.Read(r, binary.LittleEndian, n); err != nil {
		return
	}
	utf8Name := []rune(" ")
	if len(n) != 1 || n[0] != 0 {
		utf8Name = utf16.Decode(n)
	}

	child := make([]DirectoryEntry, c)
	for i := uint32(0); i < c; i++ {
		de, err := readDirectoryEntry(r)
		if err != nil {
			return ret, err
		}
		child[i] = de
	}

	ret = DirectoryRecord{l, c, d, string(utf8Name[:len(utf8Name)-1]), child}
	return
}

// Save directory record to ggpk file
func (d DirectoryRecord) Save(f *os.File) (err error) {
	name := utf16.Encode([]rune(d.Name))
	name = append(name, 0)
	err = w(f, d.NameLength, err)
	err = w(f, d.ChildCount, err)
	err = w(f, d.Digest, err)
	err = w(f, name, err)
	for _, n := range d.Entries {
		if e := n.Save(f); e != nil {
			err = e
		}
	}
	return
}

func (d DirectoryRecord) Children(f *os.File) (ret []RecordHeader, err error) {
	for _, e := range d.Entries {
		if _, err = f.Seek(int64(e.Offset), 0); err != nil {
			return
		}
		h, err := Header(f)
		if err != nil {
			return ret, err
		}
		h.Offset = e.Offset + uint64(h.ByteLength())
		ret = append(ret, h)
	}
	return
}

// ByteLength returns how many bytes occupied in ggpk file
func (d DirectoryRecord) ByteLength() (ret int) {
	ret = 4 + 4 + 32 + int(d.NameLength*2)
	for _, e := range d.Entries {
		ret += e.ByteLength()
	}
	return
}

// FreeRecord is free space
type FreeRecord uint64

// Free reads FreeRecord from stream
func Free(r io.Reader) (ret FreeRecord, err error) {
	err = binary.Read(r, binary.LittleEndian, &ret)
	return
}

func (n FreeRecord) Next(f *os.File) (ret FreeRecord, err error) {
	if n == 0 {
		return
	}

	if _, err = f.Seek(int64(n), 0); err != nil {
		return
	}

	ret, err = Free(f)
	return
}

// ByteLength returns how many bytes occupied in ggpk file
func (f FreeRecord) ByteLength() int {
	return 8
}
