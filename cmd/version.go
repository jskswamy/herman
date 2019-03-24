package cmd

var version string
var defaultVersion = "1.0-dev"

// SetBuildVersion set the build version of the application
func SetBuildVersion(v string) {
	version = v
}

// BuildVersion return the build version
func BuildVersion() string {
	return version
}
