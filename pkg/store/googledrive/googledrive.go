package googledrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"

	"github.com/cnnrznn/gamesaves/pkg/store"
)

type driveStore struct {
	client *http.Client
}

func (s *driveStore) Authorize() error {
	return nil
}

func (s *driveStore) Upload(ctx context.Context, fn string, data []byte) error {
	return nil
}

func (s *driveStore) Download(ctx context.Context, fn string) ([]byte, error) {
	return nil, nil
}

func New(ctx context.Context, accessToken string) (store.Store, error) {
	bs, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read client secret file: %v",
			err,
		)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(bs, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to parse client secret file to config: %v",
			err,
		)
	}

	if len(accessToken) == 0 {
		return nil, fmt.Errorf("missing access token")
	}

	token := &oauth2.Token{}
	err = json.Unmarshal([]byte(accessToken), token)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client token: %v", err)
	}

	client := config.Client(ctx, token)

	return &driveStore{
		client: client,
	}, nil
}
