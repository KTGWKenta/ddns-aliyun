package version

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/kentalee/log"
)

type Info struct {
	GitTag      string `json:"gitTag"`
	GitCommit   string `json:"gitCommit"`
	BuildStatus string `json:"buildStatus"`
	BuildDate   string `json:"buildDate"`
	BuildUser   string `json:"buildUser"`
	BuildHost   string `json:"buildHost"`
	GoVersion   string `json:"goVersion"`
	Compiler    string `json:"compiler"`
	Platform    string `json:"platform"`
}

func (info Info) String() string {
	return info.GitTag
}

func Get() Info {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		if gitCommit == "" {
			gitCommit = buildInfo.Main.Sum
		}
		if gitTag == "" {
			gitTag = buildInfo.Main.Version
		}
	}
	return Info{
		GitTag:      gitTag,
		GitCommit:   gitCommit,
		BuildStatus: buildStatus,
		BuildDate:   buildDate,
		BuildUser:   buildUser,
		BuildHost:   buildHost,
		GoVersion:   runtime.Version(),
		Compiler:    runtime.Compiler,
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func Print() {
	info := Get()
	log.Info("GitTag:", info.GitTag)
	log.Info("GitCommit:", info.GitCommit)
	log.Info("BuildStatus:", info.BuildStatus)
	log.Info("BuildDate:", info.BuildDate)
	log.Info("BuildUser:", info.BuildUser)
	log.Info("BuildHost:", info.BuildHost)
	log.Info("GoVersion:", info.GoVersion)
	log.Info("Compiler:", info.Compiler)
	log.Info("Platform:", info.Platform)
}
