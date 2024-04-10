package main

import (
	"os"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func TestCopy(t *testing.T) {
	t.Run("when from is empty", func(t *testing.T) {
		err := Copy("", "out.txt", 0, 0)

		require.EqualError(t, err, ErrPathIsEmpty.Error())
	})

	t.Run("when to is empty", func(t *testing.T) {
		err := Copy("testdata/input.txt", "", 0, 0)

		require.EqualError(t, err, ErrPathIsEmpty.Error())
	})

	t.Run("when offset below 0", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", -3, 0)

		require.EqualError(t, err, ErrOffsetOrLimitBellowZero.Error())
	})

	t.Run("when limit below 0", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 0, -5)

		require.EqualError(t, err, ErrOffsetOrLimitBellowZero.Error())
	})

	t.Run("when offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 100000000000, 0)

		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("when file with unsupported type", func(t *testing.T) {
		err := Copy("testdata/input.zip", "out.txt", 0, 0)

		require.EqualError(t, err, ErrUnsupportedFile.Error())

		err = Copy("testdata/input.txt", "out", 0, 0)

		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("successful copied", func(t *testing.T) {
		testFile := "test.txt"
		err := Copy("testdata/input.txt", testFile, 150, 100)

		require.NoError(t, err)

		file, err := os.OpenFile(testFile, os.O_RDONLY, 0644)
		require.NoError(t, err)

		fileStat, err := file.Stat()
		require.NoError(t, err)
		require.Equal(t, 100, int(fileStat.Size()))

		err = os.Remove(testFile)
		require.NoError(t, err)

		_, err = os.OpenFile(testFile, os.O_RDONLY, 0644)
		require.Equal(t, true, os.IsNotExist(err))
	})
}
