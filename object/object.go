package object

import "github.com/dpakach/gorkin/token"

type objectType string

// Object is a interface for all objectTypes in the Parser
type Object interface {
	Type() objectType
}

// ScenarioType is a interface for all Scenario types in the Parser
//
// This may include simple scenarios and scenario outlines
type ScenarioType interface {
	scenarioTypeObject()
}

// FeatureSet is a collection of multiple Features
type FeatureSet struct {
	Features []Feature
}

// Merge joins the features from two FeatureSet into one
func (fs *FeatureSet) Merge(newFs *FeatureSet) {
	fs.Features = append(fs.Features, newFs.Features...)
}

// Feature is the representation of each Feature
type Feature struct {
	Title      string
	Token      token.Token
	Scenarios  []ScenarioType
	Tags       []string
	Background *Background
}

// Background object represents the Background block in the Features
type Background struct {
	Steps []Step
}

// Scenario is the representation of the Scenarios
type Scenario struct {
	Steps        []Step
	Tags         []string
	ScenarioText string
	LineNumber   int
}

func (s *Scenario) scenarioTypeObject() {}

// ScenarioOutline is representation of a scenario outline object
type ScenarioOutline struct {
	Steps        []Step
	Tags         []string
	ScenarioText string
	LineNumber   int
	Table        Table
}

func (s *ScenarioOutline) scenarioTypeObject() {}

// Step is a representation of a Step in Gherkin
type Step struct {
	Token      token.Token
	StepText   string
	Table      Table
	Data       []string
	LineNumber int
}

// Table is a representation of any Table in Gherkin
type Table [][]string
