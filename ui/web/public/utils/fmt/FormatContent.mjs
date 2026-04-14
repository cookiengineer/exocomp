
export const FormatContent = (raw) => {

	if (raw === null || raw === undefined) {
		raw = "";
	}

	let tmp = raw;

	tmp = tmp.replaceAll("&", "&amp;");
	tmp = tmp.replaceAll("<", "&lt;");
	tmp = tmp.replaceAll(">", "&gt;");

	tmp = tmp.replaceAll(/```(\w+)?\n([\s\S]*?)```/g, (_, lang, code) => {

		if (lang.trim() !== "") {
			return `<pre class="${lang.trim()}">${code.trim()}</pre>`;
		} else {
			return `<pre>${code.trim()}</pre>`;
		}

    });

	tmp = tmp.replaceAll(/```([\s\S]*?)```/g, (_, code) => {
		return `<pre><code>${code.trim()}</code></pre>`;
	});

	tmp = tmp.replaceAll(/`([^\`]+?)`/g, "<code>$1</code>");

	tmp = tmp.replaceAll(/\*\*(.*?)\*\*/g, "<strong>$1</strong>");
	tmp = tmp.replaceAll(/__(.*?)__/g, "<strong>$1</strong>");

	tmp = tmp.replaceAll(/\*(.*?)\*/g, "<em>$1</em>");
	tmp = tmp.replaceAll(/_(.*?)_/g, "<em>$1</em>");

	tmp = tmp.replaceAll(/\+\+(.*?)\+\+/g, "<u>$1</u>");
	tmp = tmp.replaceAll(/~~(.*?)~~/g, "<del>$1</del>");

	tmp = tmp.replaceAll(/^###### (.*)$/gm, "<h6>$1</h6>");
	tmp = tmp.replaceAll(/^##### (.*)$/gm, "<h5>$1</h5>");
	tmp = tmp.replaceAll(/^#### (.*)$/gm, "<h4>$1</h4>");
	tmp = tmp.replaceAll(/^### (.*)$/gm, "<h3>$1</h3>");
	tmp = tmp.replaceAll(/^## (.*)$/gm, "<h2>$1</h2>");
	tmp = tmp.replaceAll(/^# (.*)$/gm, "<h1>$1</h1>");

	tmp = tmp.replaceAll(/\[([^\]]+)\]\(([^)]+)\)/g, "<a href=\"$2\" target=\"_blank\" rel=\"noopener noreferrer\">$1</a>");

	tmp = tmp.replaceAll(/^\s*[-*] (.*)$/gm, "<li>$1</li>");
	tmp = tmp.replaceAll(/^\s*\d+\. (.*)$/gm, "<li>$1</li>");
	tmp = tmp.replaceAll(/(<li>.*<\/li>)/gs, "<ul>$1</ul>");

	tmp = tmp.replaceAll("\n", "<br>");

	return tmp;

};

