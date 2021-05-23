package version

var (
	gitTag      = ""
	gitCommit   = ""                     // sha1 from git, output of $(git rev-parse HEAD)
	buildStatus = "not a git tree"       // state of git tree, either "clean" or "dirty"
	buildDate   = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildUser   = ""
	buildHost   = ""
)
