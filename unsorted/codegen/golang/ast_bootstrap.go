package golang

import "fmt"

// // SliceOfString represents a slice of a file and its 0-based start position.
// type SliceOfString struct {
// 	Text string
// 	Pos  int
// }

// Token ...
type Token struct {
	Kind string
	// Text SliceOfString
	Text string
}

// ArrayType ...
type ArrayType struct {
	Lbrack Token // "["
	Len    Expr  // Ellipsis node for [...]T array types, nil for slice types
	Elt    Expr  // element type
}

// AssignStmt ...
type AssignStmt struct {
	LHS []Expr
	Tok Token
	RHS []Expr
}

// BadDecl ...
type BadDecl struct {
	Token Token
}

// BadExpr ...
type BadExpr struct {
	Token Token
}

// BadStmt ...
type BadStmt struct {
	Token Token
}

// BasicLit ...
type BasicLit struct {
	Token Token
}

// BinaryExpr ...
type BinaryExpr struct {
	X  Expr
	Op Token
	Y  Expr
}

// BlockStmt ...
type BlockStmt struct {
	Lbrace Token // "{"
	List   []Stmt
	Rbrace Token // "}"
}

// BranchStmt ...
type BranchStmt struct {
	Tok   Token  // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
	Label *Ident // label name; or nil
}

// CallExpr ...
type CallExpr struct {
	Fun      Expr   // function expression
	Lparen   Token  // "("
	Args     []Expr // function arguments; or nil
	Ellipsis Token  // "..." (token.NoPos if there is no "...")
	Rparen   Token  // ")"
}

// CaseClause ...
type CaseClause struct {
	Case  Token  // "case" or "default" keyword
	List  []Expr // list of expressions or types; nil means default case
	Colon Token
	Body  []Stmt // statement list; or nil
}

// ChanType ...
type ChanType struct {
	Begin Token // "chan" keyword or "<-" (whichever comes first)
	Arrow Token // "<-" (token.NoPos if there is no "<-"); added in Go 1.1
	// Dir         {|"SEND", "RECV"|} powerset without {||}
	Dir   string
	Value Expr // value type
}

// CommClause ...
type CommClause struct {
	Case  Token  // "case" or "default" keyword
	Comm  Stmt   // send or receive statement; nil means default case
	Colon Token  // ":"
	Body  []Stmt // statement list; or nil
}

// Comment ...
type Comment struct {
	Token Token
}

// CommentGroup ...
type CommentGroup struct {
	List []Comment // len(List) > 0
}

// CompositeLit ...
type CompositeLit struct {
	Type       Expr   // literal type; or nil
	Lbrace     Token  // "{"
	Elts       []Expr // list of composite elements; or nil
	Rbrace     Token  // "}"
	Incomplete bool   // true if (source) expressions are missing in the Elts list; added in Go 1.11
}

// DeclStmt ...
type DeclStmt struct {
	Decl Decl // GenDecl with CONST, TYPE, or VAR token
}

// DeferStmt ...
type DeferStmt struct {
	Defer Token // "defer" keyword
	Call  CallExpr
}

// Ellipsis ...
type Ellipsis struct {
	Ellipsis Token // "..."
	Elt      Expr  // ellipsis element type (parameter lists only); or nil
}

// EmptyStmt ...
type EmptyStmt struct {
	Semicolon Token // following ";"
	Implicit  bool  // if set, ";" was omitted in the source; added in Go 1.5
}

// ExprStmt ...
type ExprStmt struct {
	X Expr // expression
}

// Field ...
type Field struct {
	Doc     *CommentGroup // associated documentation; or nil
	Names   []Ident       // field/method/parameter names; or nil
	Type    Expr          // field/method/parameter type
	Tag     *BasicLit     // field tag; or nil
	Comment *CommentGroup // line comments; or nil
}

func (n Field) WithDoc(comments ...Comment) *Field {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// FieldList ...
type FieldList struct {
	Opening Token   // opening parenthesis/brace, if any
	List    []Field // field list; or nil
	Closing Token   // closing parenthesis/brace, if any
}

// File ...
type File struct {
	Doc        *CommentGroup  // associated documentation; or nil
	Package    Token          // "package" keyword
	Name       Ident          // package name
	Decls      []Decl         // top-level declarations; or nil
	Imports    []ImportSpec   // imports in this file
	Unresolved []Ident        // unresolved identifiers in this file
	Comments   []CommentGroup // list of all comments in the source file
}

func (n File) WithDoc(comments ...Comment) *File {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// ForStmt ...
type ForStmt struct {
	For  Token // "for" keyword
	Init Stmt  // initialization statement; or nil
	Cond Expr  // condition; or nil
	Post Stmt  // post iteration statement; or nil
	Body BlockStmt
}

// FuncDecl ...
type FuncDecl struct {
	Doc  *CommentGroup // associated documentation; or nil
	Recv *FieldList    // receiver (methods); or nil (functions)
	Name Ident         // function/method name
	Type FuncType      // function signature: parameters, results, and "func" keyword
	Body *BlockStmt    // function body; or nil for external (non-Go) function
}

func (n FuncDecl) WithDoc(comments ...Comment) *FuncDecl {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// FuncLit ...
type FuncLit struct {
	Type FuncType  // function type
	Body BlockStmt // function body
}

// FuncType ...
type FuncType struct {
	Func    Token      // "func" keyword (token.NoPos if there is no "func")
	Params  FieldList  // (incoming) parameters; non-nil
	Results *FieldList // (outgoing) results; or nil
}

// GenDecl ...
type GenDecl struct {
	Doc    *CommentGroup // associated documentation; or nil
	Tok    Token         // IMPORT, CONST, TYPE, VAR
	Lparen Token         // '(', if any
	Specs  []Spec
	Rparen Token // ')', if any
}

func (n GenDecl) WithDoc(comments ...Comment) *GenDecl {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// GoStmt ...
type GoStmt struct {
	Go   Token // "go" keyword
	Call CallExpr
}

// Ident ...
type Ident struct {
	Name Token // identifier name
}

// IfStmt ...
type IfStmt struct {
	If   Token // "if" keyword
	Init Stmt  // initialization statement; or nil
	Cond Expr  // condition
	Body BlockStmt
	Else Stmt // else branch; or nil
}

// ImportSpec ...
type ImportSpec struct {
	Doc     *CommentGroup // associated documentation; or nil
	Name    *Ident        // local package name (including "."); or nil
	Path    BasicLit      // import path
	Comment *CommentGroup // line comments; or nil
	EndPos  Token         // end of spec (overrides Path.Pos if nonzero)
}

func (n ImportSpec) WithDoc(comments ...Comment) *ImportSpec {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// IncDecStmt ...
type IncDecStmt struct {
	X   Expr
	Tok Token
}

// IndexExpr ...
type IndexExpr struct {
	X      Expr
	Lbrack Token // "["
	Index  Expr
	Rbrack Token // "]"
}

// InterfaceType ...
type InterfaceType struct {
	Interface  Token     // "interface" keyword
	Methods    FieldList // list of methods
	Incomplete bool      // true if (source) methods are missing in the Methods list
}

// KeyValueExpr ...
type KeyValueExpr struct {
	Key   Expr
	Colon Token
	Value Expr
}

// LabeledStmt ...
type LabeledStmt struct {
	Label Ident
	Colon Token
	Stmt  Stmt
}

// MapType ...
type MapType struct {
	Map   Token // "map" keyword
	Key   Expr
	Value Expr
}

// ParenExpr ...
type ParenExpr struct {
	Lparen Token // "("
	X      Expr  // parenthesized expression
	Rparen Token // ")"
}

// RangeStmt ...
type RangeStmt struct {
	For   Token // "for" keyword
	Key   Expr  // Key may be nil
	Value Expr  // Value may be nil
	Tok   Token // invalid if Key == nil
	X     Expr  // value to range over
	Body  BlockStmt
}

// ReturnStmt ...
type ReturnStmt struct {
	Return  Token  // "return" keyword
	Results []Expr // result expressions; or nil
}

// SelectStmt ...
type SelectStmt struct {
	Select Token     // "select" keyword
	Body   BlockStmt // CommClauses only
}

// SelectorExpr ...
type SelectorExpr struct {
	X   Expr  // expression
	Sel Ident // field selector
}

// SendStmt ...
type SendStmt struct {
	Chan  Expr
	Arrow Token // "<-"
	Value Expr
}

// SliceExpr ...
type SliceExpr struct {
	X      Expr  // expression
	Lbrack Token // "["
	Low    Expr  // begin of slice range; or nil
	High   Expr  // end of slice range; or nil
	Max    Expr  // maximum capacity of slice; or nil; added in Go 1.2
	Slice3 bool  // true if 3-index slice (2 colons present); added in Go 1.2
	Rbrack Token // "]"
}

// StarExpr ...
type StarExpr struct {
	Star Token // "*"
	X    Expr  // operand
}

// StructType ...
type StructType struct {
	Struct     Token     // "struct" keyword
	Fields     FieldList // list of field declarations
	Incomplete bool      // true if (source) fields are missing in the Fields list
}

// SwitchStmt ...
type SwitchStmt struct {
	Switch Token     // "switch" keyword
	Init   Stmt      // initialization statement; or nil
	Tag    Expr      // tag expression; or nil
	Body   BlockStmt // CaseClauses only
}

// TypeAssertExpr ...
type TypeAssertExpr struct {
	X      Expr  // expression
	Lparen Token // "("; added in Go 1.2
	Type   Expr  // asserted type; nil means type switch X.(type)
	Rparen Token // ")"; added in Go 1.2
}

// TypeSpec ...
type TypeSpec struct {
	Doc     *CommentGroup // associated documentation; or nil
	Name    Ident         // type name
	Assign  Token         // '=', if any; added in Go 1.9
	Type    Expr          // Ident, ParenExpr, SelectorExpr, StarExpr, or any of the XxxTypes
	Comment *CommentGroup // line comments; or nil
}

func (n TypeSpec) WithDoc(comments ...Comment) *TypeSpec {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// TypeSwitchStmt ...
type TypeSwitchStmt struct {
	Switch Token     // "switch" keyword
	Init   Stmt      // initialization statement; or nil
	Assign Stmt      // x := y.(type) or y.(type)
	Body   BlockStmt // CaseClauses only
}

// UnaryExpr ...
type UnaryExpr struct {
	Op Token
	X  Expr
}

// ValueSpec ...
type ValueSpec struct {
	Doc     *CommentGroup // associated documentation; or nil
	Names   []Ident       // value names (len(Names) > 0)
	Type    Expr          // value type; or nil
	Values  []Expr        // initial values; or nil
	Comment *CommentGroup // line comments; or nil
}

func (n ValueSpec) WithDoc(comments ...Comment) *ValueSpec {
	n.Doc = &CommentGroup{List: comments}
	return &n
}

// Expr ...
type Expr interface {
	fmt.Formatter
	GoIsExpr()
}

// GoIsExpr ...
func (*BadExpr) GoIsExpr() {}

// GoIsExpr ...
func (*Ident) GoIsExpr() {}

// GoIsExpr ...
func (*Ellipsis) GoIsExpr() {}

// GoIsExpr ...
func (*BasicLit) GoIsExpr() {}

// GoIsExpr ...
func (*FuncLit) GoIsExpr() {}

// GoIsExpr ...
func (*CompositeLit) GoIsExpr() {}

// GoIsExpr ...
func (*ParenExpr) GoIsExpr() {}

// GoIsExpr ...
func (*SelectorExpr) GoIsExpr() {}

// GoIsExpr ...
func (*IndexExpr) GoIsExpr() {}

// GoIsExpr ...
func (*SliceExpr) GoIsExpr() {}

// GoIsExpr ...
func (*TypeAssertExpr) GoIsExpr() {}

// GoIsExpr ...
func (*CallExpr) GoIsExpr() {}

// GoIsExpr ...
func (*StarExpr) GoIsExpr() {}

// GoIsExpr ...
func (*UnaryExpr) GoIsExpr() {}

// GoIsExpr ...
func (*BinaryExpr) GoIsExpr() {}

// GoIsExpr ...
func (*KeyValueExpr) GoIsExpr() {}

// GoIsExpr ...
func (*ArrayType) GoIsExpr() {}

// GoIsExpr ...
func (*StructType) GoIsExpr() {}

// GoIsExpr ...
func (*FuncType) GoIsExpr() {}

// GoIsExpr ...
func (*InterfaceType) GoIsExpr() {}

// GoIsExpr ...
func (*MapType) GoIsExpr() {}

// GoIsExpr ...
func (*ChanType) GoIsExpr() {}

// Stmt ...
type Stmt interface {
	fmt.Formatter
	GoIsStmt()
}

// GoIsStmt ...
func (*AssignStmt) GoIsStmt() {}

// GoIsStmt ...
func (*BadStmt) GoIsStmt() {}

// GoIsStmt ...
func (*BlockStmt) GoIsStmt() {}

// GoIsStmt ...
func (*BranchStmt) GoIsStmt() {}

// GoIsStmt ...
func (*CaseClause) GoIsStmt() {}

// GoIsStmt ...
func (*CommClause) GoIsStmt() {}

// GoIsStmt ...
func (*DeclStmt) GoIsStmt() {}

// GoIsStmt ...
func (*DeferStmt) GoIsStmt() {}

// GoIsStmt ...
func (*EmptyStmt) GoIsStmt() {}

// GoIsStmt ...
func (*ExprStmt) GoIsStmt() {}

// GoIsStmt ...
func (*ForStmt) GoIsStmt() {}

// GoIsStmt ...
func (*GoStmt) GoIsStmt() {}

// GoIsStmt ...
func (*IfStmt) GoIsStmt() {}

// GoIsStmt ...
func (*IncDecStmt) GoIsStmt() {}

// GoIsStmt ...
func (*LabeledStmt) GoIsStmt() {}

// GoIsStmt ...
func (*RangeStmt) GoIsStmt() {}

// GoIsStmt ...
func (*ReturnStmt) GoIsStmt() {}

// GoIsStmt ...
func (*SelectStmt) GoIsStmt() {}

// GoIsStmt ...
func (*SendStmt) GoIsStmt() {}

// GoIsStmt ...
func (*SwitchStmt) GoIsStmt() {}

// GoIsStmt ...
func (*TypeSwitchStmt) GoIsStmt() {}

// Decl ...
type Decl interface {
	fmt.Formatter
	GoIsDecl()
}

// GoIsDecl ...
func (*BadDecl) GoIsDecl() {}

// GoIsDecl ...
func (*FuncDecl) GoIsDecl() {}

// GoIsDecl ...
func (*GenDecl) GoIsDecl() {}

// Spec ...
type Spec interface {
	fmt.Formatter
	GoIsSpec()
}

// GoIsSpec ...
func (*ImportSpec) GoIsSpec() {}

// GoIsSpec ...
func (*TypeSpec) GoIsSpec() {}

// GoIsSpec ...
func (*ValueSpec) GoIsSpec() {}
