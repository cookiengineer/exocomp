
export const SanitizeContent = (raw) => {

	if (raw === null || raw === undefined) {
		raw = "";
	}

	let tmp = raw;

	tmp = tmp.replaceAll("&", "&amp;");
	tmp = tmp.replaceAll("<", "&lt;");
	tmp = tmp.replaceAll(">", "&gt;");
	tmp = tmp.replaceAll('"', "&quot;");
	tmp = tmp.replaceAll("'", "&#39;");

	return tmp;

};

