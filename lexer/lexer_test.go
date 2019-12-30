package lexer

import "testing"
import "github.com/dpakach/gorkin/token"

func TestNextToken(t *testing.T) {
	input :=`
	Feature: hello world

	# this is a comment
	# this is another comment

	Background:
		Given step is parsed

	Scenario: Test Scenario
		Given hello world
		When test test "data"
		Then run test
		But not fail test

	@smoke @anotherTag
	Scenario Outline: Another Scenario
		Given hello world is "big"
		When test is 5 times test
		Then <data1> must be <data2>
		Examples:
		| data1  | data2  |
		| value1 | value2 |
		| val1   | val2   |

	Scenario Outline: Third Scenario
		Given step has some pystrings
		"""
		And some string "data" content
		And another line
		"""
		Then something happens
		| val1   |  |
		Examples:
		| value1 | value2 |
		| val1   | val2   |
		`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
		expectedLineNo int
	} {
		{token.NEW_LINE, token.NEW_LINE.String(), 1},
		{token.FEATURE, "Feature", 2},
		{token.COLON, ":", 2},
		{token.STEP_BODY, "hello world", 2},
		{token.NEW_LINE, token.NEW_LINE.String(), 2},

		{token.NEW_LINE, token.NEW_LINE.String(), 3},
		{token.COMMENT, "this is a comment", 4},
		{token.NEW_LINE, token.NEW_LINE.String(), 4},
		{token.COMMENT, "this is another comment", 5},
		{token.NEW_LINE, token.NEW_LINE.String(), 5},

		{token.NEW_LINE, token.NEW_LINE.String(), 6},
		{token.BACKGROUND, "Background", 7},
		{token.COLON, ":", 7},
		{token.NEW_LINE, token.NEW_LINE.String(), 7},
		{token.GIVEN, "Given", 8},
		{token.STEP_BODY, "step is parsed", 8},
		{token.NEW_LINE, token.NEW_LINE.String(), 8},

		{token.NEW_LINE, token.NEW_LINE.String(), 9},
		{token.SCENARIO, "Scenario", 10},
		{token.COLON, ":", 10},
		{token.STEP_BODY, "Test Scenario", 10},
		{token.NEW_LINE, token.NEW_LINE.String(), 10},
		{token.GIVEN, "Given", 11},
		{token.STEP_BODY, "hello world", 11},
		{token.NEW_LINE, token.NEW_LINE.String(), 11},
		{token.WHEN, "When", 12},
		{token.STEP_BODY, "test test", 12},
		{token.STRING, "data", 12},
		{token.NEW_LINE, token.NEW_LINE.String(), 12},
		{token.THEN, "Then", 13},
		{token.STEP_BODY, "run test", 13},
		{token.NEW_LINE, token.NEW_LINE.String(), 13},
        {token.BUT, "But", 14},
        {token.STEP_BODY, "not fail test", 14},
        {token.NEW_LINE, token.NEW_LINE.String(), 14},

        {token.NEW_LINE, token.NEW_LINE.String(), 15},
        {token.TAG, "smoke", 16},
        {token.TAG, "anotherTag", 16},
        {token.NEW_LINE, token.NEW_LINE.String(), 16},
        {token.SCENARIO, "Scenario", 17},
        {token.OUTLINE, "Outline", 17},
        {token.COLON, ":", 17},
        {token.STEP_BODY, "Another Scenario", 17},
        {token.NEW_LINE, token.NEW_LINE.String(), 17},
        {token.GIVEN, "Given", 18},
        {token.STEP_BODY, "hello world is", 18},
        {token.STRING, "big", 18},
        {token.NEW_LINE, token.NEW_LINE.String(), 18},
        {token.WHEN, "When", 19},
        {token.STEP_BODY, "test is", 19},
        {token.NUMBER, "5", 19},
        {token.STEP_BODY, "times test", 19},
        {token.NEW_LINE, token.NEW_LINE.String(), 19},

        {token.THEN, "Then", 20},
        {token.TABLE_DATA, "data1", 20},
        {token.STEP_BODY, "must be", 20},
        {token.TABLE_DATA, "data2", 20},
        {token.NEW_LINE, token.NEW_LINE.String(), 20},

        {token.EXAMPLES, "Examples", 21},
        {token.COLON, ":", 21},
        {token.NEW_LINE, token.NEW_LINE.String(), 21},
        {token.TABLE_DATA, "data1", 22},
        {token.TABLE_DATA, "data2", 22},
        {token.NEW_LINE, token.NEW_LINE.String(), 22},
        {token.TABLE_DATA, "value1", 23},
        {token.TABLE_DATA, "value2", 23},
        {token.NEW_LINE, token.NEW_LINE.String(), 23},
        {token.TABLE_DATA, "val1", 24},
        {token.TABLE_DATA, "val2", 24},
        {token.NEW_LINE, token.NEW_LINE.String(), 24},

		{token.NEW_LINE, token.NEW_LINE.String(), 25},
		{token.SCENARIO, "Scenario", 26},
		{token.OUTLINE, "Outline", 26},
		{token.COLON, ":", 26},
		{token.STEP_BODY, "Third Scenario", 26},
		{token.NEW_LINE, token.NEW_LINE.String(), 26},
		{token.GIVEN, "Given", 27},
		{token.STEP_BODY, "step has some pystrings", 27},
		{token.NEW_LINE, token.NEW_LINE.String(), 27},
		{token.PYSTRING, `And some string "data" content
		And another line`, 28},
		{token.NEW_LINE, token.NEW_LINE.String(), 31},
		{token.THEN, "Then", 32},
		{token.STEP_BODY, "something happens", 32},
		{token.NEW_LINE, token.NEW_LINE.String(), 32},
		{token.TABLE_DATA, "val1", 33},
		{token.TABLE_DATA, "", 33},
		{token.NEW_LINE, token.NEW_LINE.String(), 33},
		{token.EXAMPLES, "Examples", 34},
		{token.COLON, ":", 34},
		{token.NEW_LINE, token.NEW_LINE.String(), 34},
		{token.TABLE_DATA, "value1", 35},
		{token.TABLE_DATA, "value2", 35},
		{token.NEW_LINE, token.NEW_LINE.String(), 35},
		{token.TABLE_DATA, "val1", 36},
		{token.TABLE_DATA, "val2", 36},
		{token.NEW_LINE, token.NEW_LINE.String(), 36},
		{token.EOF, token.EOF.String(), 37},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
		if tok.LineNumber != tt.expectedLineNo {
			t.Fatalf("tests[%d] - Line Number wrong. expected=%v, got=%v", i, tt.expectedLineNo, tok.LineNumber)
		}
	}
}
