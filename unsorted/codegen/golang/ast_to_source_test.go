package golang

import (
	"fmt"
	"go/format"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFormat(t *testing.T, n interface{}, expected string) {
	expected = strings.TrimSpace(expected)
	expected = strings.Replace(expected, "\n\t\t", "\n", -1)
	actual := fmt.Sprintf("%s", n)
	formatted, err := format.Source([]byte(actual))
	if assert.NoError(t, err, "unformatted: %s", actual) {
		trimmed := strings.TrimSpace(string(formatted))
		assert.Equal(t, expected, trimmed, "unformatted: %s", actual)
	}
}

func testFormatEmpty(t *testing.T, n interface{}) {
	assert.Zero(t, fmt.Sprintf("%s", n))
}

func testFormatPanics(t *testing.T, n interface{}) {
	actual := fmt.Sprintf("%s", n)
	assert.Contains(t, actual, "%!s(PANIC=")
}

func testFormatInfix(t *testing.T, n interface{}, prefix, expected, suffix string) bool {
	expected = strings.TrimSpace(expected)
	expected = strings.Replace(expected, "\n\t\t", "\n", -1)
	output := fmt.Sprintf("%s", n)
	formatted, err := format.Source([]byte(output))
	if !assert.NoError(t, err, "unformatted: %s", output) {
		return false
	}
	actual := strings.TrimSpace(string(formatted))
	assert.True(t, strings.HasPrefix(actual, prefix), "%#v %#v", output, actual)
	assert.True(t, strings.HasSuffix(actual, suffix), "%#v %#v", output, actual)
	return assert.Equal(t,
		expected, strings.TrimSpace(actual[len(prefix):len(actual)-len(suffix)]),
		"%s", output)
}

func testFormatType(t *testing.T, typ Expr, expected string) {
	testFormatInfix(t,
		&GenDecl{
			Tok:   *T("var"),
			Specs: []Spec{&TypeSpec{Name: *I("_"), Type: typ}},
		},
		`var _ `, expected, ``,
	)
}

func TestToken(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*Token)(nil))
	testFormat(t, T("foo"), `foo`)
}

func TestArrayType(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ArrayType)(nil))
	testFormatType(t, SliceType(I("int")), `[]int`)
	testFormatType(t, ArrayN(10, I("string")), `[10]string`)
}

func TestAssignStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*AssignStmt)(nil))
	testFormat(t, Assign(I("x"))("=")(Int(42)), `x = 42`)
	testFormat(t, Init("x")(Int(42)), `x := 42`)
	testFormat(t,
		Init("a", "b")(Int(42), String("hello")),
		`a, b := 42, "hello"`,
	)
}

func TestBadDecl(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BadDecl)(nil))
	testFormatPanics(t, &BadDecl{})
}

func TestBadExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BadExpr)(nil))
	testFormatPanics(t, &BadExpr{})
}

func TestBadStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BadStmt)(nil))
	testFormatPanics(t, &BadStmt{})
}

func TestBasicLit(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BasicLit)(nil))
	testFormat(t, &BasicLit{*T("123.45")}, `123.45`)
	testFormat(t, &BasicLit{*T(`"123.45"`)}, `"123.45"`)
}

func TestBinaryExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BinaryExpr)(nil))
	// https://golang.org/ref/spec#Operator_precedence
	for _, op := range []string{
		"*", "/", "%", "<<", ">>", "&", "&^",
		"+", "-", "*", "/",
		"==", "!=", "<", "<=", ">", ">=",
		"&&",
		"||",
	} {
		testFormat(t, Binary(Int(1), op, Int(2)), fmt.Sprintf(`(1 %s 2)`, op))
	}
}

func TestBlockStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BlockStmt)(nil))
	testFormat(t,
		&BlockStmt{},
		`{
		}`,
	)
	testFormat(t,
		&BlockStmt{List: []Stmt{&ExprStmt{X: I("x")}}},
		`{
			x
		}`,
	)
	testFormat(t,
		&BlockStmt{
			List: []Stmt{Init("x")(Int(42)), Return(I("x"), Nil())},
		},
		`{
			x := 42
			return x, nil
		}`,
	)
}

func TestBranchStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*BranchStmt)(nil))
	testFormat(t, Break(), `break`)
	testFormat(t, BreakTo("fast"), `break fast`)
	testFormat(t, Continue(), `continue`)
	testFormat(t, ContinueTo("dreaming"), `continue dreaming`)
	testFormat(t, Goto("jail"), `goto jail`)
	testFormat(t, Fallthrough(), `fallthrough`)
}

func TestCallExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*CallExpr)(nil))
	testFormat(t, &CallExpr{Fun: I("f")}, `f()`)
	testFormat(t, Call(I("sin"), I("π")), `sin(π)`)
	testFormat(t, Call(Dot(I("math"), "Atan2"), I("y"), I("x")), `math.Atan2(y, x)`)
	testFormat(t,
		CallVararg(Dot(I("fmt"), "Sprintf"), String("Hi %s and %s."), I("names")),
		`fmt.Sprintf("Hi %s and %s.", names...)`,
	)
}

func TestCaseClause(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*CaseClause)(nil))

	testFormatCaseClause := func(n *CaseClause, expected string) {
		testFormatInfix(t,
			&SwitchStmt{Body: BlockStmt{List: []Stmt{n}}},
			`switch {`, expected, `}`,
		)
	}

	testFormatCaseClause(DefaultCase(), `default:`)
	testFormatCaseClause(
		DefaultCase(&BlockStmt{}),
		`default:
			{
			}`,
	)
	testFormatCaseClause(
		DefaultCase(Return(I("π"))),
		`default:
			return π`,
	)
	testFormatCaseClause(
		Case(Binary(I("π"), "<", Int(4)))(
			&ExprStmt{
				X: Call(Dot(I("log"), "Info"), String("π seems about right!")),
			},
			Return(I("π")),
		),
		`case (π < 4):
			log.Info("π seems about right!")
			return π`,
	)
}

func TestChanType(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ChanType)(nil))
	value := Struct()
	testFormatType(t, &ChanType{Value: value, Dir: ""}, `chan struct{}`)
	testFormatType(t, &ChanType{Value: value, Dir: "SEND"}, `chan<- struct{}`)
	testFormatType(t, &ChanType{Value: value, Dir: "RECV"}, `<-chan struct{}`)
	testFormatPanics(t, &ChanType{Value: value, Dir: "TRANSCEIVE"})
}

func TestCommClause(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*CommClause)(nil))

	testFormatCommClause := func(n *CommClause, expected string) {
		testFormatInfix(t,
			&SelectStmt{Body: BlockStmt{List: []Stmt{n}}},
			`select {`, expected, `}`,
		)
	}

	testFormatCommClause(DefaultComm(), `default:`)
	testFormatCommClause(
		DefaultComm(&BlockStmt{}),
		`default:
			{
			}`,
	)
	testFormatCommClause(
		DefaultComm(Return(I("π"))),
		`default:
			return π`,
	)

	body := []Stmt{
		&ExprStmt{X: Call(Dot(I("log"), "Info"), String("They need π!"))},
		Return(I("π")),
	}
	bodySource := `
			log.Info("They need π!")
			return π`

	testFormatCommClause(
		SendComm(I("ch"), I("π"), body...),
		`case ch <- π:`+bodySource)
	testFormatCommClause(
		RecvAssignComm(I("x"))("=")(I("ch"))(body...),
		`case x = <-ch:`+bodySource,
	)
	testFormatCommClause(
		RecvInitComm("x")(Call(I("ch")))(body...),
		`case x := <-ch():`+bodySource,
	)
}

func TestComment(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*Comment)(nil))

	testFormat(t, &Comment{*T("//")}, `//`)
	testFormat(t, &Comment{*T("// Hmmmm")}, `// Hmmmm`)
	testFormat(t, &Comment{*T("/* Hmmmm */")}, `/* Hmmmm */`)
}

func TestCommentGroup(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*CommentGroup)(nil))

	testFormat(t, &CommentGroup{}, ``)
	testFormat(t,
		&CommentGroup{List: []Comment{{*T("// Comment 1")}}},
		`// Comment 1`,
	)
	testFormat(t,
		&CommentGroup{
			List: []Comment{
				{*T("// Comment 1")},
				{*T("// Comment 2")},
			},
		},
		`// Comment 1
		// Comment 2`,
	)
}

func TestCompositeLit(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*CompositeLit)(nil))
	testFormat(t, Composite(I("Foo"), KV(I("X"), Int(42))), `Foo{X: 42}`)
	testFormat(t,
		Composite(I("Foo"), KV(I("X"), Int(42)), KV(I("Y"), Float(3.14159))),
		`Foo{X: 42, Y: 3.14159}`,
	)
	testFormat(t,
		Composite(ArrayEllipsis(I("string")), String("hello"), String("world")),
		`[...]string{"hello", "world"}`,
	)
}

func TestDeclStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*DeclStmt)(nil))
	testFormat(t,
		&DeclStmt{
			&GenDecl{
				Tok:   *T("var"),
				Specs: []Spec{&TypeSpec{Name: *I("x"), Type: I("int")}},
			},
		},
		`var x int`,
	)
}

func TestDeferStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*DeferStmt)(nil))
	testFormat(t, Defer(I("f"), I("x")), `defer f(x)`)
	testFormat(t,
		Defer(
			Func(
				*Fields(Field{Names: Idents("a", "b"), Type: I("int")}).Parens(),
				nil,
				&ExprStmt{
					X: Call(Dot(I("log"), "Infof"), String("(%d, %d)"), I("a"), I("b"))},
			),
			I("x"), I("y"),
		),
		`defer func(a, b int) {
			log.Infof("(%d, %d)", a, b)
		}(x, y)`,
	)
}

func TestEllipsis(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*Ellipsis)(nil))
	// Tested via ArrayType
}

func TestEmptyStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*EmptyStmt)(nil))
	testFormat(t, &EmptyStmt{}, `;`)
}

func TestExprStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ExprStmt)(nil))
	testFormat(t, &ExprStmt{X: Int(0)}, `0`)
	testFormat(t,
		&ExprStmt{X: Composite(&ArrayType{Elt: &InterfaceType{}})},
		`[]interface{}{}`,
	)
}

func TestField(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*Field)(nil))
}

func TestFieldList(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*FieldList)(nil))
}

func TestFile(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*File)(nil))
	testFormat(t,
		&File{
			Name: *I("main"),
			Imports: []ImportSpec{
				{Path: *String("fmt")},
				{Path: *String("math")},
			},
			Decls: []Decl{
				&FuncDecl{
					Name: *I("main"),
					Type: FuncType{Params: *Fields().Parens()},
					Body: &BlockStmt{
						List: []Stmt{
							&ExprStmt{
								X: Call(
									Dot(I("fmt"), "Printf"),
									String("Hello, %g\n"),
									Dot(I("math"), "PI"),
								),
							},
						},
					},
				},
			},
		},
		`package main

		import (
			"fmt"
			"math"
		)

		func main() {
			fmt.Printf("Hello, %g\n", math.PI)
		}`,
	)
}

func TestForStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ForStmt)(nil))

	testFormat(t,
		&ForStmt{},
		`for {
		}`,
	)
	testFormat(t,
		&ForStmt{Body: *Block(&ExprStmt{X: Int(0)})},
		`for {
			0
		}`,
	)
}

func TestFuncDecl(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*FuncDecl)(nil))

	testFormat(t,
		&FuncDecl{
			Name: *I("f"),
			Type: FuncType{
				Params: *Fields(
					Field{Names: Idents("a", "b"), Type: I("int")},
				).Parens(),
			},
			Body: Block(
				&ExprStmt{
					X: Call(Dot(I("log"), "Infof"),
						&BasicLit{Token{Text: `"(%d, %d)"`}}, I("a"), I("b"),
					),
				},
			),
		},
		`func f(a, b int) {
			log.Infof("(%d, %d)", a, b)
		}`,
	)
	testFormat(t,
		&FuncDecl{
			Name: *I("f"),
			Type: FuncType{
				Params: *Fields(
					Field{Names: Idents("a"), Type: I("int")},
					Field{Names: Idents("b"), Type: &Ellipsis{Elt: I("int")}},
				).Parens(),
			},
			Body: Block(),
		},
		`func f(a int, b ...int) {}`,
	)
}

func TestFuncLit(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*FuncLit)(nil))
	testFormat(t,
		Composite(
			&ArrayType{Elt: &InterfaceType{}},
			Func(
				*Fields(Field{Names: Idents("a", "b"), Type: I("int")}).Parens(),
				nil,
				&ExprStmt{
					X: Call(Dot(I("log"), "Infof"),
						&BasicLit{Token{Text: `"(%d, %d)"`}}, I("a"), I("b"),
					),
				},
			),
		),
		`[]interface{}{func(a, b int) {
			log.Infof("(%d, %d)", a, b)
		}}`,
	)
}

func TestFuncType(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*FuncType)(nil))
	// Tested via FuncDecl and FuncLit.
}

func TestGenDeclImport(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t, Import(ImportSpec{Path: *String("fmt")}), `import "fmt"`)
	testFormat(t,
		Import(
			ImportSpec{Name: I("queries"), Path: *String("database/sql")},
			ImportSpec{Path: *String("fmt")},
			ImportSpec{Path: *String("io")},
			ImportSpec{Path: *String("math")},
		),
		`import (
			queries "database/sql"
			"fmt"
			"io"
			"math"
		)`,
	)
}

func TestGenDeclVar(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t,
		Var(ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}}),
		`var x = 42`,
	)
	testFormat(t,
		Var(
			ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}},
			ValueSpec{
				Names:  Idents("r, g, b"),
				Type:   I("string"),
				Values: []Expr{String("red"), String("green"), String("blue")},
			},
		),
		`var (
			x              = 42
			r, g, b string = "red", "green", "blue"
		)`,
	)
}

func TestGenDeclConst(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t,
		Const(ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}}),
		`const x = 42`,
	)
	testFormat(t,
		Const(
			ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}},
			ValueSpec{
				Names:  Idents("r, g, b"),
				Type:   I("string"),
				Values: []Expr{String("red"), String("green"), String("blue")},
			},
		),
		`const (
			x              = 42
			r, g, b string = "red", "green", "blue"
		)`,
	)
}

func TestGenDeclType(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t,
		Types(TypeSpec{Name: *I("id"), Type: I("uint64")}),
		`type id uint64`,
	)

	coords := Idents("x", "y", "z", "w")

	testFormat(t,
		Types(
			TypeSpec{
				Name: *I("vec2"),
				Type: Struct(Field{Names: coords[:2], Type: I("float32")}),
			},
			TypeSpec{
				Name: *I("vec3"),
				Type: Struct(Field{Names: coords[:3], Type: I("float32")}),
			},
			TypeSpec{
				Name: *I("vec4"),
				Type: Struct(Field{Names: coords[:4], Type: I("float32")}),
			},
		),
		`type (
			vec2 struct {
				x, y float32
			}
			vec3 struct {
				x, y, z float32
			}
			vec4 struct {
				x, y, z, w float32
			}
		)`,
	)
}

func TestGoStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*GoStmt)(nil))
	testFormat(t, &GoStmt{Call: CallExpr{Fun: I("f")}}, `go f()`)
}

func TestIdent(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*Ident)(nil))
	testFormat(t, I("x"), `x`)
	testFormat(t, I("fooBarBaz"), `fooBarBaz`)
	testFormat(t, I("π"), `π`)
}

func TestIfStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*IfStmt)(nil))

	testFormat(t,
		If(nil, I("ok"), Return(Int(0))),
		`if ok {
			return 0
		}`,
	)
	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"), Return(Int(0)),
		),
		`if x, ok := f(); ok {
			return 0
		}`,
	)
	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"), Return(Int(0)),
		),
		`if x, ok := f(); ok {
			return 0
		}`,
	)
}

func TestIfElse(t *testing.T) {
	t.Parallel()

	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"),
			Return(Int(0)),
		).WithElse(
			Return(Int(1)),
		),
		`if x, ok := f(); ok {
			return 0
		} else {
			return 1
		}`,
	)
	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"),
			Init("y")(Binary(Int(1), "/", I("x"))),
			Return(I("y")),
		).WithElse(
			&ExprStmt{
				X: Call(
					Dot(I("fmt"), "Printf"),
					String("Oops! (%d)"),
					I("x"),
				),
			},
			Return(Int(1)),
		),
		`if x, ok := f(); ok {
			y := (1 / x)
			return y
		} else {
			fmt.Printf("Oops! (%d)", x)
			return 1
		}`,
	)
}

func TestIfElseIf(t *testing.T) {
	t.Parallel()

	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"),
			Init("y")(Binary(Int(1), "/", I("x"))),
			Return(I("y")),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		),
		`if x, ok := f(); ok {
			y := (1 / x)
			return y
		} else if panicOnError {
			panic("Ouch!")
		}`,
	)
}

func TestIfElseIfElse(t *testing.T) {
	t.Parallel()

	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"),
			Init("y")(Binary(Int(1), "/", I("x"))),
			Return(I("y")),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		).WithElse(
			&ExprStmt{
				X: Call(
					Dot(I("fmt"), "Printf"),
					String("Oops! (%d)"),
					I("x"),
				),
			},
			Return(Int(1)),
		),
		`if x, ok := f(); ok {
			y := (1 / x)
			return y
		} else if panicOnError {
			panic("Ouch!")
		} else {
			fmt.Printf("Oops! (%d)", x)
			return 1
		}`,
	)
}

func TestIfElseIfElseIf(t *testing.T) {
	t.Parallel()

	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"),
			Init("y")(Binary(Int(1), "/", I("x"))),
			Return(I("y"), Nil()),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		).WithElseIf(
			nil,
			I("defaultOnError"),
			&ExprStmt{
				X: Call(Dot(I("fmt"), "Printf"), String("Oops! (%d)"), I("x")),
			},
			Return(Int(1), Nil()),
		),
		`if x, ok := f(); ok {
			y := (1 / x)
			return y, nil
		} else if panicOnError {
			panic("Ouch!")
		} else if defaultOnError {
			fmt.Printf("Oops! (%d)", x)
			return 1, nil
		}`,
	)
}

func TestIfElseIfElseIfElse(t *testing.T) {
	t.Parallel()

	testFormat(t,
		If(
			Init("x", "ok")(Call(I("f"))),
			I("ok"),
			Init("y")(Binary(Int(1), "/", I("x"))),
			Return(I("y"), Nil()),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		).WithElseIf(
			nil,
			I("defaultOnError"),
			&ExprStmt{
				X: Call(Dot(I("fmt"), "Printf"), String("Oops! (%d)"), I("x")),
			},
			Return(Int(1), Nil()),
		).WithElse(
			Return(Int(0), Call(Dot(I("fmt"), "Errorf"), String("Oops! (%d)"), I("x"))),
		),
		`if x, ok := f(); ok {
			y := (1 / x)
			return y, nil
		} else if panicOnError {
			panic("Ouch!")
		} else if defaultOnError {
			fmt.Printf("Oops! (%d)", x)
			return 1, nil
		} else {
			return 0, fmt.Errorf("Oops! (%d)", x)
		}`,
	)
}

func TestImportSpec(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ImportSpec)(nil))
	testFormat(t, &ImportSpec{}, ``)
}

func TestIncDecStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*IncDecStmt)(nil))
	testFormat(t, Inc(I("x")), `x++`)
	testFormat(t, Dec(I("x")), `x--`)
}

func TestIndexExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*IndexExpr)(nil))
	testFormat(t, Index(I("x"), I("y")), `x[y]`)
}

// TODO: Fix this.
// func TestInterfaceType(t *testing.T) {
// 	t.Parallel()

// 	testFormatEmpty(t, (*InterfaceType)(nil))
// 	testFormatType(t, &InterfaceType{}, `interface{}`)
// 	testFormatType(t,
// 		&InterfaceType{
// 			Methods: FieldList{
// 				List: []Field{
// 					{
// 						Names: Idents("Foo"),
// 						Type: &FuncType{
// 							Params: *Fields(
// 								Field{Names: Idents("a", "b"), Type: I("int")},
// 							).Parens(),
// 							Results: Fields(
// 								Field{Type: I("int")},
// 								Field{Type: I("string")},
// 							).Parens(),
// 						},
// 					},
// 				},
// 			},
// 		},
// 		`interface {
// 			Foo(a, b int) (int, string)
// 		}`,
// 	)
// }

func TestKeyValueExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*KeyValueExpr)(nil))
}

func TestLabeledStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*LabeledStmt)(nil))

	testFormat(t,
		&LabeledStmt{Label: *I("here"), Stmt: Return(I("x"), I("y"))},
		`here:
		return x, y`,
	)
	testFormat(t,
		&LabeledStmt{Label: *I("here"), Stmt: Block(Return(I("x"), I("y")))},
		`here:
		{
			return x, y
		}`,
	)
}

func TestMapType(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*MapType)(nil))

	testFormatType(t, Map(I("int"), I("string")), `map[int]string`)
	testFormatType(t,
		Map(I("int"), Map(I("float"), I("string"))),
		`map[int]map[float]string`,
	)
}

func TestParenExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ParenExpr)(nil))
	testFormat(t, &ParenExpr{X: I("x")}, `x`)
	testFormat(t, &ParenExpr{X: Binary(I("x"), "^", I("y"))}, `(x ^ y)`)
}

func TestRangeStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*RangeStmt)(nil))

	testFormat(t,
		Range(nil, nil, "", I("records"), &IncDecStmt{X: I("x"), Tok: *T("++")}),
		`for range records {
			x++
		}`)
	testFormat(t,
		Range(I("i"), nil, ":=", I("input"),
			Assign(Index(I("output"), I("i")))("=")(Index(I("input"), I("i"))),
		),
		`for i := range input {
			output[i] = input[i]
		}`)
	testFormat(t,
		Range(I("_"), I("x"), ":=", I("input"),
			Assign(I("output"))("=")(Call(I("append"), I("output"), Call(I("f"), I("x")))),
		),
		`for _, x := range input {
			output = append(output, f(x))
		}`)
	testFormat(t,
		Range(I("i"), I("elt"), "=", Dot(I("os"), "Args"),
			If(nil, Binary(I("i"), ">", Int(0)), &ExprStmt{X: Call(I("process"), I("elt"))}),
		),
		`for i, elt = range os.Args {
			if i > 0 {
				process(elt)
			}
		}`)
}

func TestReturnStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ReturnStmt)(nil))
	testFormat(t, Return(), `return`)
	testFormat(t, Return(I("x")), `return x`)
}

func TestSelectStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*SelectStmt)(nil))

	testFormat(t, Select(), `select {}`)
	testFormat(t,
		Select(SendComm(I("ch"), I("value"))),
		`select {
		case ch <- value:
		}`,
	)
	testFormat(t,
		Select(SendComm(I("ch"), I("value"), Inc(I("x")))),
		`select {
		case ch <- value:
			x++
		}`,
	)
	testFormat(t,
		Select(
			RecvInitComm("value")(I("ch"))(
				Return(I("value")),
			),
		),
		`select {
		case value := <-ch:
			return value
		}`,
	)
	testFormat(t,
		Select(
			SendComm(I("ch"), I("value"),
				Inc(I("x")),
			),
			RecvInitComm("value")(I("ch"))(
				Return(I("value")),
			),
			DefaultComm(&ExprStmt{Call(I("panic"), String("I can't wait!"))}),
		),
		`select {
		case ch <- value:
			x++
		case value := <-ch:
			return value
		default:
			panic("I can't wait!")
		}`,
	)
}

func TestSelectorExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*SelectorExpr)(nil))
	testFormat(t, Dot(I("x"), "y"), `x.y`)
}

func TestSendStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*SendStmt)(nil))
	testFormat(t, Send(I("ch"), I("x")), `ch <- x`)
}

func TestSliceExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*SliceExpr)(nil))

	x := I("x")
	low := I("low")
	high := I("high")
	max := I("max")

	testFormat(t, Slice(x), `x[:]`)
	testFormat(t, Slice(x, low), `x[low:]`)
	testFormat(t, Slice(x, nil, high), `x[:high]`)
	testFormat(t, Slice(x, low, high), `x[low:high]`)
	testFormat(t, Slice(x, nil, high, max), `x[:high:max]`)
	testFormat(t, Slice(x, low, high, max), `x[low:high:max]`)

	assert.Panics(t, func() { Slice(x, low, high, max, I("wat")) })
}

func TestStarExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*StarExpr)(nil))
	testFormat(t, Star(I("x")), `*x`)
	testFormat(t, Star(Binary(I("x"), "+", I("y"))), `*(x + y)`)
}

func TestStructType(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*StructType)(nil))

	testFormatType(t, Struct(), `struct{}`)
	testFormatType(t, Struct(Field{Names: Idents("X"), Type: I("int")}),
		`struct {
			X int
		}`)
	testFormatType(t,
		Struct(
			Field{Names: Idents("X", "Y"), Type: I("int")},
			Field{Names: Idents("Z"), Type: I("int"), Tag: String(`json:"id"`)},
		),
		`struct {
			X, Y int
			Z    int "json:\"id\""
		}`)
}

func TestSwitchStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*SwitchStmt)(nil))
	testFormat(t,
		&SwitchStmt{},
		`switch {
		}`,
	)
	testFormat(t,
		&SwitchStmt{Body: *Block(DefaultCase(&BlockStmt{}))},
		`switch {
		default:
			{
			}
		}`,
	)
	testFormat(t,
		&SwitchStmt{Body: *Block(DefaultCase(&BlockStmt{}))},
		`switch {
		default:
			{
			}
		}`,
	)
}

func TestTypeAssertExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*TypeAssertExpr)(nil))
	testFormat(t, Assert(I("x"), I("y")), `x.(y)`)
	testFormat(t, AssertType(I("x")), `x.(type)`)
}

func TestTypeSpec(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*TypeSpec)(nil))
}

func TestTypeSwitchStmt(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*TypeSwitchStmt)(nil))

	testFormat(t,
		TypeSwitch(nil, "", I("y")),
		`switch y.(type) {
		}`,
	)
	testFormat(t,
		TypeSwitch(nil, "x", I("y")),
		`switch x := y.(type) {
		}`,
	)
	testFormat(t,
		TypeSwitch(
			Init("y")(Int(1)), "x", I("y"),
			Case(I("int"), I("float"))(
				Return(Int(0)),
			),
		),
		`switch y := 1; x := y.(type) {
		case int, float:
			return 0
		}`,
	)
}

func TestUnaryExpr(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*UnaryExpr)(nil))
	// https://golang.org/ref/spec#Operators
	for _, op := range []string{"+", "-", "!", "^", "*", "&", "<-"} {
		testFormat(t, Unary(op, I("x")), op+`x`)
	}
}

func TestValueSpec(t *testing.T) {
	t.Parallel()

	testFormatEmpty(t, (*ValueSpec)(nil))
}
