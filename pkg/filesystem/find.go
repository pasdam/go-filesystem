package filesystem

import (
	"path/filepath"
	"regexp"
)

// Find scans through the file tree and call walkFn whenever a file with a name that matches
// the pattern is found.
// If pattern is empty a match all pattern is applied.
// If recursive is true this method scans subfolders as well.
func Find(root string, pattern string, recursive bool, walkFn filepath.WalkFunc) error {
	if len(pattern) == 0 {
		pattern = ".*"
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	err = filepath.Walk(root, newFindWalkFn(root, r, recursive, walkFn))
	return err
}
