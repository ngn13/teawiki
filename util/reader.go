package util

import (
	"io"
	"os"
)

type Reader struct {
	File  *os.File
	Start int64
	End   int64
	Pos   int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	psize := int64(len(p))

	if r.End > 0 {
		if r.Pos >= r.End {
			return 0, io.EOF
		}

		if r.Pos+psize > r.End {
			psize = r.End - r.Pos
			r.Pos = r.End

			return io.LimitReader(r.File, psize).Read(p)
		}
	}

	read, err := r.File.Read(p)
	r.Pos += int64(read)

	return read, err
}

func (r *Reader) Close() {
	r.File.Close()
}

func NewReader(file *os.File, start_end ...int64) (*Reader, error) {
	reader := &Reader{
		File:  file,
		Start: 0,
		End:   0,
	}

	if len(start_end) > 0 {
		reader.Start = start_end[0]
	}

	if len(start_end) > 1 {
		reader.End = start_end[1]
	}

	reader.Pos = reader.Start
	reader.File.Seek(reader.Start, 0)

	return reader, nil
}
