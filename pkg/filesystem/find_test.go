package filesystem_test

import (
	"errors"
	"os"
	"regexp/syntax"
	"testing"

	"github.com/pasdam/go-filesystem/pkg/filesystem"
	"github.com/stretchr/testify/assert"
)

func Test_Find_ShouldReturnErrorIfPatternIsInvalid(t *testing.T) {
	pattern := "(invalid-pattern"

	err := filesystem.Find("testdata", pattern, false, nil)

	assert.Equal(t, syntax.ErrorCode("missing closing )"), err.(*syntax.Error).Code)
	assert.Equal(t, pattern, err.(*syntax.Error).Expr)
}

func Test_Find_ShouldPropagateErrorIfWalkReturnsIt(t *testing.T) {
	expected := errors.New("Expected error")
	walkFn := func(path string, info os.FileInfo, err error) error {
		return expected
	}

	err := filesystem.Find("testdata", ".*", false, walkFn)

	assert.Equal(t, expected, err)
}

func Test_Find_ShouldIgnoreFilesInSubfoldersIfRecursiveIsFalse(t *testing.T) {
	pattern := "some.*"

	var files []string

	err := filesystem.Find("testdata", pattern, false, func(path string, info os.FileInfo, err error) error {
		files = append(files, info.Name())
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, "some_file", files[0])
	assert.Equal(t, "some_folder", files[1])
}

func Test_Find_ShouldMatchAllIfPatternIsEmpty(t *testing.T) {
	pattern := ""

	var files []string

	err := filesystem.Find("testdata", pattern, true, func(path string, info os.FileInfo, err error) error {
		files = append(files, info.Name())
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, 4, len(files))
	assert.Equal(t, "testdata", files[0])
	assert.Equal(t, "some_file", files[1])
	assert.Equal(t, "some_folder", files[2])
	assert.Equal(t, "some_other_file", files[3])
}
