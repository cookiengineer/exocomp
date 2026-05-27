
export const ScheduleRenderer = function(session) {

	this.context  = null;
	this.rendered = {};
	this.running  = false;

	this.elements = {
		"canvas": document.querySelector("body > main > canvas[data-name=\"schedule\"]")
	};


};

ScheduleRenderer.prototype = {

	Clear: function() {

		if (this.elements["canvas"] !== null) {
			this.elements["canvas"]
		}

	},

	Destroy: function() {

		this.context = null;
		this.running = false;

	},

	Init: function() {

		if (this.elements["canvas"] !== null) {
			this.context = this.elements["canvas"].getContext("2d");
		}

		if (this.running === false) {
			this.running = true;
			this.RenderLoop();
		}

	},

	RenderLoop: function() {

		if (this.running === true) {

			// TODO: Render agents onto canvas

			requestAnimationFrame(() => {
				this.RenderLoop();
			});

		}

	}

};
