package lexer

import "testing"
import "gorkin/token"

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

		`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	} {
		{token.NEW_LINE, token.NEW_LINE},
		{token.FEATURE, "Feature"},
		{token.COLON, ":"},
		{token.STEP_BODY, "hello world"},
		{token.NEW_LINE, token.NEW_LINE},

		{token.NEW_LINE, token.NEW_LINE},
		{token.COMMENT, "this is a comment"},
		{token.NEW_LINE, token.NEW_LINE},
		{token.COMMENT, "this is another comment"},
		{token.NEW_LINE, token.NEW_LINE},

		{token.NEW_LINE, token.NEW_LINE},
		{token.BACKGROUND, "Background"},
		{token.COLON, ":"},
		{token.NEW_LINE, token.NEW_LINE},
		{token.GIVEN, "Given"},
		{token.STEP_BODY, "step is parsed"},
		{token.NEW_LINE, token.NEW_LINE},

		{token.NEW_LINE, token.NEW_LINE},
		{token.SCENARIO, "Scenario"},
		{token.COLON, ":"},
		{token.STEP_BODY, "Test Scenario"},
		{token.NEW_LINE, token.NEW_LINE},
		{token.GIVEN, "Given"},
		{token.STEP_BODY, "hello world"},
		{token.NEW_LINE, token.NEW_LINE},
		{token.WHEN, "When"},
		{token.STEP_BODY, "test test"},
		{token.STRING, "data"},
		{token.NEW_LINE, token.NEW_LINE},
		{token.THEN, "Then"},
		{token.STEP_BODY, "run test"},
		{token.NEW_LINE, token.NEW_LINE},
        {token.BUT, "But"},
        {token.STEP_BODY, "not fail test"},
        {token.NEW_LINE, token.NEW_LINE},

        {token.NEW_LINE, token.NEW_LINE},
        {token.TAG, "smoke"},
        {token.TAG, "anotherTag"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.SCENARIO, "Scenario"},
        {token.OUTLINE, "Outline"},
        {token.COLON, ":"},
        {token.STEP_BODY, "Another Scenario"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.GIVEN, "Given"},
        {token.STEP_BODY, "hello world is"},
        {token.STRING, "big"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.WHEN, "When"},
        {token.STEP_BODY, "test is"},
        {token.NUMBER, "5"},
        {token.STEP_BODY, "times test"},
        {token.NEW_LINE, token.NEW_LINE},

        {token.THEN, "Then"},
        {token.TABLE_DATA, "data1"},
        {token.STEP_BODY, "must be"},
        {token.TABLE_DATA, "data2"},
        {token.NEW_LINE, token.NEW_LINE},

        {token.EXAMPLES, "Examples"},
        {token.COLON, ":"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.TABLE_DATA, "data1"},
        {token.TABLE_DATA, "data2"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.TABLE_DATA, "value1"},
        {token.TABLE_DATA, "value2"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.TABLE_DATA, "val1"},
        {token.TABLE_DATA, "val2"},
        {token.NEW_LINE, token.NEW_LINE},
        {token.NEW_LINE, token.NEW_LINE},
		{token.EOF, token.EOF},
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
