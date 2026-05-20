
import { RenderSuggestion } from "../../utils/ui/RenderSuggestion.mjs";

export const CallTool = function(element, tools) {

	this.Tools    = tools;
	this.Element  = element;
	this.elements = {
		"suggestions": element.querySelector("ul[data-name=\"suggestions\"]"),
	};

	this.OnSuggest = (suggestion) => {};

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
			let tool_parameters  = prompt.split(" ").slice(1);

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

				let next_suggestion = null;
				let list_items      = tool_suggestions.map((suggestion) => RenderSuggestion(suggestion)).sort();
				let tool_schema     = this.Tools.find((schema) => {
					return schema["function"]["name"] === tool_command.substr(1, tool_command.length - 1);
				}) || null;

				if (tool_schema !== null) {

					tool_schema["function"]["parameters"]["required"].map((key) => {

						let value       = "...";
						let type        = "";
						let description = "";

						let property = tool_schema["function"]["parameters"]["properties"][key] || null;
						if (property !== null) {

							description = property["description"];

							if (property["type"] === "array") {

								if (property["items"]["type"] === "boolean") {

									return {
										key:         key,
										type:        "array-of-booleans",
										value:       "[true,...]",
										label:       key + "=[true,...]",
										description: property["description"],
									};

								} else if (property["items"]["type"] === "number") {

									return {
										key:         key,
										type:        "array-of-numbers",
										value:       "[13,37,...]",
										label:       key + "=[13,37,...]",
										description: property["description"],
									};

								} else if (property["items"]["type"] === "string") {

									return {
										key:         key,
										type:        "array-of-strings",
										value:       "[\"foo\",...]",
										label:       key + "=[\"foo\",...]",
										description: property["description"],
									};

								} else {

									console.error("Unsupported property type", property["type"], property["items"]["type"]);
									return null;

								}

							} else if (property["type"] === "boolean") {

								return {
									key:         key,
									type:        "boolean",
									value:       "true",
									label:       key + "=true",
									description: property["description"],
								};

							} else if (property["type"] === "number") {

								return {
									key:         key,
									type:        "number",
									value:       "1337",
									label:       key + "=1337",
									description: property["description"],
								};

							} else if (property["type"] === "string") {

								if (Object.prototype.toString.call(property["enum"]) === "[object Array]") {

									return {
										key:         key,
										type:        "string",
										value:       "\"" + property["enum"][0] + "\"",
										label:       key + "=\"" + property["enum"][0] + "\"",
										description: property["description"] + " (" + property["enum"].map((v) => "\"" + v + "\"").join(", ") + ")",
									};

								} else {

									return {
										key:         key,
										type:        "string",
										value:       "\"...\"",
										label:       key + "=\"...\"",
										description: property["description"],
									};

								}

							} else {

								console.error("Unsupported property type", property["type"]);
								return null;

							}

						}

						return null;

					}).forEach((suggestion) => {

						if (next_suggestion === null) {

							let found = false;

							for (let p = 0; p < tool_parameters.length; p++) {

								let parameter = tool_parameters[p];

								if (parameter !== "" && parameter.startsWith(suggestion.key + "=")) {
									found = true;
									break;
								}

							}

							if (found === false) {
								next_suggestion = suggestion;
							}

						}

						list_items.push(RenderSuggestion(suggestion));

					});

				}

				if (next_suggestion !== null) {
					this.OnSuggest(next_suggestion);
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
