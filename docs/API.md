
# API

## Status and Messages

| Verb   | Route                     | Request Schema          | Response Schema          |
|:------:|:--------------------------|:------------------------|:-------------------------|
| `GET`  | `/api/console`            |                         | `[]types.ConsoleMessage` |
| `GET`  | `/api/messages`           |                         | `[]schemas.Message`      |
| `POST` | `/api/messages/send`      | `schemas.Message`       | `[]schemas.Message`      |

## Parameters and Settings

| Verb   | Route                     | Request Schema          | Response Schema          |
|:------:|:--------------------------|:------------------------|:-------------------------|
| `GET`  | `/api/parameters/agents`  |                         | `[]string`               |
| `GET`  | `/api/parameters/models`  |                         | `[]string`               |
| `POST` | `/api/settings/agent`     | `schemas.AgentSettings` |                          |

