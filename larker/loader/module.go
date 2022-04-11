package loader

import (
	"context"
	"errors"

	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker/fs"
)

type localLocation struct {
	Path string
}

var (
	ErrUnsupportedLocation = errors.New("unsupported location")
)

func parseLocation(module string) interface{} {
	return localLocation{Path: module}
}

func findModuleFS(
	ctx context.Context,
	currentFS fs.FileSystem,
	env map[string]string,
	module string,
) (fs.FileSystem, string, error) {
	return findLocatorFS(ctx, currentFS, env, parseLocation(module))
}

func findLocatorFS(
	ctx context.Context,
	currentFS fs.FileSystem,
	env map[string]string,
	location interface{},
) (fs.FileSystem, string, error) {
	switch l := location.(type) {
	case localLocation:
		return currentFS, l.Path, nil
	default:
		return nil, "", ErrUnsupportedLocation
	}
}
