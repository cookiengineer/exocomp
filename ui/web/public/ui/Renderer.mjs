
import { RenderMessage } from "./RenderMessage.mjs";
import { Session       } from "/types/Session.mjs";



export const Renderer = function(session) {

	session = session instanceof Session ? session : null;


	this.Session  = session;
	this.rendered = 0;
	this.running  = false;

	this.elements = {
		"main":   document.querySelector("body > main"),
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

			let messages = this.Session.GetMessages(this.rendered);
			if (messages.length > 0) {

				this.RenderMessages(messages);
				this.rendered += messages.length;

			}

			let agents = this.Session.GetAgents();
			if (agents.length > 0) {
				this.RenderAgents(agents);
			}

			requestAnimationFrame(() => {
				this.RenderLoop();
			});

		}

	},

	RenderAgents: function(agents) {

		agents = Object.prototype.toString.call(agents) === "[object Object]" ? agents : {};

		Object.keys(agents).forEach((name) => {

			let agent = agents[name];

			// TODO: Render a <ul><li>...</li></ul>
			// into the sidebar

		});

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
