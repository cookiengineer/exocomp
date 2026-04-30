
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

	GetAgent: function(name) {

		if (this.Agents[name] !== undefined) {
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

	Init: function() {

		this.Update();

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

			this.Waiting = false;

			return false;

		}

	},

	SetAgent: function(agent) {

		agent = agent instanceof Agent ? agent : null;

		if (agent !== null) {

			this.Agents[agent.Name] = agent;

			if (agent.Name === this.Config.Name) {
				this.Agent = this.Agents[agent.Name];
			}

		}

	},

	Update: function() {

		this.UpdateContextUsage();

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

	}

};
