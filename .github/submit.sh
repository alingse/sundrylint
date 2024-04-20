#!bin/bash
echo $1
gh search repos $1 --limit 100 --language=go --json url --jq '.[]|.url' | xargs -I {} gh workflow run .github/workflows/check-any.yaml -F repo_url={} -F linter_url=github.com/alingse/makezero@rename-go-mod -F linter_name=makezero

