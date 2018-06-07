const SERVER_ADDR = "localhost:8080";

class ServerMessage {}

enum SimulationMode {
    Continious,
    StepByStep
}

class Application {
    private client: WsClient;
    private gameState: GameState;
    private drawer: Drawer;
    private simulationMode: SimulationMode = SimulationMode.Continious;

    constructor(ctx: CanvasRenderingContext2D) {
        this.drawer = new Drawer(ctx, 0, 0);
        this.client = new WsClient(SERVER_ADDR, randomName());
        this.gameState = new GameState(this.drawer);
    }

    start() {
        this.client.on("partial_state_sync", this.gameState.processMessage);
        this.client.on("full_state_sync", this.gameState.processMessage);
        this.client.on("control_message", this.processControlMessage);
        this.client.connect();


        if (this.simulationMode == SimulationMode.Continious) {
            requestAnimationFrame(this.onAnimationFrame);
        }
    }



    processControlMessage(): void {
        // change simulation mode, login, death, reconnect, etc.
    }

    onAnimationFrame() {
        this.do();

        if (this.simulationMode = SimulationMode.Continious) {
            requestAnimationFrame(this.onAnimationFrame);
        }
    }

    do():void {
        // get current game time
        let now = this.getCurrentGameTime();

        // process everything that we have to that time
        // (approximate if have next but not current)
        this.simulateToTime(now);

        // just show what we have on screen
        this.draw();
    }

    getCurrentGameTime(): number { return 0; }

    simulateToTime(time: number) {}

    draw(): void {}
}

class GameObject {}

// decribes visible game region and whole game state
class GameState {
    visibleObjects: GameObject[] = [];
    private drawer: Drawer;

    constructor(drawer: Drawer) {
        this.drawer = drawer;
    }

    processMessage(msg: ServerMessage): void{
        // add/remove/update visible objects (work with drawer)
        // adjust whole game state (day/night, victory, ...)
    }
}

let app:Application = new Application(new CanvasRenderingContext2D());
app.start();


