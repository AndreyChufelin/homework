package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidOffset         = errors.New("negative offset")
	ErrLimitOffset           = errors.New("negative limit")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrLimitOffset
	}

	from, err := os.OpenFile(fromPath, os.O_RDWR, 0o644)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer from.Close()

	fromInfo, err := from.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	fromSize := fromInfo.Size()

	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrUnsupportedFile
	}

	toInfo, err := os.Stat(toPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ErrUnsupportedFile
	}

	var to *os.File
	areSameFiles := os.SameFile(fromInfo, toInfo)
	if areSameFiles {
		dir := filepath.Dir(fromPath)
		to, err = os.CreateTemp(dir, "example")
		if err != nil {
			return ErrUnsupportedFile
		}

		defer func() {
			err = os.Remove(fromPath)
			fmt.Println(filepath.Join(dir, to.Name()))
			os.Rename(to.Name(), fromPath)
		}()
	} else {
		to, err = os.Create(toPath)
		if err != nil {
			return ErrUnsupportedFile
		}
	}
	defer to.Close()

	buffer := make([]byte, 1024)
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
		fmt.Print(getProgressbar(totalCopied, bytesToCopy))
	}

	return nil
}

func getProgressbar(total, current int64) string {
	percent := (100 * total) / current

	var result strings.Builder
	result.WriteString("\r\033[0K[")

	maxBarLength := int64(20)
	filledBarLength := percent / (100 / maxBarLength)

	for range filledBarLength {
		result.WriteString("#")
	}
	for range maxBarLength - filledBarLength {
		result.WriteString(" ")
	}

	result.WriteString(fmt.Sprintf("] %d%%", percent))

	return result.String()
}
