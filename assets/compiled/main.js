"use strict";
var SERVER_ADDR = "localhost:8080";
var ServerMessage = /** @class */ (function () {
    function ServerMessage() {
    }
    return ServerMessage;
}());
var SimulationMode;
(function (SimulationMode) {
    SimulationMode[SimulationMode["Continious"] = 0] = "Continious";
    SimulationMode[SimulationMode["StepByStep"] = 1] = "StepByStep";
})(SimulationMode || (SimulationMode = {}));
var Application = /** @class */ (function () {
    function Application(ctx) {
        this.simulationMode = SimulationMode.Continious;
        this.drawer = new Drawer(ctx, 0, 0);
        this.client = new WsClient(SERVER_ADDR, randomName());
        this.gameState = new GameState(this.drawer);
    }
    Application.prototype.start = function () {
        this.client.on("partial_state_sync", this.gameState.processMessage);
        this.client.on("full_state_sync", this.gameState.processMessage);
        this.client.on("control_message", this.processControlMessage);
        this.client.connect();
        if (this.simulationMode == SimulationMode.Continious) {
            requestAnimationFrame(this.onAnimationFrame);
        }
    };
    Application.prototype.processControlMessage = function () {
        // change simulation mode, login, death, reconnect, etc.
    };
    Application.prototype.onAnimationFrame = function () {
        this.do();
        if (this.simulationMode = SimulationMode.Continious) {
            requestAnimationFrame(this.onAnimationFrame);
        }
    };
    Application.prototype.do = function () {
        // get current game time
        var now = this.getCurrentGameTime();
        // process everything that we have to that time
        // (approximate if have next but not current)
        this.simulateToTime(now);
        // just show what we have on screen
        this.draw();
    };
    Application.prototype.getCurrentGameTime = function () { return 0; };
    Application.prototype.simulateToTime = function (time) { };
    Application.prototype.draw = function () { };
    return Application;
}());
var GameObject = /** @class */ (function () {
    function GameObject() {
    }
    return GameObject;
}());
// decribes visible game region and whole game state
var GameState = /** @class */ (function () {
    function GameState(drawer) {
        this.visibleObjects = [];
        this.drawer = drawer;
    }
    GameState.prototype.processMessage = function (msg) {
        // add/remove/update visible objects (work with drawer)
        // adjust whole game state (day/night, victory, ...)
    };
    return GameState;
}());
var app = new Application(new CanvasRenderingContext2D());
app.start();
//# sourceMappingURL=main.js.map