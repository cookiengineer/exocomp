
import { Agent   } from "./Agent.mjs";
import { Console } from "./Console.mjs";

export const Session = function(config) {

	// Managed by ui.Client
	this.Agent  = "";
	this.Agents = {};

	this.Config  = config;
	this.Console = new Console();
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
			return this.Agents[this.Agent] || null;
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

	GetMessages: function(from) {

		from = typeof from === "number" ? from : 0;


		let result = [];

		if (this.Agent !== "") {

			let agent = this.Agents[this.Agent] || null;
			if (agent !== null) {

				if (agent.Messages.length > 0 && from < agent.Messages.length) {

					for (let m = from; m < agent.Messages.length; m++) {
						result.push(agent.Messages[m]);
					}

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

			if (this.Agent === "" && agent.Name === this.Config.Name) {
				this.Agent = agent.Name;
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

		agent = typeof agent === "string" ? agent : null;

		if (agent !== null) {

			this.Agent = agent;

			return true;

		}

		return false;

	},

	Update: function() {

		this.UpdateTools();

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
