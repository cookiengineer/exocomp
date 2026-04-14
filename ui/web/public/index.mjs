
import { GetModels } from "./utils/GetModels.mjs";

async function init() {

	const elements = {
		"agents":      document.querySelector("select[data-name=\"agents\"]"),
		"models":      document.querySelector("select[data-name=\"models\"]"),
		"temperature": document.querySelector("input[data-name=\"temperature\"]")
	};

	const agents = await GetAgents();
	const models = await GetModels();

	console.log("agents", agents);
	console.log("models", models);

};

init();
