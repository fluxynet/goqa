package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fluxynet/goqa"
	brokers "github.com/fluxynet/goqa/broker/memory"
	caches "github.com/fluxynet/goqa/cache/memory"
	"github.com/fluxynet/goqa/emailer/smtp"
	"github.com/fluxynet/goqa/repo/flat"
	rosters "github.com/fluxynet/goqa/roster/memory"
	"github.com/fluxynet/goqa/subscriber/cachew"
	"github.com/fluxynet/goqa/subscriber/coverage"
	"github.com/fluxynet/goqa/subscriber/email"
	repos "github.com/fluxynet/goqa/subscriber/repo"
	"github.com/fluxynet/goqa/web/hook"
	"github.com/fluxynet/goqa/web/server"
)

func main() {
	var cfg, err = LoadConf()
	if err != nil {
		log.Fatalln("failed to load config: ", err.Error())
	}

	var (
		repo   = flat.New()
		broker = brokers.New()
		cache  = caches.New()
		roster = rosters.New()
		mailer = smtp.New(cfg.EmailHost, cfg.EmailPort, cfg.EmailUsr, cfg.EmailPass, cfg.EmailFrom)

		hookServer = hook.Hook{
			Broker: broker,
			SigKey: cfg.GithubSigKey,
		}

		webServer = server.Server{
			Broker:    broker,
			Cache:     cache,
			Roster:    roster,
			IndexHTML: goqa.AssetIndexHtml,
			Prefix:    "/api/",
		}
	)

	var app = App{
		cfg:        cfg,
		roster:     roster,
		cache:      cache,
		mailer:     mailer,
		broker:     broker,
		repo:       repo,
		hookServer: &hookServer,
		webServer:  &webServer,
	}

	app.Serve(context.Background())
}

type App struct {
	cfg        *Config
	roster     goqa.Roster
	cache      goqa.Cache
	mailer     goqa.Emailer
	broker     goqa.Broker
	hookServer *hook.Hook
	repo       goqa.Repo
	webServer  *server.Server
}

func (a *App) Serve(ctx context.Context) {
	var err error

	var covs []goqa.Coverage
	covs, err = a.repo.Load(ctx)
	if err != nil {
		log.Fatalln("failed to load coverage from file", err.Error())
	}

	err = a.roster.Subscribe(ctx, goqa.EventGithub, repos.New(a.repo))
	if err != nil {
		log.Fatalln("failed to subscribe repo to "+goqa.EventGithub, err.Error())
	}

	err = a.cache.Reset(covs...)
	if err != nil {
		log.Fatalln("failed to initialize cache", err.Error())
	}

	err = a.roster.Subscribe(ctx, goqa.EventGithub, cachew.New(a.cache))
	if err != nil {
		log.Fatalln("failed to subscribe cachew to "+goqa.EventGithub, err.Error())
	}

	err = a.roster.Subscribe(ctx, goqa.EventGithub, coverage.New(a.broker))
	if err != nil {
		log.Fatalln("failed to subscribe coverage to "+goqa.EventGithub, err.Error())
	}

	for i := range a.cfg.EmailSubscribers {
		var sub = email.New(a.mailer, a.cfg.EmailSubscribers[i])
		err = a.roster.Subscribe(ctx, goqa.EventCoverage, sub)
		if err != nil {
			log.Fatalln("failed to subscribe to "+goqa.EventCoverage, err.Error())
		}
	}

	go goqa.Attach(a.broker, a.roster)

	http.HandleFunc("/github", a.hookServer.Receive)
	http.HandleFunc("/api/sse", a.webServer.SSE)
	http.HandleFunc("/api/", a.webServer.Get) // slash is the difference; not best practice
	http.HandleFunc("/api", a.webServer.List) // makes life easier :(
	http.HandleFunc("/", a.webServer.Index)

	fmt.Println("Server starting on: http://" + a.cfg.ServerHost)
	http.ListenAndServe(a.cfg.ServerHost, nil)
}
