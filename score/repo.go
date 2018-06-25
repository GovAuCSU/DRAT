package score

import (
	github "github.com/google/go-github/github"
)

type RepositoryScore struct {
	weight ScoringWeight
}
type ScoringWeight struct {
	OwnerCount       int
	ContributorCount int
	LastUpdate       int
	StarCount        int
	WatcherCount     int
}

var defaultScoringWeight = ScoringWeight{
	OwnerCount:       10,
	ContributorCount: 10,
	LastUpdate:       10,
	StarCount:        10,
	WatcherCount:     10,
}

func (rs *RepositoryScore) New(c *github.Client, r *github.Repository) {

}

// This is where we check for
// func init() {
// 	Check if score.config exist
//  Load scoring weight to it
// }
