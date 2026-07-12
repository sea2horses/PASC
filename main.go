package main

import (
	"flag"
	"fmt"
	"os"
	"pseint-compiled/internal/lexer"
)

var (
	debug bool
	filename string
)

const (
	PROGRAM_NAME = "Pseint Compiler"
	VERSION = "Indev"
)

type ProgramErrors string

func (pe ProgramErrors) Error() string {
	return string(pe)
}

const (
	ErrNoFilename = ProgramErrors("expected filename to compile with -file flag")
)

func dbg(a ...any) (int, error) {
	if debug {
		return fmt.Println(a...)
	}
	return 0, nil
}

func dbgfmt(format string, a ...any) (int, error) {
	if debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Program exited successfuly.")
}

func run() error {
	flag.BoolVar(&debug, "debug", false, "Whether to print debug information about the compiler")
	flag.StringVar(&filename, "file", "", "File to compile (*.psc)")

	flag.Parse();

	dbg(PROGRAM_NAME, VERSION)

	dbg("Getting file source code...")
	if filename == "" {
		return ErrNoFilename
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	dbg("Tokenizing...")

	lexer := lexer.Lexer{Src: []rune(string(bytes))}

	tokens, err := lexer.Tokenize()
	if err != nil {
		return err
	}

	dbg("Tokens:")
	for i, token := range tokens {
		dbgfmt("#%d: %+v\n", i, token)
	}

	return nil
}
