package importer

import (
	"fmt"
	"path"
	"strings"
)

const (
	ModeOpenAPI3 = "openapi3"
	ModeOpenAPI2 = "openapi2"
	ModeOpenAPI  = "openapi"
	ModeSwagger  = "swagger"
	ModeXSD      = "xsd"
	ModeGrammar  = "grammar"
	ModeAuto     = "auto"
)

type Format struct {
	Name      string
	Signature string
	FileExt   []string
}

var SYSL = Format{
	Name:      "sysl",
	Signature: "",
	FileExt:   []string{".sysl"},
}

var XSD = Format{
	Name:      "xsd",
	Signature: ``,
	FileExt:   []string{".xsd", ".xml"},
}

var Grammar = Format{
	Name:      "grammar",
	Signature: "",
	FileExt:   []string{".g"},
}

// OpenAPI3 is identified by the openapi header. -  The value MUST be "3.x.x".
// For more details refer to https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md#oasDocument
var OpenAPI3 = Format{
	Name:      "openapi3",
	Signature: `openapi: "3`,
	FileExt:   []string{".yaml", ".json", ".yml"},
}

// Swagger only has 2.0.0 as the single valid format -  The value MUST be "2.0".
// For more details refer to https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#swaggerObject
var Swagger = Format{
	Name:      "swagger",
	Signature: `swagger: "2.0"`,
	FileExt:   []string{".yaml", ".json", ".yml"},
}

// GuessFileType detects the file based on the extension and the file itself.
// It returns the detected format if successful, or an error, if it has failed.
// It first tries to match the file extensions before checking the files for signatures such as swagger: "2.0"
func GuessFileType(filename string, data []byte, validFormats []Format) (*Format, error) {
	var matchesExt []Format
	ext := path.Ext(filename)
	for _, format := range validFormats {
		for _, formatExt := range format.FileExt {
			if formatExt == ext {
				matchesExt = append(matchesExt, format)
				break
			}
		}
	}

	var matchesSignature []Format
	for _, format := range matchesExt {
		if strings.Contains(string(data), format.Signature) {
			matchesSignature = append(matchesSignature, format)
		}
	}

	if len(matchesSignature) == 1 {
		return &matchesSignature[0], nil
	}

	// We return an error if the number of matches is less than 0 or greater than 1
	return nil, fmt.Errorf("error converting json to yaml for: %s", filename)
}
