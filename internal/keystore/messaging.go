package keystore

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// MaxMessageSize defines the maximum allowed message size (10MB)
// This prevents memory exhaustion attacks from malicious extensions
const MaxMessageSize uint32 = 10 * 1024 * 1024 // 10MB

type Messenger struct {
	in  io.Reader
	out io.Writer
}

func NewMessenger(in io.Reader, out io.Writer) *Messenger {
	return &Messenger{
		in:  in,
		out: out,
	}
}

// ReadMessage reads a length-prefixed message from the input
func (m *Messenger) ReadMessage() ([]byte, error) {
	var length uint32

	if err := binary.Read(m.in, binary.LittleEndian, &length); err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("failed to read message length: %v", err)
	}

	if length > MaxMessageSize {
		return nil, fmt.Errorf("message size %d exceeds maximum allowed size %d", length, MaxMessageSize)
	}

	if length == 0 {
		return nil, fmt.Errorf("invalid message: zero length")
	}

	msgBody := make([]byte, length)
	if _, err := io.ReadFull(m.in, msgBody); err != nil {
		return nil, fmt.Errorf("failed to read message body: %v", err)
	}

	return msgBody, nil
}

func (m *Messenger) SendResponse(resp BaseResponse) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("response serialization error: %v", err)
	}

	logSafeResponse(resp)

	// Write response length
	if err := binary.Write(m.out, binary.LittleEndian, uint32(len(respBytes))); err != nil {
		return fmt.Errorf("failed to write response length: %w", err)
	}

	// Write body
	if _, err := m.out.Write(respBytes); err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}

	return nil
}

func logSafeResponse(resp BaseResponse) {
	safeResp := BaseResponse{
		Success: resp.Success,
		Error:   resp.Error,
		Data:    "[DATA_MASKED]",
	}

	if resp.Data == nil {
		safeResp.Data = nil
	}

	safeBytes, _ := json.Marshal(safeResp)
	log.Printf("sending response: %s", string(safeBytes))
}
