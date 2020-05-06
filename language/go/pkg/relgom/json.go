package relgom

import (
	. "github.com/anz-bank/sysl/unsorted/codegen/golang" //nolint:golint,stylecheck
)

func marshalJSONMethodDecl(stmts ...Stmt) *FuncDecl {
	return &FuncDecl{
		Doc:  Comments(Commentf("// MarshalJSON implements json.Marshaler.")),
		Name: *I("MarshalJSON"),
		Type: FuncType{
			Results: Fields(
				Field{Type: &ArrayType{Elt: I("byte")}},
				Field{Type: I("error")},
			),
		},
		Body: Block(stmts...),
	}
}

func unmarshalJSONMethodDecl(stmts ...Stmt) *FuncDecl {
	return &FuncDecl{
		Doc:  Comments(Commentf("// UnmarshalJSON implements json.Unmarshaler.")),
		Name: *I("UnmarshalJSON"),
		Type: FuncType{
			Params: *Fields(
				Field{Names: Idents("data"), Type: &ArrayType{Elt: I("byte")}},
			),
			Results: Fields(
				Field{Type: I("error")},
			),
		},
		Body: Block(stmts...),
	}
}
