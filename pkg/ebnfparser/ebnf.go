// Code generated by "ωBNF gen" DO NOT EDIT.
// $ wbnf gen --grammar ebnf.wbnf --start grammar --pkg ebnfparser --output ebnf.go
package ebnfparser

import (
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

func Grammar() parser.Parsers {
	return parser.Grammar{"grammar": parser.Some(parser.Rule(`stmt`)),
		"stmt": parser.Oneof{parser.Rule(`COMMENT`),
			parser.Rule(`prod`)},
		"prod": parser.Seq{parser.Rule(`IDENT`),
			parser.S(":"),
			parser.Some(parser.Rule(`term`)),
			parser.S(";")},
		"term": parser.Stack{parser.Delim{Term: parser.Rule(`@`),
			Sep: parser.Eq("op",
				parser.S("|")),
			Assoc: parser.NonAssociative},
			parser.Some(parser.Rule(`@`)),
			parser.Seq{parser.Rule(`atom`),
				parser.Opt(parser.Eq("quant",
					parser.Oneof{parser.S("*"),
						parser.S("+"),
						parser.S("?")}))}},
		"atom": parser.Oneof{parser.Rule(`STRING`),
			parser.Eq("rule",
				parser.Rule(`IDENT`)),
			parser.Seq{parser.S("("),
				parser.Rule(`term`),
				parser.S(")")}},
		"IDENT":   parser.RE(`[A-Za-z_\.]\w*`),
		"STRING":  parser.RE(`\'([^\']*)\'`),
		"COMMENT": parser.RE(`//.*$|(?s:/\*(?:[^*]|\*+[^*/])\*/)`),
		".wrapRE": parser.RE(`[\s]*()[\s]*`)}.Compile(nil)
}

type Stopper interface {
	ExitNode() bool
	Abort() bool
}
type nodeExiter struct{}

func (n *nodeExiter) ExitNode() bool { return true }
func (n *nodeExiter) Abort() bool    { return false }

type aborter struct{}

func (n *aborter) ExitNode() bool { return true }
func (n *aborter) Abort() bool    { return true }

const (
	IdentCOMMENT = "COMMENT"
	IdentIDENT   = "IDENT"
	IdentSTRING  = "STRING"
	IdentAtom    = "atom"
	IdentOp      = "op"
	IdentProd    = "prod"
	IdentQuant   = "quant"
	IdentRule    = "rule"
	IdentStmt    = "stmt"
	IdentTerm    = "term"
)

type WrapreNode struct{ ast.Node }

func (c WrapreNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkWrapreNode(node WrapreNode, ops WalkerOps) Stopper {
	if fn := ops.EnterWrapreNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	if fn := ops.ExitWrapreNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type CommentNode struct{ ast.Node }

func (c CommentNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkCommentNode(node CommentNode, ops WalkerOps) Stopper {
	if fn := ops.EnterCommentNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	if fn := ops.ExitCommentNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type IdentNode struct{ ast.Node }

func (c IdentNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkIdentNode(node IdentNode, ops WalkerOps) Stopper {
	if fn := ops.EnterIdentNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	if fn := ops.ExitIdentNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type StringNode struct{ ast.Node }

func (c StringNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkStringNode(node StringNode, ops WalkerOps) Stopper {
	if fn := ops.EnterStringNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	if fn := ops.ExitStringNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type AtomNode struct{ ast.Node }

func (c AtomNode) Choice() int {
	return ast.Choice(c.Node)
}

func (c AtomNode) AllString() []StringNode {
	var out []StringNode
	for _, child := range ast.All(c.Node, "STRING") {
		out = append(out, StringNode{child})
	}
	return out
}

func (c AtomNode) OneString() StringNode {
	return StringNode{ast.First(c.Node, "STRING")}
}

func (c AtomNode) AllRule() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "rule") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c AtomNode) OneRule() IdentNode {
	return IdentNode{ast.First(c.Node, "rule")}
}

func (c AtomNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c AtomNode) OneTerm() TermNode {
	return TermNode{ast.First(c.Node, "term")}
}
func WalkAtomNode(node AtomNode, ops WalkerOps) Stopper {
	if fn := ops.EnterAtomNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	for _, child := range node.AllString() {
		s := WalkStringNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	for _, child := range node.AllTerm() {
		s := WalkTermNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	if fn := ops.ExitAtomNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type GrammarNode struct{ ast.Node }

func (c GrammarNode) AllStmt() []StmtNode {
	var out []StmtNode
	for _, child := range ast.All(c.Node, "stmt") {
		out = append(out, StmtNode{child})
	}
	return out
}

func (c GrammarNode) OneStmt() StmtNode {
	return StmtNode{ast.First(c.Node, "stmt")}
}
func WalkGrammarNode(node GrammarNode, ops WalkerOps) Stopper {
	if fn := ops.EnterGrammarNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	for _, child := range node.AllStmt() {
		s := WalkStmtNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	if fn := ops.ExitGrammarNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type ProdNode struct{ ast.Node }

func (c ProdNode) AllIdent() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "IDENT") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c ProdNode) OneIdent() IdentNode {
	return IdentNode{ast.First(c.Node, "IDENT")}
}

func (c ProdNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c ProdNode) OneTerm() TermNode {
	return TermNode{ast.First(c.Node, "term")}
}
func WalkProdNode(node ProdNode, ops WalkerOps) Stopper {
	if fn := ops.EnterProdNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	for _, child := range node.AllIdent() {
		s := WalkIdentNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	for _, child := range node.AllTerm() {
		s := WalkTermNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	if fn := ops.ExitProdNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type StmtNode struct{ ast.Node }

func (c StmtNode) Choice() int {
	return ast.Choice(c.Node)
}

func (c StmtNode) AllComment() []CommentNode {
	var out []CommentNode
	for _, child := range ast.All(c.Node, "COMMENT") {
		out = append(out, CommentNode{child})
	}
	return out
}

func (c StmtNode) OneComment() CommentNode {
	return CommentNode{ast.First(c.Node, "COMMENT")}
}

func (c StmtNode) AllProd() []ProdNode {
	var out []ProdNode
	for _, child := range ast.All(c.Node, "prod") {
		out = append(out, ProdNode{child})
	}
	return out
}

func (c StmtNode) OneProd() ProdNode {
	return ProdNode{ast.First(c.Node, "prod")}
}
func WalkStmtNode(node StmtNode, ops WalkerOps) Stopper {
	if fn := ops.EnterStmtNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	for _, child := range node.AllComment() {
		s := WalkCommentNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	for _, child := range node.AllProd() {
		s := WalkProdNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	if fn := ops.ExitStmtNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type TermNode struct{ ast.Node }

func (c TermNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c TermNode) OneTerm() TermNode {
	return TermNode{ast.First(c.Node, "term")}
}

func (c TermNode) Choice() int {
	return ast.Choice(c.Node)
}

func (c TermNode) AllAtom() []AtomNode {
	var out []AtomNode
	for _, child := range ast.All(c.Node, "atom") {
		out = append(out, AtomNode{child})
	}
	return out
}

func (c TermNode) OneAtom() AtomNode {
	return AtomNode{ast.First(c.Node, "atom")}
}

func (c TermNode) AllOp() []string {
	var out []string
	for _, child := range ast.All(c.Node, "op") {
		out = append(out, ast.First(child, "").Scanner().String())
	}
	return out
}

func (c TermNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return ast.First(child, "").Scanner().String()
	}
	return ""
}

func (c TermNode) AllQuant() []string {
	var out []string
	for _, child := range ast.All(c.Node, "quant") {
		out = append(out, ast.First(child, "").Scanner().String())
	}
	return out
}

func (c TermNode) OneQuant() string {
	if child := ast.First(c.Node, "quant"); child != nil {
		return ast.First(child, "").Scanner().String()
	}
	return ""
}
func WalkTermNode(node TermNode, ops WalkerOps) Stopper {
	if fn := ops.EnterTermNode; fn != nil {
		s := fn(node)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}

	for _, child := range node.AllTerm() {
		s := WalkTermNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	for _, child := range node.AllAtom() {
		s := WalkAtomNode(child, ops)
		switch {
		case s == nil:
		case s.ExitNode():
			return nil
		case s.Abort():
			return s
		}
	}
	if fn := ops.ExitTermNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

type WalkerOps struct {
	EnterWrapreNode  func(WrapreNode) Stopper
	ExitWrapreNode   func(WrapreNode) Stopper
	EnterCommentNode func(CommentNode) Stopper
	ExitCommentNode  func(CommentNode) Stopper
	EnterIdentNode   func(IdentNode) Stopper
	ExitIdentNode    func(IdentNode) Stopper
	EnterStringNode  func(StringNode) Stopper
	ExitStringNode   func(StringNode) Stopper
	EnterAtomNode    func(AtomNode) Stopper
	ExitAtomNode     func(AtomNode) Stopper
	EnterGrammarNode func(GrammarNode) Stopper
	ExitGrammarNode  func(GrammarNode) Stopper
	EnterProdNode    func(ProdNode) Stopper
	ExitProdNode     func(ProdNode) Stopper
	EnterStmtNode    func(StmtNode) Stopper
	ExitStmtNode     func(StmtNode) Stopper
	EnterTermNode    func(TermNode) Stopper
	ExitTermNode     func(TermNode) Stopper
}

func (w WalkerOps) Walk(tree GrammarNode)  { WalkGrammarNode(tree, w) }
func (c GrammarNode) GetAstNode() ast.Node { return c.Node }

func NewGrammarNode(from ast.Node) GrammarNode { return GrammarNode{from} }

func Parse(input *parser.Scanner) (GrammarNode, error) {
	p := Grammar()
	tree, err := p.Parse("grammar", input)
	if err != nil {
		return GrammarNode{nil}, err
	}
	return GrammarNode{ast.FromParserNode(p.Grammar(), tree)}, nil
}

func ParseString(input string) (GrammarNode, error) {
	return Parse(parser.NewScanner(input))
}
