
import { GetModels } from "./utils/GetModels.mjs";

async function init() {

	const models = GetModels();

	console.log(models);

};

init();
