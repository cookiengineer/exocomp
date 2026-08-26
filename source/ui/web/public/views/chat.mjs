
import { Client                       } from "../ui/Client.mjs";
import { Init      as InitHeader      } from "../ui/components/layout/Header.mjs";
import { CallTool  as CallToolPopover } from "../ui/popovers/CallTool.mjs";
import { HireAgent as HireAgentDialog } from "../ui/dialogs/HireAgent.mjs";
import { BootstrapConfig              } from "../types/Config.mjs";

async function main() {

	try {

		const config = await BootstrapConfig(name);
		const client = new Client(config);

		InitHeader();

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

				document.addEventListener("keyup", (event) => {

					if (event.key === "Escape") {

						if (dialog.IsVisible() === true) {
							dialog.Hide();
						}

					}

				});

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

	} catch (err) {
		console.error(err);
	}

};

main();
