package redis

import (
	"fmt"

	"github.com/samuraiiway/go-redis/listener"
)

var _ROOT_NAMESPACE = map[string]map[string][]byte{}
var _ROOT_INDEX = map[string]map[string][]string{}

func getNameSpace(namespace string) map[string][]byte {
	space, ok := _ROOT_NAMESPACE[namespace]

	if !ok {
		space = map[string][]byte{}
		_ROOT_NAMESPACE[namespace] = space
	}

	return space
}

func getIndex(namespace string) map[string][]string {
	index, ok := _ROOT_INDEX[namespace]

	if !ok {
		index = map[string][]string{}
		_ROOT_INDEX[namespace] = index
	}

	return index
}

func UpsertData(namespace string, id string, data []byte) {
	space := getNameSpace(namespace)
	space[id] = data
	listener.SendMessage(namespace, data)
}

func GetDataByID(namespace string, id string) []byte {
	space := getNameSpace(namespace)
	return space[id]
}

func GetDataByIndex(namespace string, index string, value string) [][]byte {
	result := [][]byte{}

	space := getNameSpace(namespace)
	indexes := getIndex(namespace)
	ids := indexes[fmt.Sprintf("%v:%v", index, value)]
	for _, id := range ids {
		data, ok := space[id]
		if ok {
			result = append(result, data)
		}
	}

	return result
}

func IndexData(namespace string, id string, keys []string) {
	index := getIndex(namespace)

	for _, key := range keys {
		ids := index[key]
		ids = append(ids, id)
		index[key] = ids
	}
}
