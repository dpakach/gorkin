package object

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dpakach/gorkin/token"
	"github.com/dpakach/gorkin/utils"
)

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
	GetTags() []string
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

// GetTags returns tags in the given scenario
func (s *Scenario) GetTags() []string { return s.Tags }

// ScenarioOutline is representation of a scenario outline object
type ScenarioOutline struct {
	Steps        []Step
	Tags         []string
	ScenarioText string
	LineNumber   int
	Table        Table
}

func (so *ScenarioOutline) scenarioTypeObject() {}

// GetTags returns tags in the given scenario outline
func (so *ScenarioOutline) GetTags() []string { return so.Tags }

func (so *ScenarioOutline) getScenarios() []Scenario {
	var scenarios []Scenario
	var steps []Step
	for i, row := range so.Table.GetHash() {
		line := so.Table[i+1][0].LineNumber
		steps = []Step{}
		for _, step := range so.Steps {
			steps = append(steps, *step.substituteExampleTable(row))
		}
		scenarios = append(
			scenarios,
			Scenario{steps, so.Tags, so.ScenarioText, line},
		)
	}
	return scenarios
}

// Step is a representation of a Step in Gherkin
type Step struct {
	Token      token.Token
	StepText   string
	Table      Table
	Data       []string
	LineNumber int
}

// TableData is a representation of a cell in a gherkin Table
type TableData struct {
	Literal    string
	LineNumber int
}

// Table is a representation of any Table in Gherkin
type Table [][]TableData

// TableFromString creates a Table type from given 2D string array
func TableFromString(strTable [][]string, startingLine int) Table {
	var res Table
	var resRow []TableData
	curLine := startingLine
	for _, row := range strTable {
		resRow = []TableData{}
		for _, item := range row {
			resRow = append(resRow, TableData{Literal: item, LineNumber: curLine})
		}
		res = append(res, resRow)
		curLine++
	}
	return res
}

func (s *Step) recompileText() {
	count := 0
	for i := 0; i < len(s.StepText); i++ {
		digits := []byte{}
		if s.StepText[i] >= '0' && s.StepText[i] <= '9' {

			// If the text contains digits add it to s.Data
			for s.StepText[i] >= '0' && s.StepText[i] <= '9' {
				digits = append(digits, s.StepText[i])
				i++
				if i == len(s.StepText) {
					break
				}
			}
			s.StepText = utils.ReplaceNth(s.StepText, string(digits), "{{d}}", 1)

			// Add new Number to the Data array
			s.Data = append(s.Data, "")
			copy(s.Data[count+1:], s.Data[count:])
			s.Data[count] = string(digits)

			// Track `{{` in the text to track the position in Data array
			if s.StepText[i] == '{' && s.StepText[i+1] == '{' {
				count++
			}
		}
	}
}

func (s *Step) substituteExampleTable(row map[string]string) *Step {
	var step = &Step{}
	*step = *s
	step.Table = make([][]TableData, len(s.Table))

	// First substitute the {{<data>}} from the step.StepText
	r := regexp.MustCompile("{{<[a-zA-Z0-9_]*>}}")
	if r.MatchString(s.StepText) {
		sup := r.FindString(s.StepText)
		sup = sup[3 : len(sup)-3]
		step.StepText = strings.ReplaceAll(s.StepText, fmt.Sprintf("{{<%v>}}", sup), row[sup])
		step.recompileText()
	}

	// Then substitute <data> occurances from the step.Data
	r = regexp.MustCompile("<[a-zA-Z0-9_]*>")
	for i, data := range step.Data {
		if r.MatchString(data) {
			sup := r.FindString(data)
			sup = sup[1 : len(sup)-1]
			step.Data[i] = strings.ReplaceAll(step.Data[i], fmt.Sprintf("<%v>", sup), row[sup])
		}
	}

	// Then substitute <data> occurances from the step.Table
	for i, rw := range s.Table {
		for _, rowi := range rw {
			if r.MatchString(rowi.Literal) {
				sup := r.FindString(rowi.Literal)
				sup = sup[1 : len(sup)-1]
				step.Table[i] = append(
					step.Table[i],
					TableData{
						Literal:    strings.ReplaceAll(rowi.Literal, fmt.Sprintf("<%v>", sup), row[sup]),
						LineNumber: rowi.LineNumber,
					},
				)
			} else {
				step.Table[i] = append(
					step.Table[i],
					TableData{
						Literal:    rowi.Literal,
						LineNumber: rowi.LineNumber,
					},
				)
			}
		}
	}
	return step
}

// GetRows returns the rows of a Table as 2D array of TableData
func (t *Table) GetRows() [][]TableData {
	return *t
}

// GetRow returns a row of a Table
func (t *Table) GetRow(i int) ([]TableData, error) {
	rows := t.GetRows()
	if i >= 0 && i < len(rows) {
		return rows[i], nil
	}

	return nil, errors.New("the row you requested does not exist")
}

// GetHash returns the data from a table as array of Hash
func (t *Table) GetHash() []map[string]string {
	keys, err := t.GetRow(0)
	if err != nil {
		panic(err)
	}

	hash := []map[string]string{}

	for _, row := range t.GetRows()[1:] {
		rowMap := map[string]string{}
		for i, key := range keys {
			rowMap[key.Literal] = row[i].Literal
		}
		hash = append(hash, rowMap)
	}

	return hash
}
