
export const Requirement = function() {

	this.Type          = "";
	this.File          = "";
	this.Symbol        = "";
	this.Declaration   = "";
	this.Behavior      = "";
	this.IsImplemented = false;

};

Requirement.from = (data) => {

	let requirement = new Requirement();

	requirement.Type          = data["type"]           || "";
	requirement.File          = data["file"]           || "";
	requirement.Symbol        = data["symbol"]         || "";
	requirement.Declaration   = data["declaration"]    || "";
	requirement.Behavior      = data["behavior"]       || "";
	requirement.IsImplemented = data["is_implemented"] || false;

	return requirement;

};

Requirement.prototype = {

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

		let checkbox = document.createElement("input");
		checkbox.setAttribute("type", "checkbox");
		checkbox.setAttribute("disabled", "");
		if (this.IsImplemented === true) {
			checkbox.setAttribute("checked", "");
		}

		let label = document.createElement("span");
		label.textContent = this.File.replace(/^\.\//, "") + "#" + this.Symbol + ": " + this.Declaration;

		node.appendChild(checkbox);
		node.appendChild(label);

		return node;

	}

};
