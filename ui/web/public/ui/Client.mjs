
import { Renderer } from "./Renderer.mjs";
import { Session  } from "../types/Session.mjs";

const time_Second = 1 * 1000;

export const Client = function(config) {

	this.Session  = new Session(config);
	this.Renderer = new Renderer(this.Session);

	this.elements = {
		"prompt": document.querySelector("body > footer textarea")
	};

	this.timers = {
		agents:  0, // every 10 seconds
		label:   0, // every  1 second
		session: 0, // every  5 seconds
	};

	this.interval_id = null;

};

Client.prototype = {

	Destroy: function() {

		if (this.Renderer !== null) {
			this.Renderer.Destroy();
			this.Renderer = null;
		}

		if (this.interval_id !== null) {
			clearInterval(this.interval_id);
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

			let usage = (this.Session.GetContextUsage() | 0);
			this.Renderer.RenderLabel("Write Message\n[" + this.Session.Config.Model + " " + usage + "%]");

		}

		if (this.elements["prompt"] !== null) {

			this.elements["prompt"].addEventListener("keyup", (event) => {

				if (event.ctrlKey === true && event.key === "Enter") {

					let prompt = (this.elements["prompt"].value || "").trim();
					if (prompt !== "") {

						(async () => {

							let usage_before = (this.Session.GetContextUsage() | 0);

							this.UpdatePrompt("");

							await this.Session.SendChatRequest({
								role:    "user",
								content: prompt
							});

						})();

					}

				} else {

					// Do Nothing

				}

			});

		}

	},

	UpdateAgents: function() {

		if (this.Session !== null) {

			fetch(this.Session.Config.ResolveAPI("/api/agents").toString(), {
				method: "GET"
			}).then((response) => {
				return response.json();
			}).then((agents) => {

				if (Object.prototype.toString.call(agents) === "[object Object]") {

					Object.keys(agents).sort().forEach((name) => {
						this.Session.Agents[name] = agents[name];
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
				this.Renderer.RenderLabel("Send with [Ctrl]+[Enter]\n[" + this.Session.Config.Model + " " + usage + "%]");
			} else {
				this.Renderer.RenderLabel("Processing ...\n[" + this.Session.Config.Model + " " + usage + "%]");
			}

		} else {

			if (this.Session.Waiting === false) {
				this.Renderer.RenderLabel("Write Message\n[" + this.Session.Config.Model + " " + usage + "%]");
			} else {
				this.Renderer.RenderLabel("Processing ...\n[" + this.Session.Config.Model + " " + usage + "%]");
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

		this.elements["prompt"].value = prompt;

		this.UpdateLabel();

	}

};
