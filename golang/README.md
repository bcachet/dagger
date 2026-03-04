Build, lint, and vulnerability-scan Go projects with integrated caching for fast, repeatable CI workflows.

```sh
# Build the project
dagger --mod github.com/bcachet/dagger/golang --source ./example call build export --path ./out

# Lint with golangci-lint
dagger --mod github.com/bcachet/dagger/golang --source ./example call lint

# Vulnerability scan with govulncheck
dagger --mod github.com/bcachet/dagger/golang --source ./example call vuln-check
```

Can also be installed as a [Dagger toolchain](https://docs.dagger.io/configuration/modules#toolchain):

```sh
dagger init
dagger toolchain install github.com/bcachet/dagger/golang
dagger call golang --source . build
dagger checks  # runs lint + vuln-check
```