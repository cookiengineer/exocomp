
import { RenderSuggestion } from "../../utils/ui/RenderSuggestion.mjs";

export const CallTool = function(element, tools) {

	this.Tools    = tools;
	this.Element  = element;
	this.elements = {
		"suggestions": element.querySelector("ul[data-name=\"suggestions\"]"),
	};

	this.Init();

};

CallTool.prototype = {

	Hide: function() {

		if (this.Element !== null) {
			this.Element.hidePopover();
		}

	},

	Init: function() {

		if (this.Element !== null) {
			this.Element.setAttribute("popover", "");
		}

	},

	IsVisible: function() {

		if (this.Element !== null) {
			return this.Element.matches(":popover-open");
		}

	},

	Show: function() {

		if (this.Element !== null) {
			this.Element.showPopover();
		}

	},

	Render: function(prompt) {

		prompt = typeof prompt === "string" ? prompt : "";

		if (prompt.startsWith("/") === true) {

			if (this.IsVisible() === false) {
				this.Show();
			}

			let tool_command     = prompt.split(" ")[0];
			let tool_suggestions = this.Tools.map((schema) => {

				return {
					label:       "/" + schema["function"]["name"],
					description: schema["function"]["description"]
				};

			}).filter((suggestion) => {
				return suggestion["label"].startsWith(tool_command);
			});

			if (tool_suggestions.length > 1) {

				let list_items = tool_suggestions.map((suggestion) => RenderSuggestion(suggestion)).sort();

				if (this.elements["suggestions"] !== null) {
					this.elements["suggestions"].innerHTML = list_items.join("");
				}

			} else if (tool_suggestions.length === 1) {

				let list_items  = tool_suggestions.map((suggestion) => RenderSuggestion(suggestion)).sort();
				let tool_schema = this.Tools.find((schema) => {
					return schema["function"]["name"] === tool_command.substr(1, tool_command.length - 1);
				}) || null;

				if (tool_schema !== null) {

					tool_schema["function"]["parameters"]["required"].map((key) => {

						let value       = "...";
						let description = "";

						let property = tool_schema["function"]["parameters"]["properties"][key] || null;
						if (property !== null) {

							description = property["description"];

							if (property["type"] === "array") {

								if (property["items"]["type"] === "boolean") {
									value = "[true,...]";
								} else if (property["items"]["type"] === "number") {
									value = "[13.37,...]";
								} else if (property["items"]["type"] === "string") {
									value = "[\"foo\",...]";
								}

							} else if (property["type"] === "boolean") {
								value = "true";
							} else if (property["type"] === "number") {
								value = "13.37";
							} else if (property["type"] === "string") {

								if (Object.prototype.toString.call(property["enum"]) === "[object Array]") {

									value        = "\"" + property["enum"][0] + "\"";
									description += property["description"] + " (" + property["enum"].map((v) => "\"" + v + "\"").join(", ") + ")";

								} else {
									value = "\"...\"";
								}

							} else {
								console.error("Unsupported property type", property["type"]);
							}

						}

						return {
							label:       "&nbsp;&nbsp;" + key + "=" + value,
							description: description
						};

					}).map((suggestion) => {
						return RenderSuggestion(suggestion);
					}).forEach((list_item) => {
						list_items.push(list_item);
					});

				}

				if (this.elements["suggestions"] !== null) {
					this.elements["suggestions"].innerHTML = list_items.join("");
				}

			} else if (tool_suggestions.length === 0) {

				if (this.elements["suggestions"] !== null) {
					this.elements["suggestions"].innerHTML = "";
				}

			}

		} else {

			if (this.elements["suggestions"] !== null) {
				this.elements["suggestions"].innerHTML = "";
			}

			if (this.IsVisible() === true) {
				this.Hide();
			}

		}

	}

};
