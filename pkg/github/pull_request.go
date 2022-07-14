package github

import "time"

type PullRequest struct {
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at"`
}
