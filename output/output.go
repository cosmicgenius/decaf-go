package output

import (
	"fmt"
	"io"
	"os"
)

var writer io.Writer = os.Stderr

func SetWriter(w io.Writer) {
	writer = w
}

func Writef(format string, v ...any) {
	fmt.Fprintf(writer, format, v...)
}

func Write(v ...any) {
	fmt.Fprint(writer, v...)
}

func Writeln(v ...any) {
	fmt.Fprintln(writer, v...)
}
