
import { Requirement          } from "../types/Requirement.mjs";
import { RequirementsRenderer } from "./renderers/RequirementsRenderer.mjs";

const time_Second = 1 * 1000;

export const Requirements = function(config) {

	this.Config    = config;
	this.Reports   = [];
	this.Collapsed = false;
	this.Renderer  = new RequirementsRenderer();

	this.elements = {
		"toggle": document.querySelector("body > header button[data-action=\"toggle-collapse\"]")
	};

	this.timers = {
		reports: 0 // every 5 seconds
	};

	setTimeout(() => {
		this.UpdateReports();
	}, 500);

};

Requirements.prototype = {

	Init: function() {

		if (this.elements["toggle"] !== null) {
			this.elements["toggle"].onclick = () => this.ToggleCollapse();
		}

		let last_interval_date = Date.now();

		this.interval_id = setInterval(() => {

			let now   = Date.now();
			let delta = now - last_interval_date;

			this.UpdateLoop(delta);

			last_interval_date = now;

		}, 1000 / 4);

		this.UpdateToggle();
		this.Renderer.Render(this.Reports, this.Collapsed);

	},

	ToggleCollapse: function() {

		this.Collapsed = this.Collapsed === false;

		this.UpdateToggle();
		this.Renderer.Render(this.Reports, this.Collapsed);

	},

	UpdateLoop: function(delta) {

		this.timers.reports += delta;

		if (this.timers.reports >= 5 * time_Second) {
			this.UpdateReports();
			this.timers.reports = 0;
		}

	},

	UpdateReports: function() {

		fetch(this.Config.ResolveAPI("/api/session/requirements").toString(), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((reports) => {

			if (Object.prototype.toString.call(reports) === "[object Object]") {

				let flat = [];

				Object.values(reports).forEach((symbols) => {
					Object.values(symbols).forEach((report) => flat.push(report));
				});

				this.Reports = flat.map((report) => {
					return Requirement.from(report);
				});

				this.Renderer.Render(this.Reports, this.Collapsed);

			}

		});

	},

	UpdateToggle: function() {

		if (this.elements["toggle"] !== null) {

			if (this.Collapsed === true) {
				this.elements["toggle"].setAttribute("data-state", "collapsed");
				this.elements["toggle"].textContent = "Expand Packages";
			} else {
				this.elements["toggle"].setAttribute("data-state", "nested");
				this.elements["toggle"].textContent = "Collapse Packages";
			}

		}

	}

};
