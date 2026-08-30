
import { SanitizeContent } from "../../utils/fmt/SanitizeContent.mjs";

export const Question = function(element, config) {

	this.Config  = config;
	this.Element = element;
	this.elements = {
		"title":    element.querySelector("article > h3"),
		"question": element.querySelector("p[data-name=\"question\"]"),
		"answers":  element.querySelector("div[data-name=\"answers\"]"),
		"errors":   element.querySelector("div[data-name=\"errors\"]"),
	};

	this.question = null;

	this.OnConfirm = (data) => {};
	this.OnCancel  = (id) => {};

	this.Init();

};

Question.prototype = {

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

		if (this.question.kind === "ask") {

			let textarea = this.elements["answers"].querySelector("textarea[data-name=\"answer\"]");
			if (textarea !== null) {

				let value = textarea.value.trim();
				if (value !== "") {
					return value;
				}

			}

			return null;

		} else {

			let inputs = Array.from(this.elements["answers"].querySelectorAll("input[name=\"question-answer\"]:checked"));
			let values = inputs.map((input) => input.value);

			if (values.length > 0) {
				return values.join("\n");
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

				let id = this.question !== null ? this.question.id : null;
				this.question = null;
				this.Hide();
				this.OnCancel(id);

			});

			let close = this.Element.querySelector("button[data-action=\"close\"]");
			if (close !== null) {
				close.onclick = () => {

					let id = this.question !== null ? this.question.id : null;
					this.question = null;
					this.Hide();
					this.OnCancel(id);

				};
			}

			let confirm = this.Element.querySelector("button[data-action=\"confirm\"]");
			if (confirm !== null) {
				confirm.onclick = () => {

					let answer = this.GetAnswer();

					if (answer !== null && this.question !== null) {

						let id = this.question.id;
						this.OnConfirm({ id: id, answer: answer });

					} else {
						this.Error([ new Error("Please select an answer or type a response.") ]);
					}

				};
			}

			let cancel = this.Element.querySelector("button[data-action=\"cancel\"]");
			if (cancel !== null) {
				cancel.onclick = () => {

					let id = this.question !== null ? this.question.id : null;
					this.question = null;
					this.Hide();
					this.OnCancel(id);

				};
			}

		}

	},

	IsVisible: function() {

		if (this.Element !== null) {
			return this.Element.hasAttribute("open");
		}

		return false;

	},

	Render: function() {

		let question = this.question;

		if (question === null) {
			return;
		}

		if (this.elements["title"] !== null) {

			if (question.kind === "choices") {
				this.elements["title"].innerHTML = "Choices";
			} else {
				this.elements["title"].innerHTML = "Question";
			}

		}

		if (this.elements["question"] !== null) {
			this.elements["question"].innerHTML = SanitizeContent(question.question || "");
		}

		if (this.elements["errors"] !== null) {
			this.elements["errors"].innerHTML = "";
		}

		if (this.elements["answers"] !== null) {

			if (question.kind === "ask") {

				this.elements["answers"].innerHTML = [
					"<textarea data-name=\"answer\" placeholder=\"Your answer ...\"></textarea>"
				].join("");

			} else {

				let type    = question.multiple === true ? "checkbox" : "radio";
				let options = Object.prototype.toString.call(question.options) === "[object Array]" ? question.options : [];

				let html = options.map((option) => {
					return [
						"<label>",
						"<input type=\"" + type + "\" name=\"question-answer\" value=\"" + SanitizeContent(option) + "\">",
						"<span>" + SanitizeContent(option) + "</span>",
						"</label>"
					].join("");
				});

				this.elements["answers"].innerHTML = html.join("");

			}

		}

	},

	Reset: function() {

		this.question = null;

		if (this.elements["answers"] !== null) {
			this.elements["answers"].innerHTML = "";
		}

		if (this.elements["errors"] !== null) {
			this.elements["errors"].innerHTML = "";
		}

	},

	Show: function(question) {

		this.question = question;
		this.Render();

		if (this.Element !== null) {

			if (this.Element.hasAttribute("open") === false) {
				this.Element.showModal();
			}

		}

	}

};
