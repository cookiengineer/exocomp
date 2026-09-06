
import { ChangelogEntry   } from "../types/ChangelogEntry.mjs";
import { ChangelogRenderer } from "./renderers/ChangelogRenderer.mjs";

const time_Second = 1 * 1000;

export const Changelog = function(config) {

	this.Config    = config;
	this.Entries   = [];
	this.Collapsed = false;
	this.Renderer  = new ChangelogRenderer();

	this.elements = {
		"toggle": document.querySelector("body > header button[data-action=\"toggle-collapse\"]")
	};

	this.timers = {
		entries: 0 // every 5 seconds
	};

	setTimeout(() => {
		this.UpdateEntries();
	}, 500);

};

Changelog.prototype = {

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
		this.Renderer.Render(this.Entries, this.Collapsed);

	},

	ToggleCollapse: function() {

		this.Collapsed = this.Collapsed === false;

		this.UpdateToggle();
		this.Renderer.Render(this.Entries, this.Collapsed);

	},

	UpdateEntries: function() {

		fetch(this.Config.ResolveAPI("/api/session/changelog").toString(), {
			method: "GET"
		}).then((response) => {
			return response.json();
		}).then((entries) => {

			if (Object.prototype.toString.call(entries) === "[object Object]") {

				let flat = [];

				Object.values(entries).forEach((symbols) => {
					Object.values(symbols).forEach((list) => {
						list.forEach((entry) => flat.push(entry));
					});
				});

				this.Entries = flat.map((entry) => {
					return ChangelogEntry.from(entry);
				});

				this.Renderer.Render(this.Entries, this.Collapsed);

			}

		});

	},

	UpdateLoop: function(delta) {

		this.timers.entries += delta;

		if (this.timers.entries >= 5 * time_Second) {
			this.UpdateEntries();
			this.timers.entries = 0;
		}

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
