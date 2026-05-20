
import { Agent           } from "../types/Agent.mjs";
import { Renderer        } from "./Renderer.mjs";
import { Session         } from "../types/Session.mjs";
import { ParseParameters } from "../utils/cli/ParseParameters.mjs";

const time_Second = 1 * 1000;

export const Client = function(config) {

	this.Session  = new Session(config);
	this.Renderer = new Renderer(this.Session);
	this.role     = "user";

	this.elements = {
		"nav":    document.querySelector("body > aside > nav[aria-label=\"agents\"]"),
		"footer": document.querySelector("body > footer"),
		"prompt": document.querySelector("body > footer textarea")
	};

	this.timers = {
		agents:  0, // every 10 seconds
		label:   0, // every  1 second
		session: 0, // every  5 seconds
	};

	this.OnChange = (prompt) => {};

	this.interval_id = null;

	setTimeout(() => {
		this.UpdateAgents();
	}, 500);

	setTimeout(() => {

		let agent = this.Session.GetAgent(null);
		if (agent !== null) {
			this.Renderer.RenderHeader(agent.Name + " | " + agent.Role + " | " + agent.Model + " | " + agent.Temperature.toFixed(1))
			this.Renderer.RenderTitle("Exocomp - " + agent.Name)
		}

	}, 1000);

};

Client.prototype = {

	CallTool: function(name, method, args) {

		this.Session.CallTool(name, method, args);

	},

	Destroy: function() {

		if (this.Renderer !== null) {
			this.Renderer.Destroy();
			this.Renderer = null;
		}

		if (this.interval_id !== null) {
			clearInterval(this.interval_id);
		}

		if (this.elements["nav"] !== null) {
			this.elements["nav"].removeEventListener("click");
		}

		if (this.elements["prompt"] !== null) {
			this.elements["prompt"].removeEventListener("keyup");
		}

	},

	HireAgent: function(data) {

		let result = false;
		let errors = [];

		let name    = (data["name"]    || "").trim();
		let role    = (data["role"]    || "").trim();
		let sandbox = (data["sandbox"] || "").trim();
		let prompt  = (data["prompt"]  || "").trim();

		if (name != "") {

			let tmp = this.Session.GetAgent(name);
			if (tmp !== null) {
				errors.push(new Error("Invalid Agent Name, must be a unique Pseudonym."));
			}

		} else {
			errors.push(new Error("Invalid Agent Name, must be a unique Pseudonym."));
		}

		if (role == "") {
			errors.push(new Error("Invalid Agent Role."));
		}

		if (sandbox != "" && sandbox.startsWith("./")) {
			// Valid
		} else {
			errors.push(new Error("Invalid Agent Sandbox, must start with \"./path/to/sandbox\""));
		}

		if (prompt != "") {
			// Valid
		} else {
			errors.push(new Error("Invalid Agent Prompt, must not be empty."));
		}

		if (errors.length === 0) {

			this.CallTool("agents.Hire", "Hire", {
				"name":    name,
				"role":    role,
				"sandbox": sandbox,
				"prompt":  prompt,
			});

			return {
				result: true,
				errors: [],
			};

		} else {

			return {
				result: false,
				errors: errors,
			};

		}

	},

	Init: function() {

		if (this.Session !== null) {
			this.Session.Init();
		}

		let last_interval_date = Date.now();

		this.interval_id = setInterval(() => {

			let now   = Date.now()
			let delta = now - last_interval_date;

			this.UpdateLoop(delta);

			last_interval_date = now;

		}, 1000 / 4);

		if (this.Renderer !== null) {
			this.Renderer.Init();
		}

		if (this.elements["nav"] !== null) {

			this.elements["nav"].addEventListener("click", (event) => {

				let element = event.target || null;
				if (element !== null) {

					if (element.tagName === "LABEL") {

						let name = (element.innerText || "").trim();
						if (name !== "") {

							this.ViewAgent(name);

							event.preventDefault()
							event.stopPropagation()

						}

						return true;

					} else if (element.tagName === "A") {

						event.preventDefault()
						event.stopPropagation()

						return false;

					}

				}

			}, {
				capture: true
			});

		}

		if (this.elements["prompt"] !== null) {

			this.elements["prompt"].addEventListener("keyup", (event) => {

				if (event.ctrlKey === true && event.key === "Enter") {

					let prompt = (this.elements["prompt"].value || "").trim();
					let role   = this.role;

					if (prompt !== "" && role !== "") {

						if (prompt.startsWith("/") && prompt.includes(" ") && !prompt.includes("\n")) {

							let name = (prompt.substr(1).split(" ")[0] || "/").trim();
							if (name.includes(".") === true) {

								let method = name.split(".").pop();
								let args   = ParseParameters(prompt.substr(name.length + 1));

								(async (name, method, args) => {
									this.CallTool(name, method, args);
								})(name, method, args);

								event.preventDefault();
								this.UpdatePrompt("");
								this.OnChange("");

							}

						} else if (prompt.startsWith("/") && !prompt.includes(" ") && !prompt.includes("\n")) {

							let name = (prompt.substr(1)).trim();
							if (name.includes(".") === true) {

								let method = name.split(".").pop();
								let args   = {};

								(async (name, method, args) => {
									this.CallTool(name, method, args);
								})(name, method, args);

								event.preventDefault();
								this.UpdatePrompt("");
								this.OnChange("");

							}

						} else {

							(async () => {

								this.UpdatePrompt("");

								let result = await this.Session.SendChatRequest({
									role:    role,
									content: prompt
								});

								if (result === true) {
									this.UpdateAgents();
								}

							})();

						}

					}

				} else if (event.key === "Enter") {

					let prompt = (this.elements["prompt"].value || "").trim();
					if (prompt.startsWith("/") && prompt.includes(" ") && !prompt.includes("\n")) {

						let name = (prompt.substr(1).split(" ")[0] || "/").trim();
						if (name.includes(".") === true) {

							let method = name.split(".").pop();
							let args   = ParseParameters(prompt.substr(name.length + 1));

							(async (name, method, args) => {
								this.CallTool(name, method, args);
							})(name, method, args);

							event.preventDefault();
							this.UpdatePrompt("");
							this.OnChange("");

						}

					}

				} else {

					let prompt = (this.elements["prompt"].value || "").trim();

					if (prompt.startsWith("/")) {
						this.elements["footer"].classList.add("command");
					} else {
						this.elements["footer"].classList.remove("command");
					}

					this.OnChange(prompt);

				}

			});

		}

	},

	SetRole: function(role) {

		if (role === "assistant" || role === "user") {
			this.role = role;
		}

	},

	UpdateAgents: function() {

		if (this.Session !== null) {

			fetch(this.Session.Config.ResolveAPI("/api/session/agents").toString(), {
				method: "GET"
			}).then((response) => {
				return response.json();
			}).then((agents) => {

				if (Object.prototype.toString.call(agents) === "[object Array]") {

					agents.forEach((agent) => {
						this.Session.ReceiveAgent(Agent.from(agent));
					});

				}

			});

		}

	},

	UpdateLabel: function() {

		let prompt = "";
		let usage  = 0.0;

		if (this.elements["prompt"] !== null) {
			prompt = (this.elements["prompt"].value || "").trim();
		}

		if (this.Session.Agent !== "") {

			let agent = this.Session.GetAgent(null);
			if (agent !== null) {

				if (agent.ContextUsage.Length > 0) {
					usage = ((agent.ContextUsage.Tokens / agent.ContextUsage.Length) * 100) | 0;
				}

			}

		}

		if (prompt !== "") {

			if (this.Session.Waiting === false) {
				this.Renderer.RenderLabel(usage.toFixed(2) + "%");
			} else {
				this.Renderer.RenderLabel("Processing ...");
			}

		} else {

			if (this.Session.Waiting === false) {
				this.Renderer.RenderLabel(usage.toFixed(2) + "%");
			} else {
				this.Renderer.RenderLabel("Processing ...");
			}

		}

	},

	UpdateLoop: function(delta) {

		this.timers.agents  += delta;
		this.timers.session += delta;
		this.timers.label   += delta;

		if (this.timers.label >= 1 * time_Second) {
			this.UpdateLabel();
			this.timers.label = 0;
		}

		if (this.timers.session >= 5 * time_Second) {

			if (this.Session !== null) {
				this.Session.Update();
			}

			this.timers.session = 0;

		}

		if (this.timers.agents >= 5 * time_Second) {
			this.UpdateAgents();
			this.timers.agents = 0;
		}

	},

	UpdatePrompt: function(message) {

		let prompt = message.trim();
		if (prompt.startsWith("/")) {
			this.elements["footer"].classList.add("command");
		} else {
			this.elements["footer"].classList.remove("command");
		}

		this.elements["prompt"].value = prompt;

		this.UpdateLabel();

	},

	ViewAgent: function(name) {

		name = typeof name === "string" ? name : null;

		if (name !== null) {

			let active = this.Session.GetAgent(null);
			let agent  = this.Session.GetAgent(name);

			if (agent !== null && agent !== active) {

				this.Session.SetAgent(agent.Name);

				this.Renderer.RenderHeader(agent.Name + " | " + agent.Role + " | " + agent.Model + " | " + agent.Temperature.toFixed(1))

				this.Renderer.ClearAgents();
				this.Renderer.RenderAgents(this.Session.Agent, this.Session.GetAgents());

				this.Renderer.ClearMessages();
				this.Renderer.RenderMessages(this.Session.GetMessages(0));

				if (agent.Name === this.Session.Config.Name) {
					this.elements["prompt"].removeAttribute("disabled");
				} else {
					this.elements["prompt"].setAttribute("disabled", "");
				}

			}

		}

	}

};
