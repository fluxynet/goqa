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
	var e, ok = event.(goqa.GithubEvent)
	if !ok {
		return subscriber.ErrUnsupportedEvent
	}

	return r.repo.Save(context.Background(), e.Coverage...)
}

func (r Repo) Serialize() (string, error) {
	return "", subscriber.ErrSerializeNotSupported
}

func (r Repo) Unserialize(string) error {
	return subscriber.ErrSerializeNotSupported
}
