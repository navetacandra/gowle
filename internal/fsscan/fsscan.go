package fsscan

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/navetacandra/gowle/internal/config"
)

type Info struct {
	Path    string
	ModSize int64
	Size    int64
}

type DiffInfo struct {
	Path string
	Diff int8
}

func DiffSnapshot(newSnapshot *[]Info, oldSnapshot *[]Info, diff *[]DiffInfo) {
	*diff = (*diff)[:0]
	new := make(map[string]Info)
	old := make(map[string]Info)

	for _, f := range *oldSnapshot {
		old[f.Path] = f
	}
	for _, f := range *newSnapshot {
		new[f.Path] = f
	}

	for p, o := range old {
		n, exist := new[p]
		switch {
		case !exist:
			*diff = append(*diff, DiffInfo{Path: p, Diff: -1}) // delete
		case o.Size != n.Size || o.ModSize != n.ModSize:
			*diff = append(*diff, DiffInfo{Path: p, Diff: 0}) // modified
		}
	}

	for p := range new {
		_, exist := old[p]
		if !exist {
			*diff = append(*diff, DiffInfo{Path: p, Diff: 1}) // created
		}
	}
}

func Scan(snapshot *[]Info, cfg *config.GowleConfig) {
	cwd, _ := os.Getwd()
	ignored := make(map[string]int8)
	*snapshot = (*snapshot)[:0]

	filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()

		if d.IsDir() {
			if shouldIgnore(path, name, &cfg.Ignore, &ignored) {
				ignored[path] = 1
			}
			return nil
		}

		if shouldIgnore(path, name, &cfg.Ignore, &ignored) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(cwd, path)
		if shouldWatch(&relPath, &cfg.Watch) {
			*snapshot = append(*snapshot, Info{
				Size:    info.Size(),
				ModSize: info.ModTime().Unix(),
				Path:    relPath,
			})
		}
		return nil
	})
}

func shouldWatch(path *string, prefixes *[]string) bool {
	if len(*prefixes) == 0 {
		return true
	}

	for _, p := range *prefixes {
		if strings.HasPrefix(*path, p+"/") || *path == p {
			return true
		}
	}
	return false
}

func shouldIgnore(absPath, base string, patterns *[]*regexp.Regexp, ignored *map[string]int8) bool {
	for dir := range *ignored {
		if strings.HasPrefix(absPath, dir+"/") {
			return true
		}
	}

	for _, r := range *patterns {
		if r.MatchString(base) {
			return true
		}
	}

	return false
}
