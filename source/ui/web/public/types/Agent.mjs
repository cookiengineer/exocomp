
export const Agent = function() {

	this.Name         = "";
	this.Role         = "";
	this.Model        = "";
	this.Prompt       = "";
	this.Temperature  = 0.0;
	this.Messages     = [];
	this.Programs     = [];
	this.Tools        = [];
	this.Sandbox      = "";
	this.ContextUsage = {
		Length: 0,
		Tokens: 0,
	};

};

Agent.from = (data) => {

	let agent = new Agent();

	agent.Name        = data["name"]        || "";
	agent.Role        = data["role"]        || "";
	agent.Model       = data["model"]       || "";
	agent.Prompt      = data["prompt"]      || "";
	agent.Temperature = data["temperature"] || 0.0;
	agent.Messages    = data["messages"]    || [];
	agent.Programs    = data["programs"]    || [];
	agent.Tools       = data["tools"]       || [];
	agent.Sandbox     = data["sandbox"]     || "";

	if (Object.prototype.toString.call(data["context-usage"]) === "[object Object]") {
		agent.ContextUsage.Length = data["context-usage"]["length"] || 0;
		agent.ContextUsage.Tokens = data["context-usage"]["tokens"] || 0;
	}

	return agent;

};
