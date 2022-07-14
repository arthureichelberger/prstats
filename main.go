package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arthureichelberger/prstats/pkg/env"
	"github.com/rs/zerolog/log"
)

func main() {
	repository := env.Get("GITHUB_REPOSITORY", "")
	secret := env.Get("GH_TOKEN", "")
	prID := env.Get("PULL_REQUEST_ID", "")

	payload := map[string]any{"body": "Hello world!", "event": "COMMENT"}
	payloadJSON, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/reviews", repository, prID)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadJSON))
	req.Header.Add("Authorization", fmt.Sprintf("token %s", secret))
	call(req)
}

func call(req *http.Request) {
	log.Debug().Str("url", req.URL.RequestURI()).Msg("calling api")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", req.URL.RequestURI()).Msg("could not execute request")
		return
	}

	defer res.Body.Close()

	var body map[string]any
	json.NewDecoder(res.Body).Decode(&body)
	bodyJSON, _ := json.Marshal(body)
	log.Info().Str("url", req.URL.RequestURI()).Str("status", res.Status).RawJSON("json", bodyJSON).Msg("request done")
}
