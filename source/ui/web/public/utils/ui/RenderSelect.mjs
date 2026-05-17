
export const RenderSelect = (element, values) => {

	let html = [];

	for (let v = 0; v < values.length; v++) {

		let value      = values[v].trim();
		let attributes = "";
		let label      = value[0].toUpperCase() + value.slice(1).toLowerCase();

		if (v === 0) {
			attributes = " selected"
		}

		html.push("<option value=\"" + value + "\"" + attributes + ">" + label + "</option>");

	}

	element.innerHTML = html.join("");

};
