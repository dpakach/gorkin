package filter

import (
	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/parser"
	"github.com/dpakach/gorkin/utils"
	"testing"
)

func TestLineFilterMatchFeature(t *testing.T) {

	input := `
	# This is feature

	Feature: test

	Scenario: test
		When test
	`

	l := lexer.New(input)
	p := parser.New(l)
	feature := p.ParseFeature()

	testData := []struct {
		input         string
		expectedMatch bool
	}{
		{"100", false},
		{"1-100", true},
		{"2-5", true},
		//{"addd", false},
		//{"1-200-300", false},
		{"4", true},
		{"5", false},
	}

	for _, tt := range testData {
		lf := &LineFilter{tt.input}
		match := lf.MatchFeature(feature)
		if match != tt.expectedMatch {
			t.Fatalf("Match feature for %q incorrect, expected: %v, got: %v", tt.input, tt.expectedMatch, match)
		}
	}
}

func TestGetTags(t *testing.T) {
	testdata := []struct {
		input           string
		expectedTags    []string
		expectedNotTags []string
	}{
		{"@tag1&&@tag2&&~@tag3", []string{"tag1", "tag2"}, []string{"tag3"}},
		{"@tag1&&@tag2&&@tag3", []string{"tag1", "tag2", "tag3"}, []string{}},
		{"~@tag1&&~@tag2&&~@tag3", []string{}, []string{"tag1", "tag2", "tag3"}},
		{"~@tag1&&@tag2&&~@tag3", []string{"tag2"}, []string{"tag1", "tag3"}},
	}

	for _, tt := range testdata {
		filter := &TagFilter{tt.input}
		tags := filter.getTags()
		tagsNot := filter.getTagsNot()

		if !utils.AreArrayEqual(tags, tt.expectedTags) {
			t.Fatalf("Tags are not equal, expected: %v, got: %v", tt.expectedTags, tags)
		}
		if !utils.AreArrayEqual(tagsNot, tt.expectedNotTags) {
			t.Fatalf("Tags not are not equal, expected: %v, got: %v", tt.expectedNotTags, tagsNot)
		}
	}
}

func TestTagMatchPresent(t *testing.T) {
	testdata := []struct {
		input     string
		inputTags []string
		match     bool
	}{
		{"@tag1&&@tag2&&~@tag3", []string{"tag1", "tag2"}, true},
		{"@tag1&&@tag2&&@tag3", []string{"tag1", "tag2", "tag3"}, true},
		{"~@tag1&&~@tag2&&~@tag3", []string{}, true},
		{"~@tag1&&@tag2&&~@tag3", []string{"tag2"}, true},
		{"@tag1&&@tag2&&~@tag3", []string{"tag3"}, false},
		{"@tag1&&@tag2&&@tag3", []string{}, false},
		{"~@tag1&&~@tag2&&~@tag3", []string{"tag1", "tag2", "tag3"}, false},
		{"~@tag1&&@tag2&&~@tag3", []string{"tag1", "tag3"}, false},

		{"@tag1", []string{"tag1"}, true},
		{"@tag1", []string{}, false},
		{"~@tag1", []string{"tag1"}, false},
		{"@tag1", []string{"tag2"}, false},
	}

	for _, tt := range testdata {
		filter := TagFilter{tt.input}

		if filter.tagsMatchPresent(tt.inputTags) != tt.match {
			t.Fatalf("Invalid match for filter %q, expected: %v, got: %v", tt.input, tt.match, filter.tagsMatchPresent(tt.inputTags))
		}
	}
}
