## actions-usage

Find your GitHub Actions usage across a given organisation.

Includes total runtime of all workflow runs and workflow jobs, including where the jobs were run within inclusive, free, billed minutes, or on self-hosted runners.

This data is not available within a single endpoint in GitHub's REST or GraphQL APIs, so many different API calls are necessary to build up the usage statistics.

```
repos = ListRepos(organisation)
   for each Repo
       ListWorkflowRuns(Repo)
          for each WorkflowRun
             jobs = ListWorkflowJobs(WorkflowRun)
sum(jobs)
```

If your team has hundreds of repositories, or thousands of builds per month, then the tool may exit early due to exceeding the API rate-limit. In this case, we suggest you run with `--days=10` and multiply the value by 3 to get a rough picture of 30-day usage.

## Usage

Or create a [Classic Token](https://github.com/settings/tokens) with: repo and admin:org and save it to ~/pat.txt. Create a short lived duration for good measure.

Download a binary from the [releases page](https://github.com/self-actuated/actions-usage/releases/tag/v0.0.1)

## Output

```bash
./actions-usage --org actuated-samples --token $(cat ~/pat.txt)

Fetching last 30 days of data
Total repos: 25
Total workflow runs: 245
Total workflow jobs: 255
Total usage: 9h10m17s
```

## Development

All changes must be proposed with an Issue prior to working on them or sending a PR. Commits must have a sign-off message, i.e. `git commit -s`

```bash
git clone https://github.com/actuated/actions-usage
cd actions-usage

go run . --org actuated-samples --token $(cat ~/pat.txt)
```

## Author

This tool was created as part of [actuated.dev](https://actuated.dev) by OpenFaaS Ltd.

