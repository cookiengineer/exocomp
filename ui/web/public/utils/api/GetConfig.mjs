
import { Config } from "../../types/Config.mjs";

export const GetConfig = () => {

	return new Promise((resolve, reject) => {

		fetch("/api/session/config", {
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
