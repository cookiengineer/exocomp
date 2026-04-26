
import { RenderMessage } from "./RenderMessage.mjs";
import { Session       } from "/types/Session.mjs";



export const Renderer = function(session) {

	session = session instanceof Session ? session : null;


	this.Session  = session;
	this.rendered = {
		agents:   0,
		messages: 0
	};
	this.running  = false;

	this.elements = {
		"main":   document.querySelector("body > main"),
		"nav":    document.querySelector("body > aside > nav[aria-label=\"agents\"]"),
		"label":  document.querySelector("body > footer label")
	};

};

Renderer.prototype = {

	Destroy: function() {
		this.running = false;
	},

	Init: function() {

		if (this.running === false) {
			this.running = true;
			this.RenderLoop();
		}

	},

	RenderLabel: function(message) {

		message = typeof message === "string" ? message : "";


		let lines = message.split("\n");

		if (this.elements["label"] !== null) {
			this.elements["label"].innerHTML = lines.join("<br>");
		}

	},

	RenderLoop: function() {

		if (this.running === true) {

			let messages = this.Session.GetMessages(this.rendered.messages);
			if (messages.length > 0) {

				this.RenderMessages(messages);
				this.rendered.messages += messages.length;

			}

			let agents = this.Session.GetAgents();
			if (this.rendered.agents != Object.keys(agents).length) {

				this.RenderAgents(this.Session.Config.Name, agents);
				this.rendered.agents = Object.keys(agents).length;

			}

			requestAnimationFrame(() => {
				this.RenderLoop();
			});

		}

	},

	RenderAgents: function(active, agents) {

		agents = Object.prototype.toString.call(agents) === "[object Object]" ? agents : {};


		let html = "";

		html += "<ul>";
		html += Object.keys(agents).map((name) => {

			let agent = agents[name];

			if (name === active) {
				return "<li class=\"active\" title=\"" + agent.Name + " working in " + agent.Sandbox + "\"><label>" + agent.Name + "</label></li>";
			} else {
				return "<li title=\"" + agent.Name + " working in " + agent.Sandbox + "\"><label>" + agent.Name + "</label></li>";
			}

		}).join("");
		html += "</ul>";

		if (this.elements["nav"] !== null) {
			this.elements["nav"].innerHTML = html;
		}

	},

	RenderMessages: function(messages) {

		messages = Object.prototype.toString.call(messages) === "[object Array]" ? messages : [];


		let debug = false;

		if (this.Session !== null) {
			debug = this.Session.Config.Debug;
		}

		messages.forEach((message) => {

			let article = RenderMessage(message, debug);
			if (article !== null) {
				this.elements["main"].appendChild(article);
			}

		});

	}

};
