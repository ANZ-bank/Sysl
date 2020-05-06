package relgom

import (
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

const builderRecv = "b"

func (g *entityGenerator) goAppendBuilderDecls(decls []Decl) []Decl {
	decls = append(decls,
		g.goBuilderTypeDecl(),
	)
	for i, nt := range g.namedAttrs {
		if !g.autoinc.Contains(nt.Name) {
			decls = append(decls, g.goBuilderSetterFuncForSyslAttr(i, nt.Name, nt.Type))
		}
	}
	return append(decls,
		g.goEntityTypeStaticMetadataDecl(),
		g.goBuilderApplyDecl(),
	)
}

func (g *entityGenerator) goBuilderTypeDecl() Decl {
	s := Struct(
		Field{Type: I(g.dataName)},
		Field{Names: Idents("model"), Type: I(g.modelName)},
		Field{Names: Idents("mask"), Type: ArrayN((g.nAttrs+63)/64, I("uint64"))},
		Field{Names: Idents("apply"), Type: g.applyFuncType()},
	)
	return Types(
		TypeSpec{Name: *ExportedID(g.builderName), Type: s},
	).WithDoc(Commentf("// %s builds an instance of %s in the model.", g.builderName, g.tname))
}

func (g *entityGenerator) builderMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(builderRecv),
		Type:  Star(I(g.builderName)),
	}).Parens()
	return &f
}

func (g *entityGenerator) goBuilderSetterFuncForSyslAttr(i int, attrName string, attr *sysl.Type) Decl {
	exp := "With" + ExportedName(attrName)
	nexp := NonExportedName(attrName)
	typeInfo := g.typeInfoForSyslType(attr)
	updateMask := CallStmt(g.relgomlib("UpdateMaskForFieldButPanicIfAlreadySet"),
		AddrOf(maskForField(i)),
		Binary(Call(I("uint64"), Int(1)), "<<", Int(i%64)),
	)
	if typeInfo.fkey == nil {
		var value Expr = I("value")
		if g.requiredMask[i/64]&(uint64(1)<<uint(i%64)) == 0 {
			value = AddrOf(value)
		}
		return g.builderMethod(FuncDecl{
			Doc:  Comments(Commentf("// %s sets the %s attribute of the %s.", exp, attrName, g.builderName)),
			Name: *I(exp),
			Type: FuncType{
				Params:  *Fields(Field{Names: Idents("value"), Type: typeInfo.param}),
				Results: Fields(Field{Type: Star(I(g.builderName))}),
			},
			Body: Block(
				updateMask,
				Assign(Dot(I(builderRecv), nexp))("=")(value),
				Return(I(builderRecv)),
			),
		})
	}
	fpath := typeInfo.fkey.Ref.Path
	ambiguous := g.fkCount[strings.Join(fpath, ".")] > 1
	// relation := ExportedID(fpath[0])
	exp2 := "With" + ExportedName(fpath[0])
	if ambiguous {
		exp2 += "For" + ExportedName(attrName)
	}
	var field Expr = Dot(I("t"), NonExportedName(fpath[1]))
	if typeInfo.opt {
		field = AddrOf(field)
	}
	return g.builderMethod(FuncDecl{
		Doc:  Comments(Commentf("// %s sets the %s attribute of the %s from t.", exp2, attrName, g.builderName)),
		Name: *I(exp2),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("t"), Type: I(fpath[0])}),
			Results: Fields(Field{Type: Star(I(g.builderName))}),
		},
		Body: Block(
			updateMask,
			Assign(Dot(I(builderRecv), nexp))("=")(field),
			Return(I(builderRecv)),
		),
	})
}

func (g *entityGenerator) goEntityTypeStaticMetadataDecl() Decl {
	pkMask := []Expr{}
	for _, u := range g.pkMask {
		pkMask = append(pkMask, Lit(u))
	}
	reqMask := []Expr{}
	for _, u := range g.requiredMask {
		reqMask = append(reqMask, Lit(u))
	}
	return Var(ValueSpec{
		Names: Idents(NonExportedName(g.tname + "StaticMetadata")),
		Values: []Expr{
			AddrOf(Composite(g.relgomlib("EntityTypeStaticMetadata"),
				KV(I("PKMask"), Composite(&ArrayType{Elt: I("uint64")}, pkMask...)),
				KV(I("RequiredMask"), Composite(&ArrayType{Elt: I("uint64")}, reqMask...)),
			)),
		},
	})
}

func (g *entityGenerator) goBuilderApplyDecl() Decl {
	pkeys := make([]string, 0, len(g.relation.AttrDefs))
	for i, nt := range g.namedAttrs {
		if g.requiredMask[i/64]&(uint64(1)<<uint(i)) != 0 {
			pkeys = append(pkeys, nt.Name)
		} else {
			pkeys = append(pkeys, "")
		}
	}

	return g.builderMethod(FuncDecl{
		Doc:  Comments(Commentf("// Apply applies the built %s.", g.tname)),
		Name: *I("Apply"),
		Type: FuncType{
			Results: Fields(Field{Type: I(g.modelName)}, Field{Type: I(g.tname)}, Field{Type: I("error")}),
		},
		Body: Block(
			CallStmt(g.relgomlib("PanicIfRequiredFieldsNotSet"),
				Slice(Dot(I(builderRecv), "mask")),
				Dot(NonExportedID(g.tname+"StaticMetadata"), "RequiredMask"),
				String(strings.Join(pkeys, ",")),
			),
			Init("set", "err")(Call(Dot(I(builderRecv), "apply"), AddrOf(Dot(I(builderRecv), g.dataName)))),
			If(nil,
				Binary(I("err"), "!=", Nil()),
				Return(Composite(I(g.modelName)), Composite(I(g.tname)), I("err")),
			),
			Init("model")(
				Call(Dot(I(builderRecv), "model", "relations", "With"),
					NonExportedID(g.tname+"Key"),
					Composite(I(g.relationDataName), I("set")),
				),
			),
			Return(
				Composite(I(g.modelName), I("model")),
				Composite(I(g.tname), AddrOf(Dot(I(builderRecv), g.dataName)), Dot(I(builderRecv), "model")),
				Nil(),
			),
		),
	})
}

func (g *entityGenerator) applyFuncType() *FuncType {
	return &FuncType{
		Params:  *Fields(Field{Names: Idents("t"), Type: Star(I(g.dataName))}),
		Results: Fields(Field{Type: g.frozen("Map")}, Field{Type: I("error")}),
	}
}
