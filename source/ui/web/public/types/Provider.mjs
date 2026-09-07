
export const Provider = function() {

	this.URL   = new URL("http://0.0.0.0:0");
	this.Alias = "";
	this.Token = "";

};

Provider.from = (data) => {

	if (typeof data === "string") {
		try {
			data = JSON.parse(data);
		} catch (err) {
			data = {};
		}
	}

	if (Object.prototype.toString.call(data) !== "[object Object]") {
		data = {};
	}

	let provider = new Provider();

	provider.URL   = new URL(data["url"] || "http://0.0.0.0:0/v1");
	provider.Alias = data["alias"] || "";
	provider.Token = data["token"] || "";

	return provider;

};
