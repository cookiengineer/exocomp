
export const group_messages = (mapped) => {

	let grouped = {};

	let times = Object.keys(mapped).map((t) => {
		return {
			key: t,
			date: new Date(t)
		};
	}).sort((a, b) => {
		return a.date - b.date;
	});

	let current = [];
	let first = null;
	let last = null;

	times.forEach((time) => {

		if (current.length === 0) {

			current = [...mapped[time.key]];
			first = time.date;
			last = time.date;

		} else {

			// overlap threshold: 8 seconds
			if ((time.date - last) < 16000) {

				current.push(...mapped[time.key]);
				last = time.date;

			} else {

				let middle = new Date(
					(first.getTime() + last.getTime()) / 2
				);

				grouped[middle] = current;

				current = [...mapped[time.key]];
				first = time.date;
				last = time.date;

			}

		}

	});

	if (current.length > 0) {

		let middle = new Date(
			(first.getTime() + last.getTime()) / 2
		);

		grouped[middle] = current;

	}

	return grouped;

};

