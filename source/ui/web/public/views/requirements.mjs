
import { Init as InitHeader } from "../ui/components/layout/Header.mjs";
import { Requirements       } from "../ui/Requirements.mjs";
import { BootstrapConfig    } from "../types/Config.mjs";

async function main() {

	try {

		const config       = await BootstrapConfig("");
		const requirements = new Requirements(config);

		InitHeader();

		window.REQUIREMENTS = requirements;
		window.REQUIREMENTS.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
