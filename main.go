package main

import (
	"context"

	"github.com/arthureichelberger/prstats/pkg/env"
	"github.com/arthureichelberger/prstats/pkg/github"
	"github.com/arthureichelberger/prstats/stats/service"
	"github.com/rs/zerolog/log"
)

func main() {
	repository := env.Get("GITHUB_REPOSITORY", "")
	secret := env.Get("GH_TOKEN", "")
	prID := env.Get("PULL_REQUEST_ID", "")

	switch {
	case repository == "":
		log.Panic().Msg("repository is undefined")
	case secret == "":
		log.Panic().Msg("secret is undefined")
	case prID == "":
		log.Panic().Msg("pull request id is undefined")
	}

	githubClient := github.NewHTTPClient(repository, secret)
	ss := service.NewStatService(githubClient)
	if err := ss.Handle(context.Background(), prID); err != nil {
		log.Panic().Err(err).Msg("could not handle stats")
	}
}
