
# API

The web frontend (`exocomp web`) exposes a REST API on port `3000`. All routes
serve the planner's [Session](../source/engine/Session.go) state.

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
| `GET`  | `/api/session/bugs`          |                  | `[]`[types.Bug](../source/types/Bug.go) (sorted by file, then symbol)            |
| `GET`  | `/api/session/changelog`     |                  | `[]`[types.ChangelogEntry](../source/types/ChangelogEntry.go) (sorted by file, symbol, date) |
| `GET`  | `/api/session/requirements`  |                  | `[]`[types.Requirement](../source/types/Requirement.go) (sorted by file, then symbol) |
| `GET`  | `/api/session/config`        |                  | [types.Config](../source/types/Config.go)            |
| `GET`  | `/api/session/config/{name}` |                  | [types.Config](../source/types/Config.go)            |
| `GET`  | `/api/session/console`       |                  | `[]types.ConsoleMessage`                             |
| `GET`  | `/api/session/tools`         |                  | `[]schemas.Tool`                                     |

## Interaction

| Verb   | Route                        | Request Schema                                     | Response Schema          |
|:------:|:-----------------------------|:---------------------------------------------------|:-------------------------|
| `POST` | `/api/session/calltool`      | [schemas.ToolCall](../source/schemas/ToolCall.go)  | `schemas.Message`        |
| `POST` | `/api/session/sendchatrequest` | [schemas.Message](../source/schemas/Message.go)  | `[]schemas.Message`      |

## Humans Tool Interaction

The `/api/session/calltool` route is also used to answer blocking `humans.Ask` / `humans.Choose`
questions. The frontend posts a `humans.Answer` tool call with the same tool call id as the
previous question's tool call id. Then the `engine/Session.go` on the server-side will correctly
handle the appendix of the `schemas.Message` with `role="tool"`.

### Example Frontend Humans Tool Call:

1. The Chat History of the currently viewed `Session.Agent` contains as its latest message with
`role="assistant"` two `humans` tool calls that the human user has to respond to. The first
question is an expecteda free-form text answer, the second is a multiple choice question with
only one allowed answer (`multiple=false`).

```javascript
// Session.Agents[Session.Agent].Messages[...]
{
    "role": "assistant",
    "content": "Let's ask the user some geographical trivia questions",
    "created": "2026-01-01T08:00:00Z",
    "tool_call_id": "", // not a tool call response, MUST be empty
    "tool_name":    "", // not a tool call response, MUST be empty
    "tool_calls": [{
        "id": "tool_call_id_of_first_question",
        "type": "function",
        "function": {
            "name": "humans.Ask",
            "arguments": {
                "question": "How old are you? Enter your age."
            }
        }
    }, {
        "id": "tool_call_id_of_second_question",
        "type": "function",
        "function": {
            "name": "humans.Choose",
            "arguments": {
                "question": "What is the mascot animal of the Linux Kernel?",
                "options": [
                    "A gnu",
                    "A penguin",
                    "A turtle",
                    "A whale"
                ],
                "multiple": false
            }
        }
    }]
}
```

2. The [Web UI Client](../source/ui/web/public/ui/Client.mjs) uses `GetUnansweredQuestions()`
to determine whether a Dialog has to be displayed for (still) unanswered questions. Then it
will call `OnQuestions(questions)` and pause the client's `UpdateAgents()` interval.

3. The [Chat View](../source/ui/web/public/views/chat.mjs) will call the [AnswerQuestions](../source/ui/web/public/ui/dialogs/AnswerQuestions.mjs)
dialog with the `Show(questions)` parameter to display the questions in a step-by-step
wizard that the human must answer. The in-between UI steps call `dialog.OnNext(answered_question)`.

The final UI step calls `dialog.OnConfirm(answered_question)`, and reset/hide the dialog,
and call `client.Resume()` which will restart the `UpdateAgents()` interval.

Each `OnNext(question_with_answer)` and `OnConfirm(question_with_answer)` will result in
a tool call that looks like the following:

```javascript
// POST to /api/session/calltool
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

Note that this doesn't have a referencing tool call id, because that is done server-side
in the [engine/Session](../source/engine/Session.go) which automatically matches answers
to their corresponding questions.

