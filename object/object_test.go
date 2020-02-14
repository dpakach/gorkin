package object

import (
	"github.com/dpakach/gorkin/token"
	"testing"
)

func areArrayEqual(a, b []string) bool {
	if len(a) == 0 {
		if len(b) == 0 {
			return true
		}
		return false
	}
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

func areTablesEqual(a, b Table) bool {
	if (a == nil) && (b == nil) {
		return true
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
			if a[i][j].LineNumber != b[i][j].LineNumber {
				return false
			}
		}
	}
	return true
}

func TestTableGetRows(t *testing.T) {
	expectedTable := TableFromString([][]string{
		[]string{"with", "data"},
		[]string{"4", "5"},
		[]string{"and", "string"},
	}, 1)
	expectedHash := []map[string]string{
		map[string]string{"with": "4", "data": "5"},
		map[string]string{"with": "and", "data": "string"},
	}
	parsed := Table(expectedTable)
	if !areTablesEqual(parsed, expectedTable) {
		t.Fatalf("Expected table to be %v but got %v", expectedTable, parsed)
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
		if !areTablesEqual([][]TableData{row}, [][]TableData{parsedRow}) {
			t.Fatalf("Expected table row to be %v but got %v", row, parsedRow)
		}
	}
}

// When There is step <with> an "<data>"
//  | <with> | 5      |
//  | <data> | string |
//
// Examples:
//  | with | data   |
//  | 4    | 5      |
//  | and  | string |
//
// Result:
// When There is step 4 an "5"
//  | 4 | 5      |
//  | 5 | string |
//
// When There is step and an "string"
//  | and    | 5      |
//  | string | string |
//

func TestSubstituteTable(t *testing.T) {
	table := TableFromString([][]string{
		[]string{"<with>", "5"},
		[]string{"<data>", "string"},
	}, 1)
	hash := []map[string]string{
		map[string]string{"with": "4", "data": "5"},
		map[string]string{"with": "and", "data": "string"},
	}
	step := &Step{
		Token:      token.Token{Type: token.WHEN, Literal: "When", LineNumber: 1},
		StepText:   "There is step {{<with>}} an {{s}}",
		Table:      table,
		Data:       []string{"<data>"},
		LineNumber: 1,
	}
	expected0 := &Step{
		Token:    token.Token{Type: token.WHEN, Literal: "When", LineNumber: 1},
		StepText: "There is step {{d}} an {{s}}",
		Table: TableFromString(
			[][]string{
				[]string{"4", "5"},
				[]string{"5", "string"},
			},
			1,
		),
		Data:       []string{"4", "5"},
		LineNumber: 1,
	}

	res := step.substituteExampleTable(hash[0])
	assertStepsEqual(t, expected0, res)

	expected1 := &Step{
		Token:    token.Token{Type: token.WHEN, Literal: "When", LineNumber: 1},
		StepText: "There is step and an {{s}}",
		Table: TableFromString(
			[][]string{
				[]string{"and", "5"},
				[]string{"string", "string"},
			},
			1,
		),
		Data:       []string{"string"},
		LineNumber: 1,
	}
	res = step.substituteExampleTable(hash[1])
	assertStepsEqual(t, expected1, res)
}

func assertTokensEqual(t *testing.T, actual, expected token.Token) {
	if expected.Type != actual.Type {
		t.Fatalf("Token type does not match, expected: %v, got: %v", expected.Type, actual.Type)
	}
	if expected.Literal != actual.Literal {
		t.Fatalf("Token literal does not match, expected: %v, got: %v", expected.Literal, actual.Literal)
	}
	if expected.LineNumber != actual.LineNumber {
		t.Fatalf("Token line number does not match, expected: %v, got: %v", expected.LineNumber, actual.LineNumber)
	}
}

func assertStepsEqual(t *testing.T, expected, actual *Step) {
	assertTokensEqual(t, expected.Token, actual.Token)
	if expected.StepText != actual.StepText {
		t.Fatalf("Step Text does not match, expected: %v, got: %v", expected.StepText, actual.StepText)
	}
	if !areTablesEqual(expected.Table, actual.Table) {
		t.Fatalf("Step Table does not match, expected: %v, got: %v", expected.Table, actual.Table)
	}
	if !areArrayEqual(expected.Data, actual.Data) {
		t.Fatalf("Step Data does not match, expected: %v, got: %v", expected.Data, actual.Data)
	}
	if expected.LineNumber != actual.LineNumber {
		t.Fatalf("Step line number does not match, expected: %v, got: %v", expected.LineNumber, actual.LineNumber)
	}
}

func assertScenariosEqual(t *testing.T, expected, actual *Scenario) {
	for i := range expected.Steps {
		assertStepsEqual(t, &expected.Steps[i], &actual.Steps[i])
	}
	if !areArrayEqual(expected.Tags, actual.Tags) {
		t.Fatalf("Scenario tags does not match, expected: %v, got: %v", expected.Tags, actual.Tags)
	}

	if expected.ScenarioText != actual.ScenarioText {
		t.Fatalf("Scenario text does not match, expected: %v, got: %v", expected.ScenarioText, actual.ScenarioText)
	}

	if expected.LineNumber != actual.LineNumber {
		t.Fatalf("Scenario line number does not match, expected: %v, got: %v", expected.LineNumber, actual.LineNumber)
	}
}

// Given some test step <with>
// Then some data is "5"
// Then some "<with>" has table
//   | <with> | 5      |
//   | and    | <data> |

var stepDataProvider = []Step{
	{
		token.Token{token.GIVEN, "Given", 1},
		"some test step {{<with>}}",
		nil,
		nil,
		1,
	},
	{
		token.Token{token.THEN, "Then", 1},
		"some data is {{s}}",
		nil,
		[]string{"5"},
		2,
	},
	{
		token.Token{token.THEN, "Then", 1},
		"some {{s}} has a table",
		TableFromString([][]string{
			[]string{"<with>", "5"},
			[]string{"and", "<data>"},
		}, 4),
		[]string{"<with>"},
		3,
	},
}

// Given some test step <with>
// Then some data is "5"
// Then some "<with>" has table
//   | <with> | 5      |
//   | and    | <data> |
//
// Examples:
//    | with | data   |
//    | 4    | 5      |
//    | and  | string |

func TestGetScenarios(t *testing.T) {
	scenarioOutline := &ScenarioOutline{
		stepDataProvider,
		[]string{},
		"Test Scenario",
		1,
		TableFromString([][]string{
			[]string{"with", "data"},
			[]string{"4", "5"},
			[]string{"and", "string"},
		}, 4),
	}

	expectedScenarios := []Scenario{
		Scenario{
			Steps: []Step{
				{
					token.Token{token.GIVEN, "Given", 1},
					"some test step {{d}}",
					nil,
					[]string{"4"},
					1,
				},
				{
					token.Token{token.THEN, "Then", 1},
					"some data is {{s}}",
					nil,
					[]string{"5"},
					2,
				},
				{
					token.Token{token.THEN, "Then", 1},
					"some {{s}} has a table",
					TableFromString([][]string{
						[]string{"4", "5"},
						[]string{"and", "5"},
					}, 4),
					[]string{"4"},
					3,
				},
			},
			Tags:         []string{},
			ScenarioText: "Test Scenario",
			LineNumber:   5,
		},
		Scenario{
			Steps: []Step{
				{
					token.Token{token.GIVEN, "Given", 1},
					"some test step and",
					nil,
					nil,
					1,
				},
				{
					token.Token{token.THEN, "Then", 1},
					"some data is {{s}}",
					nil,
					[]string{"5"},
					2,
				},
				{
					token.Token{token.THEN, "Then", 1},
					"some {{s}} has a table",
					TableFromString([][]string{
						[]string{"and", "5"},
						[]string{"and", "string"},
					}, 4),
					[]string{"and"},
					3,
				},
			},
			Tags:         []string{},
			ScenarioText: "Test Scenario",
			LineNumber:   6,
		},
	}

	scenarios := scenarioOutline.GetScenarios()

	assertScenariosEqual(t, &expectedScenarios[0], &scenarios[0])
	assertScenariosEqual(t, &expectedScenarios[1], &scenarios[1])
}
