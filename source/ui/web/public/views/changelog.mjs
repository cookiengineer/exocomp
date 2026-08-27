
import { Init as InitHeader } from "../ui/components/layout/Header.mjs";
import { Changelog          } from "../ui/Changelog.mjs";
import { BootstrapConfig    } from "../types/Config.mjs";

async function main() {

	try {

		const config    = await BootstrapConfig("");
		const changelog = new Changelog(config);

		InitHeader();

		window.CHANGELOG = changelog;
		window.CHANGELOG.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
