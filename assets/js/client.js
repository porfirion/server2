"use strict";

/**
 * Client for connecting to server
 * Auto correcting latency.
 * @param {url} wsAddr WebSocket address of server
 * @param {string} name Name to use when login
 * @constructor
 */
function WsClient(wsAddr, name) {
	
	var reconnectTimeout = 1000;

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
			type: type
			// data: JSON.stringify(data),
			// Data: btoa(JSON.stringify(data)),
		};

		if (typeof data !== 'undefined' && data != null) {
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
	};

	this.off = function(eventType, handler) {
		if (eventType in _this.handlers) {
			if (_this.handlers[eventType].indexOf(handler) > -1) {
				_this.handlers[eventType].splice(_this.handlers[eventType].indexOf(handler), 1);
			}
		}
	};
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
	};

	var onopenHandler = function() {
		console.log('on open');
		_this.sendMessage(MessageType.AUTH, {name: name});

		_this.requestTime();
		setInterval(_this.requestTime.bind(this), 10000);

		_this.trigger(WsClient.NotificationOpen);
	};

	var oncloseHandler = function() {
		console.log('on close');
		_this.websocket = null;
		_this.id = null;
		setTimeout(_this.connect, reconnectTimeout);
		_this.trigger(WsClient.NotificationClose);
	};

	var onerrorHandler = function() {
		console.warn('WebSocket error:', arguments);
		_this.websocket = null;
		_this.id = null;
		_this.trigger(WsClient.NotificationError);

		console.log('Reconnecting...');
		setTimeout(connect, reconnectTimeout);
	};

	var onmessageHandler = function (event) {
		// console.log('on message');
		try {
			var wrapper = JSON.parse(event.data);
			var data = JSON.parse(wrapper.data);
		}
		catch (err) {
			console.error(err);
		}

		if (wrapper.type === MessageType.WELLCOME) {
			_this.id = data.id;
		}
		if (wrapper.type === MessageType.SYNC_TIME) {
			_this.syncTime(data);
			return;
		}

		console.log('%c' + wrapper.type + ' (' + getMessageType(wrapper.type) + ')', 'color: green; font-weight: bold;', data);

		_this.trigger(WsClient.NotificationMessage, wrapper.type, data);
		// onmessage.call(this, wrapper.MessageType, data);
	};

	this.connect = function() {
		if (!_this.websocket) {
			_this.websocket = new WebSocket(wsAddr);
			_this.websocket.onopen = onopenHandler.bind(_this);
			_this.websocket.onclose = oncloseHandler.bind(_this);
			_this.websocket.onmessage = onmessageHandler.bind(_this);
			_this.websocket.onerror = onerrorHandler.bind(_this);
		} else {
		    console.log('we are already connected to server');
        }
	};

	this.requestTime = function() {
		this.sendMessage(MessageType.SYNC_TIME, {});
		this.lastSyncTimeRequest = Date.now();
	};

	this.syncTime = function(data) {
		var now = Date.now();

		// время туда-обратно
		var latency = now - this.lastSyncTimeRequest;

		// предполагаемое текущее серверное время
		// делим на 3, потому что пакет должен был 1) дойти 2) обработаться 3) вернуться
		var assumingServerTime = (data.time + latency / 3);

		// задержка между временем на клиенте и временем на сервере
		var correction = Math.round(assumingServerTime - now);

		// console.log('sync time!\n sent        : %d\n received    : %d\n latency     : %d\n server time : %d\n correction  : %d\n assuming server time: %f',
		// 	this.lastSyncTimeRequest,
		// 	now,
		// 	latency,
		// 	data.time,
		// 	correction,
		// 	assumingServerTime
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

		this.trigger(WsClient.NotificationTimeSynced, latency, correction);
	};
}

WsClient.NotificationOpen = 'open';
WsClient.NotificationClose = 'close';
WsClient.NotificationError = 'error';
WsClient.NotificationMessage = 'message';
WsClient.NotificationTimeSynced = 'timeSynced';

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
	SYNC_OBJECTS_POSITIONS: 10003,
	SYNC_TIME: 10004,

	ACTION_MESSAGE: 1000000,
	SIMULATE_MESSAGE: 1000001,
	CHANGE_SIMULATION_MODE: 1000002
};

function getMessageType(messageTypeId) {
	for (var key in MessageType) {
		if (MessageType[key] === messageTypeId) {
			return key;
		}
	}

	return false;
}