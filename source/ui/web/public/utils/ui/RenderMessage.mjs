
import { FormatContent   } from "../fmt/FormatContent.mjs";
import { SanitizeContent } from "../fmt/SanitizeContent.mjs";

export const RenderMessage = (message) => {

	let article = document.createElement("article");

	let role = (message["role"] || "");
	if (role === "tool") {

		article.setAttribute("data-role", role);

		let lines = (message["content"] || "").split("\n");
		if (lines.length > 1) {

			let html = "";

			html += "<details>";
			html += "<summary>" + lines[0].trim() + "</summary>";
			html += "<pre>";

			for (let l = 1; l < lines.length; l++) {
				html += SanitizeContent(lines[l].trim());
				html += "<br>";
			}

			html += "</pre>";
			html += "</details>";

			article.setHTML(html);

		} else {
			article.innerText = lines[0];
		}

	} else if (role === "assistant" || role === "user") {

		article.setAttribute("data-role", role);

		let html  = "";
		let lines = (message["content"] || "").split("\n");

		for (let l = 0; l < lines.length; l++) {
			html += FormatContent(lines[l].trim());
			html += "<br>";
		}

		article.setHTML(html);

	}

	return article;

};
