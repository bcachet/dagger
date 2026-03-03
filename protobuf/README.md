# As a module

## Generate code out of protobuf file

To generate the code associated to the Protobuf definitions located in `./example/` directory

```sh
dagger --mod github.com/bcachet/dagger/protobuf --source example call generate directory --path /out export --path ./gen
```

## Lint protobuf files

To lint the Protobuf definitions

```sh
dagger --mod github.com/bcachet/dagger/protobuf call lint stdout
```
