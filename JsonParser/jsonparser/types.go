package jsonparser

type JSONValue = any
type JSONObject = map[string]JSONValue
type JSONArray = []JSONValue

type KeyValuePair struct {
	key   string
	value JSONValue
}
