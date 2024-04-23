package relpath

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// Abs returns a best-effort absolute representation of path.
// The returned path may be relative on platforms such as WebAssembly.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path.
// See [filepath.Abs] for more information.
func Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	wd, err := Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, path), nil
}

// Rel returns a best-effort relative path. If an error occurs
// trying to make target relative to base, target is returned unmodified.
func Rel(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}

// CallerRel returns a source-file relative path.
// Used for testing when PWD is not set, such as running tests under wasm/wasip1.
// This does not work in TinyGo (missing [runtime.Caller] support).
func CallerRel(path string) string {
	if !filepath.IsLocal(path) {
		return path
	}
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return path
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, path)
}

// Getwd returns best-effort path to the current working directory, even on WebAssembly.
// The path may be absolute, or it may be relative, prefixed with one or more ".." segments.
func Getwd() (string, error) {
	var path string
	var target string
	parent := "."
outer:
	for {
		target = parent
		// println("statting target:", target)
		dirInfo, err := os.Stat(target)
		if err != nil {
			return path, err
		}
		parent = filepath.Join("..", target)

		dir, err := os.Open(parent)
		if err != nil {
			return filepath.Clean(filepath.Join(target, path)), nil
		}
		names, err := dir.Readdirnames(-1)
		dir.Close()
		if err != nil {
			return filepath.Clean(filepath.Join(target, path)), nil
		}
		// println("NAMES: ", strings.Join(names, " "))
		for _, name := range names {
			info, err := os.Stat(filepath.Join(parent, name))
			if err != nil {
				continue
			}
			if os.SameFile(dirInfo, info) {
				path = filepath.Join(name, path)
				continue outer
			}
		}
		break
	}
	return filepath.Clean(filepath.Join("/", path)), nil
}

// Walk walks the files in directory dir, passing them to func f.
// Supply glob patterns (e.g. "*.wit.json") to filter files passed to f.
func Walk(dir string, f func(path string) error, patterns ...string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if len(patterns) == 0 {
			return f(path)
		}
		for _, p := range patterns {
			matched, err := filepath.Match(p, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				return f(path)
			}
		}
		return nil
	})
}
