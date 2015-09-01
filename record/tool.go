package record

import "os"

func ReadDir(f *os.File, h RecordHeader) (ret DirectoryRecord, err error) {
	if _, err = f.Seek(int64(h.Offset), 0); err != nil {
		return
	}
	ret, err = Directory(f)
	return
}

func ReadFile(f *os.File, h RecordHeader) (ret FileRecord, err error) {
	if _, err = f.Seek(int64(h.Offset), 0); err != nil {
		return
	}
	ret, err = File(f)
	return
}
