package http

import (
	"io"
	"net/http"

	"github.com/cnnrznn/gamesaves/pkg/store/googledrive"
)

func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	storeName := r.Header.Get("store")

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

func Serve() error {
	http.HandleFunc("/authorize", HandleAuthorize)
	http.HandleFunc("/upload", HandleUpload)
	return http.ListenAndServeTLS(
		":8080",
		"server.crt",
		"server.key",
		nil,
	)
}
