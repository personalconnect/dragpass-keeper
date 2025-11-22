package keystore

import "encoding/json"

type Validator interface {
	Validate() error
}

// process is catch the pattern of unmarshaling, validating, and handling requests
// [Unmarshal] -> [Validate]-> [Handler call] pattern
func process[T any](payload json.RawMessage, handler func(T) BaseResponse) BaseResponse {
	var req T

	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &req); err != nil {
			return BaseResponse{Success: false, Error: "invalid payload format"}
		}
	}

	if v, ok := any(&req).(Validator); ok {
		if err := v.Validate(); err != nil {
			return BaseResponse{Success: false, Error: err.Error()}
		}
	}

	return handler(req)
}
