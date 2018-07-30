package crawl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

// "https://raw.githubusercontent.com/%s/%s/master/%s",owner,repository,filename

func GithubDependencyCrawl(logger *log.Logger, c *github.Client, r *github.Repository, config interface{}) ([]string, error) {
	filename := "Gopkg.lock"
	var listOfDependencies []string
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
		content := string(contentbytes)
		listOfDependencies = append(listOfDependencies, ParseGoDep(content)...)

	}
	return listOfDependencies, nil
}
