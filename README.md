## gha-usage

Find your GitHub Actions usage across a given organisation.

Includes total runtime of all workflow runs and workflow jobs, including where the jobs were run within inclusive, free, billed minutes, or on self-hosted runners.

## Usage

Create a Fine-grained personal access token with access to the given organisation, a regular personal access token doesn't have privileges to list repositories within an organisation.

```bash
go run . --org actuated-samples --token $(cat ~/pat.txt)
```

## Output

```bash
Total usage: 72h40m24s
```
