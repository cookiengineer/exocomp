
import { RenderSelect } from "../../utils/ui/RenderSelect.mjs";

export const HireAgent = function(element, config) {

	this.Config   = config;
	this.Element  = element;
	this.elements = {
		"name":    element.querySelector("input[data-name=\"name\"]"),
		"role":    element.querySelector("select[data-name=\"role\"]"),
		"sandbox": element.querySelector("input[data-name=\"sandbox\"]"),
		"prompt":  element.querySelector("textarea[data-name=\"prompt\"]"),
		"errors":  element.querySelector("div[data-name=\"errors\"]"),
	};

	this.OnConfirm = (data) => {};
	this.OnCancel  = (data) => {};

	this.options = {
		"role": [ "architect", "coder", "tester" ]
	};

	this.Init();

};

HireAgent.prototype = {

	Error: function(errors) {

		let element = this.elements["errors"] || null;
		if (element !== null) {

			let items = errors.map((err) => {
				return "<b>" + err.toString() + "</b>";
			});

			element.innerHTML = items.join("");

		}

	},

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
						"name":    "",
						"role":    "planner",
						"sandbox": "",
						"prompt":  ""
					};

					let name = this.elements["name"] || null;
					if (name !== null) {

						let tmp = name.value.trim();
						if (tmp !== "") {
							data["name"] = tmp;
						}

					}

					let roles = this.elements["role"] || null;
					if (roles !== null) {

						let tmp = roles.value.trim();
						if (tmp !== "") {
							data["role"] = tmp;
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

					let roles = this.elements["role"] || null;
					if (roles !== null) {
						roles.selectedIndex = 0;
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

	IsVisible: function() {

		if (this.Element !== null) {
			return this.Element.hasAttribute("open");
		}

		return false;

	},

	Reset: function() {

		let name = this.elements["name"] || null;
		if (name !== null) {
			name.value = "";
		}

		let roles = this.elements["role"] || null;
		if (roles !== null) {
			roles.selectedIndex = 0;
		}

		let sandbox = this.elements["sandbox"] || null;
		if (sandbox !== null) {
			sandbox.value = "";
		}

		let prompt = this.elements["prompt"] || null;
		if (prompt !== null) {
			prompt.value = "";
		}

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

		fetch(this.Config.ResolveAPI("/api/parameters/roles"), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((roles) => {

			if (Object.prototype.toString.call(roles) === "[object Array]") {

				roles.forEach((role) => {

					if (role !== "planner") {

						if (this.options["role"].includes(role) === false) {
							this.options["role"].push(role);
						}

					}

				});

				this.options["role"] = this.options["role"].sort();

				let element = this.elements["role"];
				if (element !== null) {
					RenderSelect(element, this.options["role"]);
				}

			}

		});

	}

};
