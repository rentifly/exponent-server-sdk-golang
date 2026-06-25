package expo

import (
	"encoding/json"
	"testing"
)

func TestGetReceiptsResponseUnmarshal(t *testing.T) {
	raw := `{
		"data": {
			"019f0039-9cb1-7088-a76b-504af6cea9b7": {
				"status": "ok"
			}
		}
	}`

	var response ReceiptsResponse
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	receipt, ok := response.Data["019f0039-9cb1-7088-a76b-504af6cea9b7"]
	if !ok {
		t.Fatalf("expected receipt in response data")
	}
	if receipt.Status != SuccessStatus {
		t.Fatalf("receipt.Status = %q, want %q", receipt.Status, SuccessStatus)
	}
	if err := receipt.ValidateReceipt(); err != nil {
		t.Fatalf("ValidateReceipt() error = %v", err)
	}
}

func TestValidateReceiptError(t *testing.T) {
	receipt := PushReceipt{
		Status:  "error",
		Message: "Not registered",
		Details: map[string]string{"error": ErrorDeviceNotRegistered},
	}

	err := receipt.ValidateReceipt()
	if _, ok := err.(*DeviceNotRegisteredReceiptError); !ok {
		t.Fatalf("expected DeviceNotRegisteredReceiptError, got %T", err)
	}
}
