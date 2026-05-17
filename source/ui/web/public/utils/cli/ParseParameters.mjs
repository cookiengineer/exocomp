
const seekNested = (raw, start, start_token, close_token) => {

	let result    = start_token;
	let depth     = 0;
	let in_string = false;
	let escaped   = false;

	for (let index = start + 1; index < raw.length; index++) {

		result += raw[index];

		let chr = raw[index];

		if (chr === '\\' && in_string && escaped === false) {
			escaped = true;
			continue;
		}

		if (chr === '"' && escaped === false) {
			in_string = !in_string;
		}

		if (in_string === false) {

			if (chr === start_token) {

				depth++;

			} else if (chr === close_token) {

				depth--;

				if (depth === 0) {
					break;
				}

			}

		}

		escaped = false;

	}

	return result;

};

const seekString = (raw, start, token) => {

	let result  = token;
	let escaped = false;

	for (let index = start + 1; index < raw.length; index++) {

		result += raw[index];

		let chr = raw[index];

		if (chr === '\\' && escaped === false) {
			escaped = true;
			continue;
		}

		if (chr === token && escaped === false) {
			break;
		}

		escaped = false;

	}

	return result;

};

const skipWhiteSpace = (raw, index) => {

	while (index < raw.length && /\s/.test(raw[index])) {
		index++
	}

	return index;

};

const parseValue = (buffer) => {

	let value = undefined;

	try {

		value = JSON.parse(buffer);

		return value;

	} catch (err) {

		if (buffer.length >= 2 && (buffer[0] === '-' || buffer[0] === '+')) {

			let num = parseInt(buffer, 10);

			if (Number.isNaN(num) === false) {
				return num;
			}

		}

		if (buffer.length >= 2 && buffer[0] === '"' && buffer[buffer.length - 1] === '"') {
			return buffer.slice(1, buffer.length - 1);
		}

		if (buffer.length >= 2 && buffer[0] === '\'' && buffer[buffer.length - 1] === '\'') {
			return buffer.slice(1, buffer.length - 1);
		}

		if (buffer === "true") {
			return true;
		}

		if (buffer === "false") {
			return false;
		}

		if (buffer === "null") {
			return null;
		}

		let num1 = Number.parseInt(buffer, 10);

		if (Number.isNaN(num1) === false) {
			return num1;
		}

		let num2 = Number.parseFloat(buffer)

		if (Number.isNaN(num2) === false) {
			return num2;
		}

	}

	return null;

};

export const ParseParameters = (raw) => {

	let parameters = {};
	let index      = 0;

	while (index < raw.length) {

		index = skipWhiteSpace(raw, index);

		if (index >= raw.length) {
			break
		}

		// key
		let key_start = index;

		while (index < raw.length && raw[index] !== '=') {
			index++;
		}

		if (index >= raw.length) {
			break
		}

		let key = raw.slice(key_start, index);

		// Skip '='
		index++

		index = skipWhiteSpace(raw, index);

		let value = ""

		switch (raw[index]) {

		case '{':
			value = seekNested(raw, index, '{', '}');
			break;

		case '[':
			value = seekNested(raw, index, '[', ']');
			break;

		case '\'':
			value = seekString(raw, index, '\'');
			break;

		case '"':
			value = seekString(raw, index, '"');
			break;

		default:
			index = skipWhiteSpace(raw, index);
			break;

		}

		if (value.length > 0) {
			parameters[key] = parseValue(value);
			index += value.length;
		}

	}

	return parameters;

};

