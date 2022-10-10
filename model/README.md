# Model

## Tasks

### snapshot

Snapshot all of the tests.

```sh
ls -d tests/* | xargs -I '{}' go run snapshot/main.go -pkg="github.com/a-h/gidl/model/{}" -op="./{}/snapshot.json"
```

