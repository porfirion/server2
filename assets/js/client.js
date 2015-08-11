"use strict";

function WsClient(wsAddr, name) {
	
	var reconnectTimeout = 1;

	var _this = this;

	this.name = name;
	this.id = null;
	this.websocket = false;
	this.handlers = {};

	this.sendMessage = function(type, data) {
		var msg = {
			MessageType: type,
			Data: JSON.stringify(data),
			// Data: btoa(JSON.stringify(data)),
		}
		try {
			this.websocket.send(JSON.stringify(msg))
		} catch (err) {
			console.error(err);
		}
	};

	this.on = function(eventType, handler) {
		if (!(eventType in _this.handlers)) {
			_this.handlers[eventType] = [];
		}
		_this.handlers[eventType].push(handler);

		return _this;
	}

	this.off = function(eventType, handler) {
		if (eventType in _this.handlers) {
			if (_this.handlers[eventType].indexOf(handler) > -1) {
				_this.handlers[eventType].splice(_this.handlers[eventType].indexOf(handler), 1);
			}
		}
	}
	this.trigger = function(eventType) {
		if (eventType in _this.handlers) {
			for (var i = 0; i < _this.handlers[eventType].length; i++) {
				if (arguments.length > 1) {
					_this.handlers[eventType][i].apply(null, Array.prototype.slice.call(arguments, 1));
				} else {
					_this.handlers[eventType][i].call(null);
				}
			}
		}
	}

	var onopenHandler = function() {
		console.log('on open');
		_this.sendMessage(MessageType.AUTH, {Name: name});

		_this.trigger('open');
	};

	var oncloseHandler = function() {
		console.log('on close');
		_this.websocket = null;
		_this.id = null;
		// window.setTimeout(connect, reconnectTimeout * 1000);
		_this.trigger('close');
	};

	var onerrorHandler = function() {
		console.warn('WebSocket error:', arguments);
		//_this.websocket = null;
		_this.trigger('error');

		// console.log('Reconnecting...');
		// connect()
	}

	var onmessageHandler = function (event) {
		console.log('on message');
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

		console.log('Message type: ', wrapper.MessageType);

		_this.trigger('message', wrapper.MessageType, data);
		// onmessage.call(this, wrapper.MessageType, data);
	};

	var connect = function() {
		if (!_this.websocket) {
			_this.websocket = new WebSocket(wsAddr);
			_this.websocket.onopen = onopenHandler.bind(_this);
			_this.websocket.onclose = oncloseHandler.bind(_this);
			_this.websocket.onmessage = onmessageHandler.bind(_this);
			_this.websocket.onerror = onerrorHandler.bind(_this);
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
	SYNC_TIME: 10004,
}