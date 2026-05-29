
import { group_messages } from "./group_messages.mjs";

export const MeasureAgentTimeline = (agent) => {

	if (agent.Messages.length > 0) {

		let start  = agent.Messages[0].Created;
		let end    = agent.Messages[agent.Messages.length - 1].Created;
		let width  = ((end - start) / 1000) | 0;
		let mapped = {};

		if (width < 128) {
			width = 128;
		}

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

					let x = ((message.Created - start) / 1000) | 0;
					if (x - 4 < 0) {
						x = 4;
					}

					width = Math.max(width, x + offset + 8 + 8);
					offset += 8;

				});

			}

		});

		return width;

	} else {
		return 128;
	}

};
