
import { Init as InitHeader         } from "/ui/components/layout/Header.mjs";
import { Changelog as ChangelogGrid } from "/ui/grids/Changelog.mjs";
import { BootstrapConfig            } from "/types/Config.mjs";

async function main() {

	try {

		const config    = await BootstrapConfig("");
		const changelog = new ChangelogGrid(config);

		InitHeader();

		window.CHANGELOG = changelog;
		window.CHANGELOG.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
