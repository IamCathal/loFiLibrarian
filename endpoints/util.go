package endpoints

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func DisallowFileBrowsing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = filepath.Clean(r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

func SendBasicInvalidResponse(w http.ResponseWriter, req *http.Request, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	response := struct {
		Error string `json:"error"`
	}{
		msg,
	}
	json.NewEncoder(w).Encode(response)
}

func isValidInt(strInt string) bool {
	if strInt == "" {
		return false
	}
	_, err := strconv.Atoi(strInt)
	return err == nil
}

func getRealIP(r *http.Request) string {
	if r.Header.Get("X-Real-Ip") != "" {
		return r.Header.Get("X-Real-Ip")
	} else if r.Header.Get("RemoteAddr") != "" {
		return r.Header.Get("RemoteAddr")
	}
	return r.RemoteAddr
}

func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
