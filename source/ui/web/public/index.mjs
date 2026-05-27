
import { Client                       } from "./ui/Client.mjs";
import { CallTool  as CallToolPopover } from "./ui/popovers/CallTool.mjs";
import { HireAgent as HireAgentDialog } from "./ui/dialogs/HireAgent.mjs";
import { BootstrapConfig              } from "./types/Config.mjs";

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

		((button) => {

			if (button !== null) {

				button.onclick = () => window.location.assign("/schedule.html");
				button.removeAttribute("disabled");

			}

		})(document.querySelector("header button[data-action=\"show-schedule\"]"));

		((element, button) => {

			if (element !== null && button !== null) {

				let dialog = new HireAgentDialog(element, config);

				dialog.OnConfirm = (data) => {

					let { result, errors } = client.HireAgent(data);

					if (result === true && errors.length === 0) {

						dialog.Reset();
						dialog.Hide();

					} else {

						if (errors.length > 0) {
							dialog.Error(errors);
						}

					}

				};

				button.onclick = () => dialog.Show();
				button.removeAttribute("disabled");

			}

		})(document.querySelector("dialog#hire-agent"), document.querySelector("header button[data-action=\"hire-agent\"]"));

		((element) => {

			if (element !== null) {

				let popover = new CallToolPopover(element, client.Session.Tools);

				client.OnChange   = (prompt)     => popover.Render(prompt);
				popover.OnSuggest = (suggestion) => client.Suggest(suggestion);

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
