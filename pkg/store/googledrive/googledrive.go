package googledrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"

	"github.com/cnnrznn/gamesaves/pkg/store"
)

type driveStore struct {
	client *http.Client
}

func (s *driveStore) Upload(ctx context.Context, fn string, data []byte) error {
	return nil
}

func (s *driveStore) Download(ctx context.Context, fn string) ([]byte, error) {
	return nil, nil
}

func New(ctx context.Context, accessToken string) (store.Store, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if len(accessToken) == 0 {
		return nil, fmt.Errorf("missing access token")
	}

	token := &oauth2.Token{}
	err = json.Unmarshal([]byte(accessToken), token)
	if err != nil {
		return nil, fmt.Errorf("unable to parse gamesaves token: %v", err)
	}

	client := config.Client(ctx, token)

	return &driveStore{
		client: client,
	}, nil
}

func Authorize() (*oauth2.Token, error) {
	var token *oauth2.Token
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// listen on localhost and receive code
		srv := &http.Server{Addr: ":80"}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL)
			token = &oauth2.Token{}
			_ = srv.Shutdown(context.Background())
		})
		_ = srv.ListenAndServe()
	}()

	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	url := config.AuthCodeURL("googledrive", oauth2.AccessTypeOffline)
	fmt.Println(
		"Visit this URL in your browser to authorize gamesaves: %s",
		url,
	)

	// Wait for server to receive
	wg.Wait()

	return token, nil
}

func Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	token, err := config.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token: %w", err)
	}

	return token, nil
}

func loadConfig() (*oauth2.Config, error) {
	bs, err := json.Marshal(clientConfig)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(bs, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to parse gamesaves secret file to config: %v",
			err,
		)
	}

	return config, nil
}
