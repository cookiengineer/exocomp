
import { RenderAgentTimeline } from "/utils/svg/RenderAgentTimeline.mjs";

export const ScheduleRenderer = function(session) {

	this.Session  = session;
	this.rendered = {};
	this.running  = false;

	this.elements = {
		"svg": document.querySelector("body > main > svg[data-name=\"schedule\"]")
	};

};

ScheduleRenderer.prototype = {

	Clear: function() {

		if (this.elements["svg"] !== null) {

			Array.from(this.elements["svg"].querySelector("g")).forEach((element) => {
				element.parentNode.removeChild(element);
			});

		}

	},

	Destroy: function() {

		this.context = null;
		this.running = false;

	},

	Init: function() {

		if (this.running === false) {
			this.running = true;
			this.RenderLoop();
		}

	},

	RenderAgent: function(agent, planner, index) {

		let start = new Date();
		let end   = new Date();

		if (planner !== null) {

			if (planner.Messages.length > 0) {
				start = agent.Messages[0].Created;
			}

		} else {

			if (agent.Messages.length > 0) {
				start = agent.Messages[0].Created;
			}

		}

		// TODO: Draw a line from agent to planner if planner is not null

		let cached = this.rendered[agent.Name] || null;
		if (cached !== null) {

			if (agent.Messages.length > cached.length) {

				this.rendered[agent.Name] = {
					element: RenderAgentTimeline(agent, start, end),
					length:  agent.Messages.length
				};

				this.elements["svg"].replaceChild(this.rendered[agent.Name].element, cached.element);

			}

		} else {

			this.rendered[agent.Name] = {
				element: RenderAgentTimeline(agent, start, end),
				length:  agent.Messages.length
			};

			this.elements["svg"].appendChild(this.rendered[agent.Name].element);

		}

		let offset_x = 0;
		let offset_y = (32 + (index * 96));

		this.rendered[agent.Name].element.setAttribute("transform", "translate(" + offset_x.toString() + " " + offset_y.toString() + ")");

	},

	RenderLoop: function() {

		if (this.running === true) {

			let planner     = null;
			let contractors = [];

			Object.keys(this.Session.Agents).forEach((name) => {

				let agent = this.Session.Agents[name];

				if (agent.Role === "planner") {
					planner = agent;
				} else {
					contractors.push(agent);
				}

			});

			if (planner !== null) {

				let start = planner.Messages[0].Created;
				let end   = new Date("0001-01-01T00:00:01Z");

				planner.Messages.forEach((message) => {

					if (message.Created > end) {
						end = message.Created;
					}

				});

				contractors.forEach((agent) => {

					agent.Messages.forEach((message) => {

						if (message.Created > end) {
							end = message.Created;
						}

					});

				})

				this.RenderAgent(planner, null, 0);

				let width  = ((end - start) / 1000) | 0;
				let height = 32 + 96 + 96 * contractors.length;

				this.elements["svg"].setAttribute("width",  width  + 8);
				this.elements["svg"].setAttribute("height", height + 8);

			}

			if (contractors.length > 0) {

				contractors.forEach((agent, a) => {
					this.RenderAgent(agent, planner, a + 1);
				});

			}

			requestAnimationFrame(() => {
				this.RenderLoop();
			});

		}

	}

};
