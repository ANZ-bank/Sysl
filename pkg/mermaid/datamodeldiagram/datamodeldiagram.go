package datamodeldiagram

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
)

//cause linter..
const ref = "ref"
const enum = "enum"
const mapS = "map"
const relation = "relation"
const tuple = "tuple"

//externalLink keeps track of the dependencies between two types
type externalLink struct {
	firstType, secondType string
}

//GenerateDataModelDiagram accepts a sysl module and returns a string (and an error if any)
//The resulting string is the mermaid code for the data model for that module
func GenerateDataModelDiagram(m *sysl.Module) (string, error) {
	mapper := syslwrapper.MakeAppMapper(m)
	mapper.IndexTypes()
	mapper.ConvertTypes()
	return generateDataModelDiagramHelper(mapper.SimpleTypes, &[]externalLink{})
}

//generateDataModelDiagramHelper is a helper which has additional arguments which need not be entered by the user
func generateDataModelDiagramHelper(simpleTypes map[string]*syslwrapper.Type,
	externalLinks *[]externalLink) (string, error) {
	var result string
	result = "%% AUTOGENERATED CODE -- DO NOT EDIT!\n\nclassDiagram\n"
	for appType, value := range simpleTypes {
		result += fmt.Sprintf(" class %s {\n", mermaid.CleanString(appType))
		switch value.Type {
		case relation, tuple:
			result += printProperties(value.Properties, appType, externalLinks)
		case enum:
			result += printEnum(value.Enum)
		case mapS:
			result += printMap(value.Properties, appType, value.PrimaryKey, externalLinks)
		default:
			result += ""
		}
		result += " }\n"
	}
	for _, eLink := range *externalLinks {
		result += fmt.Sprintf(" %s <-- %s\n", mermaid.CleanString(eLink.firstType),
			mermaid.CleanString(eLink.secondType))
	}
	return result, nil
}

//printProperties prints appropriate mermaid code for values inside a property
func printProperties(properties map[string]*syslwrapper.Type, appType string, externalLinks *[]externalLink) string {
	var result string
	//TODO: sort elements in map for uniformity
	for typeName, value := range properties {
		switch value.Type {
		case "int", "bool", "float", "decimal", "string", "string_8", "bytes", "date", "datetime", "xml", "uuid":
			result += fmt.Sprintf("  %s %s\n", value.Type, mermaid.CleanString(typeName))
		case ref:
			pair := externalLink{appType, value.Reference}
			if !externalLinksContain(*externalLinks, pair) {
				*externalLinks = append(*externalLinks, pair)
			}
			result += fmt.Sprintf("  %s %s\n", mermaid.CleanString(value.Reference),
				mermaid.CleanString(typeName))
		case "list", "set":
			//TODO: check if list of lists is possible
			if value.Items[0].Type == ref {
				pair := externalLink{appType,
					value.Items[0].Reference}
				if !externalLinksContain(*externalLinks, pair) {
					*externalLinks = append(*externalLinks, pair)
				}
				result += fmt.Sprintf("  List<%s> %s\n", mermaid.CleanString(value.Items[0].Reference),
					mermaid.CleanString(typeName))
			} else {
				result += fmt.Sprintf("  List<%s> %s\n", mermaid.CleanString(value.Items[0].Type),
					mermaid.CleanString(typeName))
			}
		default:
			result += ""
		}
	}
	return result
}

//printMap prints appropriate mermaid code for values inside a map
func printMap(properties map[string]*syslwrapper.Type,
	appType string, primaryKey string, externalLinks *[]externalLink) string {
	var result string
	var type1, type2, name, name1, name2 string
	name = appType
	name1 = primaryKey
	//TODO: sort elements in map for uniformity
	for typeName, value := range properties {
		if typeName == primaryKey {
			type1 = value.Type
		} else {
			name2 = typeName
			type2 = value.Type
		}
		switch value.Type {
		case "int", "bool", "float", "decimal", "string", "string_8", "bytes", "date", "datetime", "xml", "uuid":
			break
		case ref:
			pair := externalLink{mermaid.CleanString(appType), mermaid.CleanString(value.Reference)}
			if !externalLinksContain(*externalLinks, pair) {
				*externalLinks = append(*externalLinks, pair)
			}
		case "list", "set":
			//TODO: check if list of lists is possible
			if value.Items[0].Type == ref {
				pair := externalLink{mermaid.CleanString(appType),
					mermaid.CleanString(value.Items[0].Reference)}
				if !externalLinksContain(*externalLinks, pair) {
					*externalLinks = append(*externalLinks, pair)
				}
			}
		default:
			result += ""
		}
	}
	result += fmt.Sprintf("  Map<%s, %s> %s<%s, %s>\n", type1, type2, name, name1, name2)
	return result
}

//printEnum prints appropriate mermaid code for values inside an enum
func printEnum(enum map[int64]string) string {
	var result string
	//TODO: sort keys in map for uniformity
	for key, value := range enum {
		result += fmt.Sprintf("  %s %d\n", value, key)
	}
	return result
}

//externalLinksContain checks if a connection between two data types exists or not
func externalLinksContain(f []externalLink, s externalLink) bool {
	for _, a := range f {
		if a == s {
			return true
		}
	}
	return false
}
