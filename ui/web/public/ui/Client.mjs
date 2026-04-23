
import { Renderer } from "./Renderer.mjs";
import { Session  } from "../types/Session.mjs";

export const Client = function(config) {

	// TODO: Figure out if Session or Config makes more sense
	this.Session  = new Session(config);
	this.Renderer = new Renderer(this.Session);

};

Client.prototype = {

	Destroy: function() {

	},

	Init: function() {

		document.querySelector("textarea#prompt").addEventListener("keyup", (event) => {

			// TODO: this.Session.SendChatRequest()
			console.log(event);

		});

	}

};
