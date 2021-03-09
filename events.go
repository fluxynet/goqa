package goqa

import (
	"fmt"
	"strings"
)

// GithubEvent on the system
type GithubEvent struct {
	Event      string `json:"event"`
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
	Ref        string `json:"ref"`
	Head       string `json:"head"`
	Workflow   string `json:"workflow"`
	Coverage   []Coverage
}

// Name of the event
func (e GithubEvent) Name() string {
	return EventCoverage
}

// String representation of the event
func (e GithubEvent) String() string {
	var b strings.Builder
	b.WriteString(
		`Event = "` + e.Event + `"\n` +
			`Repository = "` + e.Repository + `"\n` +
			`Commit = "` + e.Commit + `"\n` +
			`Ref = "` + e.Ref + `"\n` +
			`Head = "` + e.Head + `"\n` +
			`Workflow = "` + e.Workflow + `"\n` +
			`Coverage =\n`,
	)

	for i := range e.Coverage {
		b.WriteString(e.Coverage[i].String() + "\n")
	}

	return b.String()
}

// CoverageEvent
type CoverageEvent Coverage

func (c CoverageEvent) Name() string {
	return EventCoverage
}

func (c CoverageEvent) String() string {
	return fmt.Sprintf("pkg: %s; percentage: %s; time: %s", c.Pkg, c.Pkg, c.Time)
}
