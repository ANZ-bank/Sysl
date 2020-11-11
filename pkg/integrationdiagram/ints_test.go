package integrationdiagram

import (
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/diagrams"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const plantumlHeader = `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}`

func TestGenerateIntegrations(t *testing.T) {
	t.Parallel()

	m, err := parse.NewParser().ParseFromFs("demo/simple/sysl-ints.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	require.NoError(t, err)
	require.NotNil(t, m)

	stmt := &sysl.Statement{}
	args := &Args{"", "Project", false, false}
	apps := []string{"System1", "IntegratedSystem", "System2"}
	highlights := syslutil.MakeStrSet("IntegratedSystem", "System1", "System2")
	s1 := AppElement{"IntegratedSystem", "integrated_endpoint_1"}
	t1 := AppElement{"System1", "endpoint"}
	dep1 := AppDependency{
		Self:      s1,
		Target:    t1,
		Statement: stmt,
	}
	deps := []AppDependency{
		dep1,
	}
	endpt := &sysl.Endpoint{
		Name: "_",
		Stmt: []*sysl.Statement{
			{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "IntegratedSystem"}}},
			{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "System1"}}},
		},
	}
	intsParam := &IntsParam{apps, highlights, deps, m.GetApps()["Project"], endpt}
	r := GenerateView(args, intsParam, m)

	require.NotNil(t, r)

	expected := plantumlHeader + `
[IntegratedSystem] as _0 <<highlight>>
[System1] as _1 <<highlight>>
_0 --> _1
@enduml`

	assert.Equal(t, expected, r)
}

func TestGenerateIntegrationsWithTestFile(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "indirect_1.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	expected := map[string]string{
		"all.png":            filepath.Join(testDir, "indirect_1-all-golden.puml"),
		"indirect_arrow.png": filepath.Join(testDir, "indirect_1-indirect_arrow-golden.puml"),
		"my_callers.png":     filepath.Join(testDir, "indirect_1-my_callers-golden.puml"),
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithTestFileAndFilters(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_test.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
		Filter:    "test",
	}
	expected := map[string]string{}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	// Then
	assert.Equal(t, expected, result)
}

func TestGenerateIntegrationsWithImmediatePredecessors(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_immediate_predecessors_test.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}
	expected := map[string]string{
		"immediate_predecessors.png": filepath.Join(testDir, "immediate_predecessors-golden.puml"),
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithExclude(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_excludes_test.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	expected := map[string]string{
		"excludes.png": filepath.Join(testDir, "excludes-golden.puml"),
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithPassthrough(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_passthrough_test.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	expected := map[string]string{
		"passthrough.png": filepath.Join(testDir, "passthrough-golden.puml"),
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithCluster(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_with_cluster.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
		Clustered: true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"cluster.png": filepath.Join(testDir, "cluster-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

// func TestGenerateIntegrationsWithDeepCluster(t *testing.T) {
// 	// Given
// 	args := &IntsArg{
// 		RootModel: testDir,
// 		Modules:   "integration_with_deep_cluster.sysl",
// 		Output:    "%(epname).png",
// 		Project:   "Project",
// 		Clustered: true,
// 	}

// 	// When
// 	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
// 		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
// 		args.Epa)
// 	require.NoError(t, err)

// 	expected := map[string]string{
// 		"cluster_deep.png": filepath.Join(testDir, "cluster_deep-golden.puml"),
// 	}

// 	// Then
// 	diagrams.ComparePUML(t, expected, result)
// }

func TestGenerateIntegrationsWithEpa(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_with_epa.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
		Epa:       true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"epa.png": filepath.Join(testDir, "epa-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithIndirectArrow(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "indirect_2.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"all_indirect_2.png":  filepath.Join(testDir, "all_indirect_2-golden.puml"),
		"no_passthrough.png":  filepath.Join(testDir, "no_passthrough-golden.puml"),
		"passthrough_b.png":   filepath.Join(testDir, "passthrough_b-golden.puml"),
		"passthrough_c.png":   filepath.Join(testDir, "passthrough_c-golden.puml"),
		"passthrough_d.png":   filepath.Join(testDir, "passthrough_d-golden.puml"),
		"passthrough_c_e.png": filepath.Join(testDir, "passthrough_c_e-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithRestrictBy(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_with_restrict_by.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
		Epa:       true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"with_restrict_by.png":    filepath.Join(testDir, "with_restrict_by-golden.puml"),
		"without_restrict_by.png": filepath.Join(testDir, "without_restrict_by-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithFilter(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_with_filter.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
		Filter:    "matched",
	}

	expected := map[string]string{
		"matched.png": filepath.Join(testDir, "matched-golden.puml"),
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationWithOrWithoutPassThrough(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_with_or_without_passthrough.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"with_passthrough.png":    filepath.Join(testDir, "with_passthrough-golden.puml"),
		"without_passthrough.png": filepath.Join(testDir, "without_passthrough-golden.puml"),
		"with_systema.png":        filepath.Join(testDir, "with_systema-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestPassthrough2(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "passthrough_1.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"pass_1_all.png":   filepath.Join(testDir, "pass_1_all-golden.puml"),
		"pass_1_sys_a.png": filepath.Join(testDir, "pass_1_sys_a-golden.puml"),
		"pass_b.png":       filepath.Join(testDir, "pass_b-golden.puml"),
		"pass_b_c.png":     filepath.Join(testDir, "pass_b_c-golden.puml"),
		"pass_f.png":       filepath.Join(testDir, "pass_f-golden.puml"),
		"pass_D.png":       filepath.Join(testDir, "pass_D-golden.puml"),
		"pass_e.png":       filepath.Join(testDir, "pass_e-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithPubSub(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "integration_with_pubsub.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
		Epa:       true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"pubsub.png": filepath.Join(testDir, "pubsub-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestAllStmts(t *testing.T) {
	t.Parallel()

	// Given
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "ints_stmts.sysl",
		Output:    "%(epname).png",
		Project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)

	expected := map[string]string{
		"all_stmts.png": filepath.Join(testDir, "all_stmts-golden.puml"),
	}

	// Then
	diagrams.ComparePUML(t, expected, result)
}

func TestHidden(t *testing.T) {
	t.Parallel()
	args := &IntsArg{
		RootModel: testDir,
		Modules:   "hidden.sysl",
		Output:    "out.puml",
		Project:   "Project",
	}
	result, err := GenerateIntegrationsWithParams(args.RootModel, args.Title, args.Output,
		args.Project, args.Filter, args.Modules, args.Exclude, args.Clustered,
		args.Epa)
	require.NoError(t, err)
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}
[ThirdApp] as _0 <<highlight>>
[SecondApp] as _1 <<highlight>>
_0 --> _1
[App] as _2 <<highlight>>
_0 --> _2
@enduml`
	assert.Equal(t, expected, result["out.puml"])
}

func GenerateIntegrationsWithParams(
	rootModel, title, output, project, filter, modules string,
	exclude []string,
	clustered, epa bool,
) (map[string]string, error) {
	cmdContextParamIntgen := &cmdutils.CmdContextParamIntgen{
		Title:     title,
		Output:    output,
		Project:   project,
		Filter:    filter,
		Exclude:   exclude,
		Clustered: clustered,
		EPA:       epa,
	}

	logger, _ := test.NewNullLogger()
	mod, _, err := loader.LoadSyslModule(rootModel, modules, afero.NewOsFs(), logger)
	if err != nil {
		return nil, err
	}
	return GenerateIntegrations(cmdContextParamIntgen, mod, logger)
}
