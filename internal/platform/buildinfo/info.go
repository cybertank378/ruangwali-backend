package buildinfo

import (
	"runtime"
	"strings"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
	builtBy   = "unknown"
)

const (
	name        = "RuangWali"
	description = "Platform SaaS manajemen dan pengawasan pendidikan"
	author      = "Rahman"
)

type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildTime   string `json:"buildTime"`
	BuiltBy     string `json:"builtBy"`
	GoVersion   string `json:"goVersion"`
}

func Current() Info {
	return Info{
		Name:        name,
		Description: description,
		Author:      author,
		Version:     normalize(version),
		Commit:      normalize(commit),
		BuildTime:   normalize(buildTime),
		BuiltBy:     normalize(builtBy),
		GoVersion:   runtime.Version(),
	}
}

func normalize(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return "unknown"
	}

	return value
}
