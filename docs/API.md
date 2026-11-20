# API

The web frontend (`exocomp-web web`) exposes a REST API on port `3000`. The
routes expose Session, tool and parameter state.

Successful responses use `application/json`.

## Parameters

| Verb | Route | Response |
|:-----|:------|:---------|
| `GET` | `/api/parameters/roles` | `[]string` (available agent roles) |
| `GET` | `/api/parameters/models` | `[]string` (models from `/v1/models`) |

## Session

| Verb | Route | Request | Response |
|:-----|:------|:--------|:---------|
| `GET` | `/api/session/agent` | optional `?name=` | [types.Agent](../source/types/Agent.go); empty `name` selects the planner |
| `GET` | `/api/session/agents` | | `[]*types.Agent` (planner and hired agents) |
| `GET` | `/api/session/bugs` | | `map[path]map[symbol]`[types.Bug](../source/types/Bug.go) |
| `GET` | `/api/session/changelog` | | `map[path]map[symbol][]`[types.ChangelogEntry](../source/types/ChangelogEntry.go) |
| `GET` | `/api/session/requirements` | | `map[path]map[symbol]`[types.Requirement](../source/types/Requirement.go) |
| `GET` | `/api/session/findings` | | `[]`[types.FindingSummary](../source/types/FindingSummary.go) |
| `GET` | `/api/session/config` | | [types.Config](../source/types/Config.go) |
| `GET` | `/api/session/config/{name}` | path `{name}` | [types.Config](../source/types/Config.go) for the planner or a hired agent |
| `GET` | `/api/session/console` | | `[]types.ConsoleMessage` |
| `GET` | `/api/session/tools` | | `[]schemas.Tool` |

## Interaction

| Verb | Route | Request | Response |
|:-----|:------|:--------|:---------|
| `POST` | `/api/session/calltool` | [schemas.ToolCall](../source/schemas/ToolCall.go) | `schemas.Message` (the resulting tool message) |
| `POST` | `/api/session/sendchatrequest` | [schemas.Message](../source/schemas/Message.go) | `[]*schemas.Message` (messages appended by the request) |

## Humans Tool Interaction

`humans.Ask`, `humans.Choose` and `humans.AskPermission` block until answered.
The web client detects unanswered questions with
[`GetUnansweredQuestions`](../app/utils/chat/GetUnansweredQuestions.go), renders
them in the [Chat controller](../app/controllers/Chat.go) and posts one
`humans.Answer` tool call per answer to `/api/session/calltool`.

The server matches the answer to the question text and appends the result with
`role="tool"`. The request itself does not carry a tool-call ID:

```json
{
    "type": "function",
    "function": {
        "name": "humans.Answer",
        "arguments": {
            "question": "<original_question>",
            "answer": "<selected_choice or free-form text>"
        }
    }
}
```

## Notes

`source/ui/web/routes/session/Messages.go` defines a message handler, but
`Server.Listen` does not register it. The web client reads messages from
`/api/session/agents`.
