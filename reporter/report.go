package reporter

import (
	"fmt"
	"github.com/dpakach/gorkin/object"
	"github.com/dpakach/gorkin/parser"
	"io"
	"strconv"
)

// ParseAndReport Parses the input Parser and writes the output in given writer
func ParseAndReport(p *parser.Parser, out io.Writer) {
	res := p.Parse()
	if len(p.Errors()) != 0 {
		io.WriteString(out, "Parser Errors: \n")
		for _, err := range p.Errors() {
			io.WriteString(out, err.GetMessage()+"\n")
		}
	} else {
		PrintResult(out, res)
	}
}

// PrintResult Parses the input FeatureSet and writes the output in given writer
func PrintResult(out io.Writer, featureSet *object.FeatureSet) {
	for _, feature := range featureSet.Features {
		io.WriteString(out, "\n")
		io.WriteString(out, "Feature:\n")
		io.WriteString(out, "\tTitle: ")
		io.WriteString(out, feature.Title)
		io.WriteString(out, "\n\t")
		io.WriteString(out, "Tags: ")
		io.WriteString(out, "[")
		for _, tag := range feature.Tags {
			io.WriteString(out, " "+tag+" ")
		}
		io.WriteString(out, "]")
		io.WriteString(out, "\n\n\t")
		io.WriteString(out, "Background:\n\t")
		if feature.Background != nil {
			io.WriteString(out, "\t")
			PrintSteps(out, feature.Background.Steps, 2)
		}
		io.WriteString(out, "\n\t")
		var titleString string
		var steps []object.Step
		var table object.Table
		var tags []string
		var lineNumber int
		for _, scenario := range feature.Scenarios {
			io.WriteString(out, "\n\tScenario")
			outlineObj, ok := scenario.(*object.ScenarioOutline)
			isOutline := ok
			if isOutline {
				io.WriteString(out, " Outline:\n\t\t")
				tags = outlineObj.Tags
				titleString = outlineObj.ScenarioText
				io.WriteString(out, "Title: ")
				io.WriteString(out, titleString)
				io.WriteString(out, "\n")
				for _, s := range outlineObj.GetScenarios() {
					steps = outlineObj.Steps
					table = outlineObj.Table
					lineNumber = s.LineNumber
					io.WriteString(out, "\t\tTitle: ")
					io.WriteString(out, titleString)
					io.WriteString(out, ":")
					io.WriteString(out, strconv.Itoa(lineNumber))
					io.WriteString(out, "\n\t\t\t\t")
					io.WriteString(out, "Tags: ")
					io.WriteString(out, "[")
					for _, tag := range tags {
						io.WriteString(out, " "+tag+" ")
					}
					io.WriteString(out, "]")
					io.WriteString(out, "\n\t\t\t\t")
					PrintSteps(out, steps, 4)
					PrintTable(out, table)
					io.WriteString(out, "\n\t")
				}
			} else {
				scenarioObj := scenario.(*object.Scenario)
				io.WriteString(out, ":\n\t\t")
				tags = scenarioObj.Tags
				titleString = scenarioObj.ScenarioText
				steps = scenarioObj.Steps
				lineNumber = scenarioObj.LineNumber

				io.WriteString(out, "Title: ")
				io.WriteString(out, titleString)
				io.WriteString(out, ":")
				io.WriteString(out, strconv.Itoa(lineNumber))
				io.WriteString(out, "\n\t\t")
				io.WriteString(out, "Tags: ")
				io.WriteString(out, "[")
				for _, tag := range tags {
					io.WriteString(out, " "+tag+" ")
				}
				io.WriteString(out, "]")
				io.WriteString(out, "\n\t\t")
				PrintSteps(out, steps, 2)
				PrintTable(out, table)
				io.WriteString(out, "\n\t")
			}
		}
	}
}

// PrintSteps Parses the collection of Steps and writes the output in given writer
func PrintSteps(out io.Writer, steps []object.Step, tab int) {
	for _, step := range steps {
		io.WriteString(out, "\n")
		i := 0
		for i < tab {
			io.WriteString(out, "\t")
			i++
		}

		io.WriteString(out, "Step: ")
		io.WriteString(out, step.Token.Literal)
		io.WriteString(out, " - ")
		io.WriteString(out, step.StepText)

		if step.Table != nil {
			PrintTable(out, step.Table)
		}
	}
}

// PrintTable Parses the Table and writes the output in given writer
func PrintTable(out io.Writer, table object.Table) {
	fmt.Println(table)
	if len(table) > 0 {
		io.WriteString(out, "\n\t\tTable:")
		for _, row := range table {
			io.WriteString(out, "\n\t\t\t")
			for _, data := range row {
				io.WriteString(out, data.Literal)
				io.WriteString(out, "\t")
			}
		}
	}
}
