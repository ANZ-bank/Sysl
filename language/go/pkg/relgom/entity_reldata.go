package relgom

import (
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

const relationDataRecv = "d"

func (g *entityGenerator) goAppendRelationDataDecls(decls []Decl) []Decl {
	return append(decls,
		g.relationDataTypeDecl(),
		g.relationDataCountMethodDecl(),
		g.relationDataMarshalJSONMethodDecl(),
		g.relationDataUnmarshalJSONMethodDecl(),
	)
}

// type ${.name}RelationData struct {
//     set   *seq.HashMap
// }
func (g *entityGenerator) relationDataTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.relationDataName),
		Type: Struct(
			Field{Names: Idents("set"), Type: g.frozen("Map")},
		),
	}).WithDoc(Commentf("// %s represents a set of %s.", g.relationDataName, g.tname))
}

func (g *entityGenerator) relationDataMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(relationDataRecv),
		Type:  I(g.relationDataName),
	}).Parens()
	return &f
}

func (g *entityGenerator) relationDataMethodPointerRecv(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(relationDataRecv),
		Type:  Star(I(g.relationDataName)),
	}).Parens()
	return &f
}

func (g *entityGenerator) relationDataCountMethodDecl() Decl {
	return g.relationDataMethod(FuncDecl{
		Doc:  Comments(Commentf("// Count returns the number of tuples in %s.", relationDataRecv)),
		Name: *I("Count"),
		Type: FuncType{
			Results: Fields(Field{Type: I("int")}),
		},
		Body: Block(
			Return(Call(Dot(I(relationDataRecv), "set", "Count"))),
		),
	})
}

// func (r *${.name}RelationData) MarshalJSON() ([]byte, error) {
//     a := make([]${.name}Data, 0, r.set.Count())
//     for kv, m, has := r.set.Range(); has; kv, m, has = r.set.Range() {
//         a = append(a, kv.Val.(${.name}Data))
//     }
//     return json.Marshal(a)
// }
func (g *entityGenerator) relationDataMarshalJSONMethodDecl() Decl {
	return g.relationDataMethod(*marshalJSONMethodDecl(
		Init("a")(Call(I("make"),
			&ArrayType{Elt: Star(I(g.dataName))},
			Int(0),
			Call(Dot(I(relationDataRecv), "set", "Count")),
		)),
		&ForStmt{
			Init: Init("i")(Call(Dot(I(relationDataRecv), "set", "Range"))),
			Cond: Call(Dot(I("i"), "Next")),
			Body: *Block(
				Append(I("a"), Assert(Call(Dot(I("i"), "Value")), Star(I(g.dataName)))),
			),
		},
		Return(Call(g.json("Marshal"), I("a"))),
	))
}

// func (r *${.name}RelationData) UnmarshalJSON(data []byte) error {
//     a := []${.name}Data{}
//     if err := json.Unmarshal(data, &a); err != nil {
//         return err
//     }
//     set := seq.NewMap()
//     for _, e := range a {
//         set, _ = set.Set(e.${.name}PK, e)
//     }
//     *d = ${.name}RelationData{set}
//     return nil
// }
func (g *entityGenerator) relationDataUnmarshalJSONMethodDecl() Decl {
	var i, key Expr
	if g.haveKeys {
		i, key = I("_"), Dot(I("e"), g.pkName)
	} else {
		i, key = I("i"), I("i")
	}
	return g.relationDataMethodPointerRecv(*unmarshalJSONMethodDecl(
		Init("a")(&ArrayType{Elt: Composite(Star(I(g.dataName)))}),
		If(
			Init("err")(Call(g.json("Unmarshal"), I("data"), AddrOf(I("a")))),
			Binary(I("err"), "!=", Nil()),
			Return(I("err")),
		),
		Init("set")(Call(g.frozen("NewMap"))),
		Range(i, I("e"), ":=", I("a"),
			Assign(I("set"))("=")(Call(Dot(I("set"), "With"), key, I("e"))),
		),
		Assign(Star(I(relationDataRecv)))("=")(Composite(I(g.relationDataName), I("set"))),
		Return(Nil()),
	))
}
