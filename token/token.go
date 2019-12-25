package token

type TokenType int

type Token struct {
	Type TokenType
	Literal string
	LineNumber int
}

const (
	ILLEGAL TokenType = iota
	STRING
	STEP_BODY
	NUMBER

	WHEN
	THEN
	GIVEN
	AND
	BUT

	FEATURE
	SCENARIO
	OUTLINE
	EXAMPLES
	BACKGROUND
	TAG
	EXAMPLE_VALUE
	TABLE_DATA
	LINE_TEXT

	COLON
	COMMENT
	NEW_LINE

	PYSTRING

	EOF
)

func (token TokenType) String() string {
	switch token {
	case FEATURE:
		return "Feature"
	case SCENARIO:
		return "Scenario"
	case BACKGROUND:
		return "Background"
	case WHEN:
		return "When"
	case THEN:
		return "Then"
	case GIVEN:
		return "Given"
	case AND:
		return "And"
	case BUT:
		return "But"
	case OUTLINE:
		return "Outline"
	case EXAMPLES:
		return "Examples"
	case EOF:
		return "EOF"
	case NEW_LINE:
		return "NEW_LINE"
	}

	return "Illegal"
}

var keywords = map[string]TokenType {
	"Feature": FEATURE,
	"Scenario": SCENARIO,
	"When": WHEN,
	"Given": GIVEN,
	"Then": THEN,
	"But": BUT,
	"And": AND,
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

func IsStepToken(t TokenType) bool {
	stepTokens := []TokenType{GIVEN, WHEN, THEN, AND, BUT}
	for _, step := range stepTokens {
        if step == t {
            return true
        }
    }
    return false
}
