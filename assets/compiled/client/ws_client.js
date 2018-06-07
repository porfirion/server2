"use strict";
/**
 * Client for connecting to server
 * Auto correcting latency.
 */
var WsClient = /** @class */ (function () {
    function WsClient(wsAddr, name) {
        this.reconnectTimer = 0;
        this.reconnectTimeout = 1000;
        this.id = null;
        this.websocket = null;
        this.handlers = new Map();
        this.lastSyncTimeRequest = 0;
        this.latencies = [];
        this.timeCorrections = [];
        this.requestTimeTimer = 0;
        this.addr = wsAddr;
        this.name = name;
    }
    WsClient.prototype.connect = function () {
        if (this.reconnectTimer !== 0) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = 0;
        }
        if (this.websocket === null) {
            console.log('Connecting...');
            this.websocket = new WebSocket(this.addr);
            this.websocket.onopen = this.onopenHandler;
            this.websocket.onclose = this.oncloseHandler;
            this.websocket.onmessage = this.onmessageHandler;
            this.websocket.onerror = this.onerrorHandler;
        }
        else {
            console.log('we are already connected to server');
        }
    };
    WsClient.prototype.sendMessage = function (type, data) {
        if (this.websocket == null) {
            console.error('no websocket connection');
            return;
        }
        var msg;
        if (typeof data !== 'undefined' && data != null) {
            msg = {
                type: type,
                data: JSON.stringify(data)
            };
        }
        else {
            msg = {
                type: type,
            };
        }
        try {
            this.websocket.send(JSON.stringify(msg));
        }
        catch (err) {
            console.error(err);
        }
    };
    WsClient.prototype.on = function (eventType, handler) {
        var handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            handlers.push(handler);
        }
        else {
            this.handlers.set(eventType, [handler]);
        }
        return this;
    };
    WsClient.prototype.off = function (eventType, handler) {
        var handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            var index = handlers.indexOf(handler);
            if (index > -1) {
                handlers.splice(index, 1);
            }
        }
        return this;
    };
    WsClient.prototype.trigger = function (eventType, data) {
        var handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            handlers.forEach(function (handler) {
                handler(eventType, data);
            });
        }
    };
    WsClient.prototype.onopenHandler = function () {
        console.log('on open');
        this.sendMessage(MessageType.AUTH, { name: name });
        this.requestTime();
        this.requestTimeTimer = setInterval(this.requestTime.bind(this), 10000);
        this.trigger("open" /* NotificationOpen */);
    };
    WsClient.prototype.oncloseHandler = function (ev) {
        console.log('on close');
        this.websocket = null;
        this.id = null;
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect, this.reconnectTimeout);
        }
        this.trigger("close" /* NotificationClose */);
    };
    WsClient.prototype.onerrorHandler = function () {
        console.log('WebSocket error:', arguments);
        this.websocket = null;
        this.id = null;
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }
        this.trigger("error" /* NotificationError */);
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect, this.reconnectTimeout);
        }
    };
    WsClient.prototype.onmessageHandler = function (event) {
        // console.log('on message');
        try {
            var wrapper = JSON.parse(event.data);
            var data = JSON.parse(wrapper.data);
            if (wrapper.type === MessageType.WELCOME) {
                this.id = data.id;
            }
            if (wrapper.type === MessageType.SYNC_TIME) {
                this.syncTime(data);
                return;
            }
            console.log('%c%d (%s): %o', 'color: green; font-weight: bold;', wrapper.type, MessageType[wrapper.type], data);
            this.trigger("message" /* NotificationMessage */, { type: wrapper.type, data: data });
            // onmessage.call(this, wrapper.MessageType, data);
        }
        catch (err) {
            console.error(err);
        }
    };
    WsClient.prototype.requestTime = function () {
        this.sendMessage(MessageType.SYNC_TIME, {});
        this.lastSyncTimeRequest = Date.now();
    };
    WsClient.prototype.syncTime = function (data) {
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
        this.trigger("timeSynced" /* NotificationTimeSynced */, { latency: latency, correction: correction });
    };
    return WsClient;
}());
var MessageType;
(function (MessageType) {
    MessageType[MessageType["AUTH"] = 1] = "AUTH";
    MessageType[MessageType["WELCOME"] = 2] = "WELCOME";
    MessageType[MessageType["LOGIN"] = 10] = "LOGIN";
    MessageType[MessageType["LOGOUT"] = 11] = "LOGOUT";
    MessageType[MessageType["ERROR"] = 100] = "ERROR";
    MessageType[MessageType["DATA"] = 1000] = "DATA";
    MessageType[MessageType["TEXT"] = 1001] = "TEXT";
    MessageType[MessageType["USER_LIST"] = 10000] = "USER_LIST";
    MessageType[MessageType["USER_LOGGEDIN"] = 10001] = "USER_LOGGEDIN";
    MessageType[MessageType["USER_LOGGEDOUT"] = 10002] = "USER_LOGGEDOUT";
    MessageType[MessageType["SYNC_OBJECTS_POSITIONS"] = 10003] = "SYNC_OBJECTS_POSITIONS";
    MessageType[MessageType["SYNC_TIME"] = 10004] = "SYNC_TIME";
    MessageType[MessageType["SERVER_STATE"] = 10005] = "SERVER_STATE";
    // special messages
    MessageType[MessageType["ACTION_MESSAGE"] = 1000000] = "ACTION_MESSAGE";
    MessageType[MessageType["SIMULATE_MESSAGE"] = 1000001] = "SIMULATE_MESSAGE";
    MessageType[MessageType["CHANGE_SIMULATION_MODE"] = 1000002] = "CHANGE_SIMULATION_MODE";
})(MessageType || (MessageType = {}));
//# sourceMappingURL=ws_client.js.map