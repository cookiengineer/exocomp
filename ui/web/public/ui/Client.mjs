
import { Renderer } from "./Renderer.mjs";
import { Session  } from "../types/Session.mjs";

export const Client = function(config) {

	this.Session  = new Session(config);
	this.Renderer = new Renderer(this.Session);

};

Client.prototype = {

	Destroy: function() {

	},

	Init: function() {

		this.Session.Init();

		document.querySelector("textarea#prompt").addEventListener("keyup", (event) => {

			if (event.ctrlKey === true && event.key === "Enter") {

				let prompt = (document.querySelector("textarea#prompt").value || "").trim();
				if (prompt !== "") {

					this.Session.SendChatRequest({
						role:    "user",
						content: prompt
					});

				}

			} else {
				// Do Nothing
			}

		});

	}

};
