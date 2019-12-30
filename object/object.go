package object

import "github.com/dpakach/gorkin/token"

type objectType string

const (
	FEATURE_OBJ = "FEATURE_OBJ"
	SCENARIO_OBJ = "SCENARIO_OBJ"
	WHEN_STEP = "WHEN_STEP"
	GIVEN_STEP = "GIVEN_STEP"
	THEN_STEP = "THEN_STEP"
)

type Object interface {
	Type() objectType
}

type ScenarioType interface {
	ScenarioTypeObject()
}

type FeatureSet struct {
	Features []Feature
}

func (fs *FeatureSet) Merge(newFs *FeatureSet) {
	fs.Features = append(fs.Features, newFs.Features...)
}

type Feature struct {
	Title string
	Token token.Token
	Scenarios []ScenarioType
	Tags []string
	Background *Background
}

type Background struct {
	Steps []Step
}

type Scenario struct {
	Steps []Step
	Tags []string
	ScenarioText string
	LineNumber int
}
func (s *Scenario) ScenarioTypeObject() {}

type ScenarioOutline struct {
	Steps []Step
	Tags []string
	ScenarioText string
	LineNumber int
	Table Table
}
func (s *ScenarioOutline) ScenarioTypeObject() {}

type Step struct {
	Token token.Token
	StepText string
	Table Table
	Data []string
	LineNumber int
}

type Table [][]string

