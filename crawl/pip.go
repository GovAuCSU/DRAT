package crawl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
)

const (
	PypiURLprefix = "https://pypi.org/project"
	//Libraries.IO is a great tool but rate limiting is cap at 60 requests per minute so unless
	// there is a way to up the limit, we cant use this for now.
	LibrariesIOprefix = "https://libraries.io/api/pypi"
)

func ParsePip(contentbytes []byte) ([]string, []Dependencyproblem) {
	var retlst []string
	var problems = []Dependencyproblem{}
	content := string(contentbytes)
	rexgit := regexp.MustCompile(`git\+(.*)\.git`)
	rexlib := regexp.MustCompile(`(?m:^)[\s+]?([a-zA-Z0-9]+[a-zA-Z0-9-_.]+( +)?)(([><=]=)|(?m:$))`)
	listofgit := rexgit.FindAllStringSubmatch(content, -1)
	listoflib := rexlib.FindAllStringSubmatch(content, -1)
	if len(listofgit) > 0 {
		for _, v := range listofgit {
			retlst = append(retlst, v[1])
		}
	}
	for _, r := range listoflib {
		libname := r[1]
		resp, err := http.Get(fmt.Sprintf("%s/%s", PypiURLprefix, libname))
		if err != nil {
			continue
		}
		if resp.StatusCode != 200 {
			problems = append(problems, Dependencyproblem{
				Name:      libname,
				URL:       fmt.Sprintf("%s/%s", PypiURLprefix, libname),
				RiskNotes: []string{"[MEDIUM] Could not retrieve information from pypi website for this library"},
			})
			continue
		}
		dp := Dependencyproblem{
			Name:      libname,
			URL:       fmt.Sprintf("%s/%s", PypiURLprefix, libname),
			RiskNotes: []string{},
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		doc := soup.HTMLParse(string(b))
		p := doc.Find("div", "id", "description").FindAll("p")
		if len(p) > 0 {
			if strings.TrimSpace(p[0].Text()) == "The author of this package has not provided a project description" {
				dp.RiskNotes = append(dp.RiskNotes, "[MEDIUM] This project does not have a description on Pypi")
			}
		}

		// Lets try to find a github,bitbucket and gitlab project URL
		rexgh := regexp.MustCompile(fmt.Sprintf("(((github.com)|(bitbucket.org)|(gitlab.com))/(repos/)?[a-zA-Z0-9-_]+/((%s)|(%s)))\"", libname, strings.ToLower(libname)))
		ghlinks := rexgh.FindAllStringSubmatch(string(b), -1)
		if len(ghlinks) == 0 {
			dp.RiskNotes = append(dp.RiskNotes, "[MEDIUM] Pypi page does not have any reference to a repository")
		} else {
			retlst = append(retlst, ghlinks[0][1])
		}
		// If we have a risknote related to Pypi, add it to the list of dependencies problems
		if len(dp.RiskNotes) > 0 {
			problems = append(problems, dp)
		}

	}
	return retlst, problems
}
