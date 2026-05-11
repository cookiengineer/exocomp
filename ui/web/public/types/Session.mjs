
import { Agent   } from "./Agent.mjs";
import { Console } from "./Console.mjs";

export const Session = function(config) {

	// Managed by ui.Client
	this.Agent  = null;
	this.Agents = {};

	this.Config  = config;
	this.Console = new Console();
	this.Context = {
		Length: 0,
		Tokens: 0
	};
	this.Tools   = [];
	this.Waiting = false;

};

Session.prototype = {

	CallTool: async function(name, method, args) {

		if (this.Waiting === false) {

			this.Waiting = true;

			try {

				let response = await fetch(this.Config.ResolveAPI("/api/session/calltool").toString(), {
					method:  "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify({
						"type":     "function",
						"function": {
							"name":      name,
							"arguments": args
						}
					})
				});

				if (response.ok === true) {

					this.Waiting = false;

					return true;

				} else {

					this.Waiting = false;

					return false;

				}

			} catch (err) {

				this.Waiting = false;

				return false;

			}

		} else {

			return false;

		}

	},

	GetAgent: function(name) {

		name = typeof name === "string" ? name : null;

		if (name === null) {
			return this.Agent;
		} else if (this.Agents[name] !== undefined) {
			return this.Agents[name];
		}

		return null;

	},

	GetAgents: function() {
		return this.Agents;
	},

	GetConsoleMessages: function(from) {

		if (this.Console !== null) {
			return this.Console.GetMessages(from);
		} else {
			return [];
		}

	},

	GetContextUsage: function() {

		if (this.Context.Length > 0) {
			return (this.Context.Tokens / this.Context.Length) * 100.0;
		} else {
			return 0.0;
		}

	},

	GetMessages: function(from) {

		from = typeof from === "number" ? from : 0;


		let result = [];

		if (this.Agent !== null) {

			if (this.Agent.Messages.length > 0 && from < this.Agent.Messages.length) {

				for (let m = from; m < this.Agent.Messages.length; m++) {
					result.push(this.Agent.Messages[m]);
				}

			}

		}

		return result;

	},

	GetToolNames: function() {

		let result = [];

		this.Tools.forEach((tool) => {
			result.push(tool["function"]["name"]);
		});

		return result;

	},

	GetToolSchema: function(name) {

		let found = null;

		for (let t = 0; t < this.Tools.length; t++) {

			let tool = this.Tools[t];

			if (tool["function"]["name"] === name) {
				found = tool;
				break;
			}

		}

		return found;

	},

	Init: function() {

		this.Update();

	},

	ReceiveAgent: function(agent) {

		agent = agent instanceof Agent ? agent : null;

		if (agent !== null) {

			this.Agents[agent.Name] = agent;

			if (agent.Name === this.Config.Name) {
				this.Agent = this.Agents[agent.Name];
			}

			return true;

		}

		return false;

	},

	SendChatRequest: async function(message) {

		if (this.Waiting === false) {

			this.Waiting = true;

			try {

				let response = await fetch(this.Config.ResolveAPI("/api/session/sendchatrequest").toString(), {
					method:  "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(message)
				});

				if (response.ok === true) {

					this.Waiting = false;

					return true;

				} else {

					this.Waiting = false;

					return false;

				}

			} catch (err) {

				this.Waiting = false;

				return false;

			}

		} else {

			return false;

		}

	},

	SetAgent: function(agent) {

		agent = Object.prototype.toString.call(agent) === "[object Object]" ? agent : null;

		if (agent !== null) {

			this.Agent = agent;

			return true;

		}

		return false;

	},

	Update: function() {

		this.UpdateContextUsage();
		this.UpdateTools();

	},

	UpdateContextUsage: function() {

		fetch(this.Config.ResolveAPI("/api/session/context").toString(), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((context) => {

			if (Object.prototype.toString.call(context) === "[object Object]") {

				this.Context.Length = context["length"] || 0;
				this.Context.Tokens = context["tokens"] || 0;

			}

		});

	},

	UpdateTools: function() {

		// Empty Tools without losing references
		this.Tools.splice(0, this.Tools.length);

		fetch(this.Config.ResolveAPI("/api/session/tools").toString(), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((tools) => {

			tools.forEach((tool) => {
				this.Tools.push(tool);
			});

		});

	}

};
