const enum ClientEvent {
    Open = 'open',
    Close = 'close',
    Error = 'error',
    Message = 'message',
    TimeSynced = 'timeSynced'
}

type EventCallback = (eventType: string, data?: any) => void;

interface EventEmitter<T> {
    on(eventType: T, handler: EventCallback): void
    off(eventType: T, handler: EventCallback): void
}

interface Client extends EventEmitter<ClientEvent> {
    sendMessage(msgType: number, data?: any): void;
}