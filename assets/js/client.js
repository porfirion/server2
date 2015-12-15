"use strict";

function WsClient(wsAddr, name) {
	
	var reconnectTimeout = 1;

	var _this = this;

	this.name = name;
	this.id = null;
	this.websocket = false;
	this.handlers = {};

	this.lastSyncTimeRequest = null;

	this.latencies = [];
	this.timeCorrections = [];

	this.sendMessage = function(type, data) {
		if (this.websocket == null) {
			console.error('no websocket connection');
			return;
		}

		var msg = {
			type: type,
			// data: JSON.stringify(data),
			// Data: btoa(JSON.stringify(data)),
		}

		if (typeof data != 'undefined' && data != null) {
			msg.data = JSON.stringify(data);
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
		_this.sendMessage(MessageType.AUTH, {name: name});

		setInterval(this.requestTime.bind(this), 1000);

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
		// console.log('on message');
		try {
			var wrapper = JSON.parse(event.data);
			var data = JSON.parse(wrapper.data);
		}
		catch (err) {
			console.error(err);
		}

		if (wrapper.type == MessageType.WELLCOME) {
			this.id = data.id;
		}
		if (wrapper.type == MessageType.SYNC_TIME) {
			this.syncTime(data);
			return;
		}

		console.log('Message type: ', wrapper.type);

		_this.trigger('message', wrapper.type, data);
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

	this.requestTime = function() {
		this.sendMessage(MessageType.SYNC_TIME, {})
		this.lastSyncTimeRequest = Date.now()
	}

	this.syncTime = function(data) {
		var now = Date.now();
		var latency = now - this.lastSyncTimeRequest;
		var assumingNow = (data.time + latency / 3);
		var correction = Math.round(assumingNow - now);
		// console.log('sync time!\n sent        : %d\n received    : %d\n latency     : %d\n server time : %d\n correction  : %d\n assuming now: %f', 
		// 	this.lastSyncTimeRequest,
		// 	now,
		// 	latency,
		// 	data.time,
		// 	correction,
		// 	assumingNow
		// );

		// console.log('time correction : ' + correction);
		 
		this.latencies.push(latency);
		if (this.latencies.length > 10) {
			this.latencies.shift();
		}

		this.timeCorrections.push(correction);
		if (this.timeCorrections.length > 10) {
			this.timeCorrections.shift();
		}

		this.trigger('syncTime');
	}

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