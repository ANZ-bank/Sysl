package golang

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
)

type genFormatter struct {
	I interface{}
}

func (f genFormatter) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), f.I)
}

// Format formats a node as a snippet of Go source code.
func Format(node fmt.Formatter) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", node)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	return formatted
}

// fixupFormat replaces `"%c"` with `"%"+c`, but skips `/(?<!%)(%%)+c/`.
func fixupFormat(format string, c rune) string {
	var b strings.Builder
	b.Grow(len(format))
	percent := false
	for _, r := range format {
		if percent {
			percent = false
			if r == 'c' {
				b.WriteRune(c)
				continue
			}
		} else {
			percent = r == '%'
		}
		b.WriteRune(r)
	}
	return b.String()
}

func scPrintf(w io.Writer, c rune, format string, args ...fmt.Formatter) {
	args2 := make([]interface{}, 0, len(args))
	for _, arg := range args {
		args2 = append(args2, arg)
	}
	fmt.Fprintf(w, fixupFormat(format, c), args2...)
}

func cSprintf(c rune, format string, args ...fmt.Formatter) string {
	var b strings.Builder
	scPrintf(&b, c, format, args...)
	return b.String()
}

func separator(w io.Writer, i int, sep string) {
	if i > 0 {
		if _, err := w.Write([]byte(sep)); err != nil {
			panic(err)
		}
	}
}

func sepFormat(s fmt.State, c rune, i int, f fmt.Formatter, sep string) {
	separator(s, i, sep)
	if f != nil {
		f.Format(s, c)
	}
}

func emptyNil(n interface{}) fmt.Formatter {
	if n != nil {
		return genFormatter{n}
	}
	return genFormatter{""}
}

// Format formats a Token as a snippet of Go source code.
func (n *Token) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c", genFormatter{n.Text})
}

// Format formats an ArrayType as a snippet of Go source code.
func (n *ArrayType) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "[%c]%c", emptyNil(n.Len), n.Elt)
}

// Format formats an AssignStmt as a snippet of Go source code.
func (n *AssignStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	for i, expr := range n.LHS {
		sepFormat(s, c, i, expr, ", ")
	}
	n.Tok.Format(s, c)
	for i, expr := range n.RHS {
		sepFormat(s, c, i, expr, ", ")
	}
}

// Format formats a BadDecl as a snippet of Go source code.
func (n *BadDecl) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	panic("BadDecl.Format Not implemented")
}

// Format formats a BadExpr as a snippet of Go source code.
func (n *BadExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	panic("BadExpr.Format Not implemented")
}

// Format formats a BadStmt as a snippet of Go source code.
func (n *BadStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	panic("BadStmt.Format Not implemented")
}

// Format formats a BasicLit as a snippet of Go source code.
func (n *BasicLit) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Token.Format(s, c)
}

// Format formats a BinaryExpr as a snippet of Go source code.
func (n *BinaryExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "(%c %c %c)", n.X, &n.Op, n.Y)
}

// Format formats a BlockStmt as a snippet of Go source code.
func (n *BlockStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	fmt.Fprint(s, "{")
	for _, stmt := range n.List {
		scPrintf(s, c, "\n%c", stmt)
	}
	fmt.Fprint(s, "}")
}

// Format formats a BranchStmt as a snippet of Go source code.
func (n *BranchStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c %c", &n.Tok, n.Label)
}

// Format formats a CallExpr as a snippet of Go source code.
func (n *CallExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c(", n.Fun)
	for i, expr := range n.Args {
		sepFormat(s, c, i, expr, ", ")
	}
	scPrintf(s, c, "%c)", &n.Ellipsis)
}

// Format formats a CaseClause as a snippet of Go source code.
func (n *CaseClause) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	if len(n.List) == 0 {
		fmt.Fprint(s, "default:")
	} else {
		fmt.Fprint(s, "case ")
		for i, expr := range n.List {
			sepFormat(s, c, i, expr, ", ")
		}
		fmt.Fprint(s, ":")
	}
	for _, stmt := range n.Body {
		scPrintf(s, c, "\n%c", stmt)
	}
}

// Format formats a ChanType as a snippet of Go source code.
func (n *ChanType) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	switch n.Dir {
	case "":
		scPrintf(s, c, "chan %c", n.Value)
	case "SEND":
		scPrintf(s, c, "chan<- %c", n.Value)
	case "RECV":
		scPrintf(s, c, "<-chan %c", n.Value)
	default:
		panic(fmt.Sprintf("Unknown channel type %#v", n.Dir))
	}
}

// Format formats a CommClause as a snippet of Go source code.
func (n *CommClause) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	if n.Comm == nil {
		fmt.Fprint(s, "default:")
	} else {
		scPrintf(s, c, "case %c:", n.Comm)
	}
	for _, stmt := range n.Body {
		scPrintf(s, c, "\n%c", stmt)
	}
}

// Format formats a Comment as a snippet of Go source code.
func (n *Comment) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c\n", &n.Token)
}

// Format formats a CommentGroup as a snippet of Go source code.
func (n *CommentGroup) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	for _, comment := range n.List {
		comment := comment
		scPrintf(s, c, "%c", &comment)
	}
}

// Format formats a CompositeLit as a snippet of Go source code.
func (n *CompositeLit) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c{", emptyNil(n.Type))
	for i, expr := range n.Elts {
		sepFormat(s, c, i, expr, ", ")
	}
	fmt.Fprint(s, "}")
}

// Format formats a DeclStmt as a snippet of Go source code.
func (n *DeclStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Decl.Format(s, c)
}

// Format formats a DeferStmt as a snippet of Go source code.
func (n *DeferStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "defer %c", &n.Call)
}

// Format formats an Ellipsis as a snippet of Go source code.
func (n *Ellipsis) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "...%c", emptyNil(n.Elt))
}

// Format formats an EmptyStmt as a snippet of Go source code.
func (n *EmptyStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	fmt.Fprint(s, ";")
}

// Format formats an ExprStmt as a snippet of Go source code.
func (n *ExprStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.X.Format(s, c)
}

// Format formats a Field as a snippet of Go source code.
func (n *Field) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Doc.Format(s, c)
	for i, name := range n.Names {
		name := name
		sepFormat(s, c, i, &name, ", ")
	}
	prefix := ""
	if _, ok := n.Type.(*FuncType); ok {
		prefix = "func"
	}
	scPrintf(s, c, " %c%c %c", genFormatter{prefix}, n.Type, n.Tag)
}

// Format formats a FieldList as a snippet of Go source code.
func (n *FieldList) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Opening.Format(s, c)
	for i, field := range n.List {
		field := field
		sepFormat(s, c, i, &field, ", ")
	}
	n.Closing.Format(s, c)
}

// NoCommaSepFieldList ...
type NoCommaSepFieldList FieldList

// Format ...
func (n *NoCommaSepFieldList) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Opening.Format(s, c)
	for _, field := range n.List {
		field := field
		scPrintf(s, c, "%c\n", &field)
	}
	n.Closing.Format(s, c)
}

// Format formats a File as a snippet of Go source code.
func (n *File) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Doc.Format(s, c)
	scPrintf(s, c, "package %c\n", &n.Name)
	if len(n.Imports) == 1 {
		scPrintf(s, c, "\nimport %c\n\n", &n.Imports[0])
	} else if len(n.Imports) > 1 {
		scPrintf(s, c, "\nimport (\n")
		for _, importSpec := range n.Imports {
			importSpec := importSpec
			scPrintf(s, c, "%c\n", &importSpec)
		}
		fmt.Fprintf(s, ")\n\n")
	}
	for i, decl := range n.Decls {
		sepFormat(s, c, i, decl, "\n\n")
	}
}

// Format formats a ForStmt as a snippet of Go source code.
func (n *ForStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "for %c ; %c ; %c %c",
		genFormatter{strings.TrimSpace(cSprintf(c, "%c", emptyNil(n.Init)))},
		emptyNil(n.Cond),
		genFormatter{strings.TrimSpace(cSprintf(c, "%c", emptyNil(n.Post)))},
		&n.Body,
	)
}

// Format formats a FuncDecl as a snippet of Go source code.
func (n *FuncDecl) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%cfunc %c %c %c %c", n.Doc, n.Recv, &n.Name, &n.Type, n.Body)
}

// Format formats a FuncLit as a snippet of Go source code.
func (n *FuncLit) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "func %c %c", &n.Type, &n.Body)
}

// Format formats a FuncType as a snippet of Go source code.
func (n *FuncType) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	results := n.Results
	if results != nil {
		results = results.Parens()
	}
	scPrintf(s, c, "%c %c", n.Params.Parens(), results)
}

// Format formats a GenDecl as a snippet of Go source code.
func (n *GenDecl) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	if len(n.Specs) == 1 {
		scPrintf(s, c, "%c%c %c", n.Doc, &n.Tok, n.Specs[0])
	} else {
		scPrintf(s, c, "%c%c (\n", n.Doc, &n.Tok)
		for _, spec := range n.Specs {
			scPrintf(s, c, "%c\n", spec)
		}
		fmt.Fprint(s, ")")
	}
}

// Format formats a GoStmt as a snippet of Go source code.
func (n *GoStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "go %c", &n.Call)
}

// Format formats an Ident as a snippet of Go source code.
func (n *Ident) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Name.Format(s, c)
}

// Format formats an IfStmt as a snippet of Go source code.
func (n *IfStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	ifBlock := cSprintf(c, "if %c ; %c %c",
		emptyNil(n.Init), emptyNil(n.Cond), &n.Body,
	)
	if n.Else == nil {
		fmt.Fprint(s, ifBlock)
	} else if _, ok := n.Else.(*BlockStmt); ok {
		fmt.Fprint(s)
		scPrintf(s, c, "%s else %c",
			genFormatter{strings.TrimSpace(ifBlock)}, n.Else,
		)
	} else {
		scPrintf(s, c, "%c else\n%c",
			genFormatter{strings.TrimSpace(ifBlock)}, n.Else,
		)
	}
}

// Format formats an ImportSpec as a snippet of Go source code.
func (n *ImportSpec) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c %c %c %c", n.Doc, n.Name, &n.Path, n.Comment)
}

// Format formats an IncDecStmt as a snippet of Go source code.
func (n *IncDecStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c%c", n.X, &n.Tok)
}

// Format formats an IndexExpr as a snippet of Go source code.
func (n *IndexExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c[%c]", n.X, n.Index)
}

// Format formats an InterfaceType as a snippet of Go source code.
func (n *InterfaceType) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	if len(n.Methods.List) == 0 {
		scPrintf(s, c, "interface{}")
	} else {
		scPrintf(s, c, "interface {\n")
		for _, m := range n.Methods.List {
			f := m.Type.(*FuncType)
			scPrintf(s, c, "%c %c %c\n", &m.Names[0], &f.Params, f.Results)
		}
		scPrintf(s, c, "}\n")
	}
}

// Format formats a KeyValueExpr as a snippet of Go source code.
func (n *KeyValueExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c: %c", n.Key, n.Value)
}

// Format formats a LabeledStmt as a snippet of Go source code.
func (n *LabeledStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c:\n%c", &n.Label, n.Stmt)
}

// Format formats a MapType as a snippet of Go source code.
func (n *MapType) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "map[%c]%c", n.Key, n.Value)
}

// Format formats a ParenExpr as a snippet of Go source code.
func (n *ParenExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	// All Exprs protect themselves with parens. Not required here.
	n.X.Format(s, c)
}

// Format formats a RangeStmt as a snippet of Go source code.
func (n *RangeStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	switch {
	case n.Key == nil:
		scPrintf(s, c, "for range %c %c", emptyNil(n.X), &n.Body)
	case n.Value == nil:
		scPrintf(s, c, "for %c %c range %c %c",
			emptyNil(n.Key), &n.Tok, emptyNil(n.X), &n.Body,
		)
	default:
		scPrintf(s, c, "for %c, %c %c range %c %c",
			emptyNil(n.Key), emptyNil(n.Value), &n.Tok, emptyNil(n.X), &n.Body,
		)
	}
}

// Format formats a ReturnStmt as a snippet of Go source code.
func (n *ReturnStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "return ")
	for i, expr := range n.Results {
		sepFormat(s, c, i, expr, ", ")
	}
}

// Format formats a SelectStmt as a snippet of Go source code.
func (n *SelectStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "select %c", &n.Body)
}

// Format formats a SelectorExpr as a snippet of Go source code.
func (n *SelectorExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c.%c", n.X, &n.Sel)
}

// Format formats a SendStmt as a snippet of Go source code.
func (n *SendStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c <- %c", n.Chan, n.Value)
}

// Format formats a SliceExpr as a snippet of Go source code.
func (n *SliceExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c[%c:%c", n.X, emptyNil(n.Low), emptyNil(n.High))
	if n.Max != nil {
		scPrintf(s, c, ":%c", emptyNil(n.Max))
	}
	fmt.Fprint(s, "]")
}

// Format formats a StarExpr as a snippet of Go source code.
func (n *StarExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "*%c", n.X)
}

// Format formats a StructType as a snippet of Go source code.
func (n *StructType) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	if len(n.Fields.List) == 0 {
		scPrintf(s, c, "struct{}")
	} else {
		f := NoCommaSepFieldList(n.Fields)
		scPrintf(s, c, "struct{\n%c}", &f)
	}
}

// Format formats a SwitchStmt as a snippet of Go source code.
func (n *SwitchStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "switch %c ; %c %c", emptyNil(n.Init), emptyNil(n.Tag), &n.Body)
}

// Format formats a TypeAssertExpr as a snippet of Go source code.
func (n *TypeAssertExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	if n.Type == nil {
		scPrintf(s, c, "%c.(type)", n.X)
	} else {
		scPrintf(s, c, "%c.(%c)", n.X, n.Type)
	}
}

// Format formats a TypeSpec as a snippet of Go source code.
func (n *TypeSpec) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c%c %c %c", n.Doc, &n.Name, emptyNil(n.Type), n.Comment)
}

// Format formats a TypeSwitchStmt as a snippet of Go source code.
func (n *TypeSwitchStmt) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "switch %c ; %c %c",
		genFormatter{strings.TrimSpace(cSprintf(c, "%c", emptyNil(n.Init)))},
		genFormatter{strings.TrimSpace(cSprintf(c, "%c", n.Assign))},
		&n.Body,
	)
}

// Format formats an UnaryExpr as a snippet of Go source code.
func (n *UnaryExpr) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	scPrintf(s, c, "%c%c", &n.Op, n.X)
}

// Format formats a ValueSpec as a snippet of Go source code.
func (n *ValueSpec) Format(s fmt.State, c rune) {
	if n == nil {
		return
	}
	n.Doc.Format(s, c)
	for i, name := range n.Names {
		name := name
		sepFormat(s, c, i, &name, ", ")
	}
	if len(n.Values) == 0 {
		scPrintf(s, c, " %c", emptyNil(n.Type))
	} else {
		scPrintf(s, c, " %c = ", emptyNil(n.Type))
		for i, value := range n.Values {
			sepFormat(s, c, i, value, ", ")
		}
	}
}
