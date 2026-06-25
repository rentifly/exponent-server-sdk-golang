package expo

import (
	"encoding/json"
	"testing"
)

func TestPushDetailsUnmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    PushDetails
		wantErr bool
	}{
		{
			name: "string error",
			raw:  `{"error":"DeviceNotRegistered"}`,
			want: PushDetails{"error": ErrorDeviceNotRegistered},
		},
		{
			name: "nested error object",
			raw:  `{"error":{"code":"DeviceNotRegistered"}}`,
			want: PushDetails{"error": ErrorDeviceNotRegistered},
		},
		{
			name: "fcm details object",
			raw:  `{"error":"DeviceNotRegistered","fcm":{"error":"NotRegistered"}}`,
			want: PushDetails{
				"error": ErrorDeviceNotRegistered,
				"fcm":   "NotRegistered",
			},
		},
		{
			name: "apns details object",
			raw:  `{"error":"DeviceNotRegistered","apns":{"reason":"BadDeviceToken"}}`,
			want: PushDetails{
				"error": ErrorDeviceNotRegistered,
				"apns":  "BadDeviceToken",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var details PushDetails
			err := json.Unmarshal([]byte(tt.raw), &details)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			for key, want := range tt.want {
				if got := details[key]; got != want {
					t.Fatalf("details[%q] = %q, want %q", key, got, want)
				}
			}
		})
	}
}

func TestGetReceiptsResponseUnmarshalWithProviderDetails(t *testing.T) {
	raw := `{
		"data": {
			"019f0039-9cb1-7088-a76b-504af6cea9b7": {
				"status": "error",
				"message": "The recipient device is not registered with FCM.",
				"details": {
					"error": "DeviceNotRegistered",
					"fcm": {
						"error": "NotRegistered"
					}
				}
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
	if receipt.Details["error"] != ErrorDeviceNotRegistered {
		t.Fatalf("receipt.Details[error] = %q, want %q", receipt.Details["error"], ErrorDeviceNotRegistered)
	}
	if receipt.Details["fcm"] != "NotRegistered" {
		t.Fatalf("receipt.Details[fcm] = %q, want %q", receipt.Details["fcm"], "NotRegistered")
	}
}
