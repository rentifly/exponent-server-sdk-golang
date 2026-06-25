package expo

import "testing"

func TestCountRecipients(t *testing.T) {
	messages := []PushMessage{
		{To: []ExponentPushToken{"ExponentPushToken[a]", "ExponentPushToken[b]"}},
		{To: []ExponentPushToken{"ExponentPushToken[c]"}},
	}

	if got := countRecipients(messages); got != 3 {
		t.Fatalf("countRecipients() = %d, want 3", got)
	}
}

func TestAttachPushMessages(t *testing.T) {
	messages := []PushMessage{
		{
			To:    []ExponentPushToken{"ExponentPushToken[a]", "ExponentPushToken[b]"},
			Title: "hello",
		},
		{
			To:    []ExponentPushToken{"ExponentPushToken[c]"},
			Title: "world",
		},
	}
	responses := make([]PushResponse, 3)

	attachPushMessages(messages, responses)

	if responses[0].PushMessage.Title != "hello" || responses[1].PushMessage.Title != "hello" {
		t.Fatalf("expected first two responses to reference first message")
	}
	if responses[2].PushMessage.Title != "world" {
		t.Fatalf("expected third response to reference second message")
	}
}
