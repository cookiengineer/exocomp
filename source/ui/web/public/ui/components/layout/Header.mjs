
export const Init = () => {

	const header = document.querySelector("header");

	let path = window.location.pathname;

	if (path === "/") {
		path = "/index.html";
	}

	if (header !== null) {

		let list = header.querySelector("div ul");
		if (list !== null) {

			Array.from(list.childNodes).forEach((node) => {

				if (node.nodeType === Node.TEXT_NODE) {
					node.remove();
				}

			});

			Array.from(list.querySelectorAll("li")).forEach((item) => {

				let link = item.querySelector("a");
				if (link !== null) {

					let href = link.getAttribute("href");
					if (href === path) {
						item.setAttribute("data-state", "active");
					} else {
						item.removeAttribute("data-state");
					}

				}

			});

		}

	}

};
