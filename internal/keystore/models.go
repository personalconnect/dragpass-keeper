package keystore

import (
	"encoding/json"
	"errors"
)

type BaseRequest struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type PingRequest struct{}
type GetDeviceKeyRequest struct{}
type DeleteDeviceKeyRequest struct{}
type GetSessionCodeRequest struct{}
type GetPublicKeyRequest struct{}
type GetServerPublicKeyRequest struct{}
type SaveDeviceKeyResponseData struct{}
type DeleteDeviceKeyResponseData struct{}

type GenerateKeypairRequest struct {
	ChallengeToken string `json:"challenge_token"`
	Signature      string `json:"signature"`
}

func (r GenerateKeypairRequest) Validate() error {
	if r.ChallengeToken == "" {
		return errors.New("challenge_token is required")
	}
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

type SaveDeviceKeyRequest struct {
	Key string `json:"key"`
}

func (r SaveDeviceKeyRequest) Validate() error {
	if r.Key == "" {
		return errors.New("key is required")
	}
	return nil
}

type SaveSessionCodeRequest struct {
	EncryptedSessionCode string `json:"encrypted_session_code"`
	Signature            string `json:"signature"`
}

func (r SaveSessionCodeRequest) Validate() error {
	if r.EncryptedSessionCode == "" {
		return errors.New("encrypted_session_code is required")
	}
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

type SignAliasRequest struct {
	Alias string `json:"alias"`
}

func (r SignAliasRequest) Validate() error {
	if r.Alias == "" {
		return errors.New("alias is required")
	}
	return nil
}

type SignAliasWithTimestampRequest struct {
	Alias string `json:"alias"`
}

func (r SignAliasWithTimestampRequest) Validate() error {
	if r.Alias == "" {
		return errors.New("alias is required")
	}
	return nil
}

type SignChallengeTokenRequest struct {
	ChallengeToken string `json:"challenge_token"`
	Signature      string `json:"signature"`
}

func (r SignChallengeTokenRequest) Validate() error {
	if r.ChallengeToken == "" {
		return errors.New("challenge_token is required")
	}
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

type BaseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type PingResponseData struct {
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Path    string `json:"path"`
}

type GenerateKeypairResponseData struct {
	PublicKey string `json:"publickey"`
}

type GetDeviceKeyResponseData struct {
	Key string `json:"key"`
}

type SaveSessionCodeResponseData struct {
	SessionCode string `json:"session_code"`
}

type GetSessionCodeResponseData struct {
	SessionCode string `json:"session_code"`
}

type GetPublicKeyResponseData struct {
	PublicKey string `json:"publickey"`
}

type GetServerPublicKeyResponseData struct {
	PublicKey string `json:"publickey"`
}

type SignAliasResponseData struct {
	Signature string `json:"signature"`
	PublicKey string `json:"publickey"`
}

type SignAliasWithTimestampResponseData struct {
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

type SignChallengeTokenResponseData struct {
	Signature string `json:"signature"`
}
