package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

var whiteColor *color.Color
var redColor *color.Color
var yellowColor *color.Color
var greenColor *color.Color

func init() {
	whiteColor = color.New(color.FgWhite)
	redColor = color.New(color.FgRed)
	yellowColor = color.New(color.FgYellow)
	greenColor = color.New(color.FgGreen)
}

// DieIf invoke os.Exit if there err is not nil
func DieIf(err error) {
	if err != nil {
		Fatalf("\nCommand failed: %v", err)
	}
}

// RenderAsYaml convert value as yaml and render into writer
func RenderAsYaml(writer io.Writer, value interface{}) {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return
	}
	content := string(bytes)
	_, _ = fmt.Fprint(writer, content)
}

// PrintYaml convert value as yaml and render into OutputStream()
func PrintYaml(value interface{}) {
	Info("---")
	RenderAsYaml(OutputStream(), value)
}

// Info print values to OutputStream() in white color
func Info(value ...interface{}) {
	stream := OutputStream()
	_, _ = whiteColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Info print values to OutputStream() in white color
func Success(value ...interface{}) {
	stream := OutputStream()
	_, _ = greenColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Error print values to os.Stderr in red color
func Error(value ...interface{}) {
	stream := os.Stderr
	_, _ = redColor.Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Warn print values to OutputStream() in yellow color
func Warn(value ...interface{}) {
	stream := OutputStream()
	_, _ = yellowColor.Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Fatal print values to os.Stderr in red color
// and invoke os.Exit with non zero status code
func Fatal(value ...interface{}) {
	stream := os.Stderr
	_, _ = redColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
	os.Exit(1)
}

// Infof format and print values to OutputStream() in white color
func Infof(format string, value ...interface{}) {
	Info(fmt.Sprintf(format, value...))
}

// Errorf format and print values to os.Stderr in red color
func Errorf(format string, value ...interface{}) {
	Error(fmt.Sprintf(format, value...))
}

// Warnf format and print values to OutputStream() in yellow color
func Warnf(format string, value ...interface{}) {
	Warn(fmt.Sprintf(format, value...))
}

// Fatalf format and print values to os.Stderr in red color
// and invoke os.Exit with non zero status code
func Fatalf(format string, value ...interface{}) {
	Fatal(fmt.Sprintf(format, value...))
}

// IsTTY returns true if the terminal is an interactive tty
// and false if not
func IsTTY() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

// ColorForced returns true if "FORCE_COLOR" is enabled
// Eg:
// export FORCE_COLOR=true
// (or)
// export FORCE_COLOR=1
func ColorForced() bool {
	colorForced := os.Getenv("FORCE_COLOR")
	if colorForced == "true" || colorForced == "1" {
		return true
	}

	return false
}

// OutputStream gives io.Writer based on interactive terminal
func OutputStream() io.Writer {
	if IsTTY() {
		return os.Stdout
	}
	return os.Stderr
}
