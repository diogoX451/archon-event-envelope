// Package envelope provides versioned event wrappers with dual-read support
// for legacy bare JSON payloads.
//
// Extracted from the Archon agent platform (https://github.com/diogoX451).
package envelope

import (
	"encoding/json"
	"fmt"
	"time"
)

// CurrentSchemaVersion is the envelope version this package understands fully.
const CurrentSchemaVersion = 1

// Envelope is the versioned wrapper for cross-process events and needs.
// Dual-read consumers: try UnmarshalFlexible; if no envelope fields, treat
// raw JSON as the payload itself (legacy).
type Envelope struct {
	SchemaVersion int             `json:"schema_version"`
	Type          string          `json:"type"`
	TenantID      string          `json:"tenant_id,omitempty"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	WorkflowID    string          `json:"workflow_id,omitempty"`
	OccurredAt    time.Time       `json:"occurred_at"`
	Payload       json.RawMessage `json:"payload"`
}

// MarshalV1 wraps payload as a schema_version=1 envelope.
func MarshalV1(eventType, tenantID, workflowID, correlationID string, payload any) ([]byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	env := Envelope{
		SchemaVersion: CurrentSchemaVersion,
		Type:          eventType,
		TenantID:      tenantID,
		CorrelationID: correlationID,
		WorkflowID:    workflowID,
		OccurredAt:    time.Now().UTC(),
		Payload:       raw,
	}
	return json.Marshal(env)
}

// UnmarshalFlexible decodes either a versioned envelope or legacy bare payload.
// When bare, schemaVersion is 0 and eventType is empty.
func UnmarshalFlexible(data []byte) (env Envelope, bare bool, err error) {
	if len(data) == 0 {
		return Envelope{}, true, fmt.Errorf("empty payload")
	}
	var probe struct {
		SchemaVersion int             `json:"schema_version"`
		Type          string          `json:"type"`
		Payload       json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return Envelope{}, true, err
	}
	if probe.SchemaVersion > 0 && len(probe.Payload) > 0 {
		if err := json.Unmarshal(data, &env); err != nil {
			return Envelope{}, false, err
		}
		return env, false, nil
	}
	env = Envelope{
		SchemaVersion: 0,
		Payload:       append(json.RawMessage(nil), data...),
	}
	return env, true, nil
}

// DecodePayload unmarshals env.Payload into dest.
func DecodePayload(env Envelope, dest any) error {
	if len(env.Payload) == 0 {
		return fmt.Errorf("empty envelope payload")
	}
	return json.Unmarshal(env.Payload, dest)
}

// SupportedSchema reports whether this binary can process the version.
func SupportedSchema(version int) bool {
	return version == 0 || version == CurrentSchemaVersion
}
