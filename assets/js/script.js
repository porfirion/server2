"use strict";

function generateUUID(){
	var d = new Date().getTime();
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
		var r = (d + Math.random()*16)%16 | 0;
		d = Math.floor(d/16);
		return (c=='x' ? r : (r&0x7|0x8)).toString(16);
	});
}

function getName() {
	var names = [
		"Ivan", "Mikhail", "Ilya", "Sergey", "Alexander", "Egor", "Diman", "Alexey"
	];
	var suffix = [
		"Eagle eye", "Morning star", "Black hoof", "Three finger", "Yellow wolf", "Wasp nest", "Short bull", "Hard needle", "Fetid datura", "Curved horn"
	];

	var rndNameId = Math.round(Math.random() * (names.length - 1));
	var rndSuffixId = Math.round(Math.random() * (suffix.length - 1));
	// var rndNum = Math.round(Math.random() * 1000000);

	// console.log(rndNameId, rndNum, names[rndNameId]);

	return names[rndNameId] + " " + suffix[rndSuffixId];
}

var members = {};
var myName = null;
var myId = null;
var client = null;
var syncTimeTimer = null;
var map = null;

function onmessage(messageType, data) {
	switch (messageType) {
		case MessageType.TEXT:
			if (data.sender == 0) {
				showMessage(data.text, "text-primary");
			}
			else if (data.Sender == myId) {
				showMessage(data.text, "text-success");
			}
			else {
				var username = data.sender in members ? members[data.sender].name : 'Unknown sender';
				showMessage(username + ": " + data.text);
			}
			break;
		case MessageType.WELLCOME:
			myId = data.id;
			newMember(myId, myName);
			break;
		case MessageType.USER_LIST:
			for (var i = 0, user; user = data.users[i]; i++) {
				newMember(user.id, user.name);
			}
			break;
		case MessageType.USER_LOGGEDIN:
			showMessage(data.name + " logged in", "text-muted");
			newMember(data.id, data.name);
			break;
		case MessageType.USER_LOGGEDOUT:
			removeMember(data.id);
			break;
		case MessageType.SYNC_OBJECTS_POSITIONS:
			updateObjectsPositions(data.positions);
		 	break;
		case MessageType.ERROR:
			showMessage('Error: ' + data.description);
			break;
		default:
			for (var key in MessageType) {
                if (MessageType.hasOwnProperty(key)) {
                    if (MessageType[key] == messageType) {
                        showMessage('Not implemented ' + key.toUpperCase());
                        return;
                    }
                }
			}
			showMessage("Unknown message type: " + messageType + data, "text-danger");
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
	showMessage('disconnected');
}

function updateObjectsPositions(positions) {
	for (var objectId in positions) {
        if (positions.hasOwnProperty(objectId)) {
            map.updateObjectPosition(positions[objectId]);
        }
	}
}

function newMember(id, name) {
	if (!(id in members)) {
		console.log('adding player #' + id + ' ('+ name+ ')');
		var member = new Player(id, name);
		member.anchor = $('<div class="member" aria-hidden="true" data-id="' + id + '">'+name+'</div>');
		$('.chat_members').append(member.anchor);
		if (id == myId) {
			member.isMe = true;
			member.anchor.css('font-weight', 'bold');
		} else {
			member.isMe = false;
		}
		members[id] = member;

		map.addPlayer(member);

		return member;
	}
	else {
		return members[id];
	}
}
function removeMember(id) {
	if (id in members) {
		showMessage(members[id].name + " logged out");
		members[id].anchor.remove();
		map.removePlayer(id);
		delete members[id];
	}
}

function showMessage(text, messageType) {
	if (typeof messageType == 'undefined' || messageType == null) {
		messageType = "";
	}

	$('.chat_window').append('<div class="message ' + messageType + '">' + text + '</div>');
}


var Player = function(id, name) {
	this. id = id;
	this.name = name;

	this.state = {
		position: {x: 0, y: 0}
	};
};

Player.prototype.setPosition = function(position) {
	this.state.position = position;

	$(this).trigger('change.position');
};

jQuery(document).ready(function() {
	myName = getName();
	$('.playerName').html(myName);
	client = new WsClient("ws://" + window.location.host + "/ws", myName);
	client.on('message', onmessage);
	client.on('close', onclose);
	client.on('open', function() {
		syncTimeTimer = setInterval(function() {
			//client.sendMessage(MessageType.SYNC_TIME, {time: 0});
			//console.log('sent');
		}, 10000);
	});
	client.on(WsClient.NotificationTimeSynced, function(latency, timeCorrection) {
		$('.latency .value').html(client.latencies[client.latencies.length - 1]);
		// Коррекцию выбираем как среднее из последних полученных
		var currentCorrection = client.timeCorrections.reduce(function(sum, a) { return sum + a }, 0)/(client.timeCorrections.length||1);
		$('.timeCorrection .value').html(currentCorrection)
		map.latency = latency;
		map.timeCorrection = currentCorrection;
	});

	$('#chat_form').submit(function(event) {
		event.preventDefault();

		var inp = $('.chat_input');
		var text = inp.val();
		try {
			client.sendMessage(MessageType.TEXT, {Text: text});
		} catch (err) {
			showMessage("Unable to send " + text, "text-danger");
			console.error(err);
		}

		inp.val('');

		return false;
	});

	var elem = document.getElementById("map");
	var wrapper = document.getElementById('map-wrapper');
	// console.log(wrapper);
	//elem.width = wrapper.clientWidth;
	//elem.height = wrapper.clientHeight;

	map = new Map(elem);

	// запуск анимации, если она ещё не была начата
	$(document.body).on('click', '.drawButton', function() {
		map.draw();
		return false;
	});

	// центрирование вьюпорта (0:0)
	$(document.body).on('click', '.centrateButton', function() {
		map.viewport.x = 0;
		map.viewport.y = 0;
		map.viewportAdjustPoint = null;

		return false;
	});

	// перемещение вьюпорта при помощи кнопок навигации
	$(document).on('click', '.floatingButton', function() {
		var x = $(this).data('x');
		var y = $(this).data('y');

		map.viewport.x += Number(x) * map.viewport.scale * 20;
		map.viewport.y += Number(y) * map.viewport.scale * 20;
	});

	$(map).on('game:click', function(event, data) {
		console.log('clicked at ', data);

		client.sendMessage(MessageType.ACTION_MESSAGE, { 
			actionType: 'move',
			actionData: data,
		});
	});

	// setInterval(function() {
	// 	console.log('send last pos ', map.lastCursorPosition, map.lastCursorPositionReal);
	// 	if (map.lastCursorPositionReal != null) {
	// 		client.sendMessage(MessageType.ACTION_MESSAGE, {
	// 			actionType: 'accelerate',
	// 			actionData: map.lastCursorPositionReal,
	// 		});
	// 	}
	// }, 1000);
	
	map.draw();
});