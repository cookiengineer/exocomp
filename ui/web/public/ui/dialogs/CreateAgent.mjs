
import { RenderSelect } from "../../utils/ui/RenderSelect.mjs";

export const CreateAgent = function(element, config) {

	this.Config   = config;
	this.Element  = element;
	this.elements = {
		"agent":   element.querySelector("select[data-name=\"agent\"]"),
		"model":   element.querySelector("select[data-name=\"model\"]"),
		"prompt":  element.querySelector("input[data-name=\"prompt\"]"),
		"sandbox": element.querySelector("input[data-name=\"sandbox\"]"),
	};

	this.OnConfirm = (data) => {};
	this.OnCancel  = (data) => {};

	this.options = {
		"agent": [ "planner", "architect", "coder", "tester" ],
		"model": [ "qwen3-coder:30b", "gemma4:31b" ],
	};

	this.Init();

};

CreateAgent.prototype = {

	Hide: function() {

		if (this.Element !== null) {
			this.Element.close();
		}

	},

	Init: function() {

		if (this.Element !== null) {

			let button = this.Element.querySelector("button[data-action=\"close\"]");
			if (button !== null) {
				button.onclick = () => {
					this.Hide();
				};
			}

			let confirm = this.Element.querySelector("button[data-action=\"confirm\"]");
			if (confirm !== null) {
				confirm.onclick = () => {

					let data = {
						"agent":   "planner",
						"model":   "qwen3-coder:30b",
						"prompt":  "",
						"sandbox": ""
					};

					let agent = this.elements["agent"] || null;
					if (agent !== null) {

						let tmp = agent.value.trim();
						if (tmp !== "") {
							data["agent"] = tmp;
						}

					}

					let model = this.elements["model"] || null;
					if (model !== null) {

						let tmp = model.value.trim();
						if (tmp !== "") {
							data["model"] = tmp;
						}

					}

					let prompt = this.elements["prompt"] || null;
					if (prompt !== null) {

						let tmp = prompt.value.trim();
						if (tmp !== "") {
							data["prompt"] = tmp;
						}

					}

					let sandbox = this.elements["sandbox"] || null;
					if (sandbox !== null) {

						let tmp = sandbox.value.trim();
						if (tmp !== "") {
							data["sandbox"] = tmp;
						}

					}

					this.OnConfirm(data);

				};
			}

			let cancel = this.Element.querySelector("button[data-action=\"cancel\"]");
			if (cancel !== null) {
				cancel.onclick = () => {

					let agent = this.elements["agent"] || null;
					if (agent !== null) {
						agent.selectedIndex = 0;
					}

					let model = this.elements["model"] || null;
					if (model !== null) {
						model.selectedIndex = 0;
					}

					let prompt = this.elements["prompt"] || null;
					if (prompt !== null) {
						prompt.value = "";
					}

					let sandbox = this.elements["sandbox"] || null;
					if (sandbox !== null) {
						sandbox.value = "";
					}

					this.OnCancel();

				};
			}

		}

		this.Update();

	},

	Show: function() {

		if (this.Element !== null) {
			this.Element.show();
		}

	},

	Update: function() {

		this.UpdateParameters();

	},

	UpdateParameters: async function() {

		fetch(this.Config.ResolveAPI("/api/parameters/agents"), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((agents) => {

			if (Object.prototype.toString.call(agents) === "[object Array]") {

				agents.forEach((agent) => {

					if (this.options["agent"].includes(agent) === false) {
						this.options["agent"].push(agent);
					}

				});

				let element = this.elements["agent"];
				if (element !== null) {
					RenderSelect(element, this.options["agent"]);
				}

			}

		});

		fetch(this.Config.ResolveAPI("/api/parameters/models"), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((agents) => {

			if (Object.prototype.toString.call(agents) === "[object Array]") {

				models.forEach((model) => {

					if (this.options["model"].includes(model) === false) {
						this.options["model"].push(agent);
					}

				});

				this.options["model"] = this.options["model"].sort();

				let element = this.elements["model"];
				if (element !== null) {
					RenderSelect(element, this.options["model"]);
				}

			}

		});

	}

};
