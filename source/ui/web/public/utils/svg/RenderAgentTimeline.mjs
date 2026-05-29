
import { RenderMessage  } from "./RenderMessage.mjs";
import { group_messages } from "./group_messages.mjs";

export const RenderAgentTimeline = (agent, index, schedule_start, schedule_end) => {

	let timeline = document.createElementNS("http://www.w3.org/2000/svg", "g");

	timeline.setAttribute("data-name", agent.Name);
	timeline.setAttribute("data-role", agent.Role);
	timeline.setAttribute("data-model", agent.Model);

	if (agent.Messages.length > 0) {

		let start    = agent.Messages[0].Created;
		let end      = agent.Messages[agent.Messages.length - 1].Created;
		let width    = ((end - start) / 1000) | 0;
		let offset_x = ((start - schedule_start) / 1000) | 0;
		let offset_y = 32 + (index * 96);
		let mapped   = {};

		if (width < 128) {
			width = 128;
		}

		let background = document.createElementNS("http://www.w3.org/2000/svg", "rect");

		background.setAttribute("x", 0);
		background.setAttribute("y", 0);
		background.setAttribute("width",  width);
		background.setAttribute("height", 48);
		background.setAttribute("rx", 2);
		background.setAttribute("ry", 2);

		timeline.setAttribute("transform", "translate(" + offset_x.toString() + " " + offset_y.toString() + ")");
		timeline.appendChild(background);

		agent.Messages.forEach((message) => {

			if (message.Content != "") {

				if (mapped[message.Created] === undefined) {
					mapped[message.Created] = [];
				}

				mapped[message.Created].push(message);

			}

		});

		mapped = group_messages(mapped);

		Object.keys(mapped).sort().forEach((time_str) => {

			let messages = mapped[time_str];
			if (messages.length > 0) {

				let offset = -1/2 * (messages.length * 8);

				messages.forEach((message) => {

					let element = RenderMessage(message, offset, start);
					let x       = ((message.Created - start) / 1000) | 0;

					width = Math.max(width, x + offset + 8);
					timeline.appendChild(element);
					offset += 8;

				});

			}

		});

		background.setAttribute("width", width);

		let label = document.createElementNS("http://www.w3.org/2000/svg", "text");

		label.textContent = agent.Name;
		label.setAttribute("x", (width / 2) | 0);
		label.setAttribute("y", -16);

		label.setAttribute("text-anchor", "middle");
		label.setAttribute("dominant-baseline", "middle");

		timeline.appendChild(label);

	}

	return timeline;

};
