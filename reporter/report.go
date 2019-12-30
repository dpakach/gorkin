package reporter

import (
	"io"
	"github.com/dpakach/gorkin/parser"
	"github.com/dpakach/gorkin/object"
)

func ParseAndReport(p *parser.Parser, out io.Writer) {
	res := p.Parse()
	if len(p.Errors()) != 0 {
		io.WriteString(out, "Parser Errors: \n")
		for _, err := range p.Errors() {
			io.WriteString(out, err.GetMessage() + "\n")
		}
	} else {
		PrintResult(out, res)
	}
}

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
			io.WriteString(out, " " + tag + " ")
		}
		io.WriteString(out, "]")
		io.WriteString(out, "\n\n\t")
		io.WriteString(out, "Background:\n\t")
		if feature.Background != nil {
			io.WriteString(out, "\t")
			PrintSteps(out, feature.Background.Steps)
		}
		io.WriteString(out, "\n\t")
		var titleString string
		var steps []object.Step
		var table object.Table
		var tags []string
		for _, scenario := range feature.Scenarios {
			io.WriteString(out, "\n\tScenario")
			outlineObj, ok := scenario.(*object.ScenarioOutline)
			isOutline := ok
			if isOutline {
				io.WriteString(out, " Outline:\n\t\t")
				tags = outlineObj.Tags
				titleString = outlineObj.ScenarioText
				steps = outlineObj.Steps
				table = outlineObj.Table

			} else {
				scenarioObj := scenario.(*object.Scenario)
				io.WriteString(out, ":\n\t\t")
				tags = scenarioObj.Tags
				titleString = scenarioObj.ScenarioText
				steps = scenarioObj.Steps
			}
			io.WriteString(out, "Title: ")
			io.WriteString(out, titleString)
			io.WriteString(out, "\n\t\t")
			io.WriteString(out, "Tags: ")
			io.WriteString(out, "[")
			for _, tag := range tags {
				io.WriteString(out, " " + tag + " ")
			}
			io.WriteString(out, "]")
			io.WriteString(out, "\n\t\t")
			PrintSteps(out, steps)
			PrintTable(out, table)
			io.WriteString(out, "\n\t")
		}
	}
}

func PrintSteps(out io.Writer, steps []object.Step) {
	for _, step := range steps {
		io.WriteString(out, "\n\t\t")
		io.WriteString(out, "Step: ")
		io.WriteString(out, step.Token.Literal)
		io.WriteString(out, " - ")
		io.WriteString(out, step.StepText)

		if step.Table != nil {
			PrintTable(out, step.Table)
		}
	}
}

func PrintTable(out io.Writer, table object.Table) {
	if len(table) > 0 {
		io.WriteString(out, "\n\t\tTable:")
		for _, row := range table {
			io.WriteString(out, "\n\t\t\t")
			for _, data := range row {
				io.WriteString(out, data)
				io.WriteString(out, "\t")
			}
		}
	}
}
