
# API

## Session

| Verb   | Route                          | Request Schema          | Response Schema          |
|:------:|:-------------------------------|:------------------------|:-------------------------|
| `GET`  | `/api/session/agents`          |                         | [[]schemas.Agent](../schemas/Agent.go)               |
| `GET`  | `/api/session/console`         |                         | [[]types.ConsoleMessage](../types/ConsoleMessage.go) |
| `GET`  | `/api/session/context`         |                         | [types.SessionContext](../types/SessionContext.go)   |
| `POST` | `/api/session/sendchatrequest` | `schemas.Message`       | [[]schemas.Message](../schemas/Message.go)           |

## Parameters and Settings

| Verb   | Route                     | Request Schema          | Response Schema          |
|:------:|:--------------------------|:------------------------|:-------------------------|
| `GET`  | `/api/parameters/agents`  |                         | `[]string`               |
| `GET`  | `/api/parameters/models`  |                         | `[]string`               |
| `POST` | `/api/settings/agent`     | `schemas.AgentSettings` |                          |

