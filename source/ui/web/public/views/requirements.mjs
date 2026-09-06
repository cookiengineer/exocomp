
import { Init as InitHeader               } from "/ui/components/layout/Header.mjs";
import { Requirements as RequirementsGrid } from "/ui/grids/Requirements.mjs";
import { BootstrapConfig                  } from "/types/Config.mjs";

async function main() {

	try {

		const config       = await BootstrapConfig("");
		const requirements = new RequirementsGrid(config);

		InitHeader();

		window.REQUIREMENTS = requirements;
		window.REQUIREMENTS.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
