
export const ChangelogRenderer = function() {

	this.elements = {
		"grid": document.querySelector("body > main > div[data-name=\"changelog\"]")
	};

};

const getCollapsedPackagePath = (package_path) => {

	let segments = package_path.split("/").filter((segment) => segment !== "");

	if (segments.length > 2) {
		segments = segments.slice(0, 2);
	}

	return segments.join("/");

};

ChangelogRenderer.prototype = {

	Clear: function() {

		if (this.elements["grid"] !== null) {

			Array.from(this.elements["grid"].childNodes).forEach((node) => {
				node.parentNode.removeChild(node);
			});

		}

	},

	Render: function(entries, collapsed) {

		if (this.elements["grid"] === null) {
			return;
		}

		this.Clear();

		let grouped = new Map();

		entries.forEach((entry) => {

			let package_path = entry.getPackagePath();

			if (collapsed === true) {
				package_path = getCollapsedPackagePath(package_path);
			}

			if (grouped.has(package_path) === false) {
				grouped.set(package_path, []);
			}

			grouped.get(package_path).push(entry);

		});

		let package_paths = Array.from(grouped.keys()).sort();

		package_paths.forEach((package_path) => {
			this.elements["grid"].appendChild(this.RenderPackage(package_path, grouped.get(package_path)));
		});

	},

	RenderPackage: function(package_path, entries) {

		let tile = document.createElement("article");
		tile.setAttribute("class", "changelog-package");
		tile.setAttribute("data-package", package_path);

		let heading = document.createElement("h3");
		heading.textContent = package_path !== "" ? package_path : ".";

		let list = document.createElement("ul");

		entries.forEach((entry) => {
			list.appendChild(entry.toListItem());
		});

		tile.appendChild(heading);
		tile.appendChild(list);

		return tile;

	}

};
