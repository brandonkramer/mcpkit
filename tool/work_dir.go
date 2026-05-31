package tool

import "github.com/brandonkramer/mcpkit/summarize"

// ResolveWorkDir returns argWorkDir when set, otherwise defaultWorkDir.
func ResolveWorkDir(defaultWorkDir, argWorkDir string) string {
	if argWorkDir == "" {
		return defaultWorkDir
	}
	return argWorkDir
}

// PrepareWorkDir fills an empty work_dir field from defaultWorkDir.
func PrepareWorkDir[T any](defaultWorkDir string, args T, set func(*T, string), get func(T) string) (T, summarize.Meta, error) {
	set(&args, ResolveWorkDir(defaultWorkDir, get(args)))
	return args, summarize.Meta{}, nil
}
