/*
Gist is a client for creating GitHub Gists.

	usage: gist -f file1.txt,file2.txt,file3.txt [-d "gist description"]

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
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

var (
	fileFlag   = flag.String("f", "", "`comma separated list of file(s)` to upload as a Gist. first file names Gist")
	descFlag   = flag.String("d", "", "description for Gist")
	publicFlag = flag.Bool("p", false, "create a public Gist")
	tokenFile  = flag.String("token", "", "read GitHub personal access token from `file` (default $HOME/.github-gist-token)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gist -f file,file2 [-d string] [-p]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("gist: ")

	if *fileFlag == "" {
		usage()
	}

	files := make(map[github.GistFilename]github.GistFile)
	filenames := strings.FieldsFunc(*fileFlag, func(c rune) bool {
		return c == ','
	})
	for _, file := range filenames {
		f := string(file)
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		content := string(buf)
		gistFile := github.GistFile{
			Filename: &f,
			Content:  &content,
		}
		files[github.GistFilename(file)] = gistFile
	}

	loadAuth()

	gist := &github.Gist{
		Description: descFlag,
		Public:      publicFlag,
		Files:       files,
	}
	g, _, err := client.Gists.Create(gist)
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
