MessageTypeAuth =    1;
MessageTypeData = 1000;
MessageTypeText = 1001;

var members = {};

function generateUUID(){
	var d = new Date().getTime();
	var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
		var r = (d + Math.random()*16)%16 | 0;
		d = Math.floor(d/16);
		return (c=='x' ? r : (r&0x7|0x8)).toString(16);
	});
	return uuid;
}

var websocket;
var uuid;

function onopenHandler() {
	SendMessage(MessageTypeAuth, {Uuid: uuid});
}
function oncloseHandler() {
	ShowMessage("server is down", "error");
}
function onmessageHandler(event) {
	var wrapper = JSON.parse(event.data);
	try {
		console.log(wrapper.Data);
		var data = JSON.parse(wrapper.Data);
	}
	catch (ex) {
		console.error(ex);
	}
	window.console.log("msg type " + wrapper.MessageType, data);
	
	if (wrapper.MessageType == MessageTypeText) {
		if (data.Sender == 0) {
			ShowMessage(" server: \"" + data.Text + "\"", "text-muted");
		}
		else if (data.Sender == uuid) {
			ShowMessage(" me: \"" + data.Text + "\"", "text-primary");
		}
		else {
			ShowMessage(data.Sender + " says: \"" + data.Text + "\"");
		}
	}
	else {
		ShowMessage("Error parsing " + event.data, "text-danger");
	}
}

function NewMember(uuid) {
	if (!(uuid in members)) {
		var member = {
			uuid: uuid,
			anchor: $('<div class="member">'+uuid+'</div>'),
		};
		$('.chat_members').append(member.anchor);
		members[uuid] = member;

		return member;
	}
	else {
		return members[uuid];
	}
}

function ShowMessage(text, messageType) {
	if (typeof messageType == 'undefined' || messageType == null) {
		messageType = "";
	}

	$('.chat_window').append('<div class="message ' + messageType + '">' + text + '</div>');
}

function SendMessage(type, data) {
	var msg = {
		MessageType: type,
		 Data: JSON.stringify(data)
		// Data: btoa(JSON.stringify(data))
	}
	websocket.send(JSON.stringify(msg))
}

jQuery(document).ready(function() {
	uuid = generateUUID();
	$('h1').html(uuid);
	websocket = new WebSocket("ws://" + window.location.host + "/ws");
	websocket.onopen = onopenHandler;
	websocket.onclose = oncloseHandler;
	websocket.onmessage = onmessageHandler;

	$('#chat_form').submit(function(event) {
		event.preventDefault();

		var text = $('.chat_input').val();
		try {
			SendMessage(MessageTypeText, {Text: text});
		} catch (err) {
			ShowMessage("Unable to send " + text, "text-danger");
		}

		$('.chat_input').val('');
	})
});
