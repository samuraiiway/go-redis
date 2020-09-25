package redis

import (
	"fmt"

	"github.com/samuraiiway/go-redis/listener"
)

type Empty struct {}

var _ROOT_NAMESPACE = map[string]map[string][]byte{}
var _ROOT_INDEX = map[string]map[string]map[string]Empty{}
var _ROOT_SCHEMA = map[string]map[string][]string{}

func getNameSpace(namespace string) map[string][]byte {
	space, ok := _ROOT_NAMESPACE[namespace]

	if !ok {
		space = map[string][]byte{}
		_ROOT_NAMESPACE[namespace] = space
	}

	return space
}

func getIndex(namespace string) map[string]map[string]Empty {
	index, ok := _ROOT_INDEX[namespace]

	if !ok {
		index = map[string]map[string]Empty{}
		_ROOT_INDEX[namespace] = index
	}

	return index
}

func getSchema(namespace string) map[string][]string {
	schema, ok := _ROOT_SCHEMA[namespace]
	if !ok {
		schema = map[string][]string{}
		_ROOT_SCHEMA[namespace] = schema
	}

	return schema
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
	for key, _ := range ids {
		data, ok := space[key]
		if ok {
			result = append(result, data)
		}
	}

	return result
}

func IndexData(namespace string, id string, keys []string) {
	index := getIndex(namespace)
	removeIndexBySchema(namespace, id, keys)
	updateSchema(namespace, id, keys)

	for _, key := range keys {
		val, ok := index[key]
		if !ok {
			val = map[string]Empty{}
			index[key] = val
		}

		val[id] = Empty{}
	}
}

func RemoveIndexById(namespace string, id string) {
	indexes := getIndex(namespace)

	for index, keys := range indexes {
		for idx, _ := range keys {
			if idx == id {
				delete(keys, idx)
			}
		}
		if len(keys) == 0 {
			delete(indexes, index)
		}
	}
}

func updateSchema(namespace string, id string, keys []string) {
	schema := getSchema(namespace)
	
	oldKeys, ok := schema[id]
	if ok {
		removeIndexBySchema(namespace, id, oldKeys)
	}
	
	schema[id] = keys
}

func removeIndexBySchema(namespace string, id string, keys []string) {
	schema := getSchema(namespace)
	schema[id] = keys

	indexes := getIndex(namespace)
	for _, key := range keys {
		index, ok := indexes[key]

		if ok {
			_, ok := index[id]
			if ok {
				delete(index, id)
			}
		}

		if len(indexes) == 0 {
			delete(indexes, key)
		}
	}
}
