package relgom

import (
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

const pkRecv = "k"

func (g *entityGenerator) goAppendPKDecls(decls []Decl) []Decl {
	if g.haveKeys {
		decls = append(decls,
			g.goPKTypeDecl(),
			g.goPKSetableHash(),
			g.goPKSetableEqual(),
		)
	}
	return decls
}

func (g *entityGenerator) goPKTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.pkName),
		Type: Struct(g.pkFields()...),
	}).WithDoc(Commentf("// %s is the Key for %s.", g.pkName, g.tname))
}

func (g *entityGenerator) pkFields() []Field {
	return g.goFieldsForSyslAttrDefs(g.isPkeyAttr, false, false, nil)
}

func (g *entityGenerator) pkMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(pkRecv),
		Type:  I(g.pkName),
	}).Parens()
	return &f
}

func (g *entityGenerator) goPKSetableHash() Decl {
	fields := g.pkFields()
	stmts := make([]Stmt, 0, len(fields)+1)
	for _, field := range fields {
		stmts = append(stmts, Assign(I("seed"))("=")(
			Call(g.hash("Interface"),
				Dot(I(pkRecv), field.Names[0].Name.Text),
				I("seed"))))
	}
	stmts = append(stmts, Return(I("seed")))
	return g.pkMethod(FuncDecl{
		Name: *I("Hash"),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("seed"), Type: I("uintptr")}),
			Results: Fields(Field{Type: I("uintptr")}),
		},
		Body: Block(stmts...),
	})
}

func (g *entityGenerator) goPKSetableEqual() Decl {
	return g.pkMethod(FuncDecl{
		Name: *I("Equal"),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("i"), Type: Composite(I("interface"))}),
			Results: Fields(Field{Type: I("bool")}),
		},
		Body: Block(
			If(
				Init("l", "ok")(Assert(I("i"), I(g.pkName))),
				I("ok"),
				Return(Binary(I(pkRecv), "==", I("l"))),
			),
			Return(I("false")),
		),
	})
}
