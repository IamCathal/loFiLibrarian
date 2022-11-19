package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/goodreads"
	"go.uber.org/zap"
)

var (
	logger    *zap.Logger
	appConfig dtos.AppConfig
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func InitConfig(conf dtos.AppConfig, newLogger *zap.Logger) {
	appConfig = conf
	logger = newLogger
}

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/status", status).Methods("POST")
	r.HandleFunc("/lookup", lookUp).Methods("GET")
	r.Use(logMiddleware)

	r.Handle("/static", http.NotFoundHandler())
	fs := http.FileServer(http.Dir(""))
	r.PathPrefix("/").Handler(DisallowFileBrowsing(fs))
	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func lookUp(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if isValid := isValidInt(ID); !isValid {
		errorMsg := fmt.Sprintf("Invalid id '%s' given", ID)
		SendBasicInvalidResponse(w, r, errorMsg, http.StatusBadRequest)
		return
	}
	bookInfo := goodreads.GetBookDetails(ID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookInfo)
}

func status(w http.ResponseWriter, r *http.Request) {
	req := dtos.UptimeResponse{
		Status:      "operational",
		Uptime:      time.Duration(time.Since(appConfig.ApplicationStartUpTime).Milliseconds()),
		StartUpTime: appConfig.ApplicationStartUpTime.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isStaticContent := strings.HasPrefix(r.URL.Path, "/static/"); !isStaticContent {
			logger.Sugar().Infow(fmt.Sprintf("Served request to %v", r.URL.Path),
				zap.String("requestInfo", fmt.Sprintf("%+v", r)))
		}
		next.ServeHTTP(w, r)
	})
}
