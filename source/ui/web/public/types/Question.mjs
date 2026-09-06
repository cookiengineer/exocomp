
export const Question = function() {

	this.Type     = "";
	this.Question = "";
	this.Options  = [];
	this.Multiple = false;
	this.Answer   = "";

};

Question.from = (data) => {

	if (typeof data === "string") {
		try {
			data = JSON.parse(data);
		} catch (err) {
			data = {};
		}
	}

	if (Object.prototype.toString.call(data) !== "[object Object]") {
		data = {};
	}

	let question = new Question();

	question.Question = data["question"] || "";
	question.Options  = data["options"]  || [];
	question.Multiple = typeof data["multiple"] === "boolean" ? data["multiple"] : false;

	if (question.Question !== "" && question.Options.length > 0) {
		question.Type = "Choose";
	} else {
		question.Type = "Ask";
	}

	return question;

};

