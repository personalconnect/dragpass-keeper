package keystore

import (
	"encoding/json"
	"log"
)

// HandleRequest processes incoming requests using the BaseRequest envelope pattern
func HandleRequest(msg []byte) BaseResponse {
	var base BaseRequest
	if err := json.Unmarshal(msg, &base); err != nil {
		log.Printf("failed to unmarshal base request: %v", err)
		return BaseResponse{Success: false, Error: "invalid JSON format"}
	}

	log.Printf("received action: %s", base.Action)

	switch base.Action {
	case ActionPing:
		return process(base.Payload, HandlePing)

	case ActionGenerateKeypair:
		return process(base.Payload, HandleGenerateKeypair)

	case ActionGetDeviceKey:
		return process(base.Payload, HandleGetDeviceKey)

	case ActionSaveDeviceKey:
		return process(base.Payload, HandleSaveDeviceKey)

	case ActionDeleteDeviceKey:
		return process(base.Payload, HandleDeleteDeviceKey)

	case ActionSaveSessionCode:
		return process(base.Payload, HandleSaveSessionCode)

	case ActionGetSessionCode:
		return process(base.Payload, HandleGetSessionCode)

	case ActionGetPublicKey:
		return process(base.Payload, HandleGetPublicKey)

	case ActionSignAlias:
		return process(base.Payload, HandleSignAlias)

	case ActionSignAliasWithTimestamp:
		return process(base.Payload, HandleSignAliasWithTimestamp)

	case ActionSignChallengeToken:
		return process(base.Payload, HandleSignChallengeToken)

	default:
		log.Printf("unknown action: %s", base.Action)
		return BaseResponse{Success: false, Error: "unknown action: " + base.Action}
	}
}
