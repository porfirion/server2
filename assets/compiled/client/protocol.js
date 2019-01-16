"use strict";
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
    // messages for game state
    MessageType[MessageType["SERVER_STATE"] = 20001] = "SERVER_STATE";
    MessageType[MessageType["SYNC_TIME"] = 20002] = "SYNC_TIME";
    MessageType[MessageType["SYNC_OBJECTS_PARTIAL"] = 20003] = "SYNC_OBJECTS_PARTIAL";
    MessageType[MessageType["SYNC_OBJECTS_FULL"] = 20004] = "SYNC_OBJECTS_FULL";
    // special messages
    MessageType[MessageType["ACTION_MESSAGE"] = 1000000] = "ACTION_MESSAGE";
    MessageType[MessageType["SIMULATE_MESSAGE"] = 1000001] = "SIMULATE_MESSAGE";
    MessageType[MessageType["CHANGE_SIMULATION_MODE"] = 1000002] = "CHANGE_SIMULATION_MODE";
})(MessageType || (MessageType = {}));
var Protocol = /** @class */ (function () {
    function Protocol(client, gameState) {
        this.id = 0;
        this.lastSyncTimeRequest = 0;
        this.latencies = [];
        this.timeCorrections = [];
        this.requestTimeTimer = 0;
        this.client = client;
        this.gameState = gameState;
        this.client.on("open" /* Open */, this.onOpen.bind(this));
        this.client.on("error" /* Error */, this.onError.bind(this));
        this.client.on("timeSynced" /* TimeSynced */, this.onTimeSynced.bind(this));
    }
    Protocol.prototype.onOpen = function (eventType, data) {
        this.client.sendMessage(MessageType.AUTH, { name: randomName() });
        this.requestTime();
        this.requestTimeTimer = setInterval(this.requestTime.bind(this), 10000);
    };
    Protocol.prototype.onError = function (eventType, data) {
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }
    };
    Protocol.prototype.onClose = function (eventType, data) {
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }
    };
    Protocol.prototype.onTimeSynced = function (eventType, data) {
        console.log('onTimeSynced', arguments);
    };
    Protocol.prototype.onMessage = function (eventType, wrapper) {
        if (wrapper.type === MessageType.WELCOME) {
            this.id = wrapper.data.id;
        }
        if (wrapper.type === MessageType.SYNC_TIME) {
            this.syncTime(wrapper);
            return;
        }
        switch (wrapper.type) {
            case MessageType.WELCOME:
                this.id = wrapper.data.id;
        }
    };
    Protocol.prototype.requestTime = function () {
        this.client.sendMessage(MessageType.SYNC_TIME, { id: generateUUID() });
        this.lastSyncTimeRequest = Date.now();
    };
    Protocol.prototype.syncTime = function (data) {
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
        // TODO откуда это взялось?
        //this.trigger(ClientEvent.TimeSynced, {latency: latency, correction: correction});
    };
    return Protocol;
}());
//# sourceMappingURL=protocol.js.map