package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/GovAuCSU/drat/cmd/drat-cli/jobs"
	"github.com/santrancisco/cque"
	"github.com/santrancisco/logutils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose      = kingpin.Flag("verbose", "Enable debug mode.").Default("false").Short('v').Bool()
	repo         = kingpin.Flag("repo", "Name of repo in the following format: (github.com/(username|org)/repo").Default("").Short('r').String()
	file         = kingpin.Flag("file", "File contains the list of urls to repositories seperated by newlines").Default("").Short('f').String()
	gitauthtoken = kingpin.Flag("gitauth", "Github Authentication token").Default("").OverrideDefaultFromEnvar("GITHUB_AUTH_TOKEN").String()
	depth        = kingpin.Flag("depth", "How deep do we want to crawl the dependencies").Default("5").Short('d').Int()
	sqlitedbpath = kingpin.Flag("sqlpath", "Local sqlite cache location. Default is ~/.drat.sqlite").Default("~/.drat.sqlite").Short('s').String()
)

const (
	ORGTYPE         = "Organisation"
	CONTRIBUTORTYPE = "Contributor"
	REPOSITORYTYPE  = "Repository"
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	config := map[string]interface{}{}
	config["depth"] = *depth
	config["gitauthtoken"] = *gitauthtoken
	config["sqlitedbpath"] = *sqlitedbpath

	// Configuring our log level
	logfilter := "ERROR"
	if *verbose {
		logfilter = "DEBUG"
	}
	filteroutput := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARNING", "ERROR"},
		MinLevel: logutils.LogLevel(logfilter),
		Writer:   os.Stderr,
	}
	log.SetOutput(filteroutput)

	ctx, cancel := context.WithCancel(context.Background())
	// defering canclation of all concurence processes
	defer cancel()

	qc := cque.NewQue()
	wpool := cque.NewWorkerPool(qc, cque.WorkMap{
		jobs.KeyScoreGitHubRepo: (&jobs.JobFuncWrapper{
			QC:        qc,
			Logger:    log.New(filteroutput, "", log.LstdFlags),
			F:         jobs.ScoreGitHubRepo,
			AppConfig: config}).Run,

		jobs.KeyListFromFile: (&jobs.JobFuncWrapper{
			QC:        qc,
			Logger:    log.New(filteroutput, "", log.LstdFlags),
			F:         jobs.ListFromFile,
			AppConfig: config}).Run,
	}, 5)

	wpool.Start(ctx)
	// If we give it a file, queue a new job to crawl content of the file.
	if *file != "" {
		qc.Enqueue(cque.Job{Type: jobs.KeyListFromFile, Args: *file})
	}

	if *repo != "" {
		qc.Enqueue(cque.Job{Type: jobs.KeyListFromFile, Args: *file})
	}

	time.Sleep(2 * time.Second)
	for !qc.IsQueueEmpty {
		time.Sleep(2 * time.Second)
		// when queue is not empty we loop mainthread.
	}
}
