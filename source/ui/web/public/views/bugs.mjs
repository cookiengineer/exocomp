
import { Init as InitHeader } from "../ui/components/layout/Header.mjs";
import { Bugs                } from "../ui/Bugs.mjs";
import { BootstrapConfig     } from "../types/Config.mjs";

async function main() {

	try {

		const config = await BootstrapConfig("");
		const bugs   = new Bugs(config);

		InitHeader();

		window.BUGS = bugs;
		window.BUGS.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
