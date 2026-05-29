
export const RenderAgentTimeline = (agent, schedule_start, schedule_end) => {

	let timeline = document.createElementNS("http://www.w3.org/2000/svg", "g");

	timeline.setAttribute("data-name", agent.Name);
	timeline.setAttribute("data-role", agent.Role);
	timeline.setAttribute("data-model", agent.Model);

	if (agent.Messages.length > 0) {

		let start  = agent.Messages[0].Created;
		let end    = agent.Messages[agent.Messages.length - 1].Created;
		let width  = ((end - start) / 1000) | 0;
		let x      = ((start - schedule_start) / 1000) | 0;
		let mapped = {};

		if (width < 128) {
			width = 128;
		}

		// TODO: Debug malformed agents JSON dumps
		// console.log(agent.Name, agent.Messages, (end - start) / 1000, "seconds")

		let background = document.createElementNS("http://www.w3.org/2000/svg", "rect");

		background.setAttribute("x", x);
		background.setAttribute("y", 0);
		background.setAttribute("width",  width);
		background.setAttribute("height", 48);
		background.setAttribute("rx", 2);
		background.setAttribute("ry", 2);

		timeline.appendChild(background);

		agent.Messages.forEach((message) => {

			if (message.Content != "") {

				if (mapped[message.Created] === undefined) {
					mapped[message.Created] = [];
				}

				mapped[message.Created].push(message);

			}

		});

		Object.keys(mapped).sort().forEach((time_str) => {

			let messages = mapped[time_str];
			if (messages.length > 1) {

				let offset = 0;

				// TODO: Might be better with offset to left first, then render from left to right
				// offset = messages.length * 8;

				console.log(messages);

				messages.forEach((message) => {

					let element = document.createElementNS("http://www.w3.org/2000/svg", "rect");

					let x = ((message.Created - start) / 1000) | 0;
					if (x < 4) {
						x = 4;
					}

					element.setAttribute("data-role", message.Role);
					element.setAttribute("x",       x - 4 + offset);
					element.setAttribute("y",       0);
					element.setAttribute("width",   8);
					element.setAttribute("height", 48);
					element.setAttribute("rx",      2);
					element.setAttribute("ry",      2);

					timeline.appendChild(element);
					offset += 8;

				});

			} else {

				let message = messages[0];
				let element = document.createElementNS("http://www.w3.org/2000/svg", "rect");

				let x = ((message.Created - start) / 1000) | 0;
				if (x < 4) {
					x = 4;
				}

				element.setAttribute("data-role", message.Role);
				element.setAttribute("x",       x - 4);
				element.setAttribute("y",       0);
				element.setAttribute("width",   8);
				element.setAttribute("height", 48);
				element.setAttribute("rx",      2);
				element.setAttribute("ry",      2);

				timeline.appendChild(element);

			}

		});

		let label = document.createElementNS("http://www.w3.org/2000/svg", "text");

		label.textContent = agent.Name;
		label.setAttribute("x", (width / 2) | 0);
		label.setAttribute("y", -16);

		label.setAttribute("text-anchor", "middle");
		label.setAttribute("dominant-baseline", "middle");

		timeline.appendChild(label);

	}

	// TODO: calculate width based on message from/to


	return timeline;

};
