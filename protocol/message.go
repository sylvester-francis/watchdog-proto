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
	MsgTypeTaskCancel      = "task_cancel"
	MsgTypeError           = "error"
	MsgTypeUpdateAvailable = "update_available"
	MsgTypeDiscoveryTask   = "discovery_task"
	MsgTypeDiscoveryResult = "discovery_result"
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
	APIKey      string            `json:"api_key"`
	Version     string            `json:"version,omitempty"`
	Fingerprint map[string]string `json:"fingerprint,omitempty"`
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
	MonitorID string            `json:"monitor_id"`
	Type      string            `json:"type"`
	Target    string            `json:"target"`
	Interval  int               `json:"interval"`
	Timeout   int               `json:"timeout"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// HeartbeatPayload is sent by agent with check results.
type HeartbeatPayload struct {
	MonitorID      string            `json:"monitor_id"`
	Status         string            `json:"status"`
	LatencyMs      int               `json:"latency_ms,omitempty"`
	ErrorMessage   string            `json:"error_message,omitempty"`
	CertExpiryDays *int              `json:"cert_expiry_days,omitempty"`
	CertIssuer     string            `json:"cert_issuer,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
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

// UpdateAvailablePayload is sent by hub when a newer agent version exists.
type UpdateAvailablePayload struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
	SHA256      string `json:"sha256"`
	Signature   string `json:"signature,omitempty"`
}

// Helper functions to create common messages.

// NewAuthMessage creates an authentication message.
func NewAuthMessage(apiKey, version string) *Message {
	return MustNewMessage(MsgTypeAuth, AuthPayload{
		APIKey:  apiKey,
		Version: version,
	})
}

// NewAuthMessageWithFingerprint creates an authentication message with device fingerprint.
func NewAuthMessageWithFingerprint(apiKey, version string, fingerprint map[string]string) *Message {
	return MustNewMessage(MsgTypeAuth, AuthPayload{
		APIKey:      apiKey,
		Version:     version,
		Fingerprint: fingerprint,
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

// NewTaskMessageWithMetadata creates a task assignment message with metadata.
func NewTaskMessageWithMetadata(monitorID, monitorType, target string, interval, timeout int, metadata map[string]string) *Message {
	return MustNewMessage(MsgTypeTask, TaskPayload{
		MonitorID: monitorID,
		Type:      monitorType,
		Target:    target,
		Interval:  interval,
		Timeout:   timeout,
		Metadata:  metadata,
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

// NewUpdateAvailableMessage creates an update available message.
func NewUpdateAvailableMessage(version, downloadURL, sha256, signature string) *Message {
	return MustNewMessage(MsgTypeUpdateAvailable, UpdateAvailablePayload{
		Version:     version,
		DownloadURL: downloadURL,
		SHA256:      sha256,
		Signature:   signature,
	})
}

// DiscoveryTaskPayload is sent by hub to request network discovery.
type DiscoveryTaskPayload struct {
	TaskID      string `json:"task_id"`
	Subnet      string `json:"subnet"`
	Community   string `json:"community"`
	SNMPVersion string `json:"snmp_version"`
	Timeout     int    `json:"timeout"`
}

// DiscoveryResultPayload is sent by agent with discovery results.
type DiscoveryResultPayload struct {
	TaskID   string             `json:"task_id"`
	Status   string             `json:"status"`
	Progress int                `json:"progress"`
	Devices  []DiscoveredDevice `json:"devices,omitempty"`
	Error    string             `json:"error,omitempty"`
}

// DiscoveredDevice describes a device found during network discovery.
type DiscoveredDevice struct {
	IP            string `json:"ip"`
	Hostname      string `json:"hostname,omitempty"`
	SysDescr      string `json:"sys_descr,omitempty"`
	SysObjectID   string `json:"sys_object_id,omitempty"`
	SysName       string `json:"sys_name,omitempty"`
	SNMPReachable bool   `json:"snmp_reachable"`
	PingReachable bool   `json:"ping_reachable"`
	TemplateID    string `json:"template_id,omitempty"`
}

// NewDiscoveryTaskMessage creates a discovery task message.
func NewDiscoveryTaskMessage(taskID, subnet, community, snmpVersion string, timeout int) *Message {
	return MustNewMessage(MsgTypeDiscoveryTask, DiscoveryTaskPayload{
		TaskID:      taskID,
		Subnet:      subnet,
		Community:   community,
		SNMPVersion: snmpVersion,
		Timeout:     timeout,
	})
}

// NewDiscoveryResultMessage creates a discovery result message.
func NewDiscoveryResultMessage(taskID, status string, progress int, devices []DiscoveredDevice, errMsg string) *Message {
	return MustNewMessage(MsgTypeDiscoveryResult, DiscoveryResultPayload{
		TaskID:   taskID,
		Status:   status,
		Progress: progress,
		Devices:  devices,
		Error:    errMsg,
	})
}
