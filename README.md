# daggerverse

Personal collection of [Dagger](https://dagger.io) modules for CI/CD pipelines.

## Modules

* [`docker`](./docker)
* [`golang`](./golang)
* [`protobuf`](./protobuf)

All those modules can be used as toolchains. 
```sh
export MODULE=golang
dagger init
dagger toolchain install github.com/bcachet/dagger/$MODULE
dagger call $MODULE build
```

## Development

### Prerequisites

Install tools via [mise](https://mise.jdx.dev):

```sh
mise install
```

### Create a new module (event test module)

```sh
export MODULE=my-module
mkdir $MODULE
cd $MODULE
dagger init --sdk go --source . --name my-module
```

### Regenerate Dagger scaffolding

```sh
cd $MODULE && dagger develop --sdk go
```

### Validate a module

Our modules have tests (as defined in [Dagger module documentation](https://docs.dagger.io/reference/best-practices/modules/#module-tests))
```sh
dagger --mod ./$MODULE/tests call all
```
