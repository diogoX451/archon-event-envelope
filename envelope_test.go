package envelope

import (
	"encoding/json"
	"testing"
)

func TestMarshalV1AndUnmarshalFlexible(t *testing.T) {
	t.Parallel()
	raw, err := MarshalV1("demo.event", "tenant-a", "wf-1", "c-1", map[string]string{"message": "hi"})
	if err != nil {
		t.Fatal(err)
	}
	env, bare, err := UnmarshalFlexible(raw)
	if err != nil {
		t.Fatal(err)
	}
	if bare {
		t.Fatal("expected envelope")
	}
	if env.SchemaVersion != CurrentSchemaVersion || env.Type != "demo.event" {
		t.Fatalf("env = %+v", env)
	}
	var payload map[string]string
	if err := DecodePayload(env, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["message"] != "hi" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestUnmarshalFlexibleLegacy(t *testing.T) {
	t.Parallel()
	env, bare, err := UnmarshalFlexible([]byte(`{"workflow_id":"wf-1","status":"completed"}`))
	if err != nil || !bare || env.SchemaVersion != 0 {
		t.Fatalf("bare=%v schema=%d err=%v", bare, env.SchemaVersion, err)
	}
}

func TestSupportedSchema(t *testing.T) {
	t.Parallel()
	if !SupportedSchema(0) || !SupportedSchema(1) || SupportedSchema(99) {
		t.Fatal("schema support mismatch")
	}
	b, _ := json.Marshal(Envelope{SchemaVersion: 1, Type: "t", Payload: json.RawMessage(`{}`)})
	if !json.Valid(b) {
		t.Fatal("invalid json")
	}
}
