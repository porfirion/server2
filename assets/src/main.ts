const SERVER_ADDR = "ws://localhost:8080/ws";

class ServerMessage {
}

enum SimulationMode {
    Continuous,
    StepByStep
}

class Application {
    private canvas: HTMLCanvasElement;
    private client: WsClient;
    private gameState: GameState;
    private drawer: Drawer | null = null;
    private simulationMode: SimulationMode = SimulationMode.Continuous;
    private input: CanvasInputController;
    private protocol: Protocol;

    constructor(canvas: HTMLCanvasElement) {
        this.gameState = new GameState();

        this.canvas = canvas;
        this.drawer = new Drawer(canvas.getContext("2d"), 0, 0);
        this.input = new CanvasInputController(this.canvas, this.drawer, this.gameState);
        this.client = new WsClient(SERVER_ADDR);

        this.protocol = new Protocol(this.client, this.gameState);
    }

    start() {
        this.client.connect();

        if (this.simulationMode == SimulationMode.Continuous) {
            requestAnimationFrame(this.onAnimationFrame.bind(this));
        }
    }


    processControlMessage(): void {
        // change simulation mode, login, death, reconnect, etc.
    }

    onAnimationFrame() {
        this.do();

        if (this.simulationMode == SimulationMode.Continuous) {
            requestAnimationFrame(this.onAnimationFrame.bind(this));
        }
    }

    do(): void {
        // get current game time
        let now = this.getCurrentGameTime();

        // process everything that we have to that time
        // (approximate if have next but not current)
        this.simulateToTime(now);

        // just show what we have on screen
        this.draw();
    }

    getCurrentGameTime(): number {
        return 0;
    }

    simulateToTime(time: number) {
    }

    draw(): void {
        if (this.drawer != null) {
            this.drawer.setCanvasSize(this.canvas.width, this.canvas.height);
            this.drawer.draw();
        }
    }
}

window.addEventListener('load', function (ev: Event) {
    let canvas = window.document.getElementById("canvas") as HTMLCanvasElement;
    if (canvas != null) {
        let app: Application = new Application(canvas);
        app.start();
    } else {
        console.error("Can't find canvas");
    }
});
