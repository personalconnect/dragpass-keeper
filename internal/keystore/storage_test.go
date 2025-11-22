package keystore

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestMain(m *testing.M) {
	keyring.MockInit()
	m.Run()
}

func TestPrivateKeyOperations(t *testing.T) {
	expectedKey := "mock-private-key-12345"

	if err := savePrivateKey(expectedKey); err != nil {
		t.Fatalf("Failed to save private key: %v", err)
	}

	got, err := getPrivateKey()
	if err != nil {
		t.Fatalf("Failed to get private key: %v", err)
	}

	if got != expectedKey {
		t.Errorf("Private key mismatch.\nGot: %s\nWant: %s", got, expectedKey)
	}
}

func TestPublicKeyOperations(t *testing.T) {
	expectedKey := "mock-public-key-67890"

	if err := savePublicKey(expectedKey); err != nil {
		t.Fatalf("Failed to save public key: %v", err)
	}

	got, err := getPublicKey()
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}

	if got != expectedKey {
		t.Errorf("Public key mismatch.\nGot: %s\nWant: %s", got, expectedKey)
	}
}

func TestServerPublicKeyOperations(t *testing.T) {
	expectedKey := "mock-server-public-key-abcde"

	if err := saveServerPublicKey(expectedKey); err != nil {
		t.Fatalf("Failed to save server public key: %v", err)
	}

	got, err := getServerPublicKey()
	if err != nil {
		t.Fatalf("Failed to get server public key: %v", err)
	}

	if got != expectedKey {
		t.Errorf("Server public key mismatch.\nGot: %s\nWant: %s", got, expectedKey)
	}
}

func TestDeviceKeyOperations(t *testing.T) {
	expectedKey := "mock-device-key-secret"

	if err := saveDeviceKey(expectedKey); err != nil {
		t.Fatalf("Failed to save device key: %v", err)
	}

	got, err := getDeviceKey()
	if err != nil {
		t.Fatalf("Failed to get device key: %v", err)
	}
	if got != expectedKey {
		t.Errorf("Device key mismatch.\nGot: %s\nWant: %s", got, expectedKey)
	}

	if err := deleteDeviceKey(); err != nil {
		t.Fatalf("Failed to delete device key: %v", err)
	}

	_, err = getDeviceKey()
	if err == nil {
		t.Error("Expected error after deleting device key, but got nil")
	} else if err != keyring.ErrNotFound {
		t.Logf("Correctly received error after deletion: %v", err)
	}
}

func TestSessionCodeOperations(t *testing.T) {
	expectedCode := "mock-session-code-xyz"

	if err := saveSessionCode(expectedCode); err != nil {
		t.Fatalf("Failed to save session code: %v", err)
	}

	got, err := getSessionCode()
	if err != nil {
		t.Fatalf("Failed to get session code: %v", err)
	}
	if got != expectedCode {
		t.Errorf("Session code mismatch.\nGot: %s\nWant: %s", got, expectedCode)
	}

	if err := deleteSessionCode(); err != nil {
		t.Fatalf("Failed to delete session code: %v", err)
	}

	_, err = getSessionCode()
	if err == nil {
		t.Error("Expected error after deleting session code, but got nil")
	}
}
