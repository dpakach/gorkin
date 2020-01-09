package filter

import (
	"strconv"
	"strings"

	"github.com/dpakach/gorkin/object"
)

// Filter interface is used to filter scenarios and features
type Filter interface {
	MatchFeature(feature object.Feature)
	MatchScenario(feature object.Scenario)
}

// TagFilter filters scenarios and features based on tags
type TagFilter struct {
	filterString string
}

func (tf *TagFilter) getTagsNot() []string {
	var tags []string
	for _, tag := range strings.Split(tf.filterString, "&&") {
		if tag[0] == '~' {
			tags = append(tags, tag[2:])
		}
	}
	return tags
}

func (tf *TagFilter) getTags() []string {
	var tags []string
	for _, tag := range strings.Split(tf.filterString, "&&") {
		if tag[0] != '~' {
			tags = append(tags, tag[1:])
		}
	}
	return tags
}

// MatchFeature matches given feature against the tag filter
func (tf *TagFilter) MatchFeature(feature object.Feature) bool {
	return tf.tagsMatchPresent(feature.Tags)
}

// MatchScenario matches given Scenario against the tag filter
func (tf *TagFilter) MatchScenario(scenario object.ScenarioType) bool {
	return tf.tagsMatchPresent(scenario.GetTags())
}

func (tf *TagFilter) tagsMatchPresent(tags []string) bool {
	if areArrayEqual(tags, []string{}) && !areArrayEqual(tf.getTags(), []string{}) {
		return false
	}
	for _, tag := range tags {
		if !tf.tagMatchPresent(tag) {
			return false
		}
		if tf.tagMatchNotPresent(tag) {
			return false
		}
	}
	return true
}

func (tf *TagFilter) tagMatchNotPresent(find string) bool {
	for _, tag := range tf.getTagsNot() {
		if tag == find {
			return true
		}
	}
	return false
}

func (tf *TagFilter) tagMatchPresent(find string) bool {
	for _, tag := range tf.getTags() {
		if tag == find {
			return true
		}
	}
	return false
}

// LineFilter creates a filter based on line numbers
type LineFilter struct {
	LineString string
}

func (lf *LineFilter) getStart() int {
	var startString string
	if strings.Contains(lf.LineString, "-") {
		startString = strings.Split(lf.LineString, "-")[0]
		if len(strings.Split(lf.LineString, "-")) > 2 {
			panic("Invalid filter string provided")
		}
	} else {
		startString = lf.LineString
	}
	start, err := strconv.Atoi(startString)
	if err != nil {
		panic(err)
	}
	return start
}

func (lf *LineFilter) getEnd() int {
	var endString string
	if strings.Contains(lf.LineString, "-") {
		arr := strings.Split(lf.LineString, "-")
		endString = arr[len(arr)-1]
		if len(arr) > 2 {
			panic("Invalid filter string provided")
		}
	} else {
		endString = lf.LineString
	}
	end, err := strconv.Atoi(endString)
	if err != nil {
		panic(err)
	}
	return end
}

// MatchFeature matches given feature against the line filter
func (lf *LineFilter) MatchFeature(feature *object.Feature) bool {
	ln := feature.Token.LineNumber
	if ln >= lf.getStart() && ln <= lf.getEnd() {
		return true
	}
	return false
}

// MatchScenario matches given Scenario against the Line filter
func (lf *LineFilter) MatchScenario(scenario *object.Scenario) bool {
	ln := scenario.LineNumber
	if ln >= lf.getStart() && ln <= lf.getEnd() {
		return true
	}
	return false
}
