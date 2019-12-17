package main

import (
	"os"
	"io"
	"fmt"
	"io/ioutil"
	"gorkin/lexer"
	"gorkin/parser"
	"gorkin/object"
)

func main() {
	if len(os.Args) > 1 {
		Run(os.Args[1])
	} else {
		fmt.Println("Opps, Seems like you forgot to provide the path of the feature file")
	}
}

func Run(path string) {
	out := os.Stdout
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}


	l := lexer.New(string(dat))
	p := parser.New(l)

	res := p.Parse()
	if len(p.Errors()) != 0 {
		io.WriteString(out, "Parser Errors: \n")
		for _, err := range p.Errors() {
			io.WriteString(out, err + "\n")
		}
	} else {
		PrintResult(out, res)
	}
}

func PrintResult(out io.Writer, featureSet *object.FeatureSet) {
	for _, feature := range featureSet.Features {
		io.WriteString(out, "Feature :\n")
		io.WriteString(out, "\tTitle :")
		io.WriteString(out, feature.Title)
		io.WriteString(out, "\n\t")
		io.WriteString(out, "Tags :")
		io.WriteString(out, "[")
		for _, tag := range feature.Tags {
			io.WriteString(out, " " + tag + " ")
		}
		io.WriteString(out, "]")
		io.WriteString(out, "\n\n\t")
		io.WriteString(out, "Background :\n\t")
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
			io.WriteString(out, "Scenario")
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
			io.WriteString(out, "Title :")
			io.WriteString(out, titleString)
			io.WriteString(out, "\n\t\t")
			io.WriteString(out, "Tags :")
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
		io.WriteString(out, "Step: ")
		io.WriteString(out, step.Token.Literal)
		io.WriteString(out, " - ")
		io.WriteString(out, step.StepText)

		io.WriteString(out, "\n\t\t")
		if step.Table != nil {
			PrintTable(out, step.Table)
		}
	}
}

func PrintTable(out io.Writer, table object.Table) {
	io.WriteString(out, "\n\t\tTable:\n\t\t\t")
	for _, row := range table {
		for _, data := range row {
			io.WriteString(out, data)
			io.WriteString(out, "\t")
		}
		io.WriteString(out, "\n\t\t\t")
	}
}
