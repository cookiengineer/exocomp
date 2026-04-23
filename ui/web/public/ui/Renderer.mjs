
import { marked } from "/libs/marked.mjs";

const RenderMessage = (message) => {

	if (message["role"] === "assistant") {

		let article = document.createElement("article");

		article.setAttribute("data-role", "assistant");

		if (message["content"] !== "") {
			article.innerHTML = marked.parse(message["content"]);
		} else {
			article.innerHTML = "(no content)";
		}

		return article;

	} else if (message["role"] === "system") {

		let article = document.createElement("article");

		article.setAttribute("data-role", "system");

		if (message["content"] !== "") {
			article.innerHTML = marked.parse(message["content"]);
		} else {
			article.innerHTML = "(no content)";
		}

		return article;

	} else if (message["role"] === "tool") {

		let article = document.createElement("article");

		article.setAttribute("data-role", "tool");

		let tmp = (message["content"] || "").split("\n");
		if (tmp.length == 1) {

			article.innerHTML = [
				"<details>",
				"<summary>",
				"<pre>" + tmp[0] + "</pre>",
				"</summary>",
				"</details>",
			].join("");

		} else if (tmp.length > 1) {

			article.innerHTML = [
				"<details>",
				"<summary>",
				"<pre>" + tmp[0].trim() + "</pre>",
				"</summary>",
				"<pre>" + tmp.slice(1).join("\n") + "<pre>",
				"</details>",
			].join("");

		} else {
			article.innerHTML = "<pre>(no content)</pre>";
		}

		return article;

	} else if (message["role"] === "user") {

		let article = document.createElement("article");

		article.setAttribute("data-role", "user");

		if (message["content"] !== "") {
			article.innerHTML = marked.parse(message["content"]);
		} else {
			article.innerHTML = "(no content)";
		}

		return article;

	} else {
		return null;
	}

};

export const Renderer = function(session) {

	this.Session  = session;
	this.rendered = 0;
	this.running  = false;

	this.elements = {
		"main": document.querySelector("body > main")
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

	RenderLoop: function() {

		if (this.running === true) {

			let messages = this.Session.GetMessages(this.rendered);
			if (messages.length > 0) {

				this.RenderMessages(messages);
				this.rendered += messages.length;

			}

			requestAnimationFrame(() => {
				this.RenderLoop();
			});

		}

	},

	RenderMessages: function(messages) {

		messages.forEach((message) => {

			let article = RenderMessage(message);
			if (article !== null) {
				this.elements["main"].appendChild(article);
			}

		});

	}

};
