# archon-event-envelope

[![Go Reference](https://pkg.go.dev/badge/github.com/diogoX451/archon-event-envelope.svg)](https://pkg.go.dev/github.com/diogoX451/archon-event-envelope)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Versioned event envelopes with **dual-read** for legacy bare JSON.

Part of the **Archon open-source toolkit** by [@diogoX451](https://github.com/diogoX451).

## Install

```bash
go get github.com/diogoX451/archon-event-envelope@latest
```

## Usage

```go
import "github.com/diogoX451/archon-event-envelope"

raw, _ := envelope.MarshalV1("order.created", "tenant-1", "wf-1", "corr-1", map[string]any{
    "id": "42",
})

env, bare, err := envelope.UnmarshalFlexible(raw)
// bare == false for versioned envelopes; true for legacy payloads
```

## Why

Agent platforms and microservices evolve wire formats. Dual-read lets you ship
`schema_version` without breaking older publishers.

## Related Archon packages

| Package | Role |
|---------|------|
| [archon-nats-bus](https://github.com/diogoX451/archon-nats-bus) | JetStream bus with resilient resubscribe |
| [archon-executor-runtime](https://github.com/diogoX451/archon-executor-runtime) | Need consumer SDK |
| [archon-need-protocol](https://github.com/diogoX451/archon-need-protocol) | Need/response contracts |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache-2.0 — see [LICENSE](LICENSE).
