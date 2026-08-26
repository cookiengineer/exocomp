
import { Init as InitHeader } from "../ui/components/layout/Header.mjs";
import { Schedule           } from "../ui/Schedule.mjs";
import { BootstrapConfig    } from "../types/Config.mjs";

async function main() {

	try {

		const config   = await BootstrapConfig("");
		const schedule = new Schedule(config);

		InitHeader();

		window.SCHEDULE = schedule;
		window.SCHEDULE.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
