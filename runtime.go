package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	ErrIncorrectType = "No coinciden los tipos"
)

/* I/O */
var writer bufio.Writer = *bufio.NewWriter(os.Stdout);
var reader bufio.Reader = *bufio.NewReader(os.Stdin);

type Value interface {
	int64 | float64 | string | bool
}

func RuntimeError(msg string) {
	panic(msg)
}

func Write(values... any) {
	for _, value := range values {
		switch val := value.(type) {
		case int64:
			writer.WriteString(strconv.FormatInt(val, 10))
		case string:
			writer.WriteString(val)
		case float64:
			writer.WriteString(strconv.FormatFloat(val, 'f', 2, 64))
		case bool:
			if val {
				writer.WriteString("VERDADERO")
			} else {
				writer.WriteString("FALSO")
			}
		}
		writer.WriteRune('\n')
	}
	writer.Flush()
}

func Read[T Value](dest *T) {
	writer.WriteString("> ")
	writer.Flush()

	/* TODO: Catch error from ReadLine */
	input, _, _ := reader.ReadLine()

	switch d := any(dest).(type) {
	case *int64:
		result, err := strconv.ParseInt(string(input), 10, 64)
		if err != nil {
			RuntimeError(ErrIncorrectType)
		}
		*d = result
	case *string:
		*d = string(input)
	case *float64:
		result, err := strconv.ParseFloat(string(input), 64)
		if err != nil {
			RuntimeError(ErrIncorrectType)
		}
		*d = result
	case *bool:
		*d = func() bool {
			/* Pseint takes false when no input is given */
			if len(input) == 0 {
				return false
			}

			/* Try to match input to other strings */
			input := strings.ToLower(string(input))
			if input == "verdadero" {
				return true
			}

			if input == "falso" {
				return false
			}

			/* Finally, match to 0 or 1 */
			result, err := strconv.ParseInt(input, 10, 64)
			if err == nil {
				switch result {
				case 0:
					return false
				case 1:
					return true
				}
			}

			RuntimeError(ErrIncorrectType)
			/* UNREACHABLE! */
			return false
		}()
	}
}
