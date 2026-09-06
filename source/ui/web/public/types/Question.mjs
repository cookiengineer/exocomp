
export const Question = function() {

	this.Type     = "";
	this.Question = "";
	this.Options  = [];
	this.Multiple = false;
	this.Answer   = "";

};

Question.from = (data) => {

	let question = NewQuestion();

	question.Question = data["question"] || "";
	question.Options  = data["options"]  || [];
	question.Multiple = typeof data["multiple"] === "boolean" ? data["multiple"] : false;

	if (question.Question !== "" && question.Options.length > 0) {
		question.Type = "humans.Choose";
	} else {
		question.Type = "humans.Ask";
	}

	return question;

};

