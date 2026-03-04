# Run module tests

```sh
# Will run lint/vuln-check/build tests against golang/example Go module
dagger -m golang/tests/ --source golang/example call all
```

# As a module

## Build
```sh
dagger --mod github.com/bcachet/dagger/golang --source ./golang/example call build export --path ./out
```

## Lint

```sh
dagger --mod github.com/bcachet/dagger/golang --source ./golang/example call lint
```

## Vuln-check

```sh
dagger --mod github.com/bcachet/dagger/golang --source ./golang/example call vuln-check
```

# Toolchain

You can install this module as a toolchain

```sh
dagger init
dagger toolchain install github.com/bcachet/dagger/golang
```

Once installed the functions will be available

```sh
dagger call golang --source . build
```

and checks too
```sh
dagger checks
# Will perform lint/vuln-check checks
```
