
# gist

## Installation

```
go get kkn.fi/cmd/gist
```

## Usage

```
% gist -h
usage: gist [-d string] [-p] [-a] file ... | -f file
  -a	create anonymous Gist
  -d string
    	description for Gist
  -f string
    	file name of the Gist. Reads file contents from stdin
  -p	create a public Gist
  -token file
    	read GitHub personal access token from file (default $HOME/.github-gist-token)
```

## Dependencies

 - [Google GitHub](https://godoc.org/github.com/google/go-github/github)
 - [oauth2](https://godoc.org/golang.org/x/oauth2)

## License

BSD-3-Clause. See [LICENSE](LICENSE) for details.

