package relgom

import (
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

func (g *entityGenerator) appendIterDecls(decls []Decl) []Decl {
	return append(decls,
		g.relationIteratorMethodDecl(),
		g.iteratorIntfDecl(),
		g.iteratorStructDecl(),
		g.iteratorMoveNextMethodDecl(),
		g.iteratorCurrentMethodDecl(),
	)
}

func (g *entityGenerator) iteratorName() string {
	return NonExportedName(g.tname + "Iterator")
}

func (g *entityGenerator) iteratorInterfaceName() string {
	return ExportedName(g.tname + "Iterator")
}

// // Iterator returns an iterator over ${export(.name)} tuples in r.
// func (r ${export(.name)}Relation) Iterator() ${export(.name)}Iterator {
//     return &${unexport(.name)}Iterator{r.model, r.set, nil}
// }
func (g *entityGenerator) relationIteratorMethodDecl() Decl {
	return g.relationMethod(func(recv string, recvDot dotter) FuncDecl {
		return FuncDecl{
			Doc:  Comments(Commentf("// Iterator returns an iterator over %s tuples in r.", g.tname)),
			Name: *I("Iterator"),
			Type: FuncType{
				Results: Fields(Field{Type: I(g.iteratorInterfaceName())}),
			},
			Body: Block(
				Return(AddrOf(Composite(I(g.iteratorName()),
					KV(I("model"), recvDot("model")),
					KV(I("i"), Call(recvDot("set", "Range"))),
				))),
			),
		}
	})
}

// // ${export(.name)}Iterator provides for iteration over a set of ${export(.name)} tuples.
// type ${export(.name)}Iterator interface {
//     MoveNext() bool
//     Current() ${export(.name)}
// }
func (g *entityGenerator) iteratorIntfDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.iteratorInterfaceName()),
		Type: &InterfaceType{
			Methods: *Fields(
				Field{
					Names: Idents("MoveNext"),
					Type: &FuncType{
						Params:  *Fields().Parens(),
						Results: Fields(Field{Type: I("bool")}),
					},
				},
				Field{
					Names: Idents("Current"),
					Type: &FuncType{
						Params:  *Fields().Parens(),
						Results: Fields(Field{Type: I(g.tname)}),
					},
				},
			),
		},
	}).WithDoc(Commentf("// %s provides for iteration over a set of %[1]s tuples.", g.iteratorName()))
}

// type ${u(.name)}Iterator struct {
//     model ${.model.name}
//     set   *seq.HashMap
//     t     *${export(.name)}
// }
func (g *entityGenerator) iteratorStructDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.iteratorName()),
		Type: Struct(
			Field{Names: Idents("model"), Type: I(g.modelName)},
			Field{Names: Idents("i"), Type: Star(g.frozen("MapIterator"))},
			Field{Names: Idents("t"), Type: Star(I(g.tname))},
		),
	})
}

func (g *entityGenerator) iteratorMethodDecl(f func(recv string, recvDot dotter) FuncDecl) *FuncDecl {
	return method("i", Star(I(g.iteratorName())), f)
}

// // MoveNext implements ${X(.name)}Iterator.
// func (i *${u(.name)Iterator) MoveNext() bool {
//     kv, set, has := i.set.FirstRestKV()
//     if has {
//         i.set = set
//         i.t = &${X(.name)}{${u(.name)Data: kv.Val.(*${u(.name)Data), model: i.model}
//     }
//     return has
// }
func (g *entityGenerator) iteratorMoveNextMethodDecl() *FuncDecl {
	return g.iteratorMethodDecl(func(recv string, recvDot dotter) FuncDecl {
		return FuncDecl{
			Doc:  Comments(Commentf("// MoveNext implements seq.Setable.")),
			Name: *I("MoveNext"),
			Type: FuncType{
				Results: Fields(Field{Type: I("bool")}),
			},
			Body: Block(
				If(nil, Call(recvDot("i", "Next")),
					Assign(recvDot("t"))("=")(
						AddrOf(Composite(I(g.tname),
							KV(I(g.dataName), Assert(Call(Dot(I("i"), "i", "Value")), Star(I(g.dataName)))),
							KV(I("model"), recvDot("model")),
						)),
					),
					Return(I("true")),
				),
				Return(I("false")),
			),
		}
	})
}

// // Current implements ${X(.name)}Iterator.
// func (i *${u(.name)}Iterator) Current() ${X(.name)} {
//     if i.t == nil {
//         panic("no current ${X(.name)}")
//     }
//     return i.t
// }
func (g *entityGenerator) iteratorCurrentMethodDecl() *FuncDecl {
	return g.iteratorMethodDecl(func(recv string, recvDot dotter) FuncDecl {
		return FuncDecl{
			Doc:  Comments(Commentf("// Current implements seq.Setable.")),
			Name: *I("Current"),
			Type: FuncType{
				Results: Fields(Field{Type: I(g.tname)}),
			},
			Body: Block(
				If(nil, Binary(recvDot("t"), "==", Nil()),
					Panic(String("no current "+g.tname)),
				),
				Return(Star(recvDot("t"))),
			),
		}
	})
}
