Generate, lint, and format Protobuf definitions using [buf](https://buf.build).

```sh
# Generate code from proto files
dagger --mod github.com/bcachet/dagger/protobuf --source example call generate directory --path /out export --path ./gen

# Lint proto files
dagger --mod github.com/bcachet/dagger/protobuf call lint stdout

# Format proto files
dagger --mod github.com/bcachet/dagger/protobuf call format export --path ./

# Check formatting in CI (fails if unformatted)
dagger --mod github.com/bcachet/dagger/protobuf call format --args --exit-code stdout
```