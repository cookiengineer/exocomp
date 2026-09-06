
import { Question } from "/types/Question.mjs";

export const GetUnansweredQuestions = (agent) => {

	let questions = {};

	agent.Messages.forEach((message) => {

		if (message.Role === "assistant") {

			message.ToolCalls.forEach((tool_call) => {

				if (tool_call.Function.Name === "humans.Ask") {
					questions[tool_call.ID] = Question.from(tool_call.Function.Arguments);
				} else if (tool_call.Function.Name === "humans.Choose") {
					questions[tool_call.ID] = Question.from(tool_call.Function.Arguments);
				}

			});

		} else if (message.Role === "tool") {

			if (message.ToolName === "humans.Ask") {

				let question = questions[message.ToolCallID] || null;
				if (question !== null && question.Type === "Ask") {

					let content = (message.Content || "").trim();
					if (content.includes("Answer for Question") && content.includes("===")) {
						question.Answer = content.split("===")[1];
					}

				}

			} else if (message.ToolName === "humans.Choose") {

				let question = questions[message.ToolCallID] || null;
				if (question !== null && question.Type === "Choose") {

					let content = (message.Content || "").trim();
					if (content.includes("Answer for Question") && content.includes("===")) {
						question.Answer = content.split("===")[1];
					}

				}

			}

		}

	});

	let unanswered = [];

	Object.keys(questions).forEach((tool_call_id) => {

		let question = questions[tool_call_id];
		if (question.Answer === "") {
			unanswered.push(question);
		}

	});

	return unanswered;

};
