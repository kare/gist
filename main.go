/*
Gist is a client for creating GitHub Gists.

	usage: gist [-v] [-p] [-a] [-d string] file ... | -f file

Gist uploads local file[s] to gist.github.com and prints information
about the created Gist. Default user is the authenticated user.

Authentication

Gist expects to find a GitHub "personal access token" in
$HOME/.github-gist-token and will that token to authenticate
to Github when writing Gist data.
A token can be created by visiting https://github.com/settings/tokens/new.
The token only needs the 'gist' scope checkbox.
It does not need any other permissions.
The -token flag specifies an alternate file from which to read the token.

*/
package main // import "kkn.fi/cmd/gist"

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

var (
	version    string
	date       string
	v          = flag.Bool("v", false, "print version and exit")
	descFlag   = flag.String("d", "", "description for Gist")
	publicFlag = flag.Bool("p", false, "create a public Gist")
	anonFlag   = flag.Bool("a", false, "create anonymous Gist")
	fileFlag   = flag.String("f", "", "file name of the Gist. Reads file contents from stdin")
	tokenFile  = flag.String("token", "", "read GitHub personal access token from `file` (default $HOME/.github-gist-token)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gist [-v] [-d string] [-p] [-a] file ... | -f file\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("gist: ")

	if *v {
		fmt.Printf("%v: %v %v %v\n", path.Base(os.Args[0]), version, date, runtime.Version())
		os.Exit(0)
	}
	if *fileFlag == "" && len(flag.Args()) == 0 {
		usage()
	}
	filenames := flag.Args()
	files := make(map[github.GistFilename]github.GistFile)
	for _, f := range filenames {
		file := string(f)
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		content := string(buf)
		gistFile := github.GistFile{
			Filename: &file,
			Content:  &content,
		}
		files[github.GistFilename(file)] = gistFile
	}
	if *fileFlag != "" {
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		content := string(buf)
		gistFile := github.GistFile{
			Filename: fileFlag,
			Content:  &content,
		}
		files[github.GistFilename(*fileFlag)] = gistFile
	}

	if !*anonFlag {
		loadAuth()
	} else {
		client = github.NewClient(nil)
	}

	gist := &github.Gist{
		Description: descFlag,
		Public:      publicFlag,
		Files:       files,
	}
	timeout := time.Second * time.Duration(5)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	g, _, err := client.Gists.Create(ctx, gist)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", *g.HTMLURL)
}

var client *github.Client

// GitHub personal access token, from https://github.com/settings/applications.
var authToken string

func loadAuth() {
	const short = ".github-gist-token"
	filename := filepath.Clean(os.Getenv("HOME") + "/" + short)
	shortFilename := filepath.Clean("$HOME/" + short)
	if *tokenFile != "" {
		filename = *tokenFile
		shortFilename = *tokenFile
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("reading token: ", err, "\n\n"+
			"Please create a personal access token at https://github.com/settings/tokens/new\n"+
			"and write it to ", shortFilename, " to use this program.\n"+
			"The token only needs the repo scope, or private_repo if you want to\n"+
			"view or edit issues for private repositories.\n"+
			"The benefit of using a personal access token over using your GitHub\n"+
			"password directly is that you can limit its use and revoke it at any time.\n\n")
	}
	fi, err := os.Stat(filename)
	if fi.Mode()&0077 != 0 {
		log.Fatalf("reading token: %s mode is %#o, want %#o", shortFilename, fi.Mode()&0777, fi.Mode()&0700)
	}
	authToken = strings.TrimSpace(string(data))
	t := &oauth2.Transport{
		Source: &tokenSource{AccessToken: authToken},
	}
	client = github.NewClient(&http.Client{Transport: t})
}

type tokenSource oauth2.Token

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return (*oauth2.Token)(t), nil
}
