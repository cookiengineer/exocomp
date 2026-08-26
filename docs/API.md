
# API

The web frontend (`exocomp web`) exposes a REST API on port `3000`. All routes
serve the planner's [Session](../source/types/Session.go) state.

## Parameters

| Verb   | Route                    | Response Schema                            |
|:------:|:-------------------------|:-------------------------------------------|
| `GET`  | `/api/parameters/roles`  | `[]string` (available agent roles)         |
| `GET`  | `/api/parameters/models` | `[]string` (available models via `/v1/models`) |

## Session

| Verb   | Route                        | Request Schema   | Response Schema                                      |
|:------:|:-----------------------------|:-----------------|:-----------------------------------------------------|
| `GET`  | `/api/session/agent`         |                  | [types.Agent](../source/types/Agent.go)              |
| `GET`  | `/api/session/agents`        |                  | `[]*types.Agent` (planner + hired agents)            |
| `GET`  | `/api/session/config`        |                  | [types.Config](../source/types/Config.go)            |
| `GET`  | `/api/session/config/{name}` |                  | [types.Config](../source/types/Config.go)            |
| `GET`  | `/api/session/console`       |                  | `[]types.ConsoleMessage`                             |
| `GET`  | `/api/session/tools`         |                  | `[]schemas.Tool`                                     |

## Interaction

| Verb   | Route                        | Request Schema                                     | Response Schema          |
|:------:|:-----------------------------|:---------------------------------------------------|:-------------------------|
| `POST` | `/api/session/calltool`      | [schemas.ToolCall](../source/schemas/ToolCall.go)  | `schemas.Message`        |
| `POST` | `/api/session/sendchatrequest` | [schemas.Message](../source/schemas/Message.go)  | `[]schemas.Message`      |
