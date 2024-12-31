package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/cnnrznn/gamesaves/pkg/store/googledrive"
)

func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	storeName := r.URL.Query().Get("store")

	switch storeName {
	case "googledrive":
		googledrive.Authorize(w, r)
	default:
		http.Error(
			w,
			"store not available for authorization",
			http.StatusForbidden,
		)
	}
}

func HandleAuthCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	storeName := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	var accessToken *oauth2.Token

	switch storeName {
	case "googledrive":
		token, err := googledrive.Exchange(r.Context(), code)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error exchanging code for token: %s", err),
				http.StatusInternalServerError,
			)
		}
		accessToken = token
	default:
		http.Error(
			w,
			fmt.Errorf("store not available for authorization").Error(),
			http.StatusForbidden,
		)
	}

	bs, err := json.Marshal(accessToken)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("unable to marshal access token"),
			http.StatusInternalServerError,
		)
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to write response: %s", err),
			http.StatusInternalServerError,
		)
	}
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	accessToken := r.Header.Get("access_token")
	store := r.Header.Get("store")
	game := r.Header.Get("game")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch store {
	case "googledrive":
		store, err := googledrive.New(r.Context(), accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = store.Upload(r.Context(), game, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "unrecognized store", http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/authorize", HandleAuthorize)
	http.HandleFunc("/oauth/code", HandleAuthCode)
	http.HandleFunc("/upload", HandleUpload)
	err := http.ListenAndServe(
		":8080",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}
