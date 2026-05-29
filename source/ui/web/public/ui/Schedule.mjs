
import { Agent            } from "../types/Agent.mjs";
import { ScheduleRenderer } from "./renderers/ScheduleRenderer.mjs";

const time_Second = 1 * 1000;

export const Schedule = function(config) {

	this.Agent    = "";
	this.Agents   = {};
	this.Config   = config;
	this.Renderer = new ScheduleRenderer({
		Agent:  this.Agent,
		Agents: this.Agents,
		Config: this.Config,
	});

	this.elements = {
		"schedule": document.querySelector("body > main > canvas[data-name=\"schedule\"]")
	};

	this.timers = {
		agents: 0 // every 10 seconds
	};

	setTimeout(() => {
		this.UpdateAgents();
	}, 500);

};

Schedule.prototype = {

	Init: function() {

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

	UpdateAgents: function() {

		fetch(this.Config.ResolveAPI("/api/session/agents").toString(), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((agents) => {

			if (Object.prototype.toString.call(agents) === "[object Array]") {

				agents.forEach((agent) => {
					this.ReceiveAgent(Agent.from(agent));
				});

			}

		});

	},

	UpdateLoop: function(delta) {

		this.timers.agents += delta;

		if (this.timers.agents >= 10 * time_Second) {
			this.UpdateAgents();
			this.timers.agents = 0;
		}

	}

};
