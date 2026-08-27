
export const ChangelogEntry = function() {

	this.Date        = "";
	this.Type        = "";
	this.File        = "";
	this.Symbol      = "";
	this.Description = "";

};

ChangelogEntry.from = (data) => {

	let entry = new ChangelogEntry();

	entry.Date        = data["date"]        || "";
	entry.Type        = data["type"]        || "";
	entry.File        = data["file"]        || "";
	entry.Symbol      = data["symbol"]      || "";
	entry.Description = data["description"] || "";

	return entry;

};

ChangelogEntry.prototype = {

	getPackagePath: function() {

		let path = (this.File || "").trim();

		while (path.startsWith("./") === true) {
			path = path.substring(2);
		}

		let index = path.lastIndexOf("/");

		if (index !== -1) {
			path = path.substring(0, index);
		} else {
			path = "";
		}

		return path;

	},

	toListItem: function() {

		let node = document.createElement("li");

		let date = document.createElement("time");
		date.setAttribute("datetime", this.Date);
		date.textContent = (this.Date || "").substring(0, 10);

		let label = document.createElement("span");

		let marker = this.Type !== "" ? "[" + this.Type + "] " : "";
		label.textContent = marker + this.File.replace(/^\.\//, "") + "#" + this.Symbol + ": " + this.Description;

		node.appendChild(date);
		node.appendChild(label);

		return node;

	}

};
