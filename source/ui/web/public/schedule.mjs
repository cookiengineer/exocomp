
import { Schedule        } from "./ui/Schedule.mjs";
import { BootstrapConfig } from "./types/Config.mjs";

async function main() {

	try {

		const config   = await BootstrapConfig("");
		const schedule = new Schedule(config);

		window.SCHEDULE = schedule;
		window.SCHEDULE.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
