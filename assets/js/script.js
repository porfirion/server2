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

function getName() {
	var names = [
		"Ivan", "Mikhail", "Ilya", "Sergey", "Alexander", "Egor", "Diman", "Alexey"
	];

	var rndNameId = Math.round(Math.random() * names.length - 1);
	var rndNum = Math.round(Math.random() * 1000000);

	// console.log(rndNameId, rndNum, names[rndNameId]);

	return names[rndNameId] + rndNum;
}


var myName = null;
var myId = null;
var client = null;
var syncTimeTimer = null;


function onmessage(messageType, data) {
	switch (messageType) {
		case MessageType.TEXT:
			if (data.Sender == 0) {
				ShowMessage(data.Text, "text-primary");
			}
			else if (data.Sender == myId) {
				ShowMessage(data.Text, "text-success");
			}
			else {
				var username = data.Sender in members ? members[data.Sender].name : 'Unknown';
				ShowMessage(username + ": " + data.Text);
			}
			break;
		case MessageType.WELLCOME:
			myId = data.Id;
			NewMember(myId, myName)
			break;
		case MessageType.USER_LIST:
			for (var i = 0, user; user = data.Users[i]; i++) {
				NewMember(user.Id, user.Name);
			}
			break;
		case MessageType.USER_LOGGEDIN:
			ShowMessage(data.Name + " logged in", "text-muted");
			NewMember(data.Id, data.Name);
			break;
		case MessageType.USER_LOGGEDOUT:
			RemoveMember(data.Id);
			break;
		case MessageType.SYNC_USERS_POSITIONS:
			ShowMessage("Unimplemented sync users positions");
			break;
		default:
			ShowMessage("Unknown message type: " + messageType + data, "text-danger");
			break;
	}
}

function onclose() {
	console.log('Timer: ', syncTimeTimer);

	if (syncTimeTimer) {
		clearInterval(syncTimeTimer);
	}
	
	syncTimeTimer = null;
	$('.chat_members').empty();
	ShowMessage('disconnected');
}

function NewMember(id, name) {
	if (!(id in members)) {
		var member = {
			id: id,
			name: name,
			anchor: $('<div class="member" aria-hidden="true" data-id="' + id + '">'+name+'</div>'),
		};
		$('.chat_members').append(member.anchor);
		if (id == myId) {
			member.anchor.css('font-weight', 'bold');
		}
		members[id] = member;

		return member;
	}
	else {
		return members[id];
	}
}
function RemoveMember(id) {
	if (id in members) {
		ShowMessage(members[id].name + " logged out");
		members[id].anchor.remove();
		delete members[id];
	}
}

function ShowMessage(text, messageType) {
	if (typeof messageType == 'undefined' || messageType == null) {
		messageType = "";
	}

	$('.chat_window').append('<div class="message ' + messageType + '">' + text + '</div>');
}

jQuery(document).ready(function() {
	myName = getName();
	$('h1').html(myName);
	client = new WsClient("ws://" + window.location.host + "/ws", myName);
	client.on('message', onmessage);
	client.on('close', onclose);
	client.on('open', function() {
		syncTimeTimer = setInterval(function() {
			client.sendMessage(MessageType.SYNC_TIME, {time: 0});
			console.log('sent');
		}, 1000);
	});

	$('#chat_form').submit(function(event) {
		event.preventDefault();

		var text = $('.chat_input').val();
		try {
			client.sendMessage(MessageType.TEXT, {Text: text});
		} catch (err) {
			ShowMessage("Unable to send " + text, "text-danger");
			console.error(err);
		}

		$('.chat_input').val('');

		return false;
	})
});
