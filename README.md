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

- Use ShareX, cURL or any other client

```sh
curl -F 'file=@something.png' -u "user:pass" up.simo.sh
```

## How it works

1. The user makes a `POST` request to server with both their file and `Basic Auth` headers containing the username and password

2. The server checks if the username and password match the ones in the `AUTH` environment variable

3. The server does extra validation on the file

   - is it <= 512 MiB?
   - does it exist already? (calculates the SHA256 hash of the file and gets the k/v pair from [etcd-io/bbolt](github.com/etcd-io/bbolt))
     - if yes, check the date of upload
       - if it's been 95% of its calculated lifespan, reupload
       - otherwise, return the link

4. Upload the file in the `files/` directory, as well as add a randomly generated filename + SHA256 hash in [etcd-io/bbolt](github.com/etcd-io/bbolt), like so:

   ```sh
   sha256 hash of file -> filename
   # or
   d144c320286f5a2e20573aaa1193786e2b369d7ea26bea42583771a4b0236029 -> XVlBzgbaiC
   ```

5. Return the URL of the file
6. Every 24 hours, the server calculates the lifespan of every file uploaded using the formula and deletes any files that exceed that:

   ```py
   # shoutout to 0x0.st
   retention = min_age + (-max_age + min_age) * pow((file_size / max_size - 1), 3)
   ```

## TODO

- [x] [Basic Auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#Basic_authentication_scheme)
- [x] Automatic deletion after N amount of days (size dependent)
- [x] SHA256 hash to check for duplicates
- [x] Reupload if the file is going to be deleted soon
- [x] Dockerfile for "easy" daemonizing (super overkill, especially considering the portability of go, but it's the only thing I know /shrug - better option is simply to use systemd)
- [ ] Refactor with structs?
- [ ] Tests lol
- [ ] URL shortener
