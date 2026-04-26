
export const Agent = function() {

	this.Name        = "";
	this.Type        = "";
	this.Model       = "";
	this.Prompt      = "";
	this.Temperature = 0.0;
	this.Messages    = [];
	this.Programs    = [];
	this.Tools       = [];
	this.Sandbox     = "";

};

Agent.from = (data) => {

	let agent = new Agent();

	agent.Name        = data["name"]        || "";
	agent.Type        = data["type"]        || "";
	agent.Model       = data["model"]       || "";
	agent.Prompt      = data["prompt"]      || "";
	agent.Temperature = data["temperature"] || 0.0;
	agent.Messages    = data["messages"]    || [];
	agent.Programs    = data["programs"]    || [];
	agent.Tools       = data["tools"]       || [];
	agent.Sandbox     = data["sandbox"]     || "";

	return agent;

};
