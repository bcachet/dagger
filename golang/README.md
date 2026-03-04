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
