
import { marked } from "/libs/marked.mjs";



export const RenderMessage = (message, with_empty_content) => {

	with_empty_content = typeof with_empty_content === "boolean" ? with_empty_content : false;

	if (message["role"] === "assistant") {

		if (message["content"] !== "") {

			let article = document.createElement("article");

			article.setAttribute("data-role", "assistant");
			article.innerHTML = marked.parse(message["content"]);

			return article;

		} else if (with_empty_content === true) {

			let article = document.createElement("article");

			article.setAttribute("data-role", "assistant");
			article.innerHTML = "(no content)";

			return article;

		} else {
			return null;
		}

	} else if (message["role"] === "system") {

		if (message["content"] !== "") {

			let article = document.createElement("article");

			article.setAttribute("data-role", "system");
			article.innerHTML = marked.parse(message["content"]);

			return article;

		} else if (with_empty_content === true) {

			let article = document.createElement("article");

			article.setAttribute("data-role", "system");
			article.innerHTML = "(no content)";

			return article;

		} else {
			return null;
		}

	} else if (message["role"] === "tool") {

		let article = document.createElement("article");

		article.setAttribute("data-role", "tool");

		let tmp = (message["content"] || "").split("\n");
		if (tmp.length == 1) {

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

	} else if (message["role"] === "user") {

		if (message["content"] !== "") {

			let article = document.createElement("article");

			article.setAttribute("data-role", "user");
			article.innerHTML = marked.parse(message["content"]);

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

