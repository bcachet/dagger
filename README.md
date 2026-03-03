## dagger

Personal collection of Dagger modules

## Development

### Create a new Dagger module

```sh
mkdir my-module
cd my-module
dagger init --sdk go --source . --name my-module
```

### Regenerate dagger generated code

When you update a Dagger module, you may want to generate again the Go code that Dagger generated

```sh
dagger generate
```
