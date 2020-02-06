package parser

import (
	"fmt"
	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/object"
	"github.com/dpakach/gorkin/token"
	"strings"
)

// Parser is the Object that represents Parsing through lexer input
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []ParsingError
}

// ParsingError is error object representing the Parsing errors
type ParsingError interface {
	GetMessage() string
	parserErrorType()
}

// GeneralParserError is error object representing normal Parsing errors
type GeneralParserError struct {
	parser     Parser
	LineNumber int
	Message    string
}

// GetMessage returns the formatted error message
func (p *GeneralParserError) GetMessage() string {
	return fmt.Sprintf(
		"Parser Error: %v:%v %v",
		p.parser.l.FilePath,
		p.LineNumber,
		p.Message,
	)
}

func (p *GeneralParserError) parserErrorType() {}

// PeekError is error object representing peek errors
type PeekError struct {
	parser            *Parser
	LineNumber        int
	ExpectedTokenType token.Type
	ActualToken       token.Token
}

// GetMessage returns the formatted error message
func (p *PeekError) GetMessage() string {
	return fmt.Sprintf(
		"Parser Error: %v:%v Expected token to be %q but got %q",
		p.parser.l.FilePath,
		p.LineNumber,
		p.ExpectedTokenType,
		p.ActualToken.Type,
	)
}

func (p *PeekError) parserErrorType() {}

// Parser Helper functions

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func isStepToken(t token.Token) bool {
	steps := []token.Type{
		token.WHEN,
		token.THEN,
		token.GIVEN,
		token.AND,
		token.BUT,
	}

	for _, token := range steps {
		if token == t.Type {
			return true
		}
	}

	return false
}

// Errors returns the errors in the parser
func (p *Parser) Errors() []ParsingError {
	return p.errors
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) expectPeekTokens(tokens ...token.Type) bool {
	for _, t := range tokens {
		res := p.expectPeek(t)
		if res == false {
			return res
		}
	}
	return true
}

func (p *Parser) peekError(t token.Type) {
	p.errors = append(p.errors, &PeekError{parser: p, LineNumber: p.peekToken.LineNumber, ExpectedTokenType: t})
}

func (p *Parser) getParserErrors() []string {
	var errors []string
	for _, err := range p.Errors() {
		errors = append(errors, err.GetMessage())
	}
	return errors
}

// New creates a new Parser for given lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []ParsingError{},
	}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) skipNewLines() {
	for p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.COMMENT) {
		p.nextToken()
	}
}

// Parse Functions

// Parse Parses the parser object and returns a featureSet
func (p *Parser) Parse() *object.FeatureSet {

	// TODO:Implement Parsing Multiple features from different Files

	featureSet := &object.FeatureSet{}
	p.skipNewLines()
	if !(p.curTokenIs(token.FEATURE) || p.curTokenIs(token.TAG)) {
		p.peekError(token.FEATURE)
		return nil
	}
	feature := p.ParseFeature()
	if feature == nil {
		return nil
	}
	featureSet.Features = append(featureSet.Features, *feature)
	return featureSet
}

// ParseFeature parses Feature from the current position in the parser
func (p *Parser) ParseFeature() *object.Feature {
	feature := &object.Feature{}
	var tags []string
	p.skipNewLines()
	if p.curTokenIs(token.TAG) {
		tags = p.ParseTags()
		feature.Tags = tags
		if tags == nil {
			return nil
		}
	}
	p.skipNewLines()
	if !p.curTokenIs(token.FEATURE) {
		p.peekError(token.FEATURE)
		return nil
	}
	feature.Token = p.curToken
	if !p.expectPeek(token.COLON) {
		p.peekError(token.COLON)
		return nil
	}
	p.nextToken()
	if p.curTokenIs(token.STEPBODY) {
		feature.Title = p.curToken.Literal
	}
	if !p.expectPeek(token.NEWLINE) {
		p.peekError(token.NEWLINE)
		return nil
	}
	p.skipNewLines()

	for !(p.curTokenIs(token.BACKGROUND) ||
		p.curTokenIs(token.SCENARIO) ||
		p.curTokenIs(token.TAG)) {
		p.nextToken()
		p.skipNewLines()
	}

	feature.Background = nil
	if p.curTokenIs(token.BACKGROUND) {
		background := p.ParseBackground()
		if background == nil {
			return nil
		}
		feature.Background = background
	}
	p.skipNewLines()

	var scenarios []object.ScenarioType
	if p.curTokenIs(token.SCENARIO) || p.curTokenIs(token.TAG) {
		scenarios = p.ParseScenarioTypeSet()
	}
	feature.Scenarios = scenarios
	return feature
}

// ParseBackground parses Background object from the current position in the parser
func (p *Parser) ParseBackground() *object.Background {
	p.skipNewLines()
	background := &object.Background{}
	if !p.curTokenIs(token.BACKGROUND) {
		p.peekError(token.BACKGROUND)
		return nil
	}
	if !p.expectPeekTokens(token.COLON, token.NEWLINE) {
		return nil
	}
	p.skipNewLines()
	if !isStepToken(p.curToken) {
		msg := fmt.Sprintf("Expected token to be a STEP_TYPE but got %s", p.curToken.Type)
		p.errors = append(p.errors, &GeneralParserError{parser: *p, LineNumber: p.curToken.LineNumber, Message: msg})
		return nil
	}
	steps := p.ParseBlockSteps()
	if steps == nil {
		return nil
	}
	background.Steps = steps

	return background
}

// ParseBlockSteps parses collection of Step from the current position in the parser
func (p *Parser) ParseBlockSteps() []object.Step {
	steps := []object.Step{}
	p.skipNewLines()
	if !isStepToken(p.curToken) {
		msg := fmt.Sprintf("Expected token to be a STEP_TYPE but got %s", p.curToken.Type)
		p.errors = append(p.errors, &GeneralParserError{parser: *p, LineNumber: p.curToken.LineNumber, Message: msg})
		return nil
	}
	for isStepToken(p.curToken) {
		steps = append(steps, *p.ParseStep())
		p.skipNewLines()
	}
	return steps
}

// ParseTags parses collection of Tag from the current position in the parser
func (p *Parser) ParseTags() []string {
	tags := []string{}
	for p.curTokenIs(token.TAG) {
		tags = append(tags, p.curToken.Literal)
		p.nextToken()
	}
	return tags
}

// ParseScenarioTypeSet parses collection of ScenarioType from the current position in the parser
func (p *Parser) ParseScenarioTypeSet() []object.ScenarioType {
	p.skipNewLines()
	if !(p.curTokenIs(token.SCENARIO) || p.curTokenIs(token.TAG)) {
		p.peekError(p.curToken.Type)
		return nil
	}
	scenarios := []object.ScenarioType{}
	for p.curTokenIs(token.SCENARIO) || p.curTokenIs(token.TAG) {
		scenarios = append(scenarios, p.ParseScenarioType())
		p.skipNewLines()
	}
	p.skipNewLines()
	return scenarios
}

// ParseScenarioType parses a ScenarioType from the current position in the parser
func (p *Parser) ParseScenarioType() object.ScenarioType {
	tags := []string{}
	p.skipNewLines()
	if p.curTokenIs(token.TAG) {
		tags = p.ParseTags()
		if tags == nil {
			return nil
		}
	}
	p.skipNewLines()
	if !p.curTokenIs(token.SCENARIO) {
		p.peekError(token.SCENARIO)
		return nil
	}
	outLineType := false
	if p.peekTokenIs(token.OUTLINE) {
		outLineType = true
		p.nextToken()
	}
	if !p.expectPeek(token.COLON) {
		p.peekError(token.COLON)
		return nil
	}
	p.nextToken()
	var title string
	for !p.curTokenIs(token.NEWLINE) {
		title += p.curToken.Literal
		p.nextToken()
	}
	p.skipNewLines()
	steps := p.ParseBlockSteps()
	if outLineType {
		p.skipNewLines()
		if !p.curTokenIs(token.EXAMPLES) {
			p.peekError(token.EXAMPLES)
			return nil
		}
		if !p.expectPeek(token.COLON) {
			p.peekError(token.COLON)
			return nil
		}
		if !p.expectPeek(token.NEWLINE) {
			p.peekError(token.NEWLINE)
			return nil
		}

		p.skipNewLines()
		if !p.curTokenIs(token.TABLEDATA) {
			p.peekError(p.curToken.Type)
			return nil
		}
		table := p.ParseTable()
		return &object.ScenarioOutline{
			Steps:        steps,
			Tags:         tags,
			ScenarioText: title,
			Table:        *table,
		}
	}
	return &object.Scenario{
		Steps:        steps,
		Tags:         tags,
		ScenarioText: title,
	}
}

func isDataType(check token.Token) bool {
	return check.Type == token.NUMBER || check.Type == token.STRING || check.Type == token.EXAMPLEVALUE
}

// ParseStep parses a Step from the current position in the parser
func (p *Parser) ParseStep() *object.Step {
	step := &object.Step{}
	p.skipNewLines()
	if token.IsStepToken(p.curToken.Type) {
		step.Token = p.curToken
		p.nextToken()
		for !(p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF)) {
			switch p.curToken.Type {
			case token.NUMBER:
				step.Data = append(step.Data, p.curToken.Literal)
				step.StepText = step.StepText + " {{d}}"
				if !isDataType(p.peekToken) {
					step.StepText += " "
				}
				p.nextToken()
			case token.STRING:
				step.Data = append(step.Data, p.curToken.Literal)
				step.StepText = step.StepText + " {{s}}"
				if !isDataType(p.peekToken) {
					step.StepText += " "
				}
				p.nextToken()
			case token.EXAMPLEVALUE:
				step.StepText = step.StepText + fmt.Sprintf(" {{<%v>}}", p.curToken.Literal)
				if !isDataType(p.peekToken) {
					step.StepText += " "
				}
				p.nextToken()
			default:
				step.StepText = step.StepText + p.curToken.Literal
				p.nextToken()
			}
		}
		step.StepText = strings.TrimSpace(step.StepText)
		p.nextToken()
	} else {
		msg := fmt.Sprintf("Expected token to be a STEP_TYPE but got %s", p.curToken.Type)
		p.errors = append(p.errors, &GeneralParserError{parser: *p, LineNumber: p.curToken.LineNumber, Message: msg})
		return nil
	}
	if p.curTokenIs(token.TABLEDATA) {
		table := p.ParseTable()
		if table == nil {
			return nil
		}
		step.Table = *table
	}
	if p.curTokenIs(token.PYSTRING) {
		step.StepText = step.StepText + "\n{{s}}"
		step.Data = append(step.Data, p.curToken.Literal)
		p.nextToken()
	}
	return step
}

// ParseTable parses a Table from the current position in the parser
func (p *Parser) ParseTable() *object.Table {
	var table object.Table
	var tmp []object.TableData
	p.skipNewLines()
	if !p.curTokenIs(token.TABLEDATA) {
		p.peekError(token.TABLEDATA)
		return nil
	}
	for p.curTokenIs(token.TABLEDATA) {
		tmp = []object.TableData{}
		for !(p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF)) {
			if !p.curTokenIs(token.TABLEDATA) {
				p.peekError(token.TABLEDATA)
				return nil
			}
			tmp = append(tmp, object.TableData{p.curToken.Literal, p.curToken.LineNumber})
			p.nextToken()
		}

		table = append(table, tmp)
		p.nextToken()

	}
	return &table
}
