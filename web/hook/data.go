package hook

import (
	"strconv"
	"strings"

	"github.com/fluxynet/goqa"
)

// Payload received from the hook
type Payload struct {
	Event      string  `json:"event"`
	Repository string  `json:"repository"`
	Commit     string  `json:"commit"`
	Ref        string  `json:"ref"`
	Head       string  `json:"head"`
	Workflow   string  `json:"workflow"`
	Data       []Datum `json:"data"`
}

// Datum is the singular of data
type Datum struct {
	Time    string `json:"Time"`
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Output  string `json:"Output"`
	Elapsed string `json:"Elapsed"`
}

// CreateGithubEvent get coverage information from a Payload
func CreateGithubEvent(p *Payload) *goqa.GithubEvent {
	if p == nil {
		return nil
	}

	var covs []goqa.Coverage

	for i := range p.Data {
		if !strings.HasPrefix(p.Data[i].Output, "coverage: ") {
			continue
		}

		var c = goqa.Coverage{
			Pkg:  p.Data[i].Package,
			Time: p.Data[i].Time,
		}

		var perc string
		if x := strings.Index(p.Data[i].Output, "%"); x == -1 || x < 10 { // 10 is the length of "coverage: "
			continue
		} else {
			perc = p.Data[i].Output[10:x]
		}

		if v, err := strconv.ParseFloat(perc, 64); err == nil {
			c.Percentage = int(v) // really, we do not need that level of precision
		} else {
			continue
		}

		covs = append(covs, c)
	}

	var event = goqa.GithubEvent{
		Event:      p.Event,
		Repository: p.Repository,
		Commit:     p.Commit,
		Ref:        p.Ref,
		Head:       p.Head,
		Workflow:   p.Workflow,
		Coverage:   covs,
	}

	return &event
}
