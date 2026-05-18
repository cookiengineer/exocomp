
import { Client                           } from "./ui/Client.mjs";
import { CallTool    as CallToolPopover   } from "./ui/popovers/CallTool.mjs";
import { CreateAgent as CreateAgentDialog } from "./ui/dialogs/CreateAgent.mjs";
import { RenderSelect                     } from "./utils/ui/RenderSelect.mjs";
import { BootstrapConfig                  } from "./types/Config.mjs";

const getAgentName = () => {

	let result = "";

	if (window.location.pathname === "/agent.html") {

		if (window.location.search.startsWith("?name=") === true) {

			let name = decodeURIComponent((window.location.search || "").substr(6).split("&").shift()).trim();
			if (name !== "") {
				result = name;
			}

		}

	}

	return result;

};

async function main() {

	try {

		const name   = getAgentName();
		const config = await BootstrapConfig(name);
		const client = new Client(config);

		((element, button) => {

			if (element !== null && button !== null) {

				let dialog = new CreateAgentDialog(element, config);

				dialog.OnConfirm = (data) => client.CreateAgent(data);
				button.onclick = () => dialog.Show();
				button.removeAttribute("disabled");

			}

		})(document.querySelector("dialog#create-agent"), document.querySelector("header button[data-action=\"create-agent\"]"));

		((element) => {

			if (element !== null) {

				let popover = new CallToolPopover(element, client.Session.Tools);

				client.OnChange = (prompt) => popover.Render(prompt);

			}

		})(document.querySelector("div#call-tool"));

		window.CLIENT = client;
		window.CLIENT.Init();

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
