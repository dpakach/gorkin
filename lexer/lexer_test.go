package lexer

import "testing"
import "github.com/dpakach/gorkin/token"

func TestNextToken(t *testing.T) {
	input := `
	Feature: hello world

	Scenario: Test Scenario
		Given hello world
		When test test
		Then run test
		But not fail test

	@smokeTest
	Scenario Outline: Another Scenario
		Given hello world is "big"
		When test is 5 times test
		Then <data1> must be <data2>
		Examples:
		| data1  | data2  |
		| value1 | value2 |
		| val1   | val2   |
	`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	} {
		{token.FEATURE, "Feature"},
		{token.COLON, ":"},
		{token.STEP_BODY, "hello world"},
		{token.SCENARIO, "Scenario"},
		{token.COLON, ":"},
		{token.STEP_BODY, "Test Scenario"},
		{token.GIVEN, "Given"},
		{token.STEP_BODY, "hello world"},
		{token.WHEN, "When"},
		{token.STEP_BODY, "test test"},
		{token.THEN, "Then"},
		{token.STEP_BODY, "run test"},
		{token.BUT, "But"},
		{token.STEP_BODY, "not fail test"},

		{token.TAG, "smokeTest"},
		{token.SCENARIO, "Scenario"},
		{token.OUTLINE, "Outline"},
		{token.COLON, ":"},
		{token.STEP_BODY, "Another Scenario"},
		{token.GIVEN, "Given"},
		{token.STEP_BODY, "hello world is"},
		{token.STRING, "big"},
		{token.WHEN, "When"},
		{token.STEP_BODY, "test is"},
		{token.NUMBER, "5"},
		{token.STEP_BODY, "times test"},

		{token.THEN, "Then"},
		{token.TABLE_DATA, "data1"},
		{token.STEP_BODY, "must be"},
		{token.TABLE_DATA, "data2"},

		{token.EXAMPLES, "Examples"},
		{token.COLON, ":"},
		{token.TABLE_DATA, "data1"},
		{token.TABLE_DATA, "data2"},
		{token.TABLE_LINE_BREAK, token.TABLE_LINE_BREAK},
		{token.TABLE_DATA, "value1"},
		{token.TABLE_DATA, "value2"},
		{token.TABLE_LINE_BREAK, token.TABLE_LINE_BREAK},
		{token.TABLE_DATA, "val1"},
		{token.TABLE_DATA, "val2"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
	}
}
