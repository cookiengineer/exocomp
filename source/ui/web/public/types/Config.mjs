
export const BootstrapConfig = (agent) => {

	return new Promise((resolve, reject) => {

		let url = "/api/session/config";

		if (agent !== "") {
			url = "/api/session/config/" + agent;
		}

		fetch(url, {
			method:  "GET",
			headers: {
				"Content-Type": "application/json"
			}
		}).then((response) => {

			if (response.ok) {

				response.json().then((data) => {
					resolve(Config.from(data));
				}).catch((err) => {
					reject(err);
				});

			} else {
				throw new Error(`HTTP error! Status ${response.status}`)
			}

		}).catch((err) => {
			reject(err);
		});

	});

};

export const Config = function() {

	this.Name        = "";
	this.Agent       = "";
	this.Model       = "";
	this.Prompt      = "";
	this.Temperature = 0.0;
	this.Playground  = "";
	this.Sandbox     = "";
	this.URL         = new URL("http://localhost:3000/");
	this.Debug       = false;

};

Config.from = (data) => {

	let config = new Config();

	config.Name        = data["name"]        || "";
	config.Agent       = data["agent"]       || "";
	config.Model       = data["model"]       || "";
	config.Prompt      = data["prompt"]      || "";
	config.Temperature = data["temperature"] || 0.0;
	config.Playground  = data["playground"]  || "";
	config.Sandbox     = data["sandbox"]     || "";
	config.URL         = new URL("/", window.location.origin);
	config.Debug       = data["debug"]       || false;

	return config;

};

Config.prototype = {

	GetPrompt: function() {
		return (this.Prompt || "").trim();
	},

	ResolveAPI: function(path) {

		let relative = new URL(path, "resolve://");

		return new URL(relative.pathname, this.URL);

	}

};

