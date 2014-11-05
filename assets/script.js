MessageTypeLogin    = 1;
MessageTypeWelcome  = 2;
MessageTypeForbidden = 3;
MessageTypeJoin     = 10;
MessageTypeLeave    = 11;

MessageTypeSyncMembers = 101;

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
	var msg = {
		MessageType: MessageTypeLogin,
		Data: JSON.stringify({
			UUID: uuid
		})
	}
	websocket.send(JSON.stringify(msg))
}
function oncloseHandler() {
	ShowMessage("server is down", "error");
}
function onmessageHandler(event) {
	var wrapper = JSON.parse(event.data);
	var data = JSON.parse(wrapper.Data);
	window.console.log("msg type " + wrapper.MessageType, data);
	if (wrapper.MessageType == MessageTypeWelcome) {
		ShowMessage("WELCOME!")
	}
	else if (wrapper.MessageType == MessageTypeJoin) {
		ShowMessage(data.UUID + " joined!", "join");
		NewMember(data.UUID)
	}
	else if (wrapper.MessageType == MessageTypeLeave) {
		if (data.UUID in members) {
			members[data.UUID].anchor.remove();	
			delete(members[data.UUID]);
		}

		ShowMessage(data.UUID + ' leaved!', 'leave');
	}
	else if (wrapper.MessageType == MessageTypeSyncMembers) {
		window.console.log('Synchronizing members...');
		for (var i = 0; i < data.Members.length; i++) {
			NewMember(data.Members[i]);
		}
	}
	else if (wrapper.MessageType == MessageTypeText) {
		if (data.UUID == uuid) {
			ShowMessage(" me: \"" + data.Text + "\"", "me");
		}
		else {
			ShowMessage(data.UUID + " says: \"" + data.Text + "\"");
		}
	}
	else {
		ShowMessage(event.data);
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

		SendMessage(MessageTypeText, {Text: $('.chat_input').val()});
		$('.chat_input').val('');
	})
});
