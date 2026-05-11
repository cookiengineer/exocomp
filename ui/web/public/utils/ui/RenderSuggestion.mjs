
export const RenderSuggestion = (suggestion) => {

	let html = "";

	html += "<li>";
	html += "<label>" + suggestion["label"] + "</label>";
	html += " ";
	html += "<span>" + suggestion["description"] + "</span>";
	html += "</li>";

	return html;

};
