# up! ⚡

⚡ Simple file uploading service written in go!

## Usage

- Download binary from releases / build yourself using the go compiler.
- Run the binary, and define an `AUTH` envioronment variable, like so:

```sh
AUTH="user:pass" ./up
# or
AUTH="user:pass" go run .
```

- Use ShareX, curl or any other client

```sh
curl -F 'file=@something.png' -u "user:pass" up.simo.sh
```

## TODO:

- [x] [Basic Auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#Basic_authentication_scheme)
- [x] Automatic deletion after N amount of days (size dependent)
- [x] SHA256 hash to check for duplicates
- [x] Reupload if the file is going to be deleted soon
- [ ] Dockerfile for "easy" daemonizing (super overkill, especially considering the portability of go, but it's the only thing I know /shrug - better option is simply to use systemd)
- [ ] Refactor with structs?
- [ ] Tests lol
- [ ] url shortener
