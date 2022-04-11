package loader

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

var (
	ErrChangesInclude     = errors.New("changes_include() failed")
	ErrChangesIncludeOnly = errors.New("changes_include_only() failed")
)

func starlarkArgsToStrings(args starlark.Tuple, kwargs []starlark.Tuple) ([]string, error) {
	var result []string

	if len(kwargs) != 0 {
		return nil, fmt.Errorf("%w: found %d keyword arguments, expected 0", ErrChangesInclude, len(kwargs))
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("%w: expected at least 1 positional argument, found 0", ErrChangesInclude)
	}

	for i, arg := range args {
		stringArgument, ok := arg.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("%w: expected %d'th argument to be string, got %s", ErrChangesInclude, i+1,
				arg.Type())
		}

		result = append(result, stringArgument.GoString())
	}

	return result, nil
}
