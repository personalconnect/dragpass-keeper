package keystore

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/personalconnect/dragpass-keeper/config"
	"github.com/zalando/go-keyring"
)

type RequestMessage struct {
	Action string `json:"action"`
	Key    string `json:"key,omitempty"`
}
type ResponseMessage struct {
	Success bool   `json:"success"`
	Key     string `json:"key,omitempty"`
	Error   string `json:"error,omitempty"`
}

const (
	ActionGetKey    = "getkey"
	ActionSaveKey   = "savekey"
	ActionDeleteKey = "deletekey"
)

func SaveKey(key string) error {
	return keyring.Set(config.Service, config.DeviceManager, key)
}
func GetKey() (string, error) {
	secret, err := keyring.Get(config.Service, config.DeviceManager)
	return secret, err
}
func DeleteKey() error {
	return keyring.Delete(config.Service, config.DeviceManager)
}

func ReadMessage() (RequestMessage, error) {
	var length uint32
	if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
		if err == io.EOF {
			return RequestMessage{}, io.EOF
		}
		log.Printf("length read error: %v", err)
		return RequestMessage{}, fmt.Errorf("length read error: %v", err)
	}

	msgBody := make([]byte, length)
	if _, err := io.ReadFull(os.Stdin, msgBody); err != nil {
		log.Printf("message body read error: %v", err)
		return RequestMessage{}, fmt.Errorf("message body read error: %v", err)
	}

	var req RequestMessage
	if err := json.Unmarshal(msgBody, &req); err != nil {
		log.Printf("message parsing error: %v. (original: %s)", err, string(msgBody))
		return RequestMessage{}, fmt.Errorf("message parsing error: %v", err)
	}

	return req, nil
}

func HandleRequest(req RequestMessage) ResponseMessage {
	switch req.Action {
	case ActionGetKey:
		log.Println("key retrieval request processing...")
		key, err := GetKey()
		if err != nil {
			log.Printf("key retrieval error: %v", err)
			return ResponseMessage{Success: false, Error: "key retrieval failed: " + err.Error()}
		}
		return ResponseMessage{Success: true, Key: key}

	case ActionSaveKey:
		log.Println("key save request processing...")
		if req.Key == "" {
			log.Println("key save error: Key is empty")
			return ResponseMessage{Success: false, Error: "key to save is empty"}
		}
		if err := SaveKey(req.Key); err != nil {
			log.Printf("key save error: %v", err)
			return ResponseMessage{Success: false, Error: "key save failed: " + err.Error()}
		}
		return ResponseMessage{Success: true}

	case ActionDeleteKey:
		log.Println("key delete request processing...")
		if err := DeleteKey(); err != nil {
			log.Printf("key delete error: %v", err)
			return ResponseMessage{Success: false, Error: "key delete failed: " + err.Error()}
		}
		return ResponseMessage{Success: true}

	default:
		log.Printf("unknown action: %s", req.Action)
		return ResponseMessage{Success: false, Error: "unknown action: " + req.Action}
	}
}

func LogRequest(req RequestMessage) {
	safeReq := req
	if safeReq.Action == ActionSaveKey && safeReq.Key != "" {
		safeReq.Key = "[KEY_MASKED]"
	}

	log.Printf("received request: %+v", safeReq)
}

func SendResponse(resp ResponseMessage) {
	safeResp := resp
	if safeResp.Key != "" {
		safeResp.Key = "[KEY_MASKED]"
	}
	safeRespBytes, _ := json.Marshal(safeResp)
	log.Printf("sending response: %s", string(safeRespBytes))

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response serialization error: %v", err)
		return
	}

	if err := binary.Write(os.Stdout, binary.LittleEndian, uint32(len(respBytes))); err != nil {
		log.Printf("response length write error: %v", err)
		return
	}

	if _, err := os.Stdout.Write(respBytes); err != nil {
		log.Printf("response body write error: %v", err)
		return
	}
}
