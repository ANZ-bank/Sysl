package endpointanalysisdiagram

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

const projectDir = "../../../"

//Keeps track of the statement-endpoint pairs we visit during execution
type externalLink struct {
	statement, endPoint string
}

//Accepts the sysl module and returns a string (and an error if any)
//The resulting string is the mermaid code for the endpoint analysis for that application and endpoint
func GenerateEndpointAnalysisDiagram(m *sysl.Module) (string, error) {
	return generateEndpointAnalysisDiagram(m, &[]externalLink{}, true)
}

//This is a helper which has additional arguments which need not be entered by the user
func generateEndpointAnalysisDiagram(m *sysl.Module, externalLinks *[]externalLink, theStart bool) (string, error) {
	var result string
	if theStart {
		result = "%% AUTOGENERATED CODE -- DO NOT EDIT!\n\ngraph TD\n"
	}
	for appName, app := range m.Apps {
		result += fmt.Sprintf(" subgraph %s\n", appName)
		for epName, endPoint := range app.Endpoints {
			statements := endPoint.Stmt
			result += printEndpointAnalysisStatements(m, statements, cleanString(epName), externalLinks)
		}
		result += " end\n"
	}
	for _, eLink := range *externalLinks {
		result += fmt.Sprintf(" %s-->%s\n", eLink.statement, eLink.endPoint)
	}
	return result, nil
}

//This function is used to print the mermaid code for different sysl statements
func printEndpointAnalysisStatements(m *sysl.Module, statements []*sysl.Statement,
	endPoint string, externalLinks *[]externalLink) string {
	var result string
	for _, statement := range statements {
		switch c := statement.Stmt.(type) {
		case *sysl.Statement_Call:
			appEndPoint := fmt.Sprintf("%s-%s", cleanString(c.Call.Target.Part[0]), cleanString(c.Call.Endpoint))
			result += fmt.Sprintf(" %s-->%s\n", endPoint, appEndPoint)
			pair := externalLink{appEndPoint, cleanString(c.Call.Endpoint)}
			if !externalLinksContain(*externalLinks, pair) {
				*externalLinks = append(*externalLinks, pair)
			}
		case *sysl.Statement_Group:
			result += printEndpointAnalysisStatements(m, c.Group.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Cond:
			result += printEndpointAnalysisStatements(m, c.Cond.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Loop:
			result += printEndpointAnalysisStatements(m, c.Loop.Stmt, endPoint, externalLinks)
		case *sysl.Statement_LoopN:
			result += printEndpointAnalysisStatements(m, c.LoopN.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Foreach:
			result += printEndpointAnalysisStatements(m, c.Foreach.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Action:
			result += fmt.Sprintf(" %s-->%s\n", endPoint, cleanString(c.Action.Action))
		case *sysl.Statement_Ret:
			result += fmt.Sprintf(" %s-->%s\n", endPoint, cleanString(c.Ret.Payload))
		default:
			panic("Unrecognised statement type")
		}
	}
	return result
}

//Checks if the statement-endpoint group have been already visited or not
func externalLinksContain(i []externalLink, ip externalLink) bool {
	for _, a := range i {
		if a == ip {
			return true
		}
	}
	return false
}

//Replaces certain characters in the string suitable for mermaid
//TODO:Replace more characters if necessary
func cleanString(temp string) string {
	temp = strings.ReplaceAll(temp, " ", "")
	temp = strings.ReplaceAll(temp, "{", "_")
	temp = strings.ReplaceAll(temp, "}", "_")
	temp = strings.ReplaceAll(temp, "[", "_")
	temp = strings.ReplaceAll(temp, "]", "_")
	temp = strings.ReplaceAll(temp, "\"", "")
	temp = strings.ReplaceAll(temp, "~", "")
	return temp
}
