package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	SupportedExtensions = map[string]struct{}{
		".txt": {},
	}

	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrOffsetExceedsFileSize   = errors.New("offset exceeds file size")
	ErrPathIsEmpty             = errors.New("path is empty")
	ErrOffsetOrLimitBellowZero = errors.New("offset and limit must be greater than or equal to zero")
)

const (
	fileBufSize   = 100
	timeoutForBar = time.Second / 10
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := validate(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	inputFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file doesn't exist: %v", err)
		}
		return fmt.Errorf("file cannot open: %v", err)
	}
	defer inputFile.Close()

	fileStat, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed getting file info: %v", err)
	}
	fileSize := fileStat.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	outputFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("file cannot create: %v", err)
	}
	defer outputFile.Close()

	if limit == 0 {
		limit = fileSize
	}

	maxBytes := limit
	if (offset+limit)-fileSize > 0 {
		maxBytes = fileSize - offset
	}

	buf := make([]byte, fileBufSize)
	writtenBytes := 0
	bar := pb.StartNew(int(maxBytes))

	for writtenBytes < int(maxBytes) {
		time.Sleep(timeoutForBar)

		read, errRead := inputFile.ReadAt(buf, offset)
		if read > int(limit) {
			read = int(limit)
		}

		bar.Add(read)
		offset += int64(read)
		writtenBytes += read

		_, err := outputFile.Write(buf[:read])
		if err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}

		if errRead == io.EOF {
			break
		}
		if errRead != nil {
			return fmt.Errorf("failed to read: %v", err)
		}
	}

	bar.Finish()
	return nil
}

func validate(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrPathIsEmpty
	}

	if offset < 0 || limit < 0 {
		return ErrOffsetOrLimitBellowZero
	}

	_, fromPathExtOk := SupportedExtensions[path.Ext(fromPath)]
	_, toPathExtOk := SupportedExtensions[path.Ext(toPath)]
	if !fromPathExtOk || !toPathExtOk {
		return ErrUnsupportedFile
	}

	return nil
}
