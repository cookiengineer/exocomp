
import { SanitizeContent } from "../../utils/fmt/SanitizeContent.mjs";

export const AnswerQuestions = function(element, config) {

	this.Config  = config;
	this.Element = element;
	this.elements = {
		"title":    element.querySelector("article > h3"),
		"question": element.querySelector("p[data-name=\"question\"]"),
		"answers":  element.querySelector("div[data-name=\"answers\"]"),
		"errors":   element.querySelector("div[data-name=\"errors\"]"),
	};

	this.questions = [];
	this.question = null;
	this.index    = 0;

	this.OnNext    = (data) => {}; // /types/Question
	this.OnConfirm = (data) => {}; // /types/Question
	this.OnCancel  = () => {};

	this.Init();

};

AnswerQuestions.prototype = {

	Error: function(errors) {

		let element = this.elements["errors"] || null;
		if (element !== null) {

			let items = errors.map((err) => {
				return "<b>" + err.toString() + "</b>";
			});

			element.innerHTML = items.join("");

		}

	},

	GetAnswer: function() {

		if (this.question === null) {
			return null;
		}

		if (this.IsChoice() === true) {

			let inputs = Array.from(this.elements["answers"].querySelectorAll("input[name=\"question-answer\"]:checked"));
			let values = inputs.map((input) => input.value);

			if (values.length > 0) {
				return values.join("\n");
			}

			return null;

		} else {

			let textarea = this.elements["answers"].querySelector("textarea[data-name=\"answer\"]");
			if (textarea !== null) {

				let value = textarea.value.trim();
				if (value !== "") {
					return value;
				}

			}

			return null;

		}

	},

	Hide: function() {

		if (this.Element !== null) {
			this.Element.close();
		}

	},

	Init: function() {

		if (this.Element !== null) {

			this.Element.addEventListener("cancel", (event) => {

				event.preventDefault();

				this.question = null;
				this.Hide();
				this.OnCancel();

			});

			let close = this.Element.querySelector("button[data-action=\"close\"]");
			if (close !== null) {
				close.onclick = () => {

					this.question = null;
					this.Hide();
					this.OnCancel();

				};
			}

			let confirm = this.Element.querySelector("button[data-action=\"confirm\"]");
			if (confirm !== null) {
				confirm.onclick = () => {

					let answer = this.GetAnswer();

					if (answer !== null && this.question !== null) {

						let data = {
							"question": this.question.Question,
							"answer":   answer,
						};

						if (this.IsLast() === true) {
							this.OnConfirm(data);
						} else {
							this.OnNext(data);
						}

					} else {
						this.Error([ new Error("Please select an answer or type a response.") ]);
					}

				};
			}

			let cancel = this.Element.querySelector("button[data-action=\"cancel\"]");
			if (cancel !== null) {
				cancel.onclick = () => {

					this.question = null;
					this.Hide();
					this.OnCancel();

				};
			}

		}

	},

	IsChoice: function() {

		if (this.question !== null) {
			return this.question.Type === "Choose";
		}

		return false;

	},

	IsLast: function() {

		return this.index >= this.questions.length - 1;

	},

	IsVisible: function() {

		if (this.Element !== null) {
			return this.Element.hasAttribute("open");
		}

		return false;

	},

	Next: function() {

		if (this.index < this.questions.length - 1) {

			this.index    += 1;
			this.question  = this.questions[this.index];

			this.Render();

		}

	},

	Render: function() {

		let question = this.question;

		if (question === null) {
			return;
		}

		if (this.elements["title"] !== null) {

			if (this.IsChoice() === true) {
				this.elements["title"].innerHTML = "Choices";
			} else {
				this.elements["title"].innerHTML = "Question";
			}

		}

		if (this.elements["question"] !== null) {
			this.elements["question"].innerHTML = SanitizeContent(question.Question || "");
		}

		if (this.elements["errors"] !== null) {
			this.elements["errors"].innerHTML = "";
		}

		if (this.elements["answers"] !== null) {

			if (this.IsChoice() === true) {

				let type    = question.Multiple === true ? "checkbox" : "radio";
				let options = Array.isArray(question.Options) ? question.Options : [];

				let html = options.map((option) => {
					return [
						"<label>",
						"<input type=\"" + type + "\" name=\"question-answer\" value=\"" + SanitizeContent(option) + "\">",
						"<span>" + SanitizeContent(option) + "</span>",
						"</label>"
					].join("");
				});

				this.elements["answers"].innerHTML = html.join("");

			} else {

				this.elements["answers"].innerHTML = [
					"<textarea data-name=\"answer\" placeholder=\"Your answer ...\"></textarea>"
				].join("");

			}

		}

	},

	Reset: function() {

		this.question  = null;
		this.questions = [];
		this.index     = 0;

		if (this.elements["answers"] !== null) {
			this.elements["answers"].innerHTML = "";
		}

		if (this.elements["errors"] !== null) {
			this.elements["errors"].innerHTML = "";
		}

	},

	Show: function(questions) {

		this.questions = Array.isArray(questions) ? questions : [];
		this.index     = 0;
		this.question  = this.questions.length > 0 ? this.questions[0] : null;

		this.Render();

		if (this.Element !== null && this.question !== null) {

			if (this.Element.hasAttribute("open") === false) {
				this.Element.showModal();
			}

		}

	}

};
