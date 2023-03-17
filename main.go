package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	units "github.com/docker/go-units"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func main() {

	var (
		orgName, userName, token string
		since                    int
	)

	flag.StringVar(&orgName, "org", "", "Organization name")
	flag.StringVar(&userName, "user", "", "User name")
	flag.StringVar(&token, "token", "", "GitHub token")
	flag.IntVar(&since, "since", 30, "Since when to fetch the data (in days)")

	flag.Parse()

	auth := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	switch {
	case orgName == "" && userName == "":
		log.Fatal("Organization name or username is required")
	case orgName != "" && userName != "":
		log.Fatal("Only org or username must be specified at the same time")
	}

	client := github.NewClient(auth)
	created := time.Now().AddDate(0, 0, -since)
	format := "2006-01-02"
	createdQuery := ">=" + created.Format(format)

	var (
		totalRuns    int
		totalJobs    int
		totalPrivate int
		totalPublic  int
		longestBuild time.Duration
		actors       map[string]bool
	)

	actors = make(map[string]bool)

	fmt.Printf("Fetching last %d days of data (created>=%s)\n", since, created.Format("2006-01-02"))

	var repos, allRepos []*github.Repository
	var res *github.Response
	var err error
	ctx := context.Background()
	page := 0
	for {
		log.Printf("Fetching repos %s page %d", orgName, page)
		if orgName != "" {
			opts := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{Page: page, PerPage: 100}, Type: "all"}
			log.Printf("Fetching repos %s page %d", orgName, page)
			repos, res, err = client.Repositories.ListByOrg(ctx, orgName, opts)
		}

		if userName != "" {
			opts := &github.RepositoryListOptions{ListOptions: github.ListOptions{Page: page, PerPage: 100}, Type: "all"}
			repos, res, err = client.Repositories.List(ctx, userName, opts)
		}

		if err != nil {
			log.Fatal(err)
		}

		if res.Rate.Remaining == 0 {
			panic("Rate limit exceeded")
		}

		allRepos = append(allRepos, repos...)

		log.Printf("Status: %d Page %d, next page: %d", res.StatusCode, page, res.NextPage)

		if len(allRepos) == 0 {
			break
		}
		if res.NextPage == 0 {
			break
		}

		// break
		page = res.NextPage
	}

	allUsage := time.Second * 0

	for _, repo := range allRepos {
		log.Printf("Found: %s", repo.GetFullName())
		if repo.GetPrivate() {
			totalPrivate++
		} else {
			totalPublic++
		}

		workflowRuns := []*github.WorkflowRun{}
		page := 0
		for {

			opts := &github.ListWorkflowRunsOptions{Created: createdQuery, ListOptions: github.ListOptions{Page: page, PerPage: 100}}

			var runs *github.WorkflowRuns
			log.Printf("Listing workflow runs for: %s", repo.GetFullName())
			if orgName != "" {
				runs, res, err = client.Actions.ListRepositoryWorkflowRuns(ctx, orgName, repo.GetName(), opts)
			}
			if userName != "" {
				realOwner := userName
				// if user is a member of repository
				if userName != *repo.Owner.Login {
					realOwner = *repo.Owner.Login
				}
				runs, res, err = client.Actions.ListRepositoryWorkflowRuns(ctx, realOwner, repo.GetName(), opts)
			}
			if err != nil {
				log.Fatal(err)
			}

			workflowRuns = append(workflowRuns, runs.WorkflowRuns...)

			if len(workflowRuns) == 0 {
				break
			}

			if res.NextPage == 0 {
				break
			}

			page = res.NextPage
		}
		totalRuns += len(workflowRuns)

		var entity string
		if orgName != "" {
			entity = orgName
		}
		if userName != "" {
			entity = userName
		}
		log.Printf("Found %d workflow runs for %s/%s", len(workflowRuns), entity, repo.GetName())

		for _, run := range workflowRuns {
			log.Printf("Fetching jobs for: run ID: %d, startedAt: %s, conclusion: %s", run.GetID(), run.GetRunStartedAt().Format("2006-01-02 15:04:05"), run.GetConclusion())
			workflowJobs := []*github.WorkflowJob{}

			if a := run.GetActor(); a != nil {
				actors[a.GetLogin()] = true
			}
			page := 0
			for {
				log.Printf("Fetching jobs for: %d, page %d", run.GetID(), page)
				jobs, res, err := client.Actions.ListWorkflowJobs(ctx, entity,
					repo.GetName(),
					run.GetID(),
					&github.ListWorkflowJobsOptions{Filter: "all", ListOptions: github.ListOptions{Page: page, PerPage: 100}})
				if err != nil {
					log.Fatal(err)
				}

				for _, job := range jobs.Jobs {
					dur := job.GetCompletedAt().Time.Sub(job.GetStartedAt().Time)
					if dur > longestBuild {
						longestBuild = dur
					}
				}

				workflowJobs = append(workflowJobs, jobs.Jobs...)

				if len(jobs.Jobs) == 0 {
					break
				}

				if res.NextPage == 0 {
					break
				}
				page = res.NextPage
			}

			totalJobs += len(workflowJobs)
			log.Printf("%d jobs for workflow run: %d", len(workflowJobs), run.GetID())
			for _, job := range workflowJobs {

				if !job.GetCompletedAt().IsZero() {
					dur := job.GetCompletedAt().Time.Sub(job.GetStartedAt().Time)
					allUsage += dur
					log.Printf("Job: %d [%s - %s] (%s): %s",
						job.GetID(), job.GetStartedAt().Format("2006-01-02 15:04:05"),
						job.GetCompletedAt().Format("2006-01-02 15:04:05"),
						humanDuration(dur), job.GetConclusion())
				}
			}
		}
	}

	fmt.Println("\nUsage report generated by self-actuated/actions-usage.\n")
	fmt.Printf("Total repos: %d\n", len(allRepos))
	fmt.Printf("Total private repos: %d\n", totalPrivate)
	fmt.Printf("Total public repos: %d\n", totalPublic)
	fmt.Println()
	fmt.Printf("Total workflow runs: %d\n", totalRuns)
	fmt.Printf("Total workflow jobs: %d\n", totalJobs)
	fmt.Println()
	fmt.Printf("Total users: %d\n", len(actors))
	fmt.Printf("Longest build: %s\n", longestBuild.Round(time.Second))
	fmt.Printf("Average build time: %s\n", (allUsage / time.Duration(totalJobs)).Round(time.Second))

	mins := fmt.Sprintf("%.0f mins", allUsage.Minutes())
	fmt.Printf("Total usage: %s (%s)\n", allUsage.String(), mins)
}

// types.HumanDuration fixes a long string for a value < 1s
func humanDuration(duration time.Duration) string {
	v := strings.ToLower(units.HumanDuration(duration))

	if v == "less than a second" {
		return fmt.Sprintf("%d ms", duration.Milliseconds())
	} else if v == "about a minute" {
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	}

	return v
}
