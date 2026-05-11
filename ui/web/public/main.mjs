
import { Client                           } from "./ui/Client.mjs";
import { CallTool    as CallToolDialog    } from "./ui/dialogs/CallTool.mjs";
import { CreateAgent as CreateAgentDialog } from "./ui/dialogs/CreateAgent.mjs";
import { RenderSelect                     } from "./utils/ui/RenderSelect.mjs";
import { BootstrapConfig                  } from "./types/Config.mjs";

async function main() {

	try {

		const config = await BootstrapConfig();
		const client = new Client(config);

		const dialog1 = new CreateAgentDialog(document.querySelector("dialog#create-agent"), config);

		dialog1.OnConfirm = (data) => {
			client.CreateAgent(data);
		};

		const dialog2 = new CallToolDialog(document.querySelector("dialog#call-tool"), client.Session.Tools);

		dialog2.OnConfirm = (data) => {
			client.CallTool(data["name"], data["method"], data["arguments"]);
		};

		window.CLIENT = client;

		client.Init();

	} catch (err) {
		console.error(err);
	}

};

main();
