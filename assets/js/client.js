"use strict";

function WsClient(wsAddr, name, onmessage) {
	var websocket = false;
	var reconnectTimeout = 1;

	var _this = this;

	this.name = name;
	this.id = null;

	this.sendMessage = function(type, data) {
		var msg = {
			MessageType: type,
			Data: JSON.stringify(data),
			// Data: btoa(JSON.stringify(data)),
		}
		try {
			websocket.send(JSON.stringify(msg))
		} catch (err) {
			console.error(err);
		}
	};

	var onopenHandler = function() {
		_this.sendMessage(MessageType.AUTH, {Name: name});
	};

	var oncloseHandler = function() {
		websocket = null;
		this.id = null;
		window.setTimeout(connect, reconnectTimeout * 1000)
	};

	var onmessageHandler = function (event) {
		try {
			var wrapper = JSON.parse(event.data);
			var data = JSON.parse(wrapper.Data);
		}
		catch (err) {
			console.error(err);
		}

		if (wrapper.MessageType == MessageType.WELLCOME) {
			this.id = data.Id;
		}

		onmessage.call(this, wrapper.MessageType, data);
	};

	var connect = function() {
		if (!websocket) {
			websocket = new WebSocket(wsAddr);
			websocket.onopen = onopenHandler;
			websocket.onclose = oncloseHandler;
			websocket.onmessage = onmessageHandler;
		}
	};

	connect();
}

var MessageType = {
	AUTH:           1,
	WELLCOME:       2,
	LOGIN:          10,
	LOGOUT:         11,
	ERROR:          100,
	DATA:           1000,
	TEXT:           1001,
	USER_LIST:      10000,
	USER_LOGGEDIN:  10001,
	USER_LOGGEDOUT: 10002,
	SYNC_USERS_POSITIONS: 10003,
}