package keystore

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// HandlePing handles ping requests
func HandlePing(req PingRequest) BaseResponse {
	log.Println("ping request processing...")
	return BaseResponse{
		Success: true,
		Data: PingResponseData{
			Version: Version,
			Hash:    BinaryHash,
			Path:    BinaryPath,
		},
	}
}

// HandleGenerateKeypair handles keypair generation requests
func HandleGenerateKeypair(req GenerateKeypairRequest) BaseResponse {
	log.Println("keypair generation request processing...")

	// Get server public key for signature verification
	serverPubKeyPEM, err := getServerPublicKey()
	if err != nil {
		log.Printf("keypair generation error: failed to get server public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to get server public key: " + err.Error()}
	}

	// Parse server public key
	serverPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(serverPubKeyPEM))
	if err != nil {
		log.Printf("keypair generation error: failed to parse server public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse server public key: " + err.Error()}
	}

	// Decode the signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		log.Printf("keypair generation error: failed to decode signature: %v", err)
		return BaseResponse{Success: false, Error: "failed to decode signature: " + err.Error()}
	}

	// Verify signature using server's public key
	if err := VerifySignature(serverPubKey, req.ChallengeToken, signatureBytes); err != nil {
		log.Printf("keypair generation error: signature verification failed: %v", err)
		return BaseResponse{Success: false, Error: "signature verification failed: " + err.Error()}
	}
	log.Println("signature verification successful")

	keyPair, err := GenerateRSAKeyPair()
	if err != nil {
		log.Printf("keypair generation error: %v", err)
		return BaseResponse{Success: false, Error: "keypair generation failed: " + err.Error()}
	}

	// Save(Overwrite) the new private key to the keystore
	if err := savePrivateKey(keyPair.PrivateKey); err != nil {
		log.Printf("private key save error: %v", err)
		return BaseResponse{Success: false, Error: "private key save failed: " + err.Error()}
	}

	// Save the new public key to the keystore
	if err := savePublicKey(keyPair.PublicKey); err != nil {
		log.Printf("public key save error: %v", err)
		return BaseResponse{Success: false, Error: "public key save failed: " + err.Error()}
	}

	// Delete existing session code if exists
	if err := deleteSessionCode(); err != nil {
		log.Printf("warning: failed to delete existing session code: %v", err)
	}

	log.Println("keypair generation and keypair save successful")
	return BaseResponse{Success: true, Data: GenerateKeypairResponseData{PublicKey: keyPair.PublicKey}}
}

// HandleGetDeviceKey handles device key retrieval requests
func HandleGetDeviceKey(req GetDeviceKeyRequest) BaseResponse {
	log.Println("key retrieval request processing...")
	key, err := getDeviceKey()
	if err != nil {
		log.Printf("key retrieval error: %v", err)
		return BaseResponse{Success: false, Error: "key retrieval failed: " + err.Error()}
	}
	return BaseResponse{Success: true, Data: GetDeviceKeyResponseData{Key: key}}
}

// HandleSaveDeviceKey handles device key save requests
func HandleSaveDeviceKey(req SaveDeviceKeyRequest) BaseResponse {
	log.Println("key save request processing...")
	if err := saveDeviceKey(req.Key); err != nil {
		log.Printf("key save error: %v", err)
		return BaseResponse{Success: false, Error: "key save failed: " + err.Error()}
	}
	return BaseResponse{Success: true}
}

// HandleDeleteDeviceKey handles device key deletion requests
func HandleDeleteDeviceKey(req DeleteDeviceKeyRequest) BaseResponse {
	log.Println("key delete request processing...")
	if err := deleteDeviceKey(); err != nil {
		log.Printf("key delete error: %v", err)
		return BaseResponse{Success: false, Error: "key delete failed: " + err.Error()}
	}
	return BaseResponse{Success: true}
}

// HandleSaveSessionCode handles session code save requests
func HandleSaveSessionCode(req SaveSessionCodeRequest) BaseResponse {
	log.Println("encrypted session code save request processing...")

	// Get server public key for signature verification
	serverPubKeyPEM, err := getServerPublicKey()
	if err != nil {
		log.Printf("session code save error: failed to get server public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to get server public key: " + err.Error()}
	}

	// Parse server public key
	serverPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(serverPubKeyPEM))
	if err != nil {
		log.Printf("session code save error: failed to parse server public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse server public key: " + err.Error()}
	}

	// Decode the signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		log.Printf("session code save error: failed to decode signature: %v", err)
		return BaseResponse{Success: false, Error: "failed to decode signature: " + err.Error()}
	}

	// Verify signature using server's public key
	if err := VerifySignature(serverPubKey, req.EncryptedSessionCode, signatureBytes); err != nil {
		log.Printf("session code save error: signature verification failed: %v", err)
		return BaseResponse{Success: false, Error: "signature verification failed: " + err.Error()}
	}
	log.Println("signature verification successful")

	// Get the Helper's private key from keystore
	privateKeyPEM, err := getPrivateKey()
	if err != nil {
		log.Printf("session code save error: failed to get private key: %v", err)
		return BaseResponse{Success: false, Error: "failed to get private key: " + err.Error()}
	}

	// Parse the private key
	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		log.Printf("session code save error: failed to parse private key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse private key: " + err.Error()}
	}

	// Decode the encrypted session code from base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(req.EncryptedSessionCode)
	if err != nil {
		log.Printf("session code save error: failed to decode encrypted session code: %v", err)
		return BaseResponse{Success: false, Error: "failed to decode encrypted session code: " + err.Error()}
	}

	// Decrypt the session code using Helper's private key
	decryptedBytes, err := DecryptData(privateKey, encryptedBytes)
	if err != nil {
		log.Printf("session code save error: failed to decrypt session code: %v", err)
		return BaseResponse{Success: false, Error: "failed to decrypt session code: " + err.Error()}
	}

	sessionCode := string(decryptedBytes)

	// Save the decrypted session code
	if err := saveSessionCode(sessionCode); err != nil {
		log.Printf("session code save error: %v", err)
		return BaseResponse{Success: false, Error: "session code save failed: " + err.Error()}
	}

	log.Println("session code decryption and save successful")
	return BaseResponse{Success: true, Data: SaveSessionCodeResponseData{SessionCode: sessionCode}}
}

// HandleGetSessionCode handles session code retrieval requests
func HandleGetSessionCode(req GetSessionCodeRequest) BaseResponse {
	log.Println("session code retrieval request processing...")
	sessionCode, err := getSessionCode()
	if err != nil {
		log.Printf("session code retrieval error: %v", err)
		return BaseResponse{Success: false, Error: "session code retrieval failed: " + err.Error()}
	}
	return BaseResponse{Success: true, Data: GetSessionCodeResponseData{SessionCode: sessionCode}}
}

// HandleGetPublicKey handles public key retrieval requests
func HandleGetPublicKey(req GetPublicKeyRequest) BaseResponse {
	log.Println("public key retrieval request processing...")
	publicKeyPEM, err := getPublicKey()
	if err != nil {
		log.Printf("public key retrieval error: %v", err)
		return BaseResponse{Success: false, Error: "public key retrieval failed: " + err.Error()}
	}
	log.Println("public key retrieval successful")
	return BaseResponse{Success: true, Data: GetPublicKeyResponseData{PublicKey: publicKeyPEM}}
}

// HandleGetServerPublicKey handles server public key retrieval requests
func HandleGetServerPublicKey(req GetServerPublicKeyRequest) BaseResponse {
	log.Println("server public key retrieval request processing...")
	serverPublicKeyPEM, err := getServerPublicKey()
	if err != nil {
		log.Printf("server public key retrieval error: %v", err)
		return BaseResponse{Success: false, Error: "server public key retrieval failed: " + err.Error()}
	}
	log.Println("server public key retrieval successful")
	return BaseResponse{Success: true, Data: GetServerPublicKeyResponseData{PublicKey: serverPublicKeyPEM}}
}

// HandleSignAlias handles alias signing requests (signup flow)
func HandleSignAlias(req SignAliasRequest) BaseResponse {
	log.Println("alias signing request processing...")

	_, keyErr := getPrivateKey()
	_, sessionErr := getSessionCode()
	if keyErr == nil && sessionErr == nil {
		log.Println("alias signing error: device already registered and session code exists")
		return BaseResponse{Success: false, Error: "device already registered. this device has already been registered for signup"}
	}

	log.Println("generating new keypair for signup...")
	keyPair, err := GenerateRSAKeyPair()
	if err != nil {
		log.Printf("keypair generation error: %v", err)
		return BaseResponse{Success: false, Error: "keypair generation failed: " + err.Error()}
	}

	// Save the new private key to the keystore
	if err := savePrivateKey(keyPair.PrivateKey); err != nil {
		log.Printf("private key save error: %v", err)
		return BaseResponse{Success: false, Error: "private key save failed: " + err.Error()}
	}

	// Save the new public key to the keystore
	if err := savePublicKey(keyPair.PublicKey); err != nil {
		log.Printf("public key save error: %v", err)
		return BaseResponse{Success: false, Error: "public key save failed: " + err.Error()}
	}

	log.Println("keypair generated successfully for signup")
	privateKeyPEM := keyPair.PrivateKey

	// Parse the private key
	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		log.Printf("alias signing error: failed to parse private key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse private key: " + err.Error()}
	}

	// Sign the alias using the Helper's private key
	signatureBytes, err := SignData(privateKey, req.Alias)
	if err != nil {
		log.Printf("alias signing error: failed to sign alias: %v", err)
		return BaseResponse{Success: false, Error: "failed to sign alias: " + err.Error()}
	}

	// Encode the signature to base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	// Get the Helper's public key from keystore
	publicKeyPEM, err := getPublicKey()
	if err != nil {
		log.Printf("alias signing error: failed to get public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to get public key: " + err.Error()}
	}

	log.Println("alias signing successful")
	return BaseResponse{Success: true, Data: SignAliasResponseData{Signature: signatureBase64, PublicKey: publicKeyPEM}}
}

// HandleSignAliasWithTimestamp handles alias with timestamp signing requests (login flow)
func HandleSignAliasWithTimestamp(req SignAliasWithTimestampRequest) BaseResponse {
	log.Println("alias with timestamp signing request processing...")

	// Generate current timestamp
	timestamp := time.Now().Unix()

	// Get the Helper's private key from keystore (must exist for login)
	privateKeyPEM, err := getPrivateKey()
	if err != nil {
		log.Printf("alias signing error: keypair not found. device not registered: %v", err)
		return BaseResponse{Success: false, Error: "device not registered. please complete signup first"}
	}

	// Parse the private key
	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		log.Printf("alias signing error: failed to parse private key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse private key: " + err.Error()}
	}

	// Create payload: Alias + ":" + Timestamp (matching server format)
	payload := fmt.Sprintf("%s:%d", req.Alias, timestamp)

	// Sign the payload using the Helper's private key
	signatureBytes, err := SignData(privateKey, payload)
	if err != nil {
		log.Printf("alias signing error: failed to sign alias with timestamp: %v", err)
		return BaseResponse{Success: false, Error: "failed to sign alias with timestamp: " + err.Error()}
	}

	// Encode the signature to base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	log.Println("alias with timestamp signing successful")
	return BaseResponse{Success: true, Data: SignAliasWithTimestampResponseData{Signature: signatureBase64, Timestamp: timestamp}}
}

// HandleSignChallengeToken handles challenge token signing requests
func HandleSignChallengeToken(req SignChallengeTokenRequest) BaseResponse {
	log.Println("challenge token signing request processing...")

	// Get server public key for signature verification
	serverPubKeyPEM, err := getServerPublicKey()
	if err != nil {
		log.Printf("challenge token signing error: failed to get server public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to get server public key: " + err.Error()}
	}

	// Parse server public key
	serverPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(serverPubKeyPEM))
	if err != nil {
		log.Printf("challenge token signing error: failed to parse server public key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse server public key: " + err.Error()}
	}

	// Decode the signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		log.Printf("challenge token signing error: failed to decode signature: %v", err)
		return BaseResponse{Success: false, Error: "failed to decode signature: " + err.Error()}
	}

	// Verify signature using server's public key
	if err := VerifySignature(serverPubKey, req.ChallengeToken, signatureBytes); err != nil {
		log.Printf("challenge token signing error: signature verification failed: %v", err)
		return BaseResponse{Success: false, Error: "signature verification failed: " + err.Error()}
	}
	log.Println("server signature verification successful")

	// Get the Helper's private key from keystore
	privateKeyPEM, err := getPrivateKey()
	if err != nil {
		log.Printf("challenge token signing error: failed to get private key: %v", err)
		return BaseResponse{Success: false, Error: "failed to get private key: " + err.Error()}
	}

	// Parse the private key
	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		log.Printf("challenge token signing error: failed to parse private key: %v", err)
		return BaseResponse{Success: false, Error: "failed to parse private key: " + err.Error()}
	}

	// Sign the challenge token using Helper's private key
	challengeSignatureBytes, err := SignData(privateKey, req.ChallengeToken)
	if err != nil {
		log.Printf("challenge token signing error: failed to sign challenge token: %v", err)
		return BaseResponse{Success: false, Error: "failed to sign challenge token: " + err.Error()}
	}

	// Encode the signature to base64
	challengeSignatureBase64 := base64.StdEncoding.EncodeToString(challengeSignatureBytes)

	log.Println("challenge token signing successful")
	return BaseResponse{Success: true, Data: SignChallengeTokenResponseData{Signature: challengeSignatureBase64}}
}
