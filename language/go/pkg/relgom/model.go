package relgom

import (
	"github.com/anz-bank/sysl/pkg/syslutil"
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

const (
	modelRecv = "m"
)

type modelGenerator struct {
	*sourceGenerator
	*modelScope
	*commonModules
	namedTypes syslutil.NamedTypes
}

func newModelGenerator(s *modelScope) *modelGenerator {
	g := &modelGenerator{
		sourceGenerator: s.newSourceGenerator(),
		modelScope:      s,
		namedTypes:      syslutil.NamedTypesInSourceOrder(s.model.Types),
	}
	g.commonModules = newCommonModules(g.sourceGenerator)
	return g
}

func (g *modelGenerator) genFileForSyslModel() error {
	decls := []Decl{}

	modelKeys := []ValueSpec{}
	for _, nt := range g.namedTypes {
		modelKeys = append(modelKeys, ValueSpec{
			Names: []Ident{*NonExportedID(nt.Name + "Key")},
		})
	}
	modelKeys[0].Type = I("int")
	modelKeys[0].Values = []Expr{Iota()}
	decls = append(decls, Const(modelKeys...))

	decls = append(decls, Types(TypeSpec{
		Name: *I(g.modelName),
		Type: Struct(Field{Names: Idents("relations"), Type: g.frozen("Map")}),
	})) //TODO: Fix: .WithDoc(Commentf("// %s is the model.", g.modelName))))

	// // New<Model> creates a new <Model>.
	// func New<Model>() *<Model> {
	// 	   return &<Model>{}
	// }
	newName := "New" + g.modelName
	newDeclFunc := &FuncDecl{
		Doc:  Comments(Commentf("// %s creates a new %s.", newName, g.modelName)),
		Name: *I(newName),
		Type: FuncType{
			Results: Fields(Field{Type: I(g.modelName)}),
		},
		Body: Block(
			Return(Composite(I(g.modelName),
				Call(g.frozen("NewMap"),
					Call(g.frozen("KV"),
						g.relgomlib("ModelMetadataKey"),
						Composite(g.relgomlib("ModelMetadata")),
					),
				),
			)),
		),
	}
	decls = append(decls, newDeclFunc)

	// // <Relation> returns the model's <Relation> relation.
	// func (m *<Model>) Get<Relation>() *<Relation>Relation {
	//     return m.relations.GetVal
	// }
	for _, nt := range g.namedTypes {
		ename := ExportedName(nt.Name)
		sname := ename + "Relation"
		rdname := NonExportedName(nt.Name) + "RelationData"

		decls = append(decls,
			g.modelMethod(FuncDecl{
				Doc:  Comments(Commentf("// %s returns the model's %[1]s relation.", ename)),
				Name: *I("Get" + ename),
				Type: FuncType{Results: Fields(Field{Type: Star(I(sname))})},
				Body: Block(
					If(
						Init("relation", "has")(Call(Dot(I(modelRecv), "relations", "Get"),
							NonExportedID(nt.Name+"Key"),
						)),
						I("has"),
						Return(AddrOf(Composite(I(sname), Assert(I("relation"), I(rdname)), I(modelRecv)))),
					),
					Return(AddrOf(Composite(I(sname), Composite(I(rdname)), I(modelRecv)))),
				),
			}))
	}

	decls = append(decls,
		g.marshalJSONModelMethod(),
		g.unmarshalJSONModelMethod(),
		g.modelMethod(FuncDecl{
			Doc:  Comments(Commentf("// newID returns a new id for the model")),
			Name: *I("newID"),
			Type: FuncType{Results: Fields(Field{Type: I(g.modelName)}, Field{Type: I("uint64")})},
			Body: Block(
				Init("relations", "id")(Call(g.relgomlib("NewID"), Dot(I(modelRecv), "relations"))),
				Return(Composite(I(g.modelName), I("relations")), I("id")),
			)}),
	)

	return g.genSourceForDecls(g.modelName, decls...)
}

func (g *modelGenerator) marshalJSONModelMethod() *FuncDecl {
	stmts := []Stmt{
		Init("b")(Call(g.relgomlib("NewRelationMapBuilder"), Dot(I(modelRecv), "relations"))),
	}
	for _, nt := range g.namedTypes {
		stmts = append(stmts,
			CallStmt(Dot(I("b"), "Set"), String(nt.Name), NonExportedID(nt.Name+"Key")),
		)
	}
	stmts = append(stmts,
		Return(Call(g.json("Marshal"), Call(Dot(I("b"), "Map")))),
	)
	return g.modelMethod(*marshalJSONMethodDecl(stmts...))
}

func (g *modelGenerator) unmarshalJSONModelMethod() *FuncDecl {
	stmts := []Stmt{
		Init("e")(Call(g.relgomlib("NewRelationMapExtractor"), Dot(I(modelRecv), "relations"))),
	}
	for _, nt := range g.namedTypes {
		stmts = append(stmts,
			CallStmt(Dot(I("e"), "Set"),
				String(nt.Name),
				NonExportedID(nt.Name+"Key"),
				AddrOf(Composite(NonExportedID(nt.Name+"RelationData"))),
			),
		)
	}
	stmts = append(stmts,
		Init("relations", "err")(Call(Dot(I("e"), "UnmarshalRelationDataJSON"), I("data"))),
		If(nil, Binary(I("err"), "==", Nil()),
			Assign(Dot(I(modelRecv), "relations"))("=")(I("relations")),
		),
		Return(I("err")),
	)
	return g.modelMethodPointerRecv(*unmarshalJSONMethodDecl(stmts...))
}

func (g *modelGenerator) modelMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(modelRecv),
		Type:  I(g.modelName),
	}).Parens()
	return &f
}

func (g *modelGenerator) modelMethodPointerRecv(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(modelRecv),
		Type:  Star(I(g.modelName)),
	}).Parens()
	return &f
}
