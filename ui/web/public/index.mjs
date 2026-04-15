
import { RenderSelect } from "./utils/ui/RenderSelect.mjs";
import { GetAgents    } from "./utils/api/GetAgents.mjs";
import { GetModels    } from "./utils/api/GetModels.mjs";

async function init() {

	const dialog   = document.querySelector("dialog");
	const elements = {
		"agents": document.querySelector("dialog select[data-name=\"agent\"]"),
		"models": document.querySelector("dialog select[data-name=\"model\"]")
	};


	console.log(elements);


	if (dialog !== null) {

		let close_dialog = dialog.querySelector("button[data-action=\"close\"]");
		if (close_dialog !== null) {
			close_dialog.onclick = () => {
				dialog.close();
			};
		}

		let create_agent = document.querySelector("header button[data-action=\"create-agent\"]");
		if (create_agent !== null) {
			create_agent.onclick = () => {
				dialog.show();
			};
		}

	}

	try {

		const agents = await GetAgents();
		const models = await GetModels();

		RenderSelect(elements["agents"], agents);
		RenderSelect(elements["models"], models);

	} catch (err) {
		console.error(err);
	}

};

init();
