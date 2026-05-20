
export const RenderSuggestion = (suggestion) => {

	let html        = "";
	let label       = suggestion["label"]       || "";
	let description = suggestion["description"] || "";

	if (!label.startsWith("/")) {
		label = "&nbsp;&nbsp;" + label;
	}

	html += "<li>";
	html += "<label>" + label + "</label>";
	html += " ";
	html += "<span>" + description + "</span>";
	html += "</li>";

	return html;

};
