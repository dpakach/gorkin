package lexer

import "testing"
import "github.com/dpakach/gorkin/token"
import "fmt"

func TestNextToken(t *testing.T) {
	input := `
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
		| value1 | value2 | # test comment
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
		| value_1 | value_2 |
		# some comment here
		| val1   | val2   |

		# more comment here
		| val1   | val2   |
		`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
		expectedLineNo  int
	}{
		{token.NEWLINE, token.NEWLINE.String(), 1},
		{token.FEATURE, "Feature", 2},
		{token.COLON, ":", 2},
		{token.STEPBODY, "hello world", 2},
		{token.NEWLINE, token.NEWLINE.String(), 2},

		{token.NEWLINE, token.NEWLINE.String(), 3},
		{token.COMMENT, "this is a comment", 4},
		{token.NEWLINE, token.NEWLINE.String(), 4},
		{token.COMMENT, "this is another comment", 5},
		{token.NEWLINE, token.NEWLINE.String(), 5},

		{token.NEWLINE, token.NEWLINE.String(), 6},
		{token.BACKGROUND, "Background", 7},
		{token.COLON, ":", 7},
		{token.NEWLINE, token.NEWLINE.String(), 7},
		{token.GIVEN, "Given", 8},
		{token.STEPBODY, "step is parsed", 8},
		{token.NEWLINE, token.NEWLINE.String(), 8},

		{token.NEWLINE, token.NEWLINE.String(), 9},
		{token.SCENARIO, "Scenario", 10},
		{token.COLON, ":", 10},
		{token.STEPBODY, "Test Scenario", 10},
		{token.NEWLINE, token.NEWLINE.String(), 10},
		{token.GIVEN, "Given", 11},
		{token.STEPBODY, "hello world", 11},
		{token.NEWLINE, token.NEWLINE.String(), 11},
		{token.WHEN, "When", 12},
		{token.STEPBODY, "test test", 12},
		{token.STRING, "data", 12},
		{token.NEWLINE, token.NEWLINE.String(), 12},
		{token.THEN, "Then", 13},
		{token.STEPBODY, "run test", 13},
		{token.NEWLINE, token.NEWLINE.String(), 13},
		{token.BUT, "But", 14},
		{token.STEPBODY, "not fail test", 14},
		{token.NEWLINE, token.NEWLINE.String(), 14},

		{token.NEWLINE, token.NEWLINE.String(), 15},
		{token.TAG, "smoke", 16},
		{token.TAG, "anotherTag", 16},
		{token.NEWLINE, token.NEWLINE.String(), 16},
		{token.SCENARIO, "Scenario", 17},
		{token.OUTLINE, "Outline", 17},
		{token.COLON, ":", 17},
		{token.STEPBODY, "Another Scenario", 17},
		{token.NEWLINE, token.NEWLINE.String(), 17},
		{token.GIVEN, "Given", 18},
		{token.STEPBODY, "hello world is", 18},
		{token.STRING, "big", 18},
		{token.NEWLINE, token.NEWLINE.String(), 18},
		{token.WHEN, "When", 19},
		{token.STEPBODY, "test is", 19},
		{token.NUMBER, "5", 19},
		{token.STEPBODY, "times test", 19},
		{token.NEWLINE, token.NEWLINE.String(), 19},

		{token.THEN, "Then", 20},
		{token.EXAMPLEVALUE, "data1", 20},
		{token.STEPBODY, "must be", 20},
		{token.EXAMPLEVALUE, "data2", 20},
		{token.NEWLINE, token.NEWLINE.String(), 20},

		{token.EXAMPLES, "Examples", 21},
		{token.COLON, ":", 21},
		{token.NEWLINE, token.NEWLINE.String(), 21},
		{token.TABLEDATA, "data1", 22},
		{token.TABLEDATA, "data2", 22},
		{token.NEWLINE, token.NEWLINE.String(), 22},
		{token.TABLEDATA, "value1", 23},
		{token.TABLEDATA, "value2", 23},
		{token.COMMENT, "test comment", 23},
		{token.NEWLINE, token.NEWLINE.String(), 23},
		{token.TABLEDATA, "val1", 24},
		{token.TABLEDATA, "val2", 24},
		{token.NEWLINE, token.NEWLINE.String(), 24},

		{token.NEWLINE, token.NEWLINE.String(), 25},
		{token.SCENARIO, "Scenario", 26},
		{token.OUTLINE, "Outline", 26},
		{token.COLON, ":", 26},
		{token.STEPBODY, "Third Scenario", 26},
		{token.NEWLINE, token.NEWLINE.String(), 26},
		{token.GIVEN, "Given", 27},
		{token.STEPBODY, "step has some pystrings", 27},
		{token.NEWLINE, token.NEWLINE.String(), 27},
		{token.PYSTRING, `And some string "data" content
		And another line`, 28},
		{token.NEWLINE, token.NEWLINE.String(), 31},
		{token.THEN, "Then", 32},
		{token.STEPBODY, "something happens", 32},
		{token.NEWLINE, token.NEWLINE.String(), 32},
		{token.TABLEDATA, "val1", 33},
		{token.TABLEDATA, "", 33},
		{token.NEWLINE, token.NEWLINE.String(), 33},
		{token.EXAMPLES, "Examples", 34},
		{token.COLON, ":", 34},
		{token.NEWLINE, token.NEWLINE.String(), 34},
		{token.TABLEDATA, "value_1", 35},
		{token.TABLEDATA, "value_2", 35},
		{token.NEWLINE, token.NEWLINE.String(), 35},
		{token.COMMENT, "some comment here", 36},
		{token.NEWLINE, token.NEWLINE.String(), 36},
		{token.TABLEDATA, "val1", 37},
		{token.TABLEDATA, "val2", 37},
		{token.NEWLINE, token.NEWLINE.String(), 37},

		{token.NEWLINE, token.NEWLINE.String(), 38},
		{token.COMMENT, "more comment here", 39},
		{token.NEWLINE, token.NEWLINE.String(), 39},
		{token.TABLEDATA, "val1", 40},
		{token.TABLEDATA, "val2", 40},
		{token.NEWLINE, token.NEWLINE.String(), 40},

		{token.EOF, token.EOF.String(), 41},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Println(tok.Type, tt)
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
