
export const Message = function() {

	this.Role    = "";
	this.Content = "";
	this.Created = new Date();

};

Message.from = (data) => {

	let message = new Message();

	message.Role    = data["role"]    || "";
	message.Content = data["content"] || "";
	message.Created = new Date((data["created"] || "").replace(" ", "T"));

	return message;

};
