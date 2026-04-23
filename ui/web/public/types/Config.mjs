
export const Config = function() {

	this.Name        = "";
	this.Agent       = "";
	this.Debug       = false;
	this.Model       = "";
	this.Playground  = "";
	this.Prompt      = "";
	this.Sandbox     = "";
	this.Temperature = 0.0;
	this.URL         = new URL("http://localhost:3000/");

};

Config.from = (data) => {

	let config = new Config();

	config.Name        = data["name"]        || "";
	config.Agent       = data["agent"]       || "";
	config.Debug       = data["debug"]       || false;
	config.Model       = data["model"]       || "";
	config.Playground  = data["playground"]  || "";
	config.Prompt      = data["prompt"]      || "";
	config.Sandbox     = data["sandbox"]     || "";
	config.Temperature = data["temperature"] || 0.0;
	config.URL         = new URL(data["url"] || "http://localhost:3000/");

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

