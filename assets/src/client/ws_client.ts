"use strict";

/**
 * Client for connecting to server
 * Auto correcting latency.
 */
class WsClient implements EventEmitter<ClientEvent> {
    private addr: string;

    private reconnectTimer: number = 0;
    private reconnectTimeout: number = 1000;
    private id: number | null = null;
    private websocket: WebSocket | null = null;
    private handlers: Map<String, EventCallback[]> = new Map<String, EventCallback[]>();



    constructor(wsAddr: string) {
        this.addr = wsAddr;
    }

    connect(): void {
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
        } else {
            console.log('we are already connected to server');
        }
    }

    sendMessage(type: number, data?: any): void {
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
    }

    on(eventType: string, handler: EventCallback): this {
        let handlers = this.handlers.get(eventType);

        if (typeof handlers != 'undefined') {
            handlers.push(handler);
        } else {
            this.handlers.set(eventType, [handler]);
        }
        return this;
    }

    off(eventType: string, handler: EventCallback): this {
        let handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            let index = handlers.indexOf(handler);
            if (index > -1) {
                handlers.splice(index, 1);
            }

        }
        return this;
    }

    trigger(eventType: string, data?: any): void {
        let handlers = this.handlers.get(eventType);
        if (typeof handlers != 'undefined') {
            handlers.forEach((handler) => {
                handler(eventType, data);
            });
        }
    }

    onopenHandler(): void {
        console.info('on open');

        this.trigger(ClientEvent.Open);
    }

    oncloseHandler(ev: CloseEvent): void {
        console.info('on close');
        this.websocket = null;
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect.bind(this), this.reconnectTimeout);
        }
        this.trigger(ClientEvent.Close);
    }

    onerrorHandler() {
        console.error('WebSocket error:', arguments);
        this.websocket = null;

        this.trigger(ClientEvent.Error);
        if (this.reconnectTimer === 0) {
            this.reconnectTimer = setTimeout(this.connect.bind(this), this.reconnectTimeout);
        }
    }

    onmessageHandler(event: MessageEvent) {
        // console.log('on message');
        try {
            let wrapper = JSON.parse(event.data);
            let data = JSON.parse(wrapper.data);

            console.log('%c%d (%s): %o', 'color: green; font-weight: bold;', wrapper.type, MessageType[wrapper.type], data);

            this.trigger(ClientEvent.Message, {type: wrapper.type, data: data});
            // onmessage.call(this, wrapper.MessageType, data);
        } catch (err) {
            console.error(err);
        }
    }
}