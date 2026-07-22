package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/generator"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/parser"
)

var (
	debug    bool
	filename string
)

/* TODO: Add .exe at the end for Windows */
const (
	PROGRAM_NAME = "Pseint Compiler"
	VERSION      = "Indev"
	OUTPUT       = "./build/"
	CG_FILENAME  = "out.go"
)

//go:embed runtime.go
var runtime []byte

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

	dbg("Program exited successfuly")
}

func run() error {
	flag.BoolVar(&debug, "debug", false, "Whether to print debug information about the compiler")
	flag.StringVar(&filename, "file", "", "File to compile (*.psc)")

	flag.Parse()

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

	src := []rune(string(bytes))
	lexer := lexer.NewLexer(src)

	tokens, err := lexer.Tokenize()
	ds := lexer.Diagnostics()

	for i, d := range ds {
		fm, err := diagnostics.FormatDiagnostic(filename, bytes, d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "<err #%d could not be shown (%s)>", i+1, err.Error())
		}
		fmt.Fprintln(os.Stderr, fm)
	}

	e, w, _ := diagnostics.Summary(ds)
	if e > 0 || w > 0 {
		fmt.Fprintf(os.Stderr, "(%d errors, %d warnings)\n", e, w)
	}

	if e > 0 {
		return diagnostics.SummaryAsError(ds)
	}

	if err != nil {
		return err
	}

	dbg("Tokens:")
	for i, token := range tokens {
		dbgfmt("#%d: %+v\n", i, token)
	}

	dbg("AST Parsing...")
	parser := parser.Parser{Tokens: tokens}
	ast, err := parser.Parse()
	if err != nil {
		return err
	}
	dbg(ast.String())

	dbg("Generating code...")
	codegen := generator.CodeGenerator{}

	code_bytes := codegen.Generate(ast)
	var data []byte = append(runtime, code_bytes...)

	dbg("Generated data\n", string(data))

	// Code path
	code_path := OUTPUT + CG_FILENAME
	out_path := OUTPUT + "out"

	if err = os.MkdirAll(OUTPUT, 0755); err != nil {
		return err
	}

	err = os.WriteFile(code_path, data, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"bash",
		"-c",
		"go build -o \"$1\" \"$2\"",
		"bash",
		out_path,
		code_path,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return err
	}

	run := exec.Command(out_path)
	run.Stderr = os.Stderr
	run.Stdout = os.Stdout
	run.Stdin = os.Stdin

	err = run.Run()
	if err != nil {
		return err
	}

	return nil
}
