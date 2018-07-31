package crawl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

// "https://raw.githubusercontent.com/%s/%s/master/%s",owner,repository,filename

func GithubDependencyCrawl(logger *log.Logger, c *github.Client, r *github.Repository, config interface{}) []string {
	supportfiles := make(map[string]ParserFunc)
	supportfiles["Gopkg.lock"] = ParseGoDep
	supportfiles["package.json"] = ParseNPM

	var listOfDependencies []string
	for filename, PF := range supportfiles {
		dependencies, err := GetFileAndParseResult(filename, PF, logger, c, r, config)
		if err != nil {
			continue
		}
		if len(dependencies) > 0 {
			listOfDependencies = append(listOfDependencies, dependencies...)
		}
	}
	return listOfDependencies
}

type ParserFunc func(filecontent []byte) []string

func GetFileAndParseResult(filename string, PF ParserFunc, logger *log.Logger, c *github.Client, r *github.Repository, config interface{}) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", *r.Owner.Login, *r.Name, filename))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		logger.Print(fmt.Sprintf("[DEBUG] Successfully downloaded https://raw.githubusercontent.com/%s/%s/master/%s", *r.Owner.Login, *r.Name, filename))
		contentbytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return PF(contentbytes), nil
	}
	return []string{}, nil
}
