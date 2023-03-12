package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/goodreads"
	"github.com/iamcathal/lofilibrarian/util"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"go.uber.org/zap"
)

var (
	logger    *zap.Logger
	appConfig dtos.AppConfig
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	InfluxDBClient influxdb2.Client

	TEXT              = websocket.TextMessage
	PING_TIMEOUT_WAIT = 1800 * time.Millisecond
)

func InitConfig(conf dtos.AppConfig, newLogger *zap.Logger) {
	appConfig = conf
	logger = newLogger
}

func InitInfluxClient(client influxdb2.Client) {
	InfluxDBClient = client
}

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/status", status).Methods("POST")
	r.HandleFunc("/lookup", lookUp).Methods("GET")
	r.HandleFunc("/lookupws", WsEndpoint).Methods("GET")
	// r.Use(logMiddleware)

	wsRouter := r.Path("/ws").Subrouter()
	wsRouter.HandleFunc("/lookup", WsEndpoint).Methods("GET")

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

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ctx := context.Background()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.Sugar().Warnf("error upgrading websocket connection (handshake error): %+v", err)
			SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
			return
		}
		logger.Sugar().Warnf("error upgrading websocket connection: %+v", err)
		SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
		return
	}

	_, msg, err := ws.ReadMessage()
	if err != nil {
		logger.Warn(err.Error())
	}

	fmt.Printf("Got message: %s\n", string(msg))
	lookUpRequest := dtos.InitLookupDto{}
	err = json.Unmarshal(msg, &lookUpRequest)
	if err != nil {
		logger.Warn(err.Error())
	}

	ctx = context.WithValue(ctx, "requestId", lookUpRequest.ID)
	ctx = context.WithValue(ctx, "bookId", lookUpRequest.BookId)
	ctx = context.WithValue(ctx, "ws", ws)

	if isValid := isValidInt(lookUpRequest.BookId); !isValid {
		errorMsg := fmt.Sprintf("Invalid id '%s' given", lookUpRequest.BookId)
		util.WriteWsError(ctx, errorMsg)
		ws.Close()
		return
	}

	util.WriteWsMessage(ctx, fmt.Sprintf("Looking up book %s", lookUpRequest.BookId))

	goodreads.GetBookDetailsWs(ctx, lookUpRequest.BookId)

	// ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(PING_TIMEOUT_WAIT))

	ws.Close()
}

// func WsLookUp(w http.ResponseWriter, r *http.Request) {
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		if _, ok := err.(websocket.HandshakeError); !ok {
// 			logger.Sugar().Warnf("error upgrading websocket connection (handshake error): %+v", err)
// 			SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
// 			return
// 		}
// 		logger.Sugar().Warnf("error upgrading websocket connection: %+v", err)
// 		SendBasicInvalidResponse(w, r, "unable to upgrade websocket", http.StatusBadRequest)
// 		return
// 	}

// 	ID := r.URL.Query().Get("id")
// 	if isValid := isValidInt(ID); !isValid {
// 		errorMsg := fmt.Sprintf("Invalid id '%s' given", ID)
// 		SendBasicInvalidResponse(w, r, errorMsg, http.StatusBadRequest)
// 		return
// 	}

// 	util.WriteWsMessage(ws, fmt.Sprintf("Looking up book %s", ID))

// 	goodreads.GetBookDetailsWs(ws, ID)

// }

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
		// if isStaticContent := strings.HasPrefix(r.URL.Path, "/static/"); !isStaticContent {
		// 	logger.Sugar().Infow(fmt.Sprintf("Served request to %v", r.URL.Path),
		// 		zap.String("requestInfo", fmt.Sprintf("%+v", r)))
		// }
		if os.Getenv("INFLUX_ENABLE") != "" {
			realIp := getRealIP(r)
			if realIp != "" {
				writeAPI := InfluxDBClient.WriteAPIBlocking(os.Getenv("ORG"), os.Getenv("LOFI_BUCKET"))
				point := influxdb2.NewPointWithMeasurement("clientIPLog").
					AddTag("clientIP", realIp).
					AddField("service", "lofilibrarian").
					SetTime(time.Now())
				writeAPI.WritePoint(context.Background(), point)
			}
		}
		next.ServeHTTP(w, r)
	})
}
