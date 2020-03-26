package parser

import (
	"strconv"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func (p *pscope) buildStatement(node StmtNode) *sysl.Statement {
	res := &sysl.Statement{}
	if stmt := node.OneIfElse(); stmt != nil {

	} else if stmt := node.OneForStmt(); stmt != nil {

	} else if stmt := node.OneRetStmt(); stmt != nil {
		res.Stmt = &sysl.Statement_Ret{Ret: &sysl.Return{Payload: stmt.OneRetVal().Scanner().String()}}
	} else if stmt := node.OneCallStmt(); stmt != nil {
		call := &sysl.Statement_Call{Call: &sysl.Call{
			Target:   &sysl.AppName{Part: []string{stmt.OneTarget().OneToken()}},
			Endpoint: stmt.OneTargetEndpoint().String(),
			Arg:      nil,
		}}
		if names := stmt.OneTarget().AllName(); len(names) > 0 {
			var target []string
			for _, n := range names {
				target = append(target, n.String())
			}
			call.Call.Target.Part = target
		}
		// todo call_args
		res.Stmt = call
	} else if stmt := node.OneOneOfStmt(); stmt != nil {

	} else if stmt := node.OneHttpMethodComment(); stmt != nil {

	} else if stmt := node.OneGroupStmt(); stmt != nil {

	} else if stmt := node.OneTextStmt(); stmt != nil {
		action := ""
		if txt := stmt.OneDocString(); txt != nil {
			action = txt.Scanner().String()
		} else if txt := stmt.OneQstring(); txt != nil {
			action, _ = strconv.Unquote(txt.Scanner().String())
		} else if txt := stmt.OneText(); txt != nil {
			action = txt.Scanner().String()
		}
		res.Stmt = &sysl.Statement_Action{Action: &sysl.Action{Action: action}}
	} else if stmt := node.OneAnnotation(); stmt != nil {

	}

	if attrs := node.OneAttribs(); attrs != nil {
		res.Attrs = p.buildAttributes(*attrs)
	}

	return res
}
