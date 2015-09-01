package record

import "os"

// ReadDir record from file
func ReadDir(f *os.File, h RecordHeader) (ret DirectoryRecord, err error) {
	if _, err = f.Seek(int64(h.Offset), 0); err != nil {
		return
	}
	ret, err = Directory(f)
	return
}

// ReadFile record from file
func ReadFile(f *os.File, h RecordHeader) (ret FileRecord, err error) {
	if _, err = f.Seek(int64(h.Offset), 0); err != nil {
		return
	}
	ret, err = File(f)
	return
}
