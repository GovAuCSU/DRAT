package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debugflag = kingpin.Flag("debug", "Enable debug mode.").Default("false").Bool()
	repo      = kingpin.Flag("repo", "Name of repo in the following format: (username|org)/repo").Default("").Short('r').String()
	gitauth   = kingpin.Flag("gitauth", "Github Authentication token").Default("").OverrideDefaultFromEnvar("GITHUB_AUTH_TOKEN").String()
	depth     = kingpin.Flag("depth", "How deep do we want to crawl the dependencies").Default("5").Short('d').Int()
	sqlitedb  = kingpin.Flag("sqlpath", "Local sqlite cache location. Default is ~/.drat.sqlite").Default("~/.drat.sqlite").Short('s').String()
)

const (
	ORGTYPE         = "Organisation"
	CONTRIBUTORTYPE = "Contributor"
	REPOSITORYTYPE  = "Repository"
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// Configuring our log level
	logfilter := "ERROR"
	if *debugflag {
		logfilter = "DEBUG"
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR", "INFO"},
		MinLevel: logutils.LogLevel(logfilter),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	ctx, cancel := context.WithCancel(context.Background())
	// defering canclation of all concurence processes
	defer cancel()

}
