
export const RenderSelect = (element, values) => {

	let html = [];

	values.forEach((value) => {
		html.push("<option value=" + value + ">" + value + "</option>");
	});

	element.innerHTML = html.join("");

};
