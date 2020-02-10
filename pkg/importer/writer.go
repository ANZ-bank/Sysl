package importer

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/sirupsen/logrus"
)

type OutputData struct {
	AppName     string
	Package     string
	SwaggerRoot string
	Mode        string
}

type SyslInfo struct {
	OutputData

	Title       string
	Description string
	OtherFields []string // Ordered key, val pair
}

type MethodEndpoints struct {
	Method    string // one of GET, POST, PUT, OPTION, etc
	Endpoints []Endpoint
}

type writer struct {
	io.Writer
	ind    *IndentWriter
	logger *logrus.Logger

	DisableJSONTags bool
}

const CommentLineLength = 80

func newWriter(out io.Writer, logger *logrus.Logger) *writer {
	return &writer{
		Writer: out,
		ind:    NewIndentWriter("    ", out),
		logger: logger,
	}
}

func (w *writer) Write(info SyslInfo, types TypeList, endpoints ...MethodEndpoints) error {
	if err := w.writeHeader(info); err != nil {
		return err
	}

	for _, method := range endpoints {
		for _, ep := range method.Endpoints {
			w.writeEndpoint(method.Method, ep)
			w.writeLines(BlankLine)
		}
	}

	w.writeDefinitions(types)

	return nil
}

func (w *writer) writeHeader(info SyslInfo) error {
	w.mustWrite(`##########################################
##                                      ##
##  AUTOGENERATED CODE -- DO NOT EDIT!  ##
##                                      ##
##########################################

`)
	title := info.Title

	pkg := ""
	if info.Package != "" {
		pkg = fmt.Sprintf("[package=%s]", quote(info.Package))
	}
	w.writeLines(spaceSeparate(info.AppName, quote(title), pkg)+":", PushIndent)

	for i := 0; i < len(info.OtherFields); i += 2 {
		key := info.OtherFields[i]
		val := info.OtherFields[i+1]
		if val != "" {
			w.writeLines(fmt.Sprintf("@%s = %s", key, quote(val)))
		}
	}
	w.writeLines("@description =:", PushIndent)
	desc := getDescription(info.Description)
	w.writeLines(buildDescriptionLines("| ", desc, CommentLineLength-w.ind.CurrentIndentLen())...)
	w.writeLines(PopIndent, BlankLine)

	return nil
}

func buildDescriptionLines(prefix, description string, wrapAt int) []string {
	var result []string
	scanner := bufio.NewScanner(strings.NewReader(description))
	for scanner.Scan() {
		line := scanner.Text()
		if len(prefix)+len(line) <= wrapAt {
			result = append(result, prefix+line)
		} else {
			words := bufio.NewScanner(strings.NewReader(line))
			words.Split(bufio.ScanWords)
			temp := strings.Builder{}
			temp.WriteString(prefix)
			for words.Scan() {
				w := words.Text()
				if temp.Len()+len(w) > wrapAt {
					result = append(result, strings.TrimSpace(temp.String()))
					temp.Reset()
					temp.WriteString(prefix)
				}
				temp.WriteString(w + " ")
			}
			result = append(result, strings.TrimSpace(temp.String()))
		}
	}
	return result
}

func appendAttributesString(prefix string, attrs []string) string {
	if len(attrs) == 0 {
		return prefix
	}
	return prefix + " [" + strings.Join(attrs, ", ") + "]"
}

func appendSizeSpec(prefix string, spec *sizeSpec) string {
	if spec == nil {
		return prefix
	}
	switch spec.MaxType {
	case MaxSpecified:
		return fmt.Sprintf("%s(%d..%d)", prefix, spec.Min, spec.Max)
	case OpenEnded:
		return fmt.Sprintf("%s(%d..)", prefix, spec.Min)
	case MinOnly:
		return fmt.Sprintf("%s(%d)", prefix, spec.Min)
	}
	return prefix
}

func buildQueryString(params []Param) string {
	query := ""
	if len(params) > 0 {
		var parts []string
		for _, p := range params {
			optional := ""
			if p.Optional {
				optional = "?"
			}
			parts = append(parts, fmt.Sprintf("%s=%s%s", url.QueryEscape(p.Name), getSyslTypeName(p.Type), optional))
		}
		query = " ?" + strings.Join(parts, "&")
	}
	return query
}

func buildResponseContentString(responses []Response) string {
	if len(responses) == 0 {
		return ""
	}

	contentsString := " [responses=["
	responses = removeResponseDuplicates(responses)
	for _, val := range responses {
		if val.Content.contentType != "" && val.Content.name != "" {
			contentsString = contentsString + "[" +
				val.Content.contentType[10:len(val.Content.contentType)-1] + "\", \"" + val.Content.name + "\"], "
		}
	}
	contentsString = contentsString[:len(contentsString)-2]
	contentsString += "]]"
	if contentsString == " [responses]]" {
		return ""
	}
	return contentsString
}

func removeResponseDuplicates(responses []Response) []Response {
	keys := make(map[Content]bool)
	list := []Response{}
	for _, entry := range responses {
		if _, value := keys[entry.Content]; !value {
			keys[entry.Content] = true
			list = append(list, entry)
		}
	}
	return list
}

func buildRequestBodyString(params []Param) string {
	body := ""
	if len(params) > 0 {
		sort.SliceStable(params, func(i, j int) bool {
			return strings.Compare(params[i].Name, params[j].Name) < 0
		})
		var parts []string
		for _, p := range params {
			attrs := appendAttributesString("", append(p.Field.Attributes, "~body"))
			parts = append(parts, fmt.Sprintf("%s <: %s%s", p.Name, getSyslTypeName(p.Type), attrs))
		}
		body = strings.Join(parts, ", ")
	}
	return body
}

func buildRequestHeadersString(params []Param) string {
	headers := ""
	if len(params) > 0 {
		var parts []string
		for _, p := range params {
			optional := map[bool]string{true: "~optional", false: "~required"}[p.Optional]

			safeName := regexp.MustCompile("( |-)+").ReplaceAll([]byte(p.Name), []byte("_"))
			text := fmt.Sprintf("%s <: %s", strings.ToLower(string(safeName)),
				appendAttributesString(getSyslTypeName(p.Type), []string{"~header", optional, "name=" + quote(p.Name)}))
			parts = append(parts, text)
		}
		headers = strings.Join(parts, ", ")
	}
	return headers
}

func buildPathString(path string, params []Param) string {
	result := path

	for _, p := range params {
		replacement := fmt.Sprintf("{%s<:%s}", p.Name, getSyslTypeName(p.Type))
		result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", p.Name), replacement)
	}

	return result
}

func (w *writer) writeEndpoint(method string, endpoint Endpoint) {
	responseContentType := ""
	header := buildRequestHeadersString(endpoint.Params.HeaderParams())
	body := buildRequestBodyString(endpoint.Params.BodyParams())
	if len(endpoint.Responses) != 0 {
		responseContentType = buildResponseContentString(endpoint.Responses)
	}
	reqStr := ""
	if len(header) > 0 && len(body) > 0 {
		reqStr = fmt.Sprintf(" (%s)", strings.Join([]string{body, header}, ", "))
	} else if len(header) > 0 || len(body) > 0 {
		reqStr = fmt.Sprintf(" (%s)", body+header)
	}

	pathStr := buildPathString(endpoint.Path, endpoint.Params.PathParams())
	desc := getDescription(endpoint.Description)
	queryStr := buildQueryString(endpoint.Params.QueryParams())

	w.writeLines(fmt.Sprintf("%s:", pathStr), PushIndent,
		fmt.Sprintf("%s%s%s%s:", method, reqStr, queryStr, responseContentType), PushIndent)
	w.writeLines(buildDescriptionLines("| ", desc, CommentLineLength-w.ind.CurrentIndentLen())...)

	if len(endpoint.Responses) > 0 {
		var outs []string
		for _, resp := range endpoint.Responses {
			var newline string
			switch {
			case resp.Type != nil && resp.Text != "":
				newline = fmt.Sprintf("return %s <: %s", resp.Text, getSyslTypeName(resp.Type))
			case resp.Type != nil:
				newline = "return " + getSyslTypeName(resp.Type)
			default:
				newline = "return " + resp.Text
			}
			if !swag.ContainsStrings(outs, newline) {
				outs = append(outs, newline)
			}
		}
		sort.Strings(outs)
		w.writeLines(outs...)
	}

	w.writeLines(PopIndent, PopIndent)
}

func (w *writer) writeDefinitions(types TypeList) {
	w.writeLines("#" + strings.Repeat("-", 75))
	w.writeLines("# definitions")
	var others []Type
	for _, t := range types.Items() {
		_, isEnum := t.(*Enum)
		switch {
		case isBuiltInType(t):
			// do nothing
		case isUnionType(t):
			w.writeLines(BlankLine)
			w.writeUnion(t)
		case isEnum:
			// We want the enum aliases listed with the real types
			w.writeLines(BlankLine)
			w.writeExternalAlias(t)
		case !isExternalAlias(t):
			w.writeLines(BlankLine)
			w.writeDefinition(t.(*StandardType))
		default:
			others = append(others, t)
		}
	}
	for _, t := range others {
		w.writeLines(BlankLine)
		w.writeExternalAlias(t)
	}
}

func (w *writer) writeDefinition(t *StandardType) {
	bangName := "type"
	w.writeLines(fmt.Sprintf("!%s %s:", bangName, appendAttributesString(getSyslTypeName(t), t.Attributes)))
	for _, prop := range t.Properties {
		suffix := ""
		if prop.Optional {
			suffix = "?"
		}
		suffix = appendAttributesString(suffix, prop.Attributes)

		name := prop.Name
		if IsBuiltIn(name) {
			name += "_"
		}

		if !w.DisableJSONTags {
			suffix += ":"
		}

		w.writeLines(PushIndent, fmt.Sprintf("%s <: %s%s", appendSizeSpec(name, prop.SizeSpec),
			getSyslTypeName(prop.Type), suffix))
		if !w.DisableJSONTags {
			w.writeLines(PushIndent, fmt.Sprintf("@json_tag = %s", quote(prop.Name)), PopIndent)
		}
		w.writeLines(PopIndent)
	}
}

func (w *writer) writeExternalAlias(item Type) {
	aliasType := "string"
	aliasName := getSyslTypeName(item)
	switch t := item.(type) {
	case *StandardType:
		if len(t.Properties) > 0 {
			aliasType = getSyslTypeName(t.Properties[0].Type)
		}
	case *ExternalAlias:
		aliasType = getSyslTypeName(t.Target)
	case *Alias:
		aliasType = getSyslTypeName(t.Target)
	case *Array:
		aliasType = getSyslTypeName(item)
		aliasName = t.name
	}
	w.writeLines(fmt.Sprintf("!alias %s:", aliasName),
		PushIndent, aliasType, PopIndent)
}

func (w *writer) writeUnion(item Type) {
	unionName := getSyslTypeName(item)
	t := item.(*Union)
	w.writeLines(fmt.Sprintf("!union %s:", unionName), PushIndent)
	for _, option := range t.Options {
		w.writeLines(option.Name)
	}
	w.writeLines(PopIndent)
}

func (w *writer) mustWrite(s string) {
	if _, err := w.Writer.Write([]byte(s)); err != nil {
		w.logger.Fatalf("failed to complete write: %s", err.Error())
	}
}

const PushIndent = "&& >>"
const PopIndent = "&& <<"
const BlankLine = "&& !!"

func (w *writer) writeLines(lines ...string) {
	for _, l := range lines {
		switch l {
		case PushIndent:
			w.ind.Push()
		case PopIndent:
			w.ind.Pop()
		case BlankLine:
			w.mustWrite("\n")
		default:
			_ = w.ind.Write() // nolint: errcheck
			w.mustWrite(l + "\n")
		}
	}
}
