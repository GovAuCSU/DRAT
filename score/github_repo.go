package score

import (
	"context"
	"fmt"

	github "github.com/google/go-github/github"
)

type GithubRepositoryScore struct {
	Weight     GithubScoringWeight
	R          *github.Repository
	c          *github.Client
	fullname   string
	Note       []string
	TotalScore int
}
type GithubScoringWeight struct {
	OwnerCountW       int
	ContributorCountW int
	LastUpdateW       int
	StarCountW        int
	WatcherCountW     int
}

var defaultGithubScoringWeight = GithubScoringWeight{
	OwnerCountW:       10,
	ContributorCountW: 10,
	LastUpdateW:       10,
	StarCountW:        10,
	WatcherCountW:     10,
}

// TODO: Write an init function to check for scoring weight configuration in a json or yml file and overwrite the default values.
// func init() {
// 	Check if score.config exist
//  Load scoring weight to it
// }

func NewGithubRepositoryScore(c *github.Client, r *github.Repository) *GithubRepositoryScore {
	rs := &GithubRepositoryScore{}
	rs.Weight = defaultGithubScoringWeight
	rs.R = r
	rs.c = c
	rs.fullname = "github.com/" + *r.Owner.Login + "/" + *r.Name
	rs.TotalScore = 100
	rs.Note = []string{}
	return rs
}

func (rs *GithubRepositoryScore) Score(ctx context.Context) error {
	//var err error

	// fmt.Printf("Owner Type: %v\n", *rs.R.Owner.Type) // This can either be Organisation or User  - Org is more trust worthy

	// branches, _, err := rs.c.Repositories.ListBranches(ctx, *rs.R.Owner.Login, *rs.R.Name, &github.ListOptions{})
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Number of branches: %v\n", len(branches))
	// rs.c.Repositories.List

	// issues, resp, err := rs.c.Issues.ListByRepo(ctx, *rs.R.Owner.Login, *rs.R.Name, &github.IssueListByRepoOptions{})
	// if err != nil {
	// 	return err
	// }
	// if resp.StatusCode != 200 {
	// 	return errors.New(string(resp.Body.Read(resp.ContentLength)))
	// }

	contributors, _, err := rs.c.Repositories.ListContributors(ctx, *rs.R.Owner.Login, *rs.R.Name, &github.ListContributorsOptions{ListOptions: github.ListOptions{PerPage: 100}})
	if err != nil {
		return err
	}
	fmt.Printf("Size of collaborator for %s: %d\n", rs.fullname, len(contributors))

	// commits, _, err := rs.c.Repositories.ListCommits(ctx, *rs.R.Owner.Login, *rs.R.Name, &github.CommitsListOptions{})
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Last commit at %v\n", commits[0].Commit.Author.Date)

	//Contributor stats
	// contributorStat, res, err := rs.c.Repositories.ListContributorsStats(ctx, *rs.R.Owner.Login, *rs.R.Name)
	// if err != nil {
	// 	return err
	// }
	// if res.StatusCode == 202 {
	// 	// Contributor stat require 2 operation at first run as it is an expensive operation, we can wait here until the data come back or move on and check back later
	// }

	// Latest release - does not work in many cases. we need to check Tag instead.
	// release, _, err := rs.c.Repositories.GetLatestRelease(ctx, *rs.R.Owner.Login, *rs.R.Name)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Test%v", release)

	return nil
}
