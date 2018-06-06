"use strict";

type EventCallback = (eventType:string, data?:any) => void;

interface EventEmitter<T> {
    on(eventType: T, handler: EventCallback):void
    off(eventType: T, handler: EventCallback):void
}

/**
 * Client for connecting to server
 * Auto correcting latency.
 */
class WsClient implements EventEmitter<WsClientEvent> {
    private addr: string;
    private name: string;

    private reconnectTimer: number = 0;
    private reconnectTimeout: number = 1000;
    private id:number | null = null;
    private websocket:WebSocket | null = null;
    private handlers: Map<String, EventCallback[]> = new Map<String, EventCallback[]>();

    private lastSyncTimeRequest:number | null = null;

    private latencies: number[] = [];
    private timeCorrections: number[] = [];
    private requestTimeTimer:number = 0;

    constructor(wsAddr:string, name: string) {
        this.addr = wsAddr;
        this.name = name;
    }

    connect(): void {
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
        } else {
            console.log('we are already connected to server');
        }
    };

    sendMessage(type: number, data?: object): void {
        if (this.websocket == null) {
            console.error('no websocket connection');
            return;
        }

        let msg: object;

        if (typeof data !== 'undefined' && data != null) {
            msg = {
                type: type,
                data: JSON.stringify(data)
            };
        } else {
            msg = {
                type: type,
            };
        }

        try {
            this.websocket.send(JSON.stringify(msg))
        } catch (err) {
            console.error(err);
        }
    };

    on(eventType: string, handler:EventCallback): this {
        let handlers = this.handlers.get(eventType);

        if (typeof handlers != 'undefined') {
            handlers.push(handler);
        } else {
            this.handlers.set(eventType, [handler]);
        }
        return this;
    };

    off(eventType: string, handler: EventCallback): this {
        let handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            let index = handlers.indexOf(handler);
            if (index > -1) {
                handlers.splice(index, 1);
            }

        }
        return this;
    };

    trigger (eventType: string, data?:any) {
        let handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            handlers.forEach((handler) => {
                handler(eventType, data);
            });
        }
    };

    onopenHandler(): void {
        console.log('on open');
        this.sendMessage(MessageType.AUTH, {name: name});

        this.requestTime();
        this.requestTimeTimer = setInterval(this.requestTime.bind(this), 10000);

        this.trigger(WsClientEvent.NotificationOpen);
    };

    oncloseHandler(): void {
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
        this.trigger(WsClientEvent.NotificationClose);
    };

    onerrorHandler() {
        console.log('WebSocket error:', arguments);
        this.websocket = null;
        this.id = null;
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }

        this.trigger(WsClientEvent.NotificationError);
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect, this.reconnectTimeout);
        }
    };

    onmessageHandler(event) {
        // console.log('on message');
        try {
            let wrapper = JSON.parse(event.data);
            let data = JSON.parse(wrapper.data);
        } catch (err) {
            console.error(err);
        }

        if (wrapper.type === MessageType.WELCOME) {
            this.id = data.id;
        }
        if (wrapper.type === MessageType.SYNC_TIME) {
            this.syncTime(data);
            return;
        }

        console.log('%c%d (%s): %o', 'color: green; font-weight: bold;', wrapper.type, getMessageType(wrapper.type), data);

        this.trigger(WsClientEvent.NotificationMessage, wrapper.type, data);
        // onmessage.call(this, wrapper.MessageType, data);
    };

    requestTime(): void {
        this.sendMessage(MessageType.SYNC_TIME, {});
        this.lastSyncTimeRequest = Date.now();
    };

    syncTime(data: any): void {
        let now = Date.now();

        // время туда-обратно
        let latency = now - this.lastSyncTimeRequest;

        // предполагаемое текущее серверное время
        // делим на 3, потому что пакет должен был 1) дойти 2) обработаться 3) вернуться
        let assumingServerTime = (data.time + latency / 3);

        // задержка между временем на клиенте и временем на сервере
        let correction = Math.round(assumingServerTime - now);

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

        this.trigger(WsClientEvent.NotificationTimeSynced, latency, correction);
    };
}

const enum WsClientEvent {
    NotificationOpen = 'open',
    NotificationClose = 'close',
    NotificationError = 'error',
    NotificationMessage = 'message',
    NotificationTimeSynced = 'timeSynced'
}

var MessageType = {
    AUTH: 1,
    WELCOME: 2,
    LOGIN: 10,
    LOGOUT: 11,
    ERROR: 100,
    DATA: 1000,
    TEXT: 1001,
    USER_LIST: 10000,
    USER_LOGGEDIN: 10001,
    USER_LOGGEDOUT: 10002,
    SYNC_OBJECTS_POSITIONS: 10003,
    SYNC_TIME: 10004,
    SERVER_STATE: 10005,
    // special messages
    ACTION_MESSAGE: 1000000,
    SIMULATE_MESSAGE: 1000001,
    CHANGE_SIMULATION_MODE: 1000002
};

/**
 * Возвращает название типа сообщения по его идентификатору
 * @param messageTypeId {number}
 * @returns {String}
 */
function getMessageType(messageTypeId) {
    for (var key in MessageType) {
        if (MessageType[key] === messageTypeId) {
            return key;
        }
    }

    return false;
}