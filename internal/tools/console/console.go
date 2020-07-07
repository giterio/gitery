package console

import "encoding/json"

// Log ...
func Log(data interface{}) (err error) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	println(string(bytes))
	return
}
