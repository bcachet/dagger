# As module

To create a container image out of a Dockerfile located in current directory
```sh
dagger --mod github.com/bcachet/dagger/docker call build
```

To create a container image out of a Dockerfile located in example directory and build arg foo=bar
```sh
dagger --mod github.com/bcachet/dagger/docker call --source . with-build-arg --name foo --value bar build --file example/Dockerfile
```

