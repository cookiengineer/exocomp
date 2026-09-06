
export const ToolCall = function() {

	this.ID       = "";
	this.Type     = "";
	this.Function = {
		Name:      "",
		Arguments: "",
	};

};

ToolCall.from = (data) => {

	let tool_call = new ToolCall();

	tool_call.ID   = data["id"]   || "";
	tool_call.Type = data["type"] || "function";

	let fn = data["function"] || null;

	if (fn !== null) {

		tool_call.Function.Name      = (fn["name"]      || "");
		tool_call.Function.Arguments = (fn["arguments"] || null);

	} else {

		tool_call.Function.Name      = "";
		tool_call.Function.Arguments = null;

	}

	return tool_call;

};
