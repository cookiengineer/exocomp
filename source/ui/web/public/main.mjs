
import { Client                           } from "./ui/Client.mjs";
import { CallTool    as CallToolPopover   } from "./ui/popovers/CallTool.mjs";
import { CreateAgent as CreateAgentDialog } from "./ui/dialogs/CreateAgent.mjs";
import { RenderSelect                     } from "./utils/ui/RenderSelect.mjs";
import { BootstrapConfig                  } from "./types/Config.mjs";

async function main() {

	try {

		const config  = await BootstrapConfig();
		const client  = new Client(config);
		const dialog  = new CreateAgentDialog(document.querySelector("dialog#create-agent"), config);
		const popover = new CallToolPopover(document.querySelector("div#call-tool"), client.Session.Tools);

		client.OnChange  = (prompt) => popover.Render(prompt);
		dialog.OnConfirm = (data)   => client.CreateAgent(data);

		window.CLIENT = client;

		let button = document.querySelector("header button[data-action=\"create-agent\"]");
		if (button !== null) {
			button.onclick = () => dialog.Show();
			button.removeAttribute("disabled");
		}

		client.Init();

		document.addEventListener("keyup", (event) => {

			if (event.key === "Escape") {

				if (dialog.IsVisible() === true) {
					dialog.Hide();
				}

			}

		});

	} catch (err) {
		console.error(err);
	}

};

main();
