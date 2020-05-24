package filesystem

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { log.Fatal(); return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { log.Fatal(); return os.ModeTemporary }
func (m *mockFileInfo) ModTime() time.Time { log.Fatal(); return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { log.Fatal(); return nil }

func Test_newFindWalkFn(t *testing.T) {
	type mocks struct {
		walkErr error
	}
	type args struct {
		root      string
		pattern   *regexp.Regexp
		recursive bool
		path      string
		info      os.FileInfo
		err       error
	}
	tests := []struct {
		name             string
		mocks            mocks
		args             args
		shouldCallWalkFn bool
		shouldSkip       bool
	}{
		{
			name: "Should propagate error if one is passed as argument",
			args: args{
				root:      "some-argroot",
				pattern:   regexp.MustCompile(""),
				recursive: false,
				path:      "",
				info:      nil,
				err:       errors.New("some-argerror"),
			},
			shouldCallWalkFn: false,
			shouldSkip:       false,
		},
		{
			name: "Should call walkFn if the file name matches the pattern",
			mocks: mocks{
				walkErr: errors.New("some-matching/some-walk-error"),
			},
			args: args{
				root:      "some-matching/root",
				pattern:   regexp.MustCompile("some-matching.*"),
				recursive: false,
				path:      "some-matching/path",
				info: &mockFileInfo{
					name:  "some-matching/name",
					isDir: false,
				},
				err: nil,
			},
			shouldCallWalkFn: true,
			shouldSkip:       false,
		},
		{
			name: "Should not call walkFn if the file name does not matches the pattern",
			args: args{
				root:      "some-not-matching/root",
				pattern:   regexp.MustCompile("some-matching.*"),
				recursive: false,
				path:      "some-not-matching/path",
				info: &mockFileInfo{
					name:  "some-not-matching/name",
					isDir: false,
				},
				err: nil,
			},
			shouldCallWalkFn: false,
			shouldSkip:       false,
		},
		{
			name: "Should not skip folder if recursive == true",
			args: args{
				root:      "some-not-skip/root",
				pattern:   regexp.MustCompile("some-matching.*"),
				recursive: true,
				path:      "some-not-skip/path",
				info: &mockFileInfo{
					name:  "some-not-skip/name",
					isDir: true,
				},
				err: nil,
			},
			shouldCallWalkFn: false,
			shouldSkip:       false,
		},
		{
			name: "Should not skip root folder",
			args: args{
				root:      "some-not-skip-root/root",
				pattern:   regexp.MustCompile("some-matching.*"),
				recursive: true,
				path:      "some-not-skip-root/root",
				info: &mockFileInfo{
					name:  "some-not-skip-root/name",
					isDir: true,
				},
				err: nil,
			},
			shouldCallWalkFn: false,
			shouldSkip:       false,
		},
		{
			name: "Should skip folder if recursive == false",
			args: args{
				root:      "some-skip-root",
				pattern:   regexp.MustCompile("some-matching.*"),
				recursive: false,
				path:      "some-skip-path",
				info: &mockFileInfo{
					name:  "some-skip-name",
					isDir: true,
				},
				err: nil,
			},
			shouldCallWalkFn: false,
			shouldSkip:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := tt.args.err

			walkCalled := false
			walkFn := func(path string, info os.FileInfo, err error) error {
				walkCalled = true
				assert.Equal(t, tt.args.path, path)
				assert.Equal(t, tt.args.info, info)
				assert.Equal(t, tt.args.err, err)
				return tt.mocks.walkErr
			}
			if wantErr == nil {
				wantErr = tt.mocks.walkErr
			}
			if wantErr == nil && tt.shouldSkip {
				wantErr = filepath.SkipDir
			}

			err := newFindWalkFn(tt.args.root, tt.args.pattern, tt.args.recursive, walkFn)(tt.args.path, tt.args.info, tt.args.err)

			assert.Equal(t, wantErr, err)
			assert.Equal(t, tt.shouldCallWalkFn, walkCalled)
		})
	}
}
