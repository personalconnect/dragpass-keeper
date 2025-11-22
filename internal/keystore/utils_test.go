package keystore

import (
	"encoding/json"
	"errors"
	"testing"
)

type MockRequestWithValidator struct {
	Field string `json:"field" validate:"required"`
}

func (m *MockRequestWithValidator) Validate() error {
	if m.Field == "fail" {
		return errors.New("validation failed by mock")
	}
	return nil
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Success: Valid JSON and Validation passes",
			payload:     `{"field": "success"}`,
			shouldError: false,
		},
		{
			name:        "Failure: Invalid JSON format",
			payload:     `{"field": "broken...`,
			shouldError: true,
			errorMsg:    "invalid payload format",
		},
		{
			name:        "Failure: Validation fails",
			payload:     `{"field": "fail"}`,
			shouldError: true,
			errorMsg:    "validation failed by mock",
		},
		{
			name:        "Success: Empty payload (should proceed)",
			payload:     ``,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := func(req MockRequestWithValidator) BaseResponse {
				return BaseResponse{Success: true}
			}

			resp := process(json.RawMessage(tt.payload), mockHandler)

			// 결과 검증
			if tt.shouldError {
				if resp.Success {
					t.Errorf("Expected error but got success")
				}
				if resp.Error != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, resp.Error)
				}
			} else {
				if !resp.Success {
					t.Errorf("Expected success but got error: %s", resp.Error)
				}
			}
		})
	}
}

type MockRequestNonValidator struct {
	Field string `json:"field"`
}

func TestProcess_NoValidator(t *testing.T) {
	payload := json.RawMessage(`{"field": "anything"}`)

	mockHandler := func(req MockRequestNonValidator) BaseResponse {
		if req.Field != "anything" {
			return BaseResponse{Success: false, Error: "data mismatch"}
		}
		return BaseResponse{Success: true}
	}

	resp := process(payload, mockHandler)

	if !resp.Success {
		t.Errorf("Expected success for struct without Validator, but got error: %s", resp.Error)
	}
}
