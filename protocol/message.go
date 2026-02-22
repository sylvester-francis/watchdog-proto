package protocol

import (
	"encoding/json"
	"time"
)

// Message types for WebSocket communication.
const (
	MsgTypeAuth      = "auth"
	MsgTypeAuthAck   = "auth_ack"
	MsgTypeAuthError = "auth_error"
	MsgTypeTask      = "task"
	MsgTypeHeartbeat = "heartbeat"
	MsgTypePing      = "ping"
	MsgTypePong      = "pong"
	MsgTypeTaskCancel = "task_cancel"
	MsgTypeError      = "error"
)

// Message represents a WebSocket message envelope.
type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// NewMessage creates a new message with the current timestamp.
func NewMessage(msgType string, payload any) (*Message, error) {
	var rawPayload json.RawMessage
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		rawPayload = data
	}

	return &Message{
		Type:      msgType,
		Payload:   rawPayload,
		Timestamp: time.Now(),
	}, nil
}

// ParsePayload unmarshals the payload into the provided type.
func (m *Message) ParsePayload(v any) error {
	if m.Payload == nil {
		return nil
	}
	return json.Unmarshal(m.Payload, v)
}

// MustNewMessage creates a new message and panics on error.
// Use only when payload is guaranteed to be serializable.
func MustNewMessage(msgType string, payload any) *Message {
	msg, err := NewMessage(msgType, payload)
	if err != nil {
		panic(err)
	}
	return msg
}

// AuthPayload is sent by agent to authenticate.
type AuthPayload struct {
	APIKey  string `json:"api_key"`
	Version string `json:"version,omitempty"`
}

// AuthAckPayload is sent by hub to confirm authentication.
type AuthAckPayload struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
}

// AuthErrorPayload is sent by hub when authentication fails.
type AuthErrorPayload struct {
	Error string `json:"error"`
}

// TaskPayload describes a monitoring task for the agent.
type TaskPayload struct {
	MonitorID string `json:"monitor_id"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	Interval  int    `json:"interval"`
	Timeout   int    `json:"timeout"`
}

// HeartbeatPayload is sent by agent with check results.
type HeartbeatPayload struct {
	MonitorID      string `json:"monitor_id"`
	Status         string `json:"status"`
	LatencyMs      int    `json:"latency_ms,omitempty"`
	ErrorMessage   string `json:"error_message,omitempty"`
	CertExpiryDays *int   `json:"cert_expiry_days,omitempty"`
	CertIssuer     string `json:"cert_issuer,omitempty"`
}

// TaskCancelPayload tells the agent to stop monitoring a specific monitor.
type TaskCancelPayload struct {
	MonitorID string `json:"monitor_id"`
}

// ErrorPayload is sent when an error occurs.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Helper functions to create common messages.

// NewAuthMessage creates an authentication message.
func NewAuthMessage(apiKey, version string) *Message {
	return MustNewMessage(MsgTypeAuth, AuthPayload{
		APIKey:  apiKey,
		Version: version,
	})
}

// NewAuthAckMessage creates an authentication acknowledgment message.
func NewAuthAckMessage(agentID, agentName string) *Message {
	return MustNewMessage(MsgTypeAuthAck, AuthAckPayload{
		AgentID:   agentID,
		AgentName: agentName,
	})
}

// NewAuthErrorMessage creates an authentication error message.
func NewAuthErrorMessage(err string) *Message {
	return MustNewMessage(MsgTypeAuthError, AuthErrorPayload{
		Error: err,
	})
}

// NewTaskMessage creates a task assignment message.
func NewTaskMessage(monitorID, monitorType, target string, interval, timeout int) *Message {
	return MustNewMessage(MsgTypeTask, TaskPayload{
		MonitorID: monitorID,
		Type:      monitorType,
		Target:    target,
		Interval:  interval,
		Timeout:   timeout,
	})
}

// NewTaskCancelMessage creates a task cancellation message.
func NewTaskCancelMessage(monitorID string) *Message {
	return MustNewMessage(MsgTypeTaskCancel, TaskCancelPayload{
		MonitorID: monitorID,
	})
}

// NewHeartbeatMessage creates a heartbeat message.
func NewHeartbeatMessage(monitorID, status string, latencyMs int, errorMsg string) *Message {
	return MustNewMessage(MsgTypeHeartbeat, HeartbeatPayload{
		MonitorID:    monitorID,
		Status:       status,
		LatencyMs:    latencyMs,
		ErrorMessage: errorMsg,
	})
}

// NewPingMessage creates a ping message.
func NewPingMessage() *Message {
	return MustNewMessage(MsgTypePing, nil)
}

// NewPongMessage creates a pong message.
func NewPongMessage() *Message {
	return MustNewMessage(MsgTypePong, nil)
}

// NewErrorMessage creates an error message.
func NewErrorMessage(code, message string) *Message {
	return MustNewMessage(MsgTypeError, ErrorPayload{
		Code:    code,
		Message: message,
	})
}
