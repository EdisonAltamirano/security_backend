package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EdisonAltamirano/security_backend/pkg/gst"
	"github.com/EdisonAltamirano/security_backend/pkg/webrtcstream"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"strconv"
    "io/ioutil"
)

func main() {
	logger := setupLogger()

	r := chi.NewRouter()

	r.Use(LogRequests(logger))

	logger.Debugw("initializing gstreamer4")
	gst.Init()

	config := loadConfig(logger)

	var allowedOrigins []string

	if (config.Cors.AllowAllOrigins) {
		allowedOrigins = config.Cors.AllowedOrigins
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"*"},
	}))

	streamStore := make(map[int]*webrtcstream.WebRTCStream)
	cameraServiceUrl := fmt.Sprintf("%s:%d/", config.CameraService.Hostname, config.CameraService.Port)

	streamCtx := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debugw("streamCtx")
			logger.Debugw("streamCtx", "url", r.URL)
			streamIdString := chi.URLParam(r, "streamID")
			
			// Read the request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Errorw("Error reading request body", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Print the request body as a string
			logger.Debugw("streamCtxR", "url", r.URL, "body", string(body))
			logger.Debugw("streamCtx", "streamIdString", streamIdString)
			cameraId, err := strconv.ParseInt(streamIdString, 10, 32)

			if err != nil {
				logger.Errorw("invalid camera id")
				logger.Errorw("invalid camera id", "err", err)
				logger.Errorw("invalid camera id", "cameraId", cameraId)
				http.Error(w, "invalid camera id", 400)
				return
			}

			stream, ok := streamStore[int(cameraId)]

			if !ok {
				logger.Debugw("creating camera stream", "camera id", cameraId)
				resp, err := http.Get("http://" + cameraServiceUrl + "cameras/" + strconv.FormatInt(cameraId, 10))
				if err != nil {
					logger.Errorw("could not get camera info from camera service", "err", err)
					w.WriteHeader(500)
					return
				}
				if resp.StatusCode == http.StatusNotFound {
					logger.Errorw("unknown camera stream requested")
					w.WriteHeader(http.StatusNotFound)
					return
				} else if resp.StatusCode != 200 {
					logger.Errorw("unknown server error")
					w.WriteHeader(500)
					return
				}

				var streamConfig webrtcstream.Config
				bodyDecoder := json.NewDecoder(resp.Body)
				if err := bodyDecoder.Decode(&streamConfig); err != nil {
					logger.Errorw("could not parse camera service response body")
					w.WriteHeader(500)
					return
				}

				newStream, err := webrtcstream.New(streamConfig)
				if err != nil {
					logger.Error("error creating stream: %w", err)
				}
				streamStore[int(cameraId)] = newStream
				stream = newStream
			}

			ctx := context.WithValue(r.Context(), "stream", stream)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r.Route("/{streamID}", func(r chi.Router) {
		r.Use(streamCtx)
		r.Get("/", makeGetStreamHandler(logger))
	})
	logger.Infow("starting web server", "port", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
	logger.Infow("server stopped")

	if !errors.Is(err, http.ErrServerClosed) {
		logger.Panicw("Fatal error", "err", err.Error())
	}

}
