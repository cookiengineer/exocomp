
export const Config = function() {

	this.Name        = "";
	this.Agent       = "";
	this.Debug       = false;
	this.Model       = "";
	this.Prompt      = "";
	this.Temperature = 0.0;
	this.Playground  = "";
	this.Sandbox     = "";
	this.URL         = new URL("http://localhost:3000/");

};

Config.from = (data) => {

	let config = new Config();

	config.Name        = data["name"]        || "";
	config.Agent       = data["agent"]       || "";
	config.Debug       = data["debug"]       || false;
	config.Model       = data["model"]       || "";
	config.Prompt      = data["prompt"]      || "";
	config.Temperature = data["temperature"] || 0.0;
	config.Playground  = data["playground"]  || "";
	config.Sandbox     = data["sandbox"]     || "";
	config.URL         = new URL("/", window.location.origin);

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

