package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arthureichelberger/prstats/pkg/github"
)

type StatService struct {
	ghClient github.HTTPClient
}

func NewStatService(ghClient github.HTTPClient) StatService {
	return StatService{
		ghClient: ghClient,
	}
}

func (ss StatService) Handle(ctx context.Context, prID string) error {
	pr, err := ss.ghClient.GetPullRequest(context.Background(), prID)
	if err != nil {
		return err
	}

	if pr.MergedAt == nil {
		return nil
	}

	mergedAt := *pr.MergedAt

	diff := func(t0, t1 time.Time) string {
		diff := t1.Sub(t0)
		diff = diff.Round(time.Minute)
		h := diff / time.Hour
		diff -= h * time.Hour
		m := diff / time.Minute
		return fmt.Sprintf("%02dh%02dmn", h, m)
	}(pr.CreatedAt, mergedAt)

	comment := fmt.Sprintf("# PRStats Bot\n\n**Time to merge:** %s.", diff)
	if err := ss.ghClient.CommentPullRequest(context.Background(), prID, comment); err != nil {
		return err
	}

	return nil
}
