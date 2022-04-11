// adapted from https://github.com/cirruslabs/cirrus-cli/pull/46/files#diff-fd961f8f67870410b5925d977385825c70a2da811309b101a719b383ee2d8a04
package larker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker/fs"
	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker/fs/local"
	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker/loader"
	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker/yamlhelper"
	slJSON "go.starlark.net/lib/json"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"gopkg.in/yaml.v3"
)

var (
	ErrLoadFailed           = errors.New("load failed")
	ErrNotFound             = errors.New("entrypoint not found")
	ErrMainFailed           = errors.New("failed to call main")
	ErrExecFailed           = errors.New("exec failed")
	ErrHookFailed           = errors.New("failed to call hook")
	ErrMainUnexpectedResult = errors.New("main returned unexpected result")
	ErrSanity               = errors.New("sanity check failed")
)

type Larker struct {
	fs            fs.FileSystem
	env           map[string]string
	affectedFiles []string
	isTest        bool
}

type MainResult struct {
	OutputLogs []byte
	YAMLConfig string
}

func New() *Larker {

	lrk := &Larker{
		env: make(map[string]string),
		fs:  local.New("."),
	}

	// weird global init by Starlark
	// we need floats at least for configuring CPUs for containers
	resolve.AllowFloat = true

	return lrk
}

func (larker *Larker) Main(ctx context.Context, inputFile string) (*MainResult, error) {
	source, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("could not read starlark file %s: %w", inputFile, err)
	}

	outputLogsBuffer := &bytes.Buffer{}
	capture := func(thread *starlark.Thread, msg string) {
		_, _ = fmt.Fprintln(outputLogsBuffer, msg)
	}

	thread := &starlark.Thread{
		Load:  loader.NewLoader(ctx, larker.fs, larker.env, larker.affectedFiles, larker.isTest).LoadFunc(larker.fs),
		Print: capture,
	}

	resCh := make(chan starlark.Value)
	errCh := make(chan error)

	// functions available to starlark
	predeclared := starlark.StringDict{
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
		"json":   slJSON.Module,
	}

	go func() {
		// Execute the source code for the main() to be visible
		globals, err := starlark.ExecFile(thread, inputFile, source, predeclared)
		if err != nil {
			errCh <- fmt.Errorf("%w: %v", ErrLoadFailed, err)
			return
		}

		main, ok := globals["module"]
		if !ok {
			errCh <- fmt.Errorf("%w: module()", ErrNotFound)
			return
		}

		// Ensure that main() is a function
		mainFunc, ok := main.(*starlark.Function)
		if !ok {
			errCh <- fmt.Errorf("%w: module is not a function", ErrMainFailed)
			return
		}

		var args starlark.Tuple

		// Prepare a context to pass to main() as it's first argument if needed
		if mainFunc.NumParams() != 0 {
			args = append(args, &Context{})
		}

		mainResult, err := starlark.Call(thread, main, args, nil)
		if err != nil {
			errCh <- fmt.Errorf("error occurred exeucting starlark: %w", err)
			return
		}

		resCh <- mainResult
	}()

	var mainResult starlark.Value

	select {
	case mainResult = <-resCh:
	case err := <-errCh:
		return nil, logsWithErrorAttachedErr(outputLogsBuffer.Bytes(), err)
	case <-ctx.Done():
		thread.Cancel(ctx.Err().Error())
		return nil, ctx.Err()
	}

	var tasksNode *yaml.Node

	switch typedMainResult := mainResult.(type) {
	case *starlark.List:
		tasksNode = convertList(typedMainResult)
		if err != nil {
			return nil, err
		}
		if tasksNode == nil {
			return &MainResult{OutputLogs: outputLogsBuffer.Bytes()}, nil
		}
	case *starlark.Dict:
		tasksNode = convertDict(typedMainResult)
		if tasksNode == nil {
			return &MainResult{OutputLogs: outputLogsBuffer.Bytes()}, nil
		}
	default:
		return nil, fmt.Errorf("%w: result is not a list or a dict", ErrMainUnexpectedResult)
	}

	formattedYaml, err := yamlhelper.PrettyPrint(tasksNode)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot marshal into YAML: %v", ErrMainUnexpectedResult, err)
	}

	return &MainResult{
		OutputLogs: outputLogsBuffer.Bytes(),
		YAMLConfig: formattedYaml,
	}, nil
}

func logsWithErrorAttached(logs []byte, err error) []byte {
	fmt.Printf("%T\n", err)

	ee, ok := errors.Unwrap(err).(*starlark.EvalError)
	if !ok {
		return logs
	}

	if len(logs) != 0 && !bytes.HasSuffix(logs, []byte("\n")) {
		logs = append(logs, byte('\n'))
	}

	logs = append(logs, []byte(ee.Backtrace())...)

	return logs
}

func logsWithErrorAttachedErr(logs []byte, err error) error {
	return fmt.Errorf("Starlark error: %w\n%v", err, logs)
}
