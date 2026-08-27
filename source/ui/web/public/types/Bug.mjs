
export const Bug = function() {

	this.IsFixed     = false;
	this.File        = "";
	this.Symbol      = "";
	this.Description = "";

};

Bug.from = (data) => {

	let bug = new Bug();

	bug.IsFixed     = data["is_fixed"]     || false;
	bug.File        = data["file"]         || "";
	bug.Symbol      = data["symbol"]       || "";
	bug.Description = data["description"]  || "";

	return bug;

};

Bug.prototype = {

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
		if (this.IsFixed === true) {
			checkbox.setAttribute("checked", "");
		}

		let label = document.createElement("span");
		label.textContent = this.File.replace(/^\.\//, "") + "#" + this.Symbol + ": " + this.Description;

		node.appendChild(checkbox);
		node.appendChild(label);

		return node;

	}

};
