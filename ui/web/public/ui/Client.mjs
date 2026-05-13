
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
				if (element !== null && element.tagName === "LABEL") {

					let name = (element.innerText || "").trim();
					if (name !== "") {
						this.ViewAgent(name);
					}

				}

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

		let prompt = (this.elements["prompt"].value || "").trim();
		let usage  = (this.Session.GetContextUsage() | 0);

		if (prompt !== "") {

			if (this.Session.Waiting === false) {
				this.Renderer.RenderLabel("[" + this.Session.Config.Model + " " + usage + "%]");
			} else {
				this.Renderer.RenderLabel("Processing ...");
			}

		} else {

			if (this.Session.Waiting === false) {
				this.Renderer.RenderLabel("[" + this.Session.Config.Model + " " + usage + "%]");
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

				console.info("ViewAgent: Viewing Agent \"" + name + "\" ...");
				console.info(agent);

				this.Session.SetAgent(agent.Name);

				this.Renderer.ClearAgents();
				this.Renderer.RenderAgents(this.Session.Agent, this.Session.GetAgents());

				this.Renderer.ClearMessages();
				this.Renderer.RenderMessages(this.Session.GetMessages(0));

				if (agent.Name === this.Session.Config.Name) {
					this.elements["prompt"].removeAttribute("disabled");
				} else {
					this.elements["prompt"].setAttribute("disabled", "");
				}

			} else {

				if (agent === active) {
					console.warn("ViewAgent: Agent \"" + name + "\" already viewed!");
				} else {
					console.warn("ViewAgent: Agent \"" + name + "\" doesn't exist?", agent);
				}

			}

		}

	}

};
