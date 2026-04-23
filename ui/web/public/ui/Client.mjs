
import { Renderer } from "./Renderer.mjs";
import { Session  } from "../types/Session.mjs";

export const Client = function(config) {

	this.Session  = new Session(config);
	this.Renderer = new Renderer(this.Session);

	this.elements = {
		"label":  document.querySelector("body > footer label"),
		"prompt": document.querySelector("body > footer textarea"),
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
		}

		this.elements["prompt"].addEventListener("keyup", (event) => {

			if (event.ctrlKey === true && event.key === "Enter") {

				let prompt = (this.elements["prompt"].value || "").trim();
				if (prompt !== "") {

					(async () => {

						this.elements["prompt"].value = "";
						this.elements["label"].innerHTML = "Waiting<br>...";

						await this.Session.SendChatRequest({
							role:    "user",
							content: prompt
						});

						this.elements["label"].innerHTML = "Type<br>Message";

					})();

				}

			} else {

				let prompt = (this.elements["prompt"].value || "").trim();
				if (prompt !== "") {

					if (this.Session.Waiting === false) {
						this.elements["label"].innerHTML = "Send with<br>[Ctrl]+[Enter]";
					}

				}

			}

		});

	}

};
