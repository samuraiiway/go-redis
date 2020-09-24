package handler

import(
	"strconv"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/samuraiiway/go-redis/repository/redis"
	"github.com/samuraiiway/go-redis/util/random"
)

const (
	ID = "id"
	NAMESPACE = "namespace"
	INDEX = "index"
	VALUE = "value"
	NUMBER = "number"
	REDIS_UPSERT_PATH = "/redis/{" + NAMESPACE + "}"
	REDIS_GET_ID_PATH = "/redis/{" + NAMESPACE + "}/{" + ID + "}"
	REDIS_GET_INDEX_PATH = "/redis/{" + NAMESPACE + "}/{" + INDEX + "}/{" + VALUE + "}"
	REDIS_GENERATE_PATH = "/redis/{" + NAMESPACE + "}/generate/{" + NUMBER + "}"
)

func RedisUpsert(w http.ResponseWriter, r *http.Request) {
	// Extract path variable
	vars := mux.Vars(r)
	namespace := vars[NAMESPACE]
	
	// Parse json body
	body, _ := ioutil.ReadAll(r.Body)
	var request map[string]interface{}
	json.Unmarshal(body, &request)
	
	// Create id
	idNode, hasId := request[ID];
	id := fmt.Sprintf("%v", idNode)
	if (idNode == nil || len(id) == 0 || !hasId) {
		id = uuid.New().String()
		request[ID] = id
	}

	// Create indexes
	indexes := []string{}
	indexNode, ok := request[INDEX]
	if (indexNode != nil) {
		switch index := indexNode.(type) {
		case []interface{}:
			for _, value := range index {
				indexes = append(indexes, fmt.Sprintf("%v", value))
			}
		}
	}

	if (ok) {
		delete(request, INDEX)
	}

	// Find index value
	keys := []string{}
	for _, key := range indexes {
		value := request[key]
		if (value != nil) {
			keys = append(keys, fmt.Sprintf("%v:%v", key, value))
		}
	}

	// Build data as json string
	response, _ := json.Marshal(&request)

	// Update data
	redis.UpsertData(namespace, id, &response)

	// Save indexes
	if (len(keys) > 0) {
		redis.IndexData(namespace, id, keys)
	}
	
	// Build response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func RedisGetID(w http.ResponseWriter, r *http.Request) {
	// Extract path variable
	vars := mux.Vars(r)
	namespace := vars[NAMESPACE]
	id := vars[ID]

	// Get data
	response := redis.GetDataByID(namespace, id)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func RedisGetIndex(w http.ResponseWriter, r *http.Request) {
	// Extract path variable
	vars := mux.Vars(r)
	namespace := vars[NAMESPACE]
	index := vars[INDEX]
	value := vars[VALUE]

	// Get data
	datas := redis.GetDataByIndex(namespace, index, value)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["));
	for i, data := range datas {
		w.Write(data)
		if (i != (len(datas) - 1)) {
			w.Write([]byte(","))
		}
	}
	w.Write([]byte("]"));
}

func RedisGenerate(w http.ResponseWriter, r * http.Request) {
	vars := mux.Vars(r)
	namespace := vars[NAMESPACE]
	number := vars[NUMBER]

	num, _ := strconv.Atoi(number)

	for i := 0; i < num; i++ {
		data := map[string]string{}
		keys := []string{}
		id := uuid.New().String()
		
		data["name"] = random.RandomString(2)
		data["password"] = random.RandomString(40)
		data["role"] = random.RandomRole()
		
		keys = append(keys, fmt.Sprintf("%v:%v", "name", data["name"]))
		keys = append(keys, fmt.Sprintf("%v:%v", "password", data["password"]))
		keys = append(keys, fmt.Sprintf("%v:%v", "role", data["role"]))

		response, _ := json.Marshal(data)
		redis.UpsertData(namespace, id, &response)
		redis.IndexData(namespace, id, keys)
	}
}
