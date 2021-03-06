/*
 * nigms-cli - nun-ist-genug-mit-schnee command line interface
 *
 * ... allows quick posting to the NIGMS tumblr
 *
 * $ cd src/github.com/spacetrack/nigms-cli
 * $ go run *.go help
 *
 */

package main

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kurrik/oauth1a"

	"gopkg.in/yaml.v2"
)

func doApiRequest(method string, url string, values url.Values) ([]byte, error) {
	contents, err := ioutil.ReadFile("CREDENTIALS")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := strings.Split(string(contents), "\n")

	service := &oauth1a.Service{
		RequestURL:   "https://www.tumblr.com/oauth/request_token",
		AuthorizeURL: "https://www.tumblr.com/oauth/authorize",
		AccessURL:    "https://www.tumblr.com/oauth/access_token",

		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    lines[0],
			ConsumerSecret: lines[1],
			CallbackURL:    "",
		},

		Signer: new(oauth1a.HmacSha1Signer),
	}

	httpClient := new(http.Client)
	//userConfig := &oauth1a.UserConfig{}
	//userConfig.GetRequestToken(service, httpClient)
	//url, err := userConfig.GetAuthorizeURL(service)

	userConfig := oauth1a.NewAuthorizedConfig(lines[2], lines[3])

	httpRequest, err := http.NewRequest(method, url, strings.NewReader(values.Encode()))

	if err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(1)
	}

	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	service.Sign(httpRequest, userConfig)
	httpResponse, err := httpClient.Do(httpRequest)

	if err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(1)
	}

	defer httpResponse.Body.Close()

	return ioutil.ReadAll(httpResponse.Body)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: please provide a command! Run \"nigms-cli help\" for getting list of commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	// help
	case "?", "-?", "-h", "--help", "help":
		fmt.Println("nigms-cli - nun-ist-genug-mit-schnee command line interface")
		fmt.Println("command: " + os.Args[1])

		os.Exit(0)

	// create a new post
	case "new", "create":
		// the default file to look into: ./post.yaml
		var f *os.File
		var err error
		var fileName = "post.yaml"

		if len(os.Args) > 2 {
			fileName = os.Args[2]
		}

		if fileName == "--" {
			// todo: we can use `fi, err := os.Stdin.Stat()` to find out, if there is someting from stdin
			// see https://coderwall.com/p/zyxyeg/golang-having-fun-with-os-stdin-and-shell-pipes
			f = os.Stdin
		} else {
			f, err = os.Open(fileName)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR 100! %v\n", err)
			os.Exit(1)
		}

		contents, err := ioutil.ReadAll(f)
		//fmt.Println(string(contents))

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR 200! %v\n", err)
			os.Exit(1)
		}

		p := Post{}

		err = yaml.Unmarshal(contents, &p)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR 300! %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("sending data: %+v\n", p)
		//os.Exit(0)

		apiRequestURL := "https://api.tumblr.com/v2/blog/nunistgenugmitschnee.tumblr.com/post"
		apiValues := p.GetTumblrApiValues()

		httpContents, err := doApiRequest("POST", apiRequestURL, apiValues)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! can't read http response body | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	// update existing posting:
	// nigms-cli update <id> <status> <time>
	case "update":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "ERROR! please provide a post id to update\n")
			os.Exit(1)
		}

		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "ERROR! please provide a post status to update\n")
			os.Exit(1)
		}

		//requestURL := "https://api.tumblr.com/v2/blog/nunistgenugmitschnee.tumblr.com/post/edit"

		os.Exit(0)

	// delete a post
	case "delete":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "ERROR! please provide a post id to delete\n")
			os.Exit(1)
		}

		requestURL := "https://api.tumblr.com/v2/blog/nunistgenugmitschnee.tumblr.com/post/delete"

		values := url.Values{}
		values.Set("id", os.Args[2])

		httpContents, err := doApiRequest("POST", requestURL, values)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! can't read http response body | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	// get list of draft posts
	case "drafts", "posts":
		requestURL := "https://api.tumblr.com/v2/blog/nunistgenugmitschnee.tumblr.com/posts"

		if os.Args[1] == "drafts" {
			requestURL = requestURL + "/draft"
		}

		httpContents, err := doApiRequest("POST", requestURL, url.Values{})

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! can't read http response body | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	case "version":
		fmt.Println("Nun ist genug mit Schnee! nigms-cli verson 0.3.1 (2016-03-14)")
		os.Exit(0)

	case "debug":
		// nothing
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "ERROR! unknown command \""+os.Args[1]+"\"\n")
		os.Exit(1)
	}
}
