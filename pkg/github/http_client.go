package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	getPullRequestEndpoint     string = "pulls/%s"
	commentPullRequestEndpoint string = "pulls/%s/reviews"
)

type HTTPClient struct {
	baseURL    string
	secret     string
	httpClient http.Client
}

func NewHTTPClient(repository, secret string) HTTPClient {
	return HTTPClient{
		baseURL:    fmt.Sprintf("https://api.github.com/repos/%s", repository),
		secret:     secret,
		httpClient: http.Client{Timeout: time.Second},
	}
}

func (hc HTTPClient) GetPullRequest(ctx context.Context, prID string) (PullRequest, error) {
	uri := fmt.Sprintf(getPullRequestEndpoint, prID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", hc.baseURL, uri), http.NoBody)
	if err != nil {
		return PullRequest{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", hc.secret))

	res, err := hc.httpClient.Do(req)
	if err != nil {
		return PullRequest{}, err
	}

	log.Info().Str("status", res.Status).Str("uri", uri).Msg("request sent")

	defer res.Body.Close()

	var pr PullRequest
	if err := json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return PullRequest{}, err
	}

	return pr, nil
}

func (hc HTTPClient) CommentPullRequest(ctx context.Context, prID string, comment string) error {
	payload := map[string]any{"body": comment, "event": "COMMENT"}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf(commentPullRequestEndpoint, prID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", hc.baseURL, uri), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", hc.secret))

	res, err := hc.httpClient.Do(req)
	if err != nil {
		return err
	}

	log.Info().Str("status", res.Status).Str("uri", uri).Msg("request sent")

	defer res.Body.Close()

	return nil
}
