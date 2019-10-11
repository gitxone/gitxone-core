
# Downloads

Get the binary suiting you from the following link:

https://github.com/gitxone/gitxone-core/releases


# Usage
- Start up the server.
  ```bash
  $ ./gitxone-core_x.y.z_darwin_amd64/gitxone-core
  ```

- Go to http://localhost:10098/
- Input a path to your repository
  - And click 'Go'.

## Repository view
TODO: write


# Development

```bash
# Clone the repositories.
$ git clone git@github.com:gitxone/gitxone-core.git --recursive

# Update submodules.
$ git submodule update --recursive

# Start Nuxt server. and go to http://localhost:3000/
$ docker-compose up

# Release.
$ git tag vX.Y.Z
$ git push origin vX.Y.Z
```


