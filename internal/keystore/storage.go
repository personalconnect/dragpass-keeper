package keystore

import (
	"github.com/personalconnect/dragpass-keeper/config"
	"github.com/zalando/go-keyring"
)

// Keypair related functions
func savePrivateKey(privateKey string) error {
	return keyring.Set(config.Service, config.DragPassKeeperPrivateKey, privateKey)
}

func getPrivateKey() (string, error) {
	return keyring.Get(config.Service, config.DragPassKeeperPrivateKey)
}

func getPublicKey() (string, error) {
	return keyring.Get(config.Service, config.DragPassKeeperPublicKey)
}

func savePublicKey(publicKey string) error {
	return keyring.Set(config.Service, config.DragPassKeeperPublicKey, publicKey)
}

// Server public key related functions
func saveServerPublicKey(serverPublicKey string) error {
	return keyring.Set(config.Service, config.DragPassServerPublicKey, serverPublicKey)
}

func getServerPublicKey() (string, error) {
	return keyring.Get(config.Service, config.DragPassServerPublicKey)
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
