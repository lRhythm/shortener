package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

// ls - line separator.
const ls byte = '\n'

// file - структура для работы с файлом как с хранилищем.
type file struct {
	r      *os.File
	reader *bufio.Reader
	w      *os.File
	writer *bufio.Writer
}

// newFile - создание file для операций чтения и записи.
func newFile(fname string) (*file, error) {
	r, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	w, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &file{
		r:      r,
		reader: bufio.NewReader(r),
		w:      w,
		writer: bufio.NewWriter(w),
	}, nil
}

// read - операция чтения из файла.
func (f *file) read() ([]byte, error) {
	b, err := f.reader.ReadBytes(ls)
	if err != nil && err == io.EOF {
		return b, nil
	}
	return b, err
}

// readRows - чтение из файла и преобразование в *[]Row.
func (f *file) readRows() (*[]Row, error) {
	b, err := f.read()
	if err != nil {
		return nil, err
	}
	var rows []Row
	if len(b) > 0 {
		err = json.Unmarshal(b, &rows)
		if err != nil {
			return nil, err
		}
	}
	return &rows, nil
}

// write - операция записи в файл.
func (f *file) write(b []byte) error {
	_, err := f.writer.Write(b)
	if err != nil {
		return err
	}
	err = f.writer.WriteByte(ls)
	if err != nil {
		return err
	}
	return f.writer.Flush()
}

// readRows - запись в файла *[]Row.
func (f *file) writeRows(rows *[]Row) error {
	b, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	err = f.write(b)
	if err != nil {
		return err
	}
	return nil
}

// close - закрытие файла после чтения из записи.
func (f *file) close() {
	_ = f.r.Close()
	_ = f.w.Close()
}
