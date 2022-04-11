package local

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker/fs"
)

type Local struct {
	root string
	cwd  string
}

func New(root string) *Local {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &Local{
		root: root,
		cwd:  wd,
	}
}

func (lfs *Local) Chdir(path string) {
	lfs.cwd = path
}

func (lfs *Local) Stat(ctx context.Context, path string) (*fs.FileInfo, error) {
	pivotedPath, err := lfs.Pivot(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(pivotedPath)
	if err != nil {
		return nil, err
	}

	return &fs.FileInfo{IsDir: fileInfo.IsDir()}, nil
}

func (lfs *Local) Get(ctx context.Context, path string) ([]byte, error) {
	pivotedPath, err := lfs.Pivot(path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(pivotedPath)
}

func (lfs *Local) ReadDir(ctx context.Context, path string) ([]string, error) {
	pivotedPath, err := lfs.Pivot(path)
	if err != nil {
		return nil, err
	}

	fileInfos, err := ioutil.ReadDir(pivotedPath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, fileInfo := range fileInfos {
		result = append(result, fileInfo.Name())
	}

	return result, nil
}

func (lfs *Local) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (lfs *Local) Pivot(path string) (string, error) {
	adaptedPath := filepath.FromSlash(path)
	return filepath.Join(lfs.cwd, adaptedPath), nil
}
