## actions-usage

Find your GitHub Actions usage across a given organisation (or user account).

![Example console output](https://pbs.twimg.com/media/FrbYxbwWwAMvQZN?format=jpg&name=large)
> Example console output for the [inlets OSS repos](https://github.com/inlets)

Includes total runtime of all workflow runs and workflow jobs, including where the jobs were run within inclusive, free, billed minutes, or on self-hosted runners.

This data is not available within a single endpoint in GitHub's REST or GraphQL APIs, so many different API calls are necessary to build up the usage statistics.

```
repos = ListRepos(organisation || user)
   for each Repo
       ListWorkflowRuns(Repo)
          for each WorkflowRun
             jobs = ListWorkflowJobs(WorkflowRun)
sum(jobs)
```

If your team has hundreds of repositories, or thousands of builds per month, then the tool may exit early due to exceeding the API rate-limit. In this case, we suggest you run with `--days=10` and multiply the value by 3 to get a rough picture of 30-day usage.

## Usage

This tool is primarily designed for use with an organisation, however you can also use it with a regular user account by changing the `--org` flag to `--user`.

Or create a [Classic Token](https://github.com/settings/tokens) with: repo and admin:org and save it to ~/pat.txt. Create a short lived duration for good measure.

Download a binary from the [releases page](https://github.com/self-actuated/actions-usage/releases)

## Output

```bash
actions-usage --org openfaas --token $(cat ~/pat.txt)

Fetching last 30 days of data (created>=2023-01-29)

Total repos: 45
Total private repos: 0
Total public repos: 45

Total workflow runs: 95
Total workflow jobs: 113
Total usage: 6h16m16s (376 mins)
```

As a user:

```bash
actions-usage --user alexellis --token $(cat ~/pat.txt)
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

