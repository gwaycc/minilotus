package version

import "github.com/gwaycc/goapp/version"

var (
	GitCommit = version.GitCommit
)
var (
	ver = "0.0.1"
)

func Version() string {
	return ver + "-" + GitCommit
}
