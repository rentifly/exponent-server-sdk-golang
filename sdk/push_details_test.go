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
			raw:  `{"error":"DeveloperError","apns":{"reason":"BadDeviceToken","statusCode":400},"errorCodeEnum":2,"sentAt":1782424694}`,
			want: PushDetails{
				"error":         "DeveloperError",
				"apns":          "BadDeviceToken",
				"errorCodeEnum": "2",
				"sentAt":        "1782424694",
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

func TestGetReceiptsResponseUnmarshalDeveloperError(t *testing.T) {
	raw := `{
		"data": {
			"019f00ca-6e6d-778f-8204-c90038832cc5": {
				"status": "error",
				"message": "The Apple Push Notification service failed to send the notification (reason: BadDeviceToken, status code: 400).",
				"messageEnum": 1006,
				"messageParamValues": ["BadDeviceToken", "400"],
				"details": {
					"apns": {
						"reason": "BadDeviceToken",
						"statusCode": 400
					},
					"error": "DeveloperError",
					"errorCodeEnum": 2,
					"sentAt": 1782424694
				}
			}
		}
	}`

	var response ReceiptsResponse
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	receipt, ok := response.Data["019f00ca-6e6d-778f-8204-c90038832cc5"]
	if !ok {
		t.Fatalf("expected receipt in response data")
	}
	if receipt.Details["error"] != "DeveloperError" {
		t.Fatalf("receipt.Details[error] = %q, want DeveloperError", receipt.Details["error"])
	}
	if receipt.Details["errorCodeEnum"] != "2" {
		t.Fatalf("receipt.Details[errorCodeEnum] = %q, want 2", receipt.Details["errorCodeEnum"])
	}
	if receipt.Details["apns"] != "BadDeviceToken" {
		t.Fatalf("receipt.Details[apns] = %q, want BadDeviceToken", receipt.Details["apns"])
	}
	if receipt.Details["sentAt"] != "1782424694" {
		t.Fatalf("receipt.Details[sentAt] = %q, want 1782424694", receipt.Details["sentAt"])
	}
	if receipt.MessageEnum == nil || *receipt.MessageEnum != 1006 {
		t.Fatalf("receipt.MessageEnum = %v, want 1006", receipt.MessageEnum)
	}
}
