// The compiledef package is used to declare variables that are reassigned at compile time.
package compiledef

// the values of the following variables will be reassigned at compile time
var (
	_appVersion  = "${app-version}"
	_goVersion   = "${go-version}"
	_gitCommitID = "${git-commit-id}"
	_gitDescribe = "${git-describe}"
)

// GetAPPVersion returns appVersion.
func GetAPPVersion() string {
	return _appVersion
}

// GetGoVersion returns goVersion.
func GetGoVersion() string {
	return _goVersion
}

// GetGitCommitID returns gitCommitID.
func GetGitCommitID() string {
	return _gitCommitID
}

// GetGitDescribe returns gitDescribe.
func GetGitDescribe() string {
	return _gitDescribe
}
