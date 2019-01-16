"use strict";
var SERVER_ADDR = "ws://localhost:8080/ws";
var ServerMessage = /** @class */ (function () {
    function ServerMessage() {
    }
    return ServerMessage;
}());
var SimulationMode;
(function (SimulationMode) {
    SimulationMode[SimulationMode["Continuous"] = 0] = "Continuous";
    SimulationMode[SimulationMode["StepByStep"] = 1] = "StepByStep";
})(SimulationMode || (SimulationMode = {}));
var Application = /** @class */ (function () {
    function Application(canvas) {
        this.drawer = null;
        this.simulationMode = SimulationMode.Continuous;
        this.gameState = new GameState();
        this.canvas = canvas;
        this.drawer = new Drawer(canvas.getContext("2d"), 0, 0);
        // @ts-ignore
        this.input = new CanvasInputController(this.canvas, this.drawer, this.gameState, jQuery);
        this.client = new WsClient(SERVER_ADDR);
        this.protocol = new Protocol(this.client, this.gameState);
    }
    Application.prototype.start = function () {
        this.client.connect();
        if (this.simulationMode == SimulationMode.Continuous) {
            requestAnimationFrame(this.onAnimationFrame.bind(this));
        }
    };
    Application.prototype.processControlMessage = function () {
        // change simulation mode, login, death, reconnect, etc.
    };
    Application.prototype.onAnimationFrame = function () {
        this.do();
        if (this.simulationMode == SimulationMode.Continuous) {
            requestAnimationFrame(this.onAnimationFrame.bind(this));
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
    Application.prototype.getCurrentGameTime = function () {
        return 0;
    };
    Application.prototype.simulateToTime = function (time) {
    };
    Application.prototype.draw = function () {
        if (this.drawer != null) {
            this.drawer.setCanvasSize(this.canvas.width, this.canvas.height);
            this.drawer.draw();
        }
    };
    return Application;
}());
window.addEventListener('load', function (ev) {
    var canvas = window.document.getElementById("canvas");
    if (canvas != null) {
        var app = new Application(canvas);
        // @ts-ignore
        window.app = app;
        app.start();
    }
    else {
        console.error("Can't find canvas");
    }
});
//# sourceMappingURL=main.js.map