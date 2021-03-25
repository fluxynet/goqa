package repo

import (
	"context"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

func New(repo goqa.Repo) *Repo {
	return &Repo{repo: repo}
}

type Repo struct {
	subscriber.Identifiable
	repo goqa.Repo
}

func (r Repo) Notify(event goqa.Event) error {
	var e *goqa.GithubEvent

	switch v := event.(type) {
	default:
		return subscriber.ErrUnsupportedEvent
	case *goqa.GithubEvent:
		e = v
	case goqa.GithubEvent:
		e = &v
	}

	return r.repo.Save(context.Background(), e.Coverage...)
}
