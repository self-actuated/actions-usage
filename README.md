## actions-usage

Find your GitHub Actions usage across a given organisation.

Includes total runtime of all workflow runs and workflow jobs, including where the jobs were run within inclusive, free, billed minutes, or on self-hosted runners.

## Usage

Or create a [Classic Token](https://github.com/settings/tokens) with: repo and admin:org and save it to ~/pat.txt. Create a short lived duration for good measure.

Download a binary from the [releases page](https://github.com/self-actuated/actions-usage/releases)

Or clone the source and run it:

```bash
git clone https://github.com/actuated/actions-usage --depth=1
cd actions-usage

go run . --org actuated-samples --token $(cat ~/pat.txt)
```

## Output

```bash
Fetching last 30 days of data
Total repos: 25
Total workflow runs: 245
Total workflow jobs: 255
Total usage: 9h10m17s
```

## Author

This tool was created as part of [actuated.dev](https://actuated.dev) by OpenFaaS Ltd.
