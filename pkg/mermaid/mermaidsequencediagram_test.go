package mermaid

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

func TestBadInputsToGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appname := "wrongname"
	epname := "wrongep"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.Empty(t, r)
	assert.Error(t, err)
}

func TestGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appname := "WebFrontend"
	epname := "RequestProfile"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>WebFrontend: RequestProfile
 WebFrontend->>+Api: GET /users/{user_id}/profile
 Api->>+Database: QueryUser
 Database-->>-Api: User
 Api-->>-WebFrontend: UserProfile
 WebFrontend-->>...: Profile Page
`
	assert.Equal(t, expected, r)
}

func TestGenerateMermaidSequenceDiagram2(t *testing.T) {
	t.Parallel()
	appname := "WebFrontend"
	epname := "RequestProfile"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>WebFrontend: RequestProfile
 WebFrontend->>+Api: GET /users/{user_id}/profile
 Api->>+Database: QueryUser
 Database-->>-Api: User [~y, x="1"]
 Api-->>-WebFrontend: UserProfile
 WebFrontend->>+WebFrontend: FooBar
 WebFrontend-->>...: Profile Page
`
	assert.Equal(t, expected, r)
}

func TestGenerateMermaidSequenceDiagram3(t *testing.T) {
	t.Parallel()
	appname := "MobileApp"
	epname := "Login"
	m, err := parse.NewParser().Parse("demo/simple/sysl-app-hyperlink.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>MobileApp: Login
 MobileApp->>+Server: LoginRequest
 Server-->>-MobileApp: MobileApp.LoginResponse
`
	assert.Equal(t, expected, r)
}

func TestGenerateMermaidSequenceDiagramWithIfElseLoopActionAndGroupStatements(t *testing.T) {
	t.Parallel()
	appname := "DragonEater"
	epname := "EatDragon"
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>DragonEater: EatDragon
 DragonEater->>+Actions: EatDragon
 Actions-->>-DragonEater: Firebreath
 alt DragonEater got Firebreath
  DragonEater->>+Actions: SpreadFirebreath
  Actions->>+TheWorld: WTFBro
  TheWorld->>TheWorld: DoNothing
  Actions->>+HouseLannister: FirstResponse
  HouseLannister->>HouseLannister: DenyEverything
  loop until further notice
   HouseLannister->>+Actions: SilenceLittleBirds
   Actions->>Actions: ...
   HouseLannister->>+Actions: BuildTemplesVeryQuickly
   Actions->>Actions: HowdYouDoThat
   Actions-->>-HouseLannister: Temple
  end
  HouseLannister-->>-Actions: Nothing
  Actions->>+HouseLannister: SecondResponse
  alt Firebreath still exists
   HouseLannister->>HouseLannister: Lockdown
  end
  HouseLannister-->>-Actions: Insult
  Actions->>+TheWorld: WereScrewed
  TheWorld->>+Actions: Giveup
  Actions->>Actions: LearnArchery
  Actions->>Actions: FeedTheWolves
  Actions->>Actions: GameOver
  Actions-->>-DragonEater: Firebreath
 else if sdsd
  DragonEater->>DragonEater: norhinf
 else
  DragonEater->>DragonEater: go back to normal
 end
 DragonEater-->>...: Firebreath
`
	assert.Equal(t, expected, r)
}