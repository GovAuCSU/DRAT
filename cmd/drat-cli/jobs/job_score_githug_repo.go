package jobs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/GovAuCSU/DRAT/crawl"
	"github.com/GovAuCSU/DRAT/score"
	"github.com/santrancisco/cque"
)

const (
	KeyScoreGitHubRepo = "score_githug_repo"
)

// Example - we can wrap the other function
func ScoreGitHubRepo(logger *log.Logger, qc *cque.Client, j *cque.Job, appconfig map[string]interface{}) error {
	repojob := j.Args.(RepoJob)
	err := IsCrawlLimitReach(repojob, appconfig)
	if err != nil {
		return err
	}

	err = ScoreGitHubRepoFunc(logger, qc, j, appconfig)
	switch err {
	// Do something about error here, eg: in some case we could return ErrImmediateReschedule
	// We could also trap error code for rate limit reach?
	case nil:
		return nil
	default:
		logger.Print("[DEBUG] Error out in ScoreGithubRepo")
		return err

	}
}

func ScoreGitHubRepoFunc(logger *log.Logger, qc *cque.Client, j *cque.Job, appconfig map[string]interface{}) error {
	repojob := j.Args.(RepoJob)
	repofullurlparts := strings.Split(repojob.fullname, "/")
	reponame := repofullurlparts[len(repofullurlparts)-1]
	owner := repofullurlparts[len(repofullurlparts)-2]
	c, err := GetConn(appconfig)
	if err != nil {
		return err
	}
	repo, resp, err := c.Repositories.Get(context.Background(), owner, reponame)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to get information about repository %s/%s", owner, reponame)
	}
	rs := score.NewGithubRepositoryScore(c, repo)
	rs.Score(context.Background())
	// fmt.Printf("Repo score: \n %v \n", rs)
	lst, err := crawl.GithubDependencyCrawl(logger, c, repo, appconfig)
	if err != nil {
		return err
	}
	for _, v := range lst {
		if strings.Contains(v, "github.com") {
			logger.Printf("[INFO] Queue scoring job for \"%s\" found in repo %s/%s\n", v, *repo.Owner.Login, *repo.Name)
			qc.Enqueue(cque.Job{
				Type: KeyScoreGitHubRepo,
				Args: RepoJob{fullname: v, currentdepth: repojob.currentdepth + 1},
			})
			continue
		}
		// If it does not match with any type of Scoring job, discard it
		logger.Printf("[INFO] No scoring job for \"%s\" found in repo %s/%s\n", v, *repo.Owner.Login, *repo.Name)
	}

	return nil
}
