
import { Renderer } from "./Renderer.mjs";
import { Session  } from "../types/Session.mjs";

export const Client = function(config) {

	this.Session  = new Session(config);
	this.Renderer = new Renderer(this.Session);

	this.elements = {
		"prompt": document.querySelector("body > footer textarea")
	};

};

Client.prototype = {

	Destroy: function() {

		if (this.Renderer !== null) {
			this.Renderer.Destroy();
			this.Renderer = null;
		}

	},

	Init: function() {

		if (this.Session !== null) {
			this.Session.Init();
		}

		if (this.Renderer !== null) {

			this.Renderer.Init();

			let usage = (this.Session.GetContextUsage() | 0);
			this.Renderer.RenderLabel("Write Message\n[" + this.Session.Config.Model + " " + usage + "%]");

		}

		this.elements["prompt"].addEventListener("keyup", (event) => {

			if (event.ctrlKey === true && event.key === "Enter") {

				let prompt = (this.elements["prompt"].value || "").trim();
				if (prompt !== "") {

					(async () => {

						let usage_before = (this.Session.GetContextUsage() | 0);

						this.UpdatePrompt("");

						// Executed after this.Session.SendChatRequest()
						setTimeout(() => {
							this.UpdateLabel();
						}, 16);

						// TODO: Remove this if above works
						// this.Renderer.RenderLabel("Processing ...\n[" + this.Session.Config.Model + " " + usage_before + "%]");

						await this.Session.SendChatRequest({
							role:    "user",
							content: prompt
						});

						this.UpdateLabel();

					})();

				}

			} else {

				this.UpdateLabel();

			}

		});

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

	UpdatePrompt: function(message) {

		let prompt = message.trim();

		this.elements["prompt"].value = prompt;

		this.UpdateLabel();

	}

};
