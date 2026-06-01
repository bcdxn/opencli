package cliutils

import (
	"bytes"
	"io"
	"os"
)

const DefaultWidth = 80

// IOStreams is an abstraction of the std in,out,err channels used for interacting with the CLI.
// This enables injecting custom implementations and testing.
type IOStreams struct {
	In     FileReader
	Out    FileWriter
	ErrOut FileWriter
}

type FileWriter interface {
	io.Writer
	Fd() uintptr
}

type FileReader interface {
	io.ReadCloser
	Fd() uintptr
}

// System returns an IOStreams that represents the system default stdin, stdout, stderr
func System() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// Test returns an IOStreams object perfect for injecting into tests
func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	io := &IOStreams{
		In: &fdReader{
			fd:         0,
			ReadCloser: io.NopCloser(in),
		},
		Out:    &fdWriter{fd: 1, Writer: out},
		ErrOut: &fdWriter{fd: 2, Writer: errOut},
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

// fdWriteCloser represents a wrapped stdout Writer that preserves the original file descriptor
type fdWriteCloser struct {
	io.WriteCloser
	fd uintptr
}

func (w *fdWriteCloser) Fd() uintptr {
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
