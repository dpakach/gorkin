package token

// Type represents the type of an token
type Type int

// Token represents each token in parsing through Gherkin
type Token struct {
	Type       Type
	Literal    string
	LineNumber int
}

// TokenTypes used in Gherkin
const (
	// Basic Types and symbols
	ILLEGAL Type = iota
	STRING
	STEPBODY
	NUMBER
	COLON
	COMMENT
	NEWLINE
	PYSTRING

	// Steps
	WHEN
	THEN
	GIVEN
	AND
	BUT

	// Data Structures types in gherkin
	FEATURE
	SCENARIO
	OUTLINE
	EXAMPLES
	BACKGROUND
	TAG
	EXAMPLEVALUE
	TABLEDATA
	LINETEXT

	// Eof token
	EOF
)

func (token Type) String() string {
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
	case NEWLINE:
		return "NEW_LINE"
	}

	return "Illegal"
}

var keywords = map[string]Type{
	"Feature":    FEATURE,
	"Scenario":   SCENARIO,
	"When":       WHEN,
	"Given":      GIVEN,
	"Then":       THEN,
	"But":        BUT,
	"And":        AND,
	"Outline":    OUTLINE,
	"Examples":   EXAMPLES,
	"Background": BACKGROUND,
}

// GherkinKeyword represents available keywords in Gherkin and their token id
var GherkinKeyword = keywords

// LookupIdent returns Type for given Identifier
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ILLEGAL
}

// IsStepToken checks if given Type is a "Step" token
func IsStepToken(t Type) bool {
	stepTokens := []Type{GIVEN, WHEN, THEN, AND, BUT}
	for _, step := range stepTokens {
		if step == t {
			return true
		}
	}
	return false
}
