package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64, ch chan int64) error {
	defer func() {
		if ch != nil {
			close(ch)
		}
	}()

	from, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer from.Close()

	fi, err := from.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	fromSize := fi.Size()

	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrUnsupportedFile
	}

	to, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer to.Close()

	buffer := make([]byte, 1)
	totalCopied := int64(0)
	bytesToCopy := limit

	if limit == 0 || offset+limit > fromSize {
		bytesToCopy = fromSize - offset
	}

	for totalCopied < bytesToCopy {
		bytesLeft := bytesToCopy - totalCopied
		if bytesLeft < int64(len(buffer)) {
			buffer = buffer[:bytesLeft]
		}

		n, err := from.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return ErrUnsupportedFile
		}
		if n == 0 {
			break
		}

		n, err = to.Write(buffer)
		if err != nil {
			return ErrUnsupportedFile
		}

		totalCopied += int64(n)

		p := (100 * totalCopied) / bytesToCopy

		if ch != nil {
			ch <- p
		}
	}

	return nil
}
