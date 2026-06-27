package gencli

import (
	"bytes"
	"io"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"golang.org/x/term"
)

type IOStreams interface {
	In() FileReader
	Out() FileWriter
	ErrOut() FileWriter
	TerminalTheme() string
	TerminalSize() (int, int, error)
}

type FileWriter interface {
	io.Writer
	Fd() uintptr
}

type FileReader interface {
	io.ReadCloser
	Fd() uintptr
}

type OcliCheckArgs struct {
	PathToSpec string
}

type OcliCheckFlags struct {
	FailOnErr bool
}

type DefaultIOStreams struct {
	in     FileReader
	out    FileWriter
	errOut FileWriter
	term   Terminal
}

func (ios DefaultIOStreams) In() FileReader {
	return ios.in
}

func (ios DefaultIOStreams) Out() FileWriter {
	return ios.out
}

func (ios DefaultIOStreams) ErrOut() FileWriter {
	return ios.errOut
}

func (ios DefaultIOStreams) TerminalTheme() string {
	return ios.term.Theme()
}

func (ios DefaultIOStreams) TerminalSize() (int, int, error) {
	return ios.term.Size()
}

func DefaultIOS() IOStreams {
	return DefaultIOStreams{
		in:     os.Stdin,
		out:    os.Stdout,
		errOut: os.Stderr,
		term: &DefaultTerminal{
			in:           os.Stdin,
			out:          os.Stdout,
			errOut:       os.Stderr,
			isTTY:        isTerminal(os.Stdout),
			is256enabled: is256ColorSupported(),
			hasTrueColor: isTrueColorSupported(),
		},
	}
}

type Terminal interface {
	IsTerminalOutput() bool
	Is256ColorSupported() bool
	IsTrueColorSupported() bool
	Theme() string
	Size() (int, int, error)
}

type DefaultTerminal struct {
	in           *os.File
	out          *os.File
	errOut       *os.File
	isTTY        bool
	colorEnabled bool
	is256enabled bool
	hasTrueColor bool
	width        int
}

func (t *DefaultTerminal) IsTerminalOutput() bool {
	return t.isTTY
}

func (t *DefaultTerminal) Is256ColorSupported() bool {
	return t.is256enabled
}

func (t *DefaultTerminal) IsTrueColorSupported() bool {
	return t.hasTrueColor
}

func (t *DefaultTerminal) Theme() string {
	if lipgloss.HasDarkBackground(t.in, t.out) {
		return "dark"
	}

	return "light"
}

func (t *DefaultTerminal) Size() (int, int, error) {
	ttyOut := t.out
	if ttyOut == nil || !isTerminal(ttyOut) {
		if f, err := openTTY(); err == nil {
			defer f.Close()
			ttyOut = f
		} else {
			return -1, -1, err
		}
	}

	width, height, err := terminalSize(ttyOut)
	return width, height, err
}

/* Types to make testing easier
------------------------------------------------------------------------------------------------- */

func TestIOS() (IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	io := DefaultIOStreams{
		in: &fdReader{
			fd:         0,
			ReadCloser: io.NopCloser(in),
		},
		out:    &fdWriter{fd: 1, Writer: out},
		errOut: &fdWriter{fd: 2, Writer: errOut},
		// term   Terminal
	}

	return io, in, out, errOut
}

// fdWriter represents a wrapped stdout Writer that preserves the original file descriptor
type fdWriter struct {
	io.Writer
	fd uintptr
}

func (w *fdWriter) Fd() uintptr {
	return w.fd
}

// fdReader represents a wrapped stdin ReadCloser that preserves the original file descriptor
type fdReader struct {
	io.ReadCloser
	fd uintptr
}

func (r *fdReader) Fd() uintptr {
	return r.fd
}

// type TestIOStreams struct{}
// type TestTerminal struct {
// 	in           *os.File
// 	out          *os.File
// 	errOut       *os.File
// 	isTTY        bool
// 	colorEnabled bool
// 	is256enabled bool
// 	hasTrueColor bool
// 	width        int
// }

// func (t *TestTerminal) IsTerminalOutput() bool {
// 	return t.isTTY
// }

// func (t *TestTerminal) Is256ColorSupported() bool {
// 	return t.is256enabled
// }

// func (t *TestTerminal) IsTrueColorSupported() bool {
// 	return t.hasTrueColor
// }

// func (t *TestTerminal) Theme() string {
// 	return "dark"
// }

// func (t *TestTerminal) Size() (int, int, error) {
// 	return 80, 100, nil
// }

/* Private helper functions
------------------------------------------------------------------------------------------------- */

// isTerminal reports whether a file descriptor is connected to a terminal.
func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func is256ColorSupported() bool {
	return isTrueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

func isTrueColorSupported() bool {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	return strings.Contains(term, "24bit") ||
		strings.Contains(term, "truecolor") ||
		strings.Contains(colorterm, "24bit") ||
		strings.Contains(colorterm, "truecolor")
}

func openTTY() (*os.File, error) {
	return os.Open("/dev/tty")
}

func terminalSize(f *os.File) (int, int, error) {
	return term.GetSize(int(f.Fd()))
}
