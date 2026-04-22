
export const GetModels = () => {

	return new Promise((resolve, reject) => {

		fetch("/api/parameters/models", {
			method:  "GET",
			headers: {
				"Content-Type": "application/json"
			}
		}).then((response) => {

			if (response.ok) {

				response.json().then((data) => {
					resolve(data);
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

