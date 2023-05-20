package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rashad-j/tiktokuploader/configs"
	logger "github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func Authorize(cfg configs.Config) (service *youtube.Service, err error) {
	ctx := context.Background()
	var b []byte
	var config *oauth2.Config
	var client *http.Client

	if b, err = ioutil.ReadFile(cfg.YoutubeClientSecret); err != nil {
		return
	}

	if config, err = google.ConfigFromJSON(b, youtube.YoutubeUploadScope, youtube.YoutubeReadonlyScope); err != nil {
		return
	}

	if client, err = getClient(ctx, config); err != nil {
		return
	}

	if service, err = youtube.New(client); err != nil {
		return
	}

	return
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) (client *http.Client, err error) {
	cacheFile := tokenCacheFile()
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return
		}
		if err = saveToken(cacheFile, tok); err != nil {
			return
		}
	}
	client = config.Client(ctx, tok)
	return
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (tok *oauth2.Token, err error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	_, err = fmt.Scan(&code)
	if err != nil {
		return
	}

	tok, err = config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return
	}
	return
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() string {
	tokenCacheDir := "./youtube/"
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube.json"))
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) error {
	logger.Info().Str("file", file).Msg("saving credential file")
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
