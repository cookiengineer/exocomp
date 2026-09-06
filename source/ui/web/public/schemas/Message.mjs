
import { ToolCall } from "/schemas/ToolCall.mjs";

export const Message = function() {

	this.Role       = "";
	this.Content    = "";
	this.Created    = new Date();
	this.ToolCallID = "";
	this.ToolName   = "";
	this.ToolCalls  = [];

};

Message.from = (data) => {

	let message = new Message();

	message.Role       = data["role"]    || "";
	message.Content    = data["content"] || "";
	message.Created    = new Date((data["created"] || "").replace(" ", "T"));
	message.ToolCallID = data["tool_call_id"] || "";
	message.ToolName   = data["tool_name"]    || "";

	if (Array.isArray(data["tool_calls"]) === true) {

		message.ToolCalls = data["tool_calls"].map((raw) => {
			return ToolCall.from(raw);
		});

	}

	return message;

};
