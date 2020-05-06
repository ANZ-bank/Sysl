package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/language/go/pkg/codegen"
	"github.com/anz-bank/sysl/language/go/pkg/relgom"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

func main() {
	root := os.Args[1]
	filename := os.Args[2]
	modelName := os.Args[3]
	out := os.Args[4]

	fsw := &codegen.OSFileSystemWriter{Root: out}
	if err := relgom.Generate(fsw, syslutil.NewChrootFs(afero.NewOsFs(), root), filename, modelName); err != nil {
		var buf bytes.Buffer

		buf.WriteString(errors.Cause(err).Error())

		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		if err, ok := err.(stackTracer); ok {
			for _, f := range err.StackTrace() {
				fmt.Fprintf(&buf, "%+s:%d\n", f, f)
			}
		} else {
			logrus.Warn("(Error stack trace missing)")
		}

		logrus.Fatal(err.Error() + "\n" + buf.String())
	}
}
