
import { Client       } from "./ui/Client.mjs";
import { RenderSelect } from "./utils/ui/RenderSelect.mjs";
import { GetAgents    } from "./utils/api/parameters/GetAgents.mjs";
import { GetModels    } from "./utils/api/parameters/GetModels.mjs";
import { GetConfig    } from "./utils/api/GetConfig.mjs";

async function main() {

	const dialog   = document.querySelector("dialog");
	const elements = {
		"agents": document.querySelector("dialog select[data-name=\"agent\"]"),
		"models": document.querySelector("dialog select[data-name=\"model\"]")
	};

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

	try {

		const config = await GetConfig();
		const client = new Client(config);

		window.CLIENT = client;

		client.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
