package expo

import "encoding/json"

// PushDetails contains optional fields from Expo push ticket/receipt details.
// Expo may return string values, numbers, or nested objects (fcm, apns).
type PushDetails map[string]string

func (d *PushDetails) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = nil
		return nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	out := make(map[string]string, len(raw))
	for key, value := range raw {
		if decoded, ok := decodeDetailsValue(value); ok {
			out[key] = decoded
		}
	}

	*d = out
	return nil
}

func decodeDetailsValue(raw json.RawMessage) (string, bool) {
	var value string
	if err := json.Unmarshal(raw, &value); err == nil {
		return value, true
	}

	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return number.String(), true
	}

	var nested map[string]json.RawMessage
	if err := json.Unmarshal(raw, &nested); err != nil {
		return "", false
	}

	for _, key := range []string{"error", "code", "reason", "message"} {
		if nestedValue, ok := nested[key]; ok {
			if decoded, ok := decodeDetailsValue(nestedValue); ok {
				return decoded, true
			}
		}
	}

	return "", false
}
