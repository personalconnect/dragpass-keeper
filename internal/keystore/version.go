package keystore

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

const (
	Version = "0.0.4"
)

var (
	BinaryHash string
	BinaryPath string
)

func LoadBinaryInfo() error {
	var err error
	BinaryPath, err = os.Executable()
	if err != nil {
		return err
	}

	file, err := os.Open(BinaryPath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}
	BinaryHash = hex.EncodeToString(hasher.Sum(nil))

	return nil
}
