
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

	Clear: function() {

		this.ClearAgents();
		this.ClearLabel();
		this.ClearMessages();

	},

	ClearAgents: function() {

		this.rendered.agents = 0;

		if (this.elements["nav"] !== null) {
			this.elements["nav"].innerHTML = "";
		}

	},

	ClearLabel: function() {

		if (this.elements["label"] !== null) {
			this.elements["label"].innerHTML = "";
		}

	},

	ClearMessages: function() {

		this.rendered.messages = 0;

		if (this.elements["main"] !== null) {

			Array.from(this.elements["main"].querySelectorAll("article")).forEach((article) => {
				article.parentNode.removeChild(article);
			});

		}

	},

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

			let messages = this.Session.GetMessages(0);

			if (this.rendered.messages != messages.length) {

				this.RenderMessages(messages);
				this.rendered.messages = messages.length;

			}

			let agents = this.Session.GetAgents();
			if (this.rendered.agents != Object.keys(agents).length) {

				this.RenderAgents(this.Session.Agent, agents);
				this.rendered.agents = Object.keys(agents).length;

			}

			requestAnimationFrame(() => {
				this.RenderLoop();
			});

		}

	},

	RenderAgents: function(active, agents) {

		active = Object.prototype.toString.call(active) === "[object Object]" ? active : null;
		agents = Object.prototype.toString.call(agents) === "[object Object]" ? agents : {};


		let html = "";

		html += "<ul>";
		html += Object.keys(agents).map((name) => {

			let agent = agents[name];

			if (active !== null && active.Name === agent.Name) {
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

		if (this.elements["main"] !== null) {
			this.elements["main"].innerHTML = "";
		}

		messages.forEach((message) => {

			let article = RenderMessage(message, debug);
			if (article !== null) {
				this.elements["main"].appendChild(article);
			}

		});

	}

};
