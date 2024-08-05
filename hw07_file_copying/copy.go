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

	inputFile, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file doesn't exist: %w", err)
		}
		return fmt.Errorf("file cannot open: %w", err)
	}
	defer inputFile.Close()

	fileStat, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed getting file info: %w", err)
	}
	fileSize := fileStat.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit+offset > fileSize {
		limit = fileSize - offset
	}

	_, err = inputFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("set offset in resource file: %w", err)
	}

	limitReader := io.LimitReader(inputFile, limit)
	bar := pb.Full.Start64(limit)
	readerWithPB := bar.NewProxyReader(limitReader)

	err = copyFile(readerWithPB, toPath)
	if err != nil {
		return err
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

	if path.Ext(fromPath) == "" || path.Ext(toPath) == "" {
		return ErrUnsupportedFile
	}

	return nil
}

func copyFile(resource io.Reader, copyPath string) error {
	file, err := os.Create(copyPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resource)
	if err != nil {
		defer os.Remove(copyPath)
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}
