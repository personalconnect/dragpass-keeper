package keystore

import (
	"encoding/base64"
	"fmt"
)

const (
	serverPubKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF3MG1NZ0FycExYVUhTemJmTGNudAowU1NhTEVhMnhCVms2SXNGTFlOVEl2NzdiZTdYdHhwZzRPd0hDc3JMMzAxV3R0Z2FEWDJBM0pYSnZEQ3FuNXJsCkZGbXNQY2RoeGxwbWdsRjNmODVSMW5KNlB6RW9Dekt1aVVjWE1pc21YSkJteGU2bEpDenZoWXJnbWpKT2xtMkUKY0xJUUpzelFvMUllRml3Mm5wN2c2TzNGSCt2aXRYSkRmV2toakV2RlFGQnd6aFp6cXZUT1o3SDNveUhGZ3RGSwpYeEJwOW5uN2N5L2RmRmVlYkRhSzBmVE1jQ2dEMWxGMjUwZDJMNDdPUmIrbkpEaklObjU4WkZxRVIvTkhWb3dpCnRyanFROU5mWG9rVVFYV2RCWHpjajZDMnNFbGRuR3B5TzFIUzhpYVEvM0RYeXZ2eG9oUWQrWTl3RDJqQnBOajkKYVFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
)

func EnsureServerPublicKey() error {
	// Check if the server public key already exists
	_, err := getServerPublicKey()
	if err == nil {
		return nil
	}

	// Decode the hardcoded server public key
	serverPubKeyBytes, err := base64.StdEncoding.DecodeString(serverPubKey)
	if err != nil {
		return fmt.Errorf("failed to decode hardcoded server public key: %v", err)
	}

	// Save to keystore
	if err := saveServerPublicKey(string(serverPubKeyBytes)); err != nil {
		return fmt.Errorf("failed to save server public key: %v", err)
	}

	return nil
}
