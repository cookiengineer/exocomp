
import { Init as InitHeader } from "/ui/components/layout/Header.mjs";
import { Bugs as BugsGrid   } from "/ui/grids/Bugs.mjs";
import { BootstrapConfig    } from "/types/Config.mjs";

async function main() {

	try {

		const config = await BootstrapConfig("");
		const bugs   = new BugsGrid(config);

		InitHeader();

		window.BUGS = bugs;
		window.BUGS.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
