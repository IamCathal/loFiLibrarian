package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/goodreads"
	"github.com/iamcathal/lofilibrarian/influxdb"
	"github.com/iamcathal/lofilibrarian/openlibrary"
	"github.com/iamcathal/lofilibrarian/util"
	"go.uber.org/zap"
)

var (
	logger    *zap.Logger
	appConfig dtos.AppConfig
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	TEXT              = websocket.TextMessage
	PING_TIMEOUT_WAIT = 1800 * time.Millisecond
)

func InitConfig(conf dtos.AppConfig, newLogger *zap.Logger) {
	appConfig = conf
	logger = newLogger
}

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/status", status).Methods("GET")
	// r.HandleFunc("/lookup", lookUp).Methods("GET")
	r.HandleFunc("/eee", wsLookUp).Methods("GET")
	r.HandleFunc("/livestatus", liveStatus).Methods("GET")
	r.HandleFunc("/openlib", openLibrary).Methods("GET")
	// r.HandleFunc("/lookup", lookUp).Methods("GET")
	// r.Use(logMiddleware)

	r.Use(corsMiddleware)

	wsRouter := r.Path("/ws").Subrouter()
	wsRouter.HandleFunc("/lookup", wsLookUp).Methods("GET")

	r.Handle("/static", http.NotFoundHandler())
	fs := http.FileServer(http.Dir(""))
	r.PathPrefix("/").Handler(DisallowFileBrowsing(fs))
	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func wsLookUp(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().UnixMilli()
	logger.Sugar().Infof("Initiated new ws connection from %s with user-agent: %s", getRealIP(r), r.Header["User-Agent"])
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ctx := context.Background()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.Sugar().Warnf("error upgrading websocket connection (handshake error): %+v", err)
			SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
			ws.Close()
			return
		}
		logger.Sugar().Warnf("error upgrading websocket connection: %+v", err)
		SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
		ws.Close()
		return
	}
	defer ws.Close()

	_, msg, err := ws.ReadMessage()
	if err != nil {
		logger.Warn(err.Error())
	}

	lookUpRequest := dtos.LookupRequest{}
	err = json.Unmarshal(msg, &lookUpRequest)
	if err != nil {
		util.WriteWsError(ctx, "Invalid lookup request given")
		return
	}
	logger.Sugar().Infof("Lookup request was: %+v", lookUpRequest)

	ctx = context.WithValue(ctx, dtos.START_TIME, startTime)
	ctx = context.WithValue(ctx, dtos.REQUEST_ID, lookUpRequest.ID)
	ctx = context.WithValue(ctx, dtos.BOOK_ID, lookUpRequest.BookId)
	ctx = context.WithValue(ctx, dtos.WS, ws)

	if isValid := isValidInt(lookUpRequest.BookId); !isValid {
		errorMsg := fmt.Sprintf("Invalid id ' %s ' given", lookUpRequest.BookId)
		util.WriteWsError(ctx, errorMsg)
		return
	}

	util.WriteWsMessage(ctx, fmt.Sprintf("Processing request to lookup ISBN: %s", lookUpRequest.BookId))

	_, err = goodreads.GetBookDetailsWs(ctx, lookUpRequest.BookId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting book details: %+v", err)
		logger.Sugar().Error(errorMessage)
		util.WriteWsError(ctx, errorMessage)
		return
	}

	logger.Sugar().Infof("Completed book search in %vms", time.Now().UnixMilli()-startTime)
}

func liveStatus(w http.ResponseWriter, r *http.Request) {
	logger.Sugar().Infof("Initiated new ws connection from %s with user-agent: %s", getRealIP(r), r.Header["User-Agent"])

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ctx := context.Background()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.Sugar().Warnf("error upgrading websocket connection (handshake error): %+v", err)
			SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
			ws.Close()
			return
		}
		logger.Sugar().Warnf("error upgrading websocket connection: %+v", err)
		SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
		ws.Close()
		return
	}
	defer ws.Close()

	ctx = context.WithValue(ctx, dtos.WS, ws)

	for {
		util.WriteWsLiveStatus(ctx, appConfig.ApplicationStartUpTime)
		time.Sleep(250 * time.Millisecond)
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	req := dtos.UptimeResponse{
		Status:      "operational",
		Uptime:      time.Duration(time.Since(appConfig.ApplicationStartUpTime).Milliseconds()),
		StartUpTime: appConfig.ApplicationStartUpTime.UnixMilli(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

func openLibrary(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().UnixMilli()
	ctx := context.Background()

	isbn := r.URL.Query().Get("id")
	if isValid := isValidInt(isbn); !isValid {
		errorMsg := fmt.Sprintf("Invalid isbn ' %s ' given", isbn)
		SendBasicInvalidResponse(w, r, errorMsg, http.StatusBadRequest)
		return
	}

	ctx = context.WithValue(ctx, dtos.START_TIME, startTime)
	ctx = context.WithValue(ctx, dtos.BOOK_ID, isbn)

	bookInfo, err := openlibrary.IsbnSearch(ctx, isbn)
	if err != nil {
		SendBasicInvalidResponse(w, r, "failed to lookup book", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookInfo)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if isStaticContent := strings.HasPrefix(r.URL.Path, "/static/"); !isStaticContent {
		// 	logger.Sugar().Infow(fmt.Sprintf("Served request to %v", r.URL.Path),
		// 		zap.String("requestInfo", fmt.Sprintf("%+v", r)))
		// }
		if influxdb.IsInfluxDBEnabled() {
			realIp := getRealIP(r)
			if realIp != "" {
				influxdb.WriteMetricPoint(realIp)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetupCORS(w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
