package service

import (
	"regexp"
	"strconv"
)

type RegexText struct {
	Pattern string
	Content string
}

type RegexResult struct {
	*RegexText
	GroupNames    []string
	Captures      []Capture
	RangeCaptures []Capture
	Params        map[string]string
	MatchCount    int
	Positions     [][]int
}
type Capture struct {
	Start   int
	End     int
	Value   string
	IsMatch bool
	Params  map[string]string
	Groups  []Capture
}

func NewRegexText(pattern, content string) *RegexText {
	return &RegexText{pattern, content}
}

func (rs *RegexText) GetMatchResult() *RegexResult {
	r := regexp.MustCompile(rs.Pattern)
	positions := r.FindAllStringSubmatchIndex(rs.Content, -1)
	result := &RegexResult{
		RegexText:  rs,
		GroupNames: r.SubexpNames(),
		Positions:  positions,
		Params:     make(map[string]string),
	}
	if len(positions) == 0 {
		return result
	}
	result.SplitBy(rs.Content, true)
	result.FillParams(rs.Content, true)

	return result
}

// split positions
// if matchOnly=true, will get matches Capture[Start:End] only.
func (rs *RegexResult) SplitBy(input string, matchOnly bool) *[]Capture {
	// no alias, because no update, reference only.
	positions := rs.Positions
	//alias result.Captures for updating itself
	captures := &rs.Captures
	if len(positions) == 0 {
		*captures = append(*captures, Capture{Start: 0, End: len(input)})
	} else {
		cpos := 0
		epos := len(input)
		for _, pos := range positions {
			// match Capture
			match := Capture{Start: pos[0], End: pos[1], IsMatch: true}
			// append the ahead of match
			if !matchOnly && cpos < epos && cpos < match.Start {
				*captures = append(*captures, Capture{Start: cpos, End: match.Start})
			}
			// append match.value
			*captures = append(*captures, match)
			cpos = match.End
		}
		// append last string
		if !matchOnly && cpos < epos {
			*captures = append(*captures, Capture{Start: cpos, End: epos})
		}
	}
	//refresh result.Captures.value
	for i := 0; i < len(rs.Captures); i++ {
		c := &rs.Captures[i]
		c.Value = input[c.Start:c.End]
	}
	return captures
}

func (rs *RegexResult) FillParams(input string, detail bool) {

	// match.Index.
	x := 0
	for i, c := range rs.Captures {
		// skip if it's not match
		if !c.IsMatch {
			continue
		}

		match := &rs.Captures[i]
		match.Groups = make([]Capture, 0)
		match.Params = make(map[string]string)
		position := rs.Positions[x]
		if detail {
			match.Params["match_index"] = strconv.Itoa(x)
			match.Params["match_start"] = strconv.Itoa(match.Start)
			match.Params["match_end"] = strconv.Itoa(match.End)
		}

		for y := 0; y < len(rs.GroupNames); y++ {
			group := &Capture{Start: position[y*2+0], End: position[y*2+1], IsMatch: true}
			group.Value = input[group.Start:group.End]
			if group.Params == nil {
				group.Params = make(map[string]string)
			}
			gname := rs.GroupNames[y]
			if y == 0 {
				gname = "match_value"
			}
			if detail {
				group.Params["group_index"] = strconv.Itoa(y)
				group.Params["group_start"] = strconv.Itoa(group.Start)
				group.Params["group_end"] = strconv.Itoa(group.End)
				group.Params["group_key"] = gname
				group.Params["group_value"] = group.Value
			}
			match.Params[gname] = group.Value
			match.Groups = append(match.Groups, *group)
		}
		x++
	}
	rs.MatchCount = x
}
func (rs *RegexText) IsMatch() bool {
	r := regexp.MustCompile(rs.Pattern)
	return r.MatchString(rs.Content)
}
