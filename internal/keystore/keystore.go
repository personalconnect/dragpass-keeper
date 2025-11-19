package keystore

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/personalconnect/dragpass-keeper/config"
	"github.com/zalando/go-keyring"
)

type RequestMessage struct {
	Action         string `json:"action"`
	Key            string `json:"key,omitempty"`
	PublicKey      string `json:"publickey,omitempty"`
	SessionCode    string `json:"session_code,omitempty"`
	ChallengeToken string `json:"challenge_token,omitempty"`
	Signature      string `json:"signature,omitempty"`
}

type ResponseMessage struct {
	Success     bool   `json:"success"`
	Key         string `json:"key,omitempty"`
	PublicKey   string `json:"publickey,omitempty"`
	SessionCode string `json:"session_code,omitempty"`
	Error       string `json:"error,omitempty"`
}

const (
	serverPubKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF3MG1NZ0FycExYVUhTemJmTGNudAowU1NhTEVhMnhCVms2SXNGTFlOVEl2NzdiZTdYdHhwZzRPd0hDc3JMMzAxV3R0Z2FEWDJBM0pYSnZEQ3FuNXJsCkZGbXNQY2RoeGxwbWdsRjNmODVSMW5KNlB6RW9Dekt1aVVjWE1pc21YSkJteGU2bEpDenZoWXJnbWpKT2xtMkUKY0xJUUpzelFvMUllRml3Mm5wN2c2TzNGSCt2aXRYSkRmV2toakV2RlFGQnd6aFp6cXZUT1o3SDNveUhGZ3RGSwpYeEJwOW5uN2N5L2RmRmVlYkRhSzBmVE1jQ2dEMWxGMjUwZDJMNDdPUmIrbkpEaklObjU4WkZxRVIvTkhWb3dpCnRyanFROU5mWG9rVVFYV2RCWHpjajZDMnNFbGRuR3B5TzFIUzhpYVEvM0RYeXZ2eG9oUWQrWTl3RDJqQnBOajkKYVFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
)

const (
	// Device key related actions
	ActionGetDeviceKey    = "getkey"
	ActionSaveDeviceKey   = "savekey"
	ActionDeleteDeviceKey = "deletekey"

	// Keypair related actions
	ActionGenerateKeypair = "generatekeypair"
	ActionGetPublicKey    = "getpublickey"
	ActionGetPrivatekey   = "getprivatekey"

	// Server public key related actions
	ActionGetServerPubkey = "getserverpubkey"

	// Session code related actions
	ActionSaveSessionCode = "savesessioncode"
	ActionGetSessionCode  = "getsessioncode"
)

// Keypair related functions
func savePrivateKey(privateKey string) error {
	return keyring.Set(config.Service, config.DragPassKeeperPrivateKey, privateKey)
}

func getPrivateKey() (string, error) {
	return keyring.Get(config.Service, config.DragPassKeeperPrivateKey)
}

func deletePrivateKey() error {
	return keyring.Delete(config.Service, config.DragPassKeeperPrivateKey)
}

func savePublicKey(publicKey string) error {
	return keyring.Set(config.Service, config.DragPassKeeperPublicKey, publicKey)
}

func getPublicKey() (string, error) {
	return keyring.Get(config.Service, config.DragPassKeeperPublicKey)
}

// Server public key related functions
func saveServerPublicKey(serverPublicKey string) error {
	return keyring.Set(config.Service, config.DragPassServerPublicKey, base64.StdEncoding.EncodeToString([]byte(serverPublicKey)))
}

func getServerPublicKey() (string, error) {
	serverPubKeyBytes, err := base64.StdEncoding.DecodeString(serverPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode server public key: %v", err)
	}
	return string(serverPubKeyBytes), nil
}

// Device key related functions
func saveDeviceKey(key string) error {
	return keyring.Set(config.Service, config.DeviceKey, key)
}

func getDeviceKey() (string, error) {
	secret, err := keyring.Get(config.Service, config.DeviceKey)
	return secret, err
}

func deleteDeviceKey() error {
	return keyring.Delete(config.Service, config.DeviceKey)
}

// Session code related functions
func saveSessionCode(sessionCode string) error {
	return keyring.Set(config.Service, config.SessionCode, sessionCode)
}

func getSessionCode() (string, error) {
	return keyring.Get(config.Service, config.SessionCode)
}

func deleteSessionCode() error {
	return keyring.Delete(config.Service, config.SessionCode)
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
	case ActionGenerateKeypair:
		log.Println("keypair generation request processing...")

		// Check if challenge token and signature are provided
		if req.ChallengeToken == "" || req.Signature == "" {
			log.Println("keypair generation error: challenge token and signature are required")
			return ResponseMessage{Success: false, Error: "challenge token and signature are required"}
		}

		// Get server public key for signature verification
		serverPubKeyPEM, err := getServerPublicKey()
		if err != nil {
			log.Printf("keypair generation error: failed to get server public key: %v", err)
			return ResponseMessage{Success: false, Error: "failed to get server public key: " + err.Error()}
		}

		// Parse server public key
		serverPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(serverPubKeyPEM))
		if err != nil {
			log.Printf("keypair generation error: failed to parse server public key: %v", err)
			return ResponseMessage{Success: false, Error: "failed to parse server public key: " + err.Error()}
		}

		// Decode the signature from base64
		signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
		if err != nil {
			log.Printf("keypair generation error: failed to decode signature: %v", err)
			return ResponseMessage{Success: false, Error: "failed to decode signature: " + err.Error()}
		}

		// Verify signature using server's public key
		if err := VerifySignature(serverPubKey, req.ChallengeToken, signatureBytes); err != nil {
			log.Printf("keypair generation error: signature verification failed: %v", err)
			return ResponseMessage{Success: false, Error: "signature verification failed: " + err.Error()}
		}
		log.Println("signature verification successful")

		// Delete existing private key if exists
		if err := deletePrivateKey(); err != nil {
			log.Printf("warning: failed to delete existing private key: %v", err)
		}

		// Delete existing session code if exists
		if err := deleteSessionCode(); err != nil {
			log.Printf("warning: failed to delete existing session code: %v", err)
		}

		keyPair, err := GenerateRSAKeyPair()
		if err != nil {
			log.Printf("keypair generation error: %v", err)
			return ResponseMessage{Success: false, Error: "keypair generation failed: " + err.Error()}
		}

		// Save the new private key to the keystore
		if err := savePrivateKey(keyPair.PrivateKey); err != nil {
			log.Printf("private key save error: %v", err)
			return ResponseMessage{Success: false, Error: "private key save failed: " + err.Error()}
		}

		// Save the new public key to the keystore
		if err := savePublicKey(keyPair.PublicKey); err != nil {
			log.Printf("public key save error: %v", err)
			return ResponseMessage{Success: false, Error: "public key save failed: " + err.Error()}
		}

		log.Println("keypair generation and keypair save successful")
		return ResponseMessage{Success: true, PublicKey: keyPair.PublicKey}

	case ActionGetPublicKey:
		log.Println("public key retrieval request processing...")
		publicKey, err := getPublicKey()
		if err != nil {
			log.Printf("public key retrieval error: %v", err)
			return ResponseMessage{Success: false, Error: "public key retrieval failed: " + err.Error()}
		}
		return ResponseMessage{Success: true, PublicKey: publicKey}

	case ActionGetPrivatekey:
		log.Println("private key retrieval request processing...")
		privateKey, err := getPrivateKey()
		if err != nil {
			log.Printf("private key retrieval error: %v", err)
			return ResponseMessage{Success: false, Error: "private key retrieval failed: " + err.Error()}
		}
		return ResponseMessage{Success: true, Key: privateKey}

	case ActionGetServerPubkey:
		log.Println("server public key retrieval request processing...")
		decodeKey, err := getServerPublicKey()
		if err != nil {
			log.Printf("server public key retrieval error: %v", err)
			return ResponseMessage{Success: false, Error: "server public key retrieval failed: " + err.Error()}
		}
		return ResponseMessage{Success: true, PublicKey: decodeKey}

	case ActionGetDeviceKey:
		log.Println("key retrieval request processing...")
		key, err := getDeviceKey()
		if err != nil {
			log.Printf("key retrieval error: %v", err)
			return ResponseMessage{Success: false, Error: "key retrieval failed: " + err.Error()}
		}
		return ResponseMessage{Success: true, Key: key}

	case ActionSaveDeviceKey:
		log.Println("key save request processing...")
		if req.Key == "" {
			log.Println("key save error: Key is empty")
			return ResponseMessage{Success: false, Error: "key to save is empty"}
		}
		if err := saveDeviceKey(req.Key); err != nil {
			log.Printf("key save error: %v", err)
			return ResponseMessage{Success: false, Error: "key save failed: " + err.Error()}
		}
		return ResponseMessage{Success: true}

	case ActionDeleteDeviceKey:
		log.Println("key delete request processing...")
		if err := deleteDeviceKey(); err != nil {
			log.Printf("key delete error: %v", err)
			return ResponseMessage{Success: false, Error: "key delete failed: " + err.Error()}
		}
		return ResponseMessage{Success: true}

	case ActionSaveSessionCode:
		log.Println("session code save request processing...")
		if req.SessionCode == "" {
			log.Println("session code save error: SessionCode is empty")
			return ResponseMessage{Success: false, Error: "session code to save is empty"}
		}
		if err := saveSessionCode(req.SessionCode); err != nil {
			log.Printf("session code save error: %v", err)
			return ResponseMessage{Success: false, Error: "session code save failed: " + err.Error()}
		}
		return ResponseMessage{Success: true}

	case ActionGetSessionCode:
		log.Println("session code retrieval request processing...")
		sessionCode, err := getSessionCode()
		if err != nil {
			log.Printf("session code retrieval error: %v", err)
			return ResponseMessage{Success: false, Error: "session code retrieval failed: " + err.Error()}
		}
		return ResponseMessage{Success: true, SessionCode: sessionCode}

	default:
		log.Printf("unknown action: %s", req.Action)
		return ResponseMessage{Success: false, Error: "unknown action: " + req.Action}
	}
}

func LogRequest(req RequestMessage) {
	safeReq := req
	if safeReq.Action == ActionSaveDeviceKey && safeReq.Key != "" {
		safeReq.Key = "[KEY_MASKED]"
	}
	if safeReq.Action == ActionSaveSessionCode && safeReq.SessionCode != "" {
		safeReq.SessionCode = "[SESSION_CODE_MASKED]"
	}

	log.Printf("received request: %+v", safeReq)
}

func SendResponse(resp ResponseMessage) {
	safeResp := resp
	if safeResp.Key != "" {
		safeResp.Key = "[KEY_MASKED]"
	}
	if safeResp.PublicKey != "" {
		safeResp.PublicKey = "[PUBLIC_KEY_MASKED]"
	}
	if safeResp.SessionCode != "" {
		safeResp.SessionCode = "[SESSION_CODE_MASKED]"
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

func init() {
	// Initialize server public key on startup
	_, err := keyring.Get(config.Service, config.DragPassServerPublicKey)
	if err == nil {
		// Key already exists, no need to initialize
		return
	}

	// Decode the hardcoded server public key
	serverPubKeyBytes, err := base64.StdEncoding.DecodeString(serverPubKey)
	if err != nil {
		log.Printf("warning: failed to decode hardcoded server public key: %v", err)
		return
	}

	// Save to keystore
	if err := saveServerPublicKey(string(serverPubKeyBytes)); err != nil {
		log.Printf("warning: failed to save server public key: %v", err)
		return
	}

	// Generate a new keypair for the client
	keyPair, err := GenerateRSAKeyPair()
	if err != nil {
		log.Printf("keypair generation error: %v", err)
		return
	}

	// Save the new private key to the keystore
	if err := savePrivateKey(keyPair.PrivateKey); err != nil {
		log.Printf("private key save error: %v", err)
		return
	}

	// Save the new public key to the keystore
	if err := savePublicKey(keyPair.PublicKey); err != nil {
		log.Printf("public key save error: %v", err)
		return
	}

	log.Println("server public key initialized successfully")
}
