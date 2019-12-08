package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL="ILLEGAL"
	STRING="STRING"
	STEP_BODY="STEP_BODY"
	NUMBER="NUMBER"

	WHEN="WHEN"
	THEN="THEN"
	GIVEN="GIVEN"
	BUT="BUT"

	FEATURE="FEATURE"
	SCENARIO="SCENARIO"
	OUTLINE="OUTLINE"
	EXAMPLES="EXAMPLES"
	BACKGROUND="BACKGROUND"
	TAG="TAG"
	EXAMPLE_VALUE="EXAMPLE_VALUE"
	TABLE_DATA="TABLE_DATA"
	LINE_TEXT="LINE_TEXT"

	COLON=":"
	COMMENT="COMMENT"
	NEW_LINE="NEW_LINE"
)

var keywords = map[string]TokenType {
	"Feature": FEATURE,
	"Scenario": SCENARIO,
	"When": WHEN,
	"Given": GIVEN,
	"Then": THEN,
	"But": BUT,
	"Outline": OUTLINE,
	"Examples": EXAMPLES,
	"Background": BACKGROUND,
}

var GherkinKeyword = keywords

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ILLEGAL
}
