package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v48/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file does not exist, will get the variables from the environment")
	}
}

func main() {
	ctx := context.Background()

	client, err := getGithubClient(ctx)
	if err != nil {
		panic(err)
	}

	repos, err := getGithubRepos(ctx, client)
	if err != nil {
		panic(err)
	}

	options, err := setProtectionOptions()
	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		fmt.Println("#################")
		fmt.Println(*repo.Name)
		err = setGithubProtection(ctx, client, *repo.Name, "main", options)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("...main branch now protected")
	}
}

func setProtectionOptions() (*github.ProtectionRequest, error) {
	var enabled bool = true

	var adminUser string
	adminUserValue, adminUserPresent := os.LookupEnv("USER")
	if adminUserPresent {
		adminUser = adminUserValue
	} else {
		return nil, errors.New("missing ENV Variable USER")
	}

	opt := &github.ProtectionRequest{
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews:          enabled,
			RequireCodeOwnerReviews:      enabled,
			RequiredApprovingReviewCount: 1,
			BypassPullRequestAllowancesRequest: &github.BypassPullRequestAllowancesRequest{
				Users: []string{adminUser},
				Teams: []string{},
				Apps:  []string{},
			},
		},
		RequiredConversationResolution: &enabled,
		EnforceAdmins:                  false,
	}

	return opt, nil
}

func getGithubClient(c context.Context) (*github.Client, error) {
	var githubToken string
	githubTokenValue, githubTokenPresent := os.LookupEnv("TOKEN")
	if githubTokenPresent {
		githubToken = githubTokenValue
	} else {
		return nil, errors.New("missing ENV Variable TOKEN")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(c, ts)

	client := github.NewClient(tc)

	return client, nil
}

func getGithubRepos(c context.Context, client *github.Client) ([]*github.Repository, error) {
	var organization string
	organizationValue, organizationPresent := os.LookupEnv("ORG")
	if organizationPresent {
		organization = organizationValue
	} else {
		return nil, errors.New("missing ENV Variable ORG")
	}

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(c, organization, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

func setGithubProtection(c context.Context, client *github.Client, repo string, branch string, opt *github.ProtectionRequest) error {
	var organization string
	organizationValue, organizationPresent := os.LookupEnv("ORG")
	if organizationPresent {
		organization = organizationValue
	} else {
		return errors.New("missing ENV Variable ORG")
	}

	_, _, err := client.Repositories.UpdateBranchProtection(c, organization, repo, branch, opt)
	if err != nil {
		return err
	}

	return nil
}
