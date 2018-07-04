"use strict";
/**
 * Client for connecting to server
 * Auto correcting latency.
 */
var WsClient = /** @class */ (function () {
    function WsClient(wsAddr) {
        this.reconnectTimer = 0;
        this.reconnectTimeout = 1000;
        this.id = null;
        this.websocket = null;
        this.handlers = new Map();
        this.addr = wsAddr;
    }
    WsClient.prototype.connect = function () {
        if (this.reconnectTimer !== 0) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = 0;
        }
        if (this.websocket === null || this.websocket.readyState === WebSocket.CLOSED) {
            console.log('Connecting...');
            this.websocket = new WebSocket(this.addr);
            this.websocket.onopen = this.onopenHandler.bind(this);
            this.websocket.onclose = this.oncloseHandler.bind(this);
            this.websocket.onmessage = this.onmessageHandler.bind(this);
            this.websocket.onerror = this.onerrorHandler.bind(this);
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
                data: data
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
                if (handler != null) {
                    handler(eventType, data);
                }
                else {
                    console.warn("handler is null");
                }
            });
        }
    };
    WsClient.prototype.onopenHandler = function () {
        console.info('on open');
        this.trigger("open" /* Open */);
    };
    WsClient.prototype.oncloseHandler = function (ev) {
        console.info('on close');
        this.websocket = null;
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect.bind(this), this.reconnectTimeout);
        }
        this.trigger("close" /* Close */);
    };
    WsClient.prototype.onerrorHandler = function () {
        console.error('WebSocket error:', arguments);
        this.websocket = null;
        this.trigger("error" /* Error */);
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect.bind(this), this.reconnectTimeout);
        }
    };
    WsClient.prototype.onmessageHandler = function (event) {
        // console.log('on message');
        try {
            var wrapper = JSON.parse(event.data);
            // let data = JSON.parse(wrapper.data);
            var data = wrapper.data;
            console.log('%c%d (%s): %o', 'color: green; font-weight: bold;', wrapper.type, MessageType[wrapper.type], data);
            this.trigger("message" /* Message */, { type: wrapper.type, data: data });
            // onmessage.call(this, wrapper.MessageType, data);
        }
        catch (err) {
            console.error(err);
        }
    };
    return WsClient;
}());
//# sourceMappingURL=ws_client.js.map