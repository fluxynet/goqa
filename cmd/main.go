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
	"github.com/fluxynet/goqa/subscriber/email"
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
		mailer = smtp.New(cfg.EmailHost, cfg.EmailUsr, cfg.EmailPass, cfg.EmailFrom)

		hookServer = hook.Hook{
			Broker: broker,
			SigKey: cfg.GithubSigKey,
		}

		webServer = server.Server{
			Broker:    broker,
			Cache:     cache,
			Roster:    roster,
			IndexHTML: goqa.AssetIndexHtml,
		}
	)

	var ctx = context.Background()

	var covs []goqa.Coverage
	covs, err = repo.Load(ctx)
	if err != nil {
		log.Fatalln("failed to load coverage from file", err.Error())
	}

	err = cache.Reset(covs...)
	if err != nil {
		log.Fatalln("failed to initialize cache", err.Error())
	}

	err = roster.Subscribe(ctx, goqa.EventGithub, cachew.New(cache))
	if err != nil {
		log.Fatalln("failed to subscribe to "+goqa.EventGithub, err.Error())
	}

	for i := range cfg.EmailSubscribers {
		err = roster.Subscribe(ctx, goqa.EventCoverage, email.New(mailer, cfg.EmailSubscribers[i]))
		if err != nil {
			log.Fatalln("failed to subscribe to "+goqa.EventCoverage, err.Error())
		}
	}

	goqa.Attach(ctx, broker, roster)

	http.HandleFunc("/github", hookServer.Receive)
	http.HandleFunc("/api/sse", webServer.SSE)
	http.HandleFunc("/api/", webServer.Get) // slash is the difference; not best practice
	http.HandleFunc("/api", webServer.List) // makes life easier :(
	http.HandleFunc("/", webServer.Index)

	fmt.Println("Server starting on: http://" + cfg.ServerHost)
	http.ListenAndServe(cfg.ServerHost, nil)
}
