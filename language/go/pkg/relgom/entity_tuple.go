package relgom

import (
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

const tupleRecv = "t"

func (g *entityGenerator) goAppendTupleDecls(decls []Decl) []Decl {
	decls = append(decls,
		g.goTupleTypeDecl(),
	)
	for _, nt := range g.namedAttrs {
		decls = append(decls,
			g.goGetterFuncForSyslAttr(nt.Name, nt.Type),
		)
	}
	return decls
}

func (g *entityGenerator) goTupleTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.tname),
		Type: Struct(
			Field{Type: Star(I(g.dataName))},
			Field{Names: Idents("model"), Type: I(g.modelName)},
		),
	}).WithDoc(Commentf("// %s is the public representation tuple in the model.", g.tname))
}

func (g *entityGenerator) tupleMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(tupleRecv),
		Type:  I(g.tname),
	}).Parens()
	return &f
}

func (g *entityGenerator) goGetterFuncForSyslAttr(attrName string, attr *sysl.Type) Decl {
	exp := ExportedID(attrName)
	nexp := NonExportedID(attrName)
	var field Expr = Dot(I(tupleRecv), nexp.Name.Text)
	typeInfo := g.typeInfoForSyslType(attr)
	if typeInfo.fkey == nil {
		return g.tupleMethod(FuncDecl{
			Doc:  Comments(Commentf("// %s gets the %s attribute from the %s.", exp.Name.Text, attrName, g.tname)),
			Name: *exp,
			Type: FuncType{
				Results: Fields(Field{Type: g.typeInfoForSyslType(attr).final}),
			},
			Body: Block(
				Return(field),
			),
		})
	}

	// FK special case
	fpath := typeInfo.fkey.Ref.Path
	ambiguous := g.fkCount[strings.Join(fpath, ".")] > 1
	relation := ExportedID(fpath[0])
	exp2 := ExportedID(fpath[0])
	if ambiguous {
		exp2.Name.Text += "Via" + exp.Name.Text
	}
	if typeInfo.opt {
		field = Star(field)
	}
	return g.tupleMethod(FuncDecl{
		Doc: Comments(Commentf(
			"// %s gets the %s corresponding to the %s attribute from t.",
			exp2.Name.Text, fpath[0], attrName,
		)),
		Name: *exp2,
		Type: FuncType{
			Results: Fields(Field{Type: relation}),
		},
		Body: Block(
			Init("u", "_")(Call(
				Dot(Call(Dot(I(tupleRecv), "model", "Get"+relation.Name.Text)), "Lookup"),
				field,
			)),
			Return(I("u")),
		),
	})
}
