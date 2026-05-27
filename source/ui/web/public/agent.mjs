
import { Client          } from "./ui/Client.mjs";
import { BootstrapConfig } from "./types/Config.mjs";

const getAgentName = () => {

	let result = "";

	if (window.location.pathname === "/agent.html") {

		if (window.location.search.startsWith("?name=") === true) {

			let name = decodeURIComponent((window.location.search || "").substr(6).split("&").shift()).trim();
			if (name !== "") {
				result = name;
			}

		}

	}

	return result;

};

async function main() {

	try {

		const name   = getAgentName();
		const config = await BootstrapConfig(name);
		const client = new Client(config);

		window.CLIENT = client;
		window.CLIENT.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
