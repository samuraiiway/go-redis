package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/go-redis/handler"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc(handler.REDIS_UPSERT_PATH, handler.RedisUpsert).Methods("POST")
	router.HandleFunc(handler.REDIS_LISTEN_PATH, handler.RedisListen).Methods("GET")
	router.HandleFunc(handler.REDIS_GET_ID_PATH, handler.RedisGetID).Methods("GET")
	router.HandleFunc(handler.REDIS_GET_INDEX_PATH, handler.RedisGetIndex).Methods("GET")
	router.HandleFunc(handler.REDIS_GENERATE_PATH, handler.RedisGenerate).Methods("POST")
	http.ListenAndServe(":8000", router)
}
