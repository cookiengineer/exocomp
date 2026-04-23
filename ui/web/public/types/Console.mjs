
export const Console = function() {

	this.Messages = [];

};

Console.prototype = {

	GetMessages: function(from) {

		let result = [];

		if (this.Messages.length > 0 && from < this.Messages.length) {

			for (let m = from; m < this.Messages.length; m++) {
				result.push(this.Messages[m]);
			}

		}

		return result;

	}

};
