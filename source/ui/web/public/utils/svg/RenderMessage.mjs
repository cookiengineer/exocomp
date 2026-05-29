
export const RenderMessage = (message, offset, start) => {

	let content = message.Content.trim();
	let element = document.createElementNS("http://www.w3.org/2000/svg", "rect");
	let title   = document.createElementNS("http://www.w3.org/2000/svg", "title");

	let x = ((message.Created - start) / 1000) | 0;
	if (x - 4 < 0) {
		x = 4;
	}

	element.setAttribute("data-role", message.Role);
	element.setAttribute("x",       x + offset);
	element.setAttribute("y",       0);
	element.setAttribute("width",   8);
	element.setAttribute("height", 48);
	element.setAttribute("rx",      2);
	element.setAttribute("ry",      2);

	if (message.Role === "tool") {

		if (content.startsWith("Error: ")) {
			title.textContent = "Tool Call: " + content.split("\n").shift().trim();
		} else if (content.includes(":")) {
			title.textContent = "Tool Call: " + content.split(":").shift().trim();
		} else {
			title.textContent = "Tool Call";
		}

	} else if (message.Role === "system") {
		title.textContent = "System Prompt";
	} else if (message.Role === "user") {
		title.textContent = "User Message";
	} else if (message.Role === "assistant") {
		title.textContent = "AI Message";
	}

	element.appendChild(title);

	return element;

};
