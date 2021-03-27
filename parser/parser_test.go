package parser

import (
	"fmt"
	"testing"

	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/object"
	"github.com/dpakach/gorkin/token"
)

type stepDataType struct {
	input         string
	expectedToken token.Type
	expectedBody  string
	expectedData  []string
	expectedTable object.Table
}

const stepInput1 = `Given some test step
					Then some data is 5
					But some "guy" has a table
						| with | data   |
						| 4    | 5      |
						| and  | string |`

const stepInput2 = `Given some other step
					Then some "string" data
					But also table with one row
						| just |
						| one  |
						| row  |`

const stepInput3 = `When running tests
					Then a basic step`

const stepInput4 = `When running tests with <example>
					"""
					This is a basic pystring
					multiline too
					"""`

var stepDataProvider = map[string][]stepDataType{
	"data1": []stepDataType{
		{
			"Given some test step",
			token.GIVEN,
			"some test step",
			nil,
			nil,
		},
		{
			"Then some data is 5",
			token.THEN,
			"some data is {{d}}",
			[]string{"5"},
			nil,
		},
		{
			`But some "guy" has a table
				| with | data   |
				| 4    | 5      |
				| and  | string |`,
			token.BUT,
			"some {{s}} has a table",
			[]string{"guy"},
			object.TableFromString([][]string{
				[]string{"with", "data"},
				[]string{"4", "5"},
				[]string{"and", "string"},
			}, 2),
		},
	},
	"data2": []stepDataType{
		{
			"Given some other step",
			token.GIVEN,
			"some other step",
			nil,
			nil,
		},
		{
			"Then some \"string\" data",
			token.THEN,
			"some {{s}} data",
			[]string{"string"},
			nil,
		},
		{
			`But also table with one row
				| just |
				| one  |
				| row  |`,
			token.BUT,
			"also table with one row",
			nil,
			object.TableFromString([][]string{
				[]string{"just"},
				[]string{"one"},
				[]string{"row"},
			}, 2),
		},
	},
	"data3": []stepDataType{
		{
			"When running tests",
			token.WHEN,
			"running tests",
			nil,
			nil,
		},
		{
			"Then a basic step",
			token.THEN,
			"a basic step",
			nil,
			nil,
		},
	},

	"data4": []stepDataType{
		{
			stepInput4,
			token.WHEN,
			"running tests with {{<example>}}\n{{s}}",
			[]string{
				`This is a basic pystring
					multiline too`,
			},
			nil,
		},
	},
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))

	for _, err := range errors {
		t.Errorf("\nparser error: %q", err.GetMessage())
	}
	t.FailNow()
}

func TestStepParsing(t *testing.T) {
	for _, dataProvider := range stepDataProvider {
		for _, tt := range dataProvider {
			l := lexer.New(tt.input)
			p := New(l)
			parsed := p.ParseStep()
			checkParserErrors(t, p)
			assertStepsEqual(t, parsed, tt)
		}
	}
}

func assertStepsEqual(t *testing.T, actual *object.Step, expected stepDataType) {
	if actual.Token.Type != expected.expectedToken {
		t.Fatalf("Expected Type to be %q, but got %q", expected.expectedToken, actual.Token.Type)
	}
	if actual.StepText != expected.expectedBody {
		t.Fatalf("Expected step text to be %q, but got %q", expected.expectedBody, actual.StepText)
	}
	if expected.expectedTable != nil && !areTablesEqual(actual.Table, expected.expectedTable) {
		t.Fatalf("Expected table to be %q, but got %q", expected.expectedTable, actual.Table)
	}
	if expected.expectedData != nil && !areArrayEqual(actual.Data, expected.expectedData) {
		t.Fatalf("Expected Data to be %q, but got %q", expected.expectedData, actual.Data)
	}
}

func areArrayEqual(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func areMapEqual(a, b map[string]string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func areTablesEqual(a, b object.Table) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if (a[i] == nil) != (b[i] == nil) {
			return false
		}

		if len(a[i]) != len(b[i]) {
			return false
		}

		for j := range a[i] {
			if a[i][j].Literal != b[i][j].Literal {
				return false
			}
			// skipping for now
			//if a[i][j].LineNumber != b[i][j].LineNumber {
			//	return false
			//}
		}
	}
	return true
}

func TestParsingBlockSteps(t *testing.T) {
	input := stepInput1

	l := lexer.New(input)
	p := New(l)

	steps := p.ParseBlockSteps()
	checkParserErrors(t, p)
	assertBlockStepsEqual(t, stepDataProvider["data1"], steps)
}

func assertBlockStepsEqual(t *testing.T, expected []stepDataType, actual []object.Step) {
	if len(expected) != len(actual) {
		t.Fatalf("Number of steps mismatch, expected %v, got %v", len(expected), len(actual))
	}
	for i, tt := range expected {
		assertStepsEqual(t, &actual[i], tt)
	}
}

func TestParsingTags(t *testing.T) {
	input := `@Tag @anotherTag @newTag`
	expected := []string{"Tag", "anotherTag", "newTag"}

	l := lexer.New(input)
	p := New(l)

	tags := p.ParseTags()
	checkParserErrors(t, p)

	if !areArrayEqual(tags, expected) {
		t.Fatalf("Tags mismatch, expedted %q, got %q", expected, tags)
	}
}

func TestParseScenario(t *testing.T) {
	input := fmt.Sprintf(`
	@testTag @randomTag
	Scenario: test Scenario
		%v
	`, stepInput1)
	l := lexer.New(input)
	p := New(l)

	res := p.ParseScenarioType([]string{})
	checkParserErrors(t, p)
	scenario, ok := res.(*object.Scenario)
	if !ok {
		t.Fatalf("Type mismatch, expected Scenario but not got")
	}
	if len(scenario.Steps) != 3 {
		t.Fatalf("Steps length mismatch, expected 3, got %v", len(scenario.Steps))
	}

	for i, tt := range stepDataProvider["data1"] {
		assertStepsEqual(t, &scenario.Steps[i], tt)
	}
	expectedTags := []string{"testTag", "randomTag"}
	if !areArrayEqual(scenario.Tags, expectedTags) {
		t.Fatalf("Tags mismatch, expected %v, got %v", expectedTags, scenario.Tags)
	}
}
func TestParseScenarioOutline(t *testing.T) {
	input := fmt.Sprintf(`
	Scenario Outline: test Scenario
		%v
	Examples:
		| data1  | data2 |
		| value1 | v1    |
		| value2 | 5     |
	`, stepInput1)
	l := lexer.New(input)
	p := New(l)

	res := p.ParseScenarioType([]string{})
	checkParserErrors(t, p)
	scenario, ok := res.(*object.ScenarioOutline)
	if !ok {
		t.Fatalf("Type mismatch, expected Scenario but not got")
	}
	if len(scenario.Steps) != 3 {
		t.Fatalf("Steps length mismatch, expected 3, got %v", len(scenario.Steps))
	}

	for i, tt := range stepDataProvider["data1"] {
		assertStepsEqual(t, &scenario.Steps[i], tt)
	}
	expectedTags := []string{}
	if !areArrayEqual(scenario.Tags, expectedTags) {
		t.Fatalf("Tags mismatch, expected %v, got %v", expectedTags, scenario.Tags)
	}
	expectedTable := object.TableFromString([][]string{
		[]string{"data1", "data2"},
		[]string{"value1", "v1"},
		[]string{"value2", "5"},
	}, 5)

	if !areTablesEqual(expectedTable, scenario.Table) {
		t.Fatalf("Tables mismatch, expected %v, got %v", expectedTable, scenario.Table)
	}
}

func TestParseBackground(t *testing.T) {
	backgrounds := []string{
		fmt.Sprintf(`
		Background:
			%v
		`, stepInput1),
		fmt.Sprintf(`
		Background: background with some text
			%v
		`, stepInput1),
	}

	for _, input := range backgrounds {
		l := lexer.New(input)
		p := New(l)

		background := p.ParseBackground()
		checkParserErrors(t, p)
		if background == nil {
			t.Fatal("Expected background but got nil.")
		}
		if len(background.Steps) != 3 {
			t.Fatalf("Number of steps mismatch, expected 3, got %v", len(background.Steps))
		}
		for i, tt := range stepDataProvider["data1"] {
			assertStepsEqual(t, &background.Steps[i], tt)
		}
	}
}

func TestParseScenarioTypeSet(t *testing.T) {
	input := fmt.Sprintf(`
		Scenario: Scenario test Case 1
			%v

		Scenario: not Outline another test
			%v

		Scenario: test new
			%v
	`, stepInput1, stepInput2, stepInput3)

	l := lexer.New(input)
	p := New(l)

	scenarios := p.ParseScenarioTypeSet()
	checkParserErrors(t, p)
	if scenarios == nil {
		t.Fatal("Expected feature but got nil")
	}

	if len(scenarios) != 3 {
		t.Fatalf("Expected number of scenarios to be 3 but got %v", len(scenarios))
	}
	expected := []struct {
		title           string
		dataProviderKey string
	}{
		{
			title:           "Scenario test Case 1",
			dataProviderKey: "data1",
		},
		{
			title:           "another test",
			dataProviderKey: "data2",
		},
		{
			title:           "test new",
			dataProviderKey: "data3",
		},
	}

	for i, data := range expected {
		assertBlockStepsEqual(t, stepDataProvider[data.dataProviderKey], scenarios[i].(*object.Scenario).Steps)
	}
}

func TestParseScenarioTypeSetWithOutline(t *testing.T) {
	input := fmt.Sprintf(`
		Scenario: Scenario test Case 1
			%v

		Scenario: not Outline another test
			%v

		Scenario Outline: test new
			%v
			Examples:
			 | data |
			 | row  |
			 | row1 |
	`, stepInput1, stepInput2, stepInput3)

	l := lexer.New(input)
	p := New(l)

	scenarios := p.ParseScenarioTypeSet()
	checkParserErrors(t, p)
	if scenarios == nil {
		t.Fatal("Expected feature but got nil")
	}

	if len(scenarios) != 3 {
		t.Fatalf("Expected number of scenarios to be 3 but got %v", len(scenarios))
	}
	expected := []struct {
		title           string
		dataProviderKey string
		outline         bool
	}{
		{
			title:           "Scenario test Case 1",
			dataProviderKey: "data1",
			outline:         false,
		},
		{
			title:           "another test",
			dataProviderKey: "data2",
			outline:         false,
		},
		{
			title:           "test new",
			dataProviderKey: "data3",
			outline:         true,
		},
	}

	for i, data := range expected {
		if !data.outline {
			assertBlockStepsEqual(t, stepDataProvider[data.dataProviderKey], scenarios[i].(*object.Scenario).Steps)
		} else {
			scenario, ok := scenarios[i].(*object.ScenarioOutline)
			if !ok {
				t.Fatalf("Type mismatch, expected Scenario Outline but not got")
			}
			if len(scenario.Steps) != len(stepDataProvider[data.dataProviderKey]) {
				t.Fatalf("Steps length mismatch, expected 3, got %v", len(scenario.Steps))
			}

			for i, tt := range stepDataProvider[data.dataProviderKey] {
				assertStepsEqual(t, &scenario.Steps[i], tt)
			}
			expectedTable := object.TableFromString([][]string{
				[]string{"data"},
				[]string{"row"},
				[]string{"row1"},
			}, 22)

			if !areTablesEqual(expectedTable, scenario.Table) {
				t.Fatalf("Tables mismatch, expected %v, got %v", expectedTable, scenario.Table)
			}
		}
	}
}

func TestParsingFeature(t *testing.T) {
	input := fmt.Sprintf(`
	@coolFeature
	Feature: This is a feature
		Some description about the feature
		Also some more description
		Plus to top it off some extra description

		Background:
			%v

		Scenario: Scenario test Case 1
			%v

		Scenario: not Outline another test
			%v

		Scenario: test new
			%v
	`, stepInput1, stepInput1, stepInput2, stepInput3)

	l := lexer.New(input)
	p := New(l)

	feature := p.ParseFeature()
	checkParserErrors(t, p)

	expectedTitle := "This is a feature"
	if feature.Title != expectedTitle {
		t.Fatalf("Title mismatch, expected %v, got %v", expectedTitle, feature.Title)
	}

	expectedTags := []string{"coolFeature"}

	if !areArrayEqual(expectedTags, feature.Tags) {
		t.Fatalf("Tags mismatch, expected %v, got %v", expectedTags, feature.Tags)
	}

	if feature.Background == nil {
		t.Fatal("Expected background to not be null but got nil")
	}

	assertBlockStepsEqual(t, stepDataProvider["data1"], feature.Background.Steps)

	scenarios := feature.Scenarios
	if scenarios == nil {
		t.Fatal("Expected feature but got nil")
	}

	if len(scenarios) != 3 {
		t.Fatalf("Expected number of scenarios to be 3 but got %v", len(scenarios))
	}
	expected := []struct {
		title           string
		dataProviderKey string
	}{
		{
			title:           "Scenario test Case 1",
			dataProviderKey: "data1",
		},
		{
			title:           "another test",
			dataProviderKey: "data2",
		},
		{
			title:           "test new",
			dataProviderKey: "data3",
		},
	}

	for i, data := range expected {
		assertBlockStepsEqual(t, stepDataProvider[data.dataProviderKey], scenarios[i].(*object.Scenario).Steps)
	}

}

func TestParsingFeatureSet(t *testing.T) {
	input := fmt.Sprintf(`
	@coolFeature
	Feature: This is a feature
		Some description about the feature
		Also some more description
		Plus to top it off some extra description

		Background:
			%v

		Scenario: Scenario test Case 1
			%v

		Scenario: not Outline another test
			%v

		Scenario: test new
			%v
	`, stepInput1, stepInput1, stepInput2, stepInput3)

	l := lexer.New(input)
	p := New(l)

	featureSet := p.Parse()
	checkParserErrors(t, p)
	if len(featureSet.Features) != 1 {
		t.Fatalf("Featureset length mismatch, expected %v, got %v", 3, len(featureSet.Features))
	}

	feature := featureSet.Features[0]

	expectedTitle := "This is a feature"
	if feature.Title != expectedTitle {
		t.Fatalf("Title mismatch, expected %v, got %v", expectedTitle, feature.Title)
	}

	expectedTags := []string{"coolFeature"}

	if !areArrayEqual(expectedTags, feature.Tags) {
		t.Fatalf("Tags mismatch, expected %v, got %v", expectedTags, feature.Tags)
	}

	if feature.Background == nil {
		t.Fatal("Expected background to not be null but got nil")
	}

	assertBlockStepsEqual(t, stepDataProvider["data1"], feature.Background.Steps)

	scenarios := feature.Scenarios
	if scenarios == nil {
		t.Fatal("Expected feature but got nil")
	}

	if len(scenarios) != 3 {
		t.Fatalf("Expected number of scenarios to be 3 but got %v", len(scenarios))
	}
	expected := []struct {
		title           string
		dataProviderKey string
	}{
		{
			title:           "Scenario test Case 1",
			dataProviderKey: "data1",
		},
		{
			title:           "another test",
			dataProviderKey: "data2",
		},
		{
			title:           "test new",
			dataProviderKey: "data3",
		},
	}

	for i, data := range expected {
		assertBlockStepsEqual(t, stepDataProvider[data.dataProviderKey], scenarios[i].(*object.Scenario).Steps)
	}
}

func TestTableGetRows(t *testing.T) {
	input := `| with | data   |
			| 4    | 5      |
			| and  | string |`
	expectedTable := object.TableFromString([][]string{
		[]string{"with", "data"},
		[]string{"4", "5"},
		[]string{"and", "string"},
	}, 1)
	expectedHash := []map[string]string{
		map[string]string{"with": "4", "data": "5"},
		map[string]string{"with": "and", "data": "string"},
	}
	l := lexer.New(input)
	p := New(l)
	parsed := p.ParseTable()
	checkParserErrors(t, p)

	if !areTablesEqual(*parsed, expectedTable) {
		t.Fatalf("Expected table to be %v but got %v", expectedTable, parsed)
	}

	if len(parsed.GetHash()) != 2 {
		t.Fatalf("Expected table hash length to be %v but got %v", 2, len(parsed.GetHash()))
	}
	for i, item := range parsed.GetHash() {
		if !areMapEqual(item, expectedHash[i]) {
			t.Fatalf("Expected table hash to be %v but got %v", expectedHash[i], item)
		}
	}

	for i, row := range expectedTable {
		parsedRow, err := parsed.GetRow(i)
		if err != nil {
			t.Fatal(err)
		}
		if !areTablesEqual([][]object.TableData{row}, [][]object.TableData{parsedRow}) {
			t.Fatalf("Expected table row to be %v but got %v", row, parsedRow)
		}
	}
}
