enum MessageType {
    AUTH = 1,
    WELCOME = 2,
    LOGIN = 10,
    LOGOUT = 11,
    ERROR = 100,
    DATA = 1000,
    TEXT = 1001, // just text messages from server or other players

    USER_LIST = 10000,
    USER_LOGGEDIN = 10001,
    USER_LOGGEDOUT = 10002,

    // messages for game state
    SERVER_STATE = 20001,
    SYNC_TIME = 20002,

    SYNC_OBJECTS_PARTIAL = 20003,
    SYNC_OBJECTS_FULL = 20004,


    // special messages
    ACTION_MESSAGE = 1000000,
    SIMULATE_MESSAGE = 1000001,
    CHANGE_SIMULATION_MODE = 1000002
}

interface ServerState {
    simulation_by_step: boolean;
    simulation_step_time: number;
    simulation_step_real_time: number;
    simulation_time: number;
    server_time: number;
}
interface SyncTime {
    latency: number;
    correction: number;
}

interface MessageWrapper {
    type: MessageType
    data?: any
}

interface Welcome {
    id: number
}

class Protocol {
    private id: number = 0;
    private client: Client;
    private gameState: GameState;

    private lastSyncTimeRequest: number = 0;

    private latencies: number[] = [];
    private timeCorrections: number[] = [];
    private requestTimeTimer: number = 0;

    constructor(client: Client, gameState: GameState) {
        this.client = client;
        this.gameState = gameState;

        this.client.on(ClientEvent.Open, this.onOpen.bind(this));
        this.client.on(ClientEvent.Error, this.onError.bind(this));
        this.client.on(ClientEvent.TimeSynced, this.onTimeSynced.bind(this));
    }

    onOpen(eventType: string, data?: any): void {
        this.client.sendMessage(MessageType.AUTH, {name: randomName()});

        this.requestTime();
        this.requestTimeTimer = setInterval(this.requestTime.bind(this), 10000);
    }
    onError(eventType: string, data?: any): void {
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }
    }
    onClose(eventType: string, data?: any): void {
        if (this.requestTimeTimer !== 0) {
            clearInterval(this.requestTimeTimer);
            this.requestTimeTimer = 0;
        }
    }

    onTimeSynced(eventType: string, data?: any): void {
        console.log('onTimeSynced', arguments);
    }

    onMessage(eventType: string, wrapper: MessageWrapper): void {
        if (wrapper.type === MessageType.WELCOME) {
            this.id = wrapper.id;
        }
        if (wrapper.type === MessageType.SYNC_TIME) {
            this.syncTime(wrapper);
            return;
        }

        switch (wrapper.type) {
            case MessageType.WELCOME:
                this.id = (wrapper.data as Welcome).id;
        }
    }

    requestTime(): void {
        this.client.sendMessage(MessageType.SYNC_TIME, {id: generateUUID()});
        this.lastSyncTimeRequest = Date.now();
    }

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

        this.trigger(ClientEvent.TimeSynced, {latency: latency, correction: correction});
    }
}