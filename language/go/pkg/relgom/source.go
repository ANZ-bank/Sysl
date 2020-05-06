package relgom

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/language/go/pkg/codegen"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

//nolint:gochecknoglobals
var forceOptional = func() bool {
	forceOptional := os.Getenv("RELGOM_ALL_FIELDS_OPTIONAL") == "1"
	if forceOptional {
		logrus.Warnf("All generated fields will be marked optional.")
	}
	return forceOptional
}()

type commonModules struct {
	json      func(name string) Expr
	relgomlib func(name string) Expr
	frozen    func(name string) Expr
	hash      func(name string) Expr
}

func newCommonModules(g *sourceGenerator) *commonModules {
	return &commonModules{
		json:      g.imported("encoding/json"),
		relgomlib: g.imported("github.com/anz-bank/sysl/language/go/pkg/relgom/relgomlib"),
		frozen:    g.imported("github.com/arr-ai/frozen"),
		hash:      g.imported("github.com/arr-ai/hash"),
	}
}

type sourceGenerator struct {
	fsw             codegen.FileSystemWriter
	packageName     string
	model           *sysl.Application
	builtinImports  map[string]struct{}
	externalImports map[string]struct{}
}

func newSourceGenerator(fsw codegen.FileSystemWriter, packageName string, model *sysl.Application) *sourceGenerator {
	return &sourceGenerator{
		fsw:             fsw,
		packageName:     packageName,
		model:           model,
		builtinImports:  map[string]struct{}{},
		externalImports: map[string]struct{}{},
	}
}

func (g *sourceGenerator) genSourceForDecls(
	basepath string, decls ...Decl,
) error {
	return g.genSourceForFile(basepath, &File{Decls: decls})
}

func (g *sourceGenerator) genSourceForFile(basepath string, file *File) error {
	{
		file := *file
		file.Doc = codegen.Prelude()
		file.Name = *I(g.packageName)
		file.Imports = ImportGroups(
			Imports(sortedSetElements(g.builtinImports)...),
			Imports(sortedSetElements(g.externalImports)...),
		)

		var buf bytes.Buffer
		fmt.Fprintln(&buf, &file)

		final, err := format.Source(buf.Bytes())
		if err != nil {
			var lines bytes.Buffer
			for i, line := range strings.Split(buf.String(), "\n") {
				fmt.Fprintf(&lines, "%3d: %s\n", 1+i, line)
			}
			logrus.Errorf("Error formatting %#v:\n%s", basepath+".go", lines.String())
			return errors.Wrap(err, "gofmt")
		}

		f, err := g.fsw.Create(basepath + ".go")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(final)
		return err
	}
}

func (g *sourceGenerator) imported(imp string, extra ...string) func(name string) Expr { //nolint:unparam
	var pkgName string
	switch len(extra) {
	case 0:
		pkgName = imp[strings.LastIndex(imp, "/")+1:]
	case 1:
		pkgName = extra[0]
	default:
		panic("Too many parameters to imported()")
	}

	var imports map[string]struct{}
	if strings.Contains(imp, ".") {
		imports = g.externalImports
	} else {
		imports = g.builtinImports
	}

	return func(name string) Expr {
		imports[imp] = struct{}{}
		return Dot(I(pkgName), name)
	}
}

type typeInfo struct {
	final   Expr
	param   Expr
	staging Expr
	stage   func(value Expr) Expr
	unstage func(value Expr) Expr
	opt     bool
	fkey    *sysl.ScopedRef
}

func (g *sourceGenerator) typeInfoForSyslType(t *sysl.Type) *typeInfo {
	var ti typeInfo
	ti.stage = func(value Expr) Expr { return value }

	switch t := t.Type.(type) {
	case *sysl.Type_Primitive_:
		switch t.Primitive {
		case sysl.Type_BOOL:
			ti.final = I("bool")
		case sysl.Type_INT:
			ti.final = I("int64")
		case sysl.Type_FLOAT:
			ti.final = I("float64")
		case sysl.Type_DECIMAL:
			ti.final = g.imported("github.com/anz-bank/decimal")("Decimal64")
		case sysl.Type_DATE, sysl.Type_DATETIME:
			relgomlib := g.imported("github.com/anz-bank/sysl/language/go/pkg/relgom/relgomlib")
			ti.final = g.imported("time")("Time")
			ti.staging = relgomlib("DateTimeString")
			ti.stage = func(value Expr) Expr {
				return Call(relgomlib("NewDateTimeString"), value)
			}
			ti.unstage = func(value Expr) Expr {
				return Call(Dot(value, "Unstage"))
			}
		case sysl.Type_STRING, sysl.Type_STRING_8:
			ti.final = I("string")
		default:
			panic(fmt.Errorf("type: %#v", t))
		}
	case *sysl.Type_TypeRef:
		fkInfo := g.typeInfoForSyslType(g.typeForScopedRef(t.TypeRef))
		ti = *fkInfo
		ti.fkey = t.TypeRef
	default:
		panic(fmt.Errorf("type: %#v", t))
	}
	ti.param = ti.final
	patterns := syslutil.MakeStrSetFromAttr("patterns", t.Attrs)
	if (t.Opt || forceOptional) && !patterns.Contains("pk") {
		ti.opt = true
		ti.final = Star(ti.final)
		if ti.staging != nil {
			ti.staging = Star(ti.staging)
		}
	}
	return &ti
}

func (g *sourceGenerator) typeForScopedRef(t *sysl.ScopedRef) *sysl.Type {
	if len(t.GetRef().GetAppname().GetPart()) > 0 {
		panic(fmt.Errorf("non-local refs not implemented: %#v", t.Ref))
	}
	if len(t.Ref.Path) != 2 {
		panic(fmt.Errorf("ScopedRef path must be length 2: %#v", t.Ref)) //nolint:golint
	}
	return g.model.Types[t.Ref.Path[0]].Type.(*sysl.Type_Relation_).Relation.AttrDefs[t.Ref.Path[1]]
}
