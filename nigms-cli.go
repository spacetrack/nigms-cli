/*
 * nigms-cli - nun-ist-genug-mit-schnee command line interface
 *
 * ... allows quick posting to the NIGMS tumblr
 *
 */

package main

import (
    // "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/kurrik/oauth1a"
)

func main() {

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

    // get list of draft posts
    case "drafts":
        httpRequest, err := http.NewRequest("GET", "https://api.tumblr.com/v2/blog/nunistgenugmitschnee.tumblr.com/posts/draft", nil)
        service.Sign(httpRequest, userConfig)
        httpResponse, err := httpClient.Do(httpRequest)

        defer httpResponse.Body.Close()

        httpContents, err := ioutil.ReadAll(httpResponse.Body)

        if err != nil {
            fmt.Println("ERROR: can't read http response body")
            os.Exit(1)
        }

        fmt.Println(string(httpContents))
        os.Exit(0)

    case "debug":
        // nothing
        os.Exit(0)

	default:
		fmt.Println("ERROR: unknown command \"" + os.Args[1] + "\"")
		os.Exit(1)
	}
}
