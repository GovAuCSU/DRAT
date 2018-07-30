package score

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GovAuCSU/DRAT/crawl"
)

// func setupTest() (r *github.Repository, u *github.User, c *github.Client, err error) {
// 	err = nil
// 	ctx := context.Background()
// 	c = github.NewClient(&http.Client{})
// 	if os.Getenv("GITHUB_AUTH_TOKEN") != "" {
// 		ts := oauth2.StaticTokenSource(
// 			&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
// 		)

// 		tc := oauth2.NewClient(ctx, ts)
// 		c = github.NewClient(tc)
// 	}
// 	r, _, err = c.Repositories.Get(ctx, "google", "go-github")
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	u, _, err = c.Users.Get(ctx, "santrancisco")
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	return r, u, c, err
// }

func TestParseGoDep(t *testing.T) {
	t.Log("Test ParseGoDep")
	b, err := ioutil.ReadFile("Gopkg.lock")
	if err != nil {
		t.Errorf("Cannot open Gopkg.lock")
	}
	content := string(b)
	lst := crawl.ParseGoDep(content)
	final := strings.Join(lst, ",")
	expected := "github.com/beorn7/perks,github.com/bgentry/que-go,github.com/cloudfoundry-community/go-cfenv,github.com/gogo/protobuf,github.com/golang/protobuf,github.com/google/certificate-transparency-go,github.com/jackc/pgx,github.com/matttproud/golang_protobuf_extensions,github.com/mitchellh/mapstructure,github.com/pkg/errors,github.com/prometheus/client_golang,github.com/prometheus/client_model,github.com/prometheus/common,github.com/prometheus/procfs,golang.org/x/crypto,golang.org/x/net"
	if final != expected {
		t.Errorf("ParseGoDep fails, output is not as expected")
	}
}
