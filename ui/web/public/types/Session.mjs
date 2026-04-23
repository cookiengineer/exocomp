
import { Console } from "./Console.mjs";

export const Session = function(config) {

	this.Config   = config;
	this.Console  = new Console();
	this.Messages = [];
	this.Tools    = [];
	this.Waiting  = false;
	this.context  = {
		length: 0,
		tokens: 0
	};

};

Session.prototype = {

	GetConsoleMessages: function(from) {

		if (this.Console !== null) {
			return this.Console.GetMessages(from);
		} else {
			return [];
		}

	},

	GetContextUsage: function() {

		if (this.context.length > 0) {
			return (this.context.tokens / this.context.length) * 100.0;
		} else {
			return 0.0;
		}

	},

	GetMessages: function(from) {

		let result = [];

		if (this.Messages.length > 0 && from < this.Messages.length) {

			for (let m = from; m < this.Messages.length; m++) {
				result.push(this.Messages[m]);
			}

		}

		return result;

	},

	Init: function() {

		fetch(this.Config.ResolveAPI("/api/messages").toString(), {
		}).then((response) => {
			return response.json();
		}).then((messages) => {

			if (Object.prototype.toString.call(messages) === "[object Array]" && messages.length > 0) {
				this.Messages = messages;
			}

		}).catch((err) => {
			console.error(err);
		});

	},

	SendChatRequest: async function(message) {

		if (this.Waiting === false) {

			this.Waiting = true;

			try {

				let response = await fetch(this.Config.ResolveAPI("/api/messages/send").toString(), {
					method:  "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(message)
				});

				if (response.ok === true) {

					let messages = await response.json();

					if (Object.prototype.toString.call(messages) === "[object Array]" && messages.length > 0) {

						messages.forEach((message) => {
							this.ReceiveChatResponse(message);
						});

					}

					this.Waiting = false;

					return null;

				} else {

					this.Waiting = false;
					return new Error("Session is busy, LLM is not responding ...");

				}

			} catch (err) {

				this.Waiting = false;
				return new Error("Session is busy, LLM is not responding ...");

			}

		} else {

			this.Waiting = false;
			return new Error("Session is busy, LLM is not responding ...");

		}

	},

	ReceiveChatResponse: function(message) {
		this.Messages.push(message);
	}

};
