
export const CallTool = function(element, tools) {

	this.Tools    = tools;
	this.Element  = element;
	this.elements = {
		"command":     element.querySelector("input[data-name=\"command\"]"),
		"suggestions": element.querySelector("ul[data-name=\"suggestions\"]"),
	};

	this.OnConfirm = (data) => {};
	this.OnCancel  = (data) => {};

	this.Init();

};

CallTool.prototype = {

	Hide: function() {

		if (this.Element !== null) {
			this.Element.close();
		}

	},

	Init: function() {

		// TODO: Render session.Tools into suggestions?
		// TODO: Execute Session.CallTool()

	},

	Show: function() {

		if (this.Element !== null) {
			this.Element.show();
		}

	},

};
