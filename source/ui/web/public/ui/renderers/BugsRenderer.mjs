
export const BugsRenderer = function() {

	this.elements = {
		"grid": document.querySelector("body > main > div[data-name=\"bugs\"]")
	};

};

const getCollapsedPackagePath = (package_path) => {

	let segments = package_path.split("/").filter((segment) => segment !== "");

	if (segments.length > 2) {
		segments = segments.slice(0, 2);
	}

	return segments.join("/");

};

BugsRenderer.prototype = {

	Clear: function() {

		if (this.elements["grid"] !== null) {

			Array.from(this.elements["grid"].childNodes).forEach((node) => {
				node.parentNode.removeChild(node);
			});

		}

	},

	Render: function(bugs, collapsed) {

		if (this.elements["grid"] === null) {
			return;
		}

		this.Clear();

		let grouped = new Map();

		bugs.forEach((bug) => {

			let package_path = bug.getPackagePath();

			if (collapsed === true) {
				package_path = getCollapsedPackagePath(package_path);
			}

			if (grouped.has(package_path) === false) {
				grouped.set(package_path, []);
			}

			grouped.get(package_path).push(bug);

		});

		let package_paths = Array.from(grouped.keys()).sort();

		package_paths.forEach((package_path) => {
			this.elements["grid"].appendChild(this.RenderPackage(package_path, grouped.get(package_path)));
		});

	},

	RenderPackage: function(package_path, bugs) {

		let tile = document.createElement("article");
		tile.setAttribute("class", "bugs-package");
		tile.setAttribute("data-package", package_path);

		let heading = document.createElement("h3");
		heading.textContent = package_path !== "" ? package_path : ".";

		let list = document.createElement("ul");

		bugs.forEach((bug) => {
			list.appendChild(bug.toListItem());
		});

		tile.appendChild(heading);
		tile.appendChild(list);

		return tile;

	}

};
