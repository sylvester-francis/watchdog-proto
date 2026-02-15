# watchdog-proto

Shared WebSocket message protocol for the [WatchDog](https://github.com/sylvester-francis/watchdog) monitoring system.

This package defines the JSON message envelope, payload types, and helper constructors used by both the [WatchDog Hub](https://github.com/sylvester-francis/watchdog) and the [WatchDog Agent](https://github.com/sylvester-francis/watchdog-agent) to communicate over WebSocket connections.

## Installation

```bash
go get github.com/sylvester-francis/watchdog-proto@latest
```

## Usage

```go
import "github.com/sylvester-francis/watchdog-proto/protocol"

// Create an authentication message (agent -> hub)
msg := protocol.NewAuthMessage("my-api-key", "1.0.0")

// Create a heartbeat message (agent -> hub)
msg := protocol.NewHeartbeatMessage("monitor-uuid", "up", 42, "")

// Create a task assignment message (hub -> agent)
msg := protocol.NewTaskMessage("monitor-uuid", "http", "https://example.com", 30, 10)

// Parse an incoming message payload
var payload protocol.HeartbeatPayload
if err := msg.ParsePayload(&payload); err != nil {
    log.Fatal(err)
}
```

## Message Envelope

Every message sent over the WebSocket connection uses the same envelope format:

```json
{
  "type": "heartbeat",
  "payload": { ... },
  "timestamp": "2025-01-15T10:30:00Z"
}
```

The `Message` struct wraps the type, a raw JSON payload, and a UTC timestamp:

```go
type Message struct {
    Type      string          `json:"type"`
    Payload   json.RawMessage `json:"payload,omitempty"`
    Timestamp time.Time       `json:"timestamp"`
}
```

## Message Types

| Type         | Direction       | Description                          |
|--------------|-----------------|--------------------------------------|
| `auth`       | Agent -> Hub    | Agent sends API key to authenticate  |
| `auth_ack`   | Hub -> Agent    | Hub confirms successful authentication |
| `auth_error` | Hub -> Agent    | Hub rejects authentication           |
| `task`       | Hub -> Agent    | Hub assigns a monitoring task        |
| `heartbeat`  | Agent -> Hub    | Agent reports check results          |
| `ping`       | Hub -> Agent    | Hub checks agent liveness            |
| `pong`       | Agent -> Hub    | Agent responds to ping               |
| `error`      | Either          | Generic error message                |

## Payload Types

### AuthPayload

Sent by the agent during the initial handshake.

```go
type AuthPayload struct {
    APIKey  string `json:"api_key"`
    Version string `json:"version,omitempty"`
}
```

### AuthAckPayload

Sent by the hub to confirm authentication and provide the agent's identity.

```go
type AuthAckPayload struct {
    AgentID   string `json:"agent_id"`
    AgentName string `json:"agent_name"`
}
```

### AuthErrorPayload

Sent by the hub when authentication fails.

```go
type AuthErrorPayload struct {
    Error string `json:"error"`
}
```

### TaskPayload

Sent by the hub to assign a monitoring check to the agent.

```go
type TaskPayload struct {
    MonitorID string `json:"monitor_id"`
    Type      string `json:"type"`       // "http", "tcp", "ping", "dns"
    Target    string `json:"target"`     // URL, host:port, or hostname
    Interval  int    `json:"interval"`   // Check interval in seconds
    Timeout   int    `json:"timeout"`    // Check timeout in seconds
}
```

### HeartbeatPayload

Sent by the agent after each monitoring check completes.

```go
type HeartbeatPayload struct {
    MonitorID    string `json:"monitor_id"`
    Status       string `json:"status"`                    // "up", "down", "timeout", "error"
    LatencyMs    int    `json:"latency_ms,omitempty"`
    ErrorMessage string `json:"error_message,omitempty"`
}
```

### ErrorPayload

Generic error message for non-auth error conditions.

```go
type ErrorPayload struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

## Helper Constructors

The package provides constructors that handle JSON serialization and timestamp assignment:

| Function                | Creates              |
|-------------------------|----------------------|
| `NewAuthMessage`        | `auth` message       |
| `NewAuthAckMessage`     | `auth_ack` message   |
| `NewAuthErrorMessage`   | `auth_error` message |
| `NewTaskMessage`        | `task` message       |
| `NewHeartbeatMessage`   | `heartbeat` message  |
| `NewPingMessage`        | `ping` message       |
| `NewPongMessage`        | `pong` message       |
| `NewErrorMessage`       | `error` message      |
| `NewMessage`            | Any type (returns error) |
| `MustNewMessage`        | Any type (panics on error) |

## Connection Lifecycle

A typical WebSocket session follows this sequence:

```
Agent                                Hub
  |                                   |
  |-- auth {api_key, version} ------->|
  |                                   |  (validate API key)
  |<------ auth_ack {agent_id} -------|
  |                                   |
  |<------ task {monitor, target} ----|  (one per enabled monitor)
  |<------ task {monitor, target} ----|
  |                                   |
  |-- heartbeat {status, latency} --->|  (after each check)
  |-- heartbeat {status, latency} --->|
  |                                   |
  |<------ ping ----------------------|  (periodic liveness check)
  |-- pong --------------------------->|
  |                                   |
```

## Dependencies

None. This package uses only the Go standard library (`encoding/json`, `time`).

## Related Repositories

| Repository | Description |
|------------|-------------|
| [watchdog](https://github.com/sylvester-francis/watchdog) | Hub server -- dashboard, API, alerting, and data storage |
| [watchdog-agent](https://github.com/sylvester-francis/watchdog-agent) | Monitoring agent binary deployed inside customer networks |

## License

MIT License. See [LICENSE](LICENSE) for details.
