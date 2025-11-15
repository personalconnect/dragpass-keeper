package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/zalando/go-keyring"
)

const serviceName = "com.blindfold.keeper"
const userName = "device_key"

func SaveKey(key string) error {
	return keyring.Set(serviceName, userName, key)
}
func GetKey() (string, error) {
	secret, err := keyring.Get(serviceName, userName)
	return secret, err
}
func DeleteKey() error {
	return keyring.Delete(serviceName, userName)
}

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

func main() {
	log.SetOutput(os.Stderr)
	log.Println("blindfold extension helper started")

	for {
		req, err := readMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("lost chrome extension connection")
				break
			}
			sendResponse(ResponseMessage{Success: false, Error: "wrong message format: " + err.Error()})
			continue
		}
		logRequest(req)
		resp := handleRequest(req)
		sendResponse(resp)
	}
}

func readMessage() (RequestMessage, error) {
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

func handleRequest(req RequestMessage) ResponseMessage {
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

func logRequest(req RequestMessage) {
	safeReq := req
	if safeReq.Action == ActionSaveKey && safeReq.Key != "" {
		safeReq.Key = "[KEY_MASKED]"
	}

	log.Printf("received request: %+v", safeReq)
}

func sendResponse(resp ResponseMessage) {
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
