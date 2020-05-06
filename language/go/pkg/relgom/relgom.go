package relgom

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/language/go/pkg/codegen"
	"github.com/anz-bank/sysl/pkg/parse"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

// Generate generates a Go model library.
func Generate(fsw codegen.FileSystemWriter, srcFS afero.Fs, filename, modelName string) error {
	module, err := parse.NewParser().Parse(filename, srcFS)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll(os.Stdin)")
	}

	model, has := module.Apps[modelName]
	if !has {
		return fmt.Errorf("model (application) %#v not found in %#v", modelName, filename)
	}

	g := generator{
		modelScope: &modelScope{
			fsw:       fsw,
			model:     model,
			modelName: strings.Join(model.Name.Part, ""),
		},
	}

	return g.generate()
}

type modelScope struct {
	fsw       codegen.FileSystemWriter
	modelName string
	model     *sysl.Application
}

func (g *modelScope) newSourceGenerator() *sourceGenerator {
	return newSourceGenerator(g.fsw, strings.ToLower(g.modelName), g.model)
}

type generator struct {
	*modelScope
}

func sortedSetElements(set map[string]struct{}) []string {
	elements := make([]string, 0, len(set))
	for e := range set {
		elements = append(elements, e)
	}
	sort.Strings(elements)
	return elements
}

func (g *generator) generate() error {
	if err := newModelGenerator(g.modelScope).genFileForSyslModel(); err != nil {
		return err
	}
	for tname, t := range g.model.Types {
		if err := genFileForSyslTypeDecl(g.modelScope, tname, t); err != nil {
			return err
		}
	}
	return nil
}
