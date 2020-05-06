package relgom

import (
	"github.com/anz-bank/sysl/pkg/syslutil"
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

const tupleDataRecv = "d"

func (g *entityGenerator) goAppendTupleDataDecls(decls []Decl) []Decl {
	return append(decls,
		g.goTupleDataTypeDecl(),
		g.marshalTupleDataJSONFunc(),
		g.unmarshalTupleDataJSONFunc(),
	)
}

func (g *entityGenerator) goTupleDataTypeDecl() Decl {
	fields := []Field{}
	if g.haveKeys {
		fields = append(fields, Field{Type: I(g.pkName)})
	}
	fields = append(fields,
		g.goFieldsForSyslAttrDefs(syslutil.NamedTypeNot(g.isPkeyAttr), false, false, nil)...,
	)
	return Types(TypeSpec{
		Name: *I(g.dataName),
		Type: Struct(fields...),
	}).WithDoc(Commentf("// %s is the internal representation of a tuple in the model.", g.dataName))
}

func (g *entityGenerator) tupleDataMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(tupleDataRecv),
		Type:  Star(I(g.dataName)),
	}).Parens()
	return &f
}

func (g *entityGenerator) marshalTupleDataJSONFunc() Decl {
	kvs := []Expr{}
	for _, nt := range g.namedAttrs {
		ti := g.typeInfoForSyslType(nt.Type)
		kvs = append(kvs, &KeyValueExpr{
			Key:   ExportedID(nt.Name),
			Value: ti.stage(Dot(I(tupleDataRecv), NonExportedName(nt.Name))),
		})
	}

	return g.tupleDataMethod(*marshalJSONMethodDecl(
		Return(Call(g.json("Marshal"), Composite(g.goExportedStruct(true), kvs...))),
	))
}

func (g *entityGenerator) unmarshalTupleDataJSONFunc() Decl {
	tmp := "u"

	pkKVs := []Expr{}
	kvs := []Expr{}
	stagings := []Stmt{}
	if g.haveKeys {
		kvs = append(kvs, nil)
	}

	for _, nt := range g.namedAttrs {
		ti := g.typeInfoForSyslType(nt.Type)
		ename := ExportedName(nt.Name)
		var value Expr = Dot(I(tmp), ename)
		if ti.unstage != nil {
			u := "unstage" + ename
			stagings = append(stagings,
				Init(u, "err")(ti.unstage(value)),
				If(nil, Binary(I("err"), "!=", I("nil")),
					Return(Call(g.imported("fmt")("Errorf"),
						String("error unstaging %s.%s: %v"), String(g.tname), String(nt.Name), I("err"),
					)),
				),
			)
			value = I(u)
		}
		if g.isPkeyAttr(nt) {
			pkKVs = append(pkKVs, &KeyValueExpr{
				Key:   NonExportedID(nt.Name),
				Value: value,
			})
		} else {
			kvs = append(kvs, &KeyValueExpr{
				Key:   NonExportedID(nt.Name),
				Value: value,
			})
		}
	}
	if g.haveKeys {
		kvs[0] = &KeyValueExpr{
			Key:   I(g.pkName),
			Value: Composite(I(g.pkName), pkKVs...),
		}
	}

	stmts := []Stmt{
		&DeclStmt{Decl: Var(ValueSpec{
			Names: Idents(tmp),
			Type:  g.goExportedStruct(true),
		})},
		If(
			Init("err")(Call(g.json("Unmarshal"), I("data"), AddrOf(I(tmp)))),
			Binary(I("err"), "!=", Nil()),
			Return(I("err")),
		),
	}
	stmts = append(stmts, stagings...)
	stmts = append(stmts,
		Assign(Star(I(tupleDataRecv)))("=")(Composite(I(g.dataName), kvs...)),
		Return(Nil()),
	)

	return g.tupleDataMethod(*unmarshalJSONMethodDecl(stmts...))
}

func (g *entityGenerator) goExportedStruct(staging bool) *StructType {
	return Struct(g.goFieldsForSyslAttrDefs(
		syslutil.NamedTypeAll,
		true,
		staging,
		func(nt syslutil.NamedType) map[string]string {
			return map[string]string{"json": nt.Name + ",omitempty"}
		},
	)...)
}
