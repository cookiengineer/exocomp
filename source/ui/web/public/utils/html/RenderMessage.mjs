
import { marked } from "/libs/marked.mjs";



export const RenderMessage = (message, with_empty_content) => {

	with_empty_content = typeof with_empty_content === "boolean" ? with_empty_content : false;

	if (message.Role === "assistant") {

		let has_content          = message.Content          !== "";
		let has_reasoning        = message.ReasoningContent !== "";

		if (has_content === true || has_reasoning === true) {

			let article = document.createElement("article");
			let html    = "";

			article.setAttribute("data-role", "assistant");

			if (has_reasoning === true) {

				html += "<details class=\"reasoning\">";
				html += "<summary>Thought</summary>";
				html += marked.parse(message.ReasoningContent);
				html += "</details>";

			}

			if (has_content === true) {
				html += marked.parse(message.Content);
			}

			article.innerHTML = html;

			return article;

		} else if (with_empty_content === true) {

			let article = document.createElement("article");

			article.setAttribute("data-role", "assistant");
			article.innerHTML = "(no content)";

			return article;

		} else {
			return null;
		}

	} else if (message.Role === "system") {

		if (message.Content !== "") {

			let article = document.createElement("article");

			article.setAttribute("data-role", "system");
			article.innerHTML = marked.parse(message.Content);

			return article;

		} else if (with_empty_content === true) {

			let article = document.createElement("article");

			article.setAttribute("data-role", "system");
			article.innerHTML = "(no content)";

			return article;

		} else {
			return null;
		}

	} else if (message.Role === "tool") {

		let article = document.createElement("article");

		article.setAttribute("data-role", "tool");

		let tmp = message.Content.split("\n");
		if (tmp.length === 1) {

			article.innerHTML = [
				"<pre>" + tmp[0] + "</pre>"
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

	} else if (message.Role === "user") {

		if (message.Content !== "") {

			let article = document.createElement("article");

			article.setAttribute("data-role", "user");
			article.innerHTML = marked.parse(message.Content);

			return article;

		} else if (with_empty_content === true) {

			let article = document.createElement("article");

			article.setAttribute("data-role", "user");
			article.innerHTML = "(no content)";

			return article;

		} else {
			return null;
		}

	} else {
		return null;
	}

};

