package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFile = "test.txt"

func TestCopy(t *testing.T) {
	t.Run("when from is empty", func(t *testing.T) {
		err := Copy("", testFile, 0, 0)

		require.EqualError(t, err, ErrPathIsEmpty.Error())

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when to is empty", func(t *testing.T) {
		err := Copy("testdata/input.txt", "", 0, 0)

		require.EqualError(t, err, ErrPathIsEmpty.Error())

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when offset below 0", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFile, -3, 0)

		require.EqualError(t, err, ErrOffsetOrLimitBellowZero.Error())

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when limit below 0", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFile, 0, -5)

		require.EqualError(t, err, ErrOffsetOrLimitBellowZero.Error())

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFile, 100000000000, 0)

		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when file without extension", func(t *testing.T) {
		err := Copy("testdata/input", testFile, 0, 0)

		require.EqualError(t, err, ErrUnsupportedFile.Error())

		err = Copy("testdata/input.txt", "out", 0, 0)

		require.EqualError(t, err, ErrUnsupportedFile.Error())

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when from file does not exist", func(t *testing.T) {
		inputFile := "testdata/input111.txt"
		err := Copy(inputFile, testFile, 150, 100)
		require.EqualError(t, err, fmt.Sprintf("file doesn't exist: open %s: no such file or directory", inputFile))
	})

	t.Run("successful copied", func(t *testing.T) {
		err := Copy("testdata/input.txt", testFile, 150, 100)

		require.NoError(t, err)

		file, err := os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.NoError(t, err)

		fileStat, err := file.Stat()
		require.NoError(t, err)
		require.Equal(t, 100, int(fileStat.Size()))

		err = os.Remove(testFile)
		require.NoError(t, err)

		_, err = os.OpenFile(testFile, os.O_RDONLY, os.ModePerm)
		require.Equal(t, true, os.IsNotExist(err))
	})
}
