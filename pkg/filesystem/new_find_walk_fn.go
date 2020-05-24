package filesystem

import (
	"os"
	"path/filepath"
	"regexp"
)

func newFindWalkFn(root string, pattern *regexp.Regexp, recursive bool, walkFn filepath.WalkFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if pattern.MatchString(info.Name()) {
			err = walkFn(path, info, nil)
		}

		if info.IsDir() && path != root && !recursive {
			return filepath.SkipDir
		}

		return err
	}
}
