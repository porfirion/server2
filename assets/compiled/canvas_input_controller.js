"use strict";
var SCALE_STEP = 1.05;
var CanvasInputController = /** @class */ (function () {
    function CanvasInputController(canvas, drawer, gameState) {
        this.canvas = canvas;
        this.drawer = drawer;
        this.gameState = gameState;
        this.initHandlers();
    }
    CanvasInputController.prototype.initHandlers = function () {
        $(this.canvas).on("mousedown", this.onMouseDown.bind(this));
        $(this.canvas).on('mousemove', this.onMouseMove.bind(this));
        $(this.canvas).on('mouseenter', this.onMouseEnter.bind(this));
        $(this.canvas).on('mouseout', this.onMouseOut.bind(this));
        $(this.canvas).on('mousewheel', this.onMouseWheel.bind(this));
        $(this.canvas).on('DOMMouseScroll', this.onMouseWheel.bind(this));
        $(this.canvas).on('contextmenu', this.onContextMenu.bind(this));
    };
    CanvasInputController.prototype.onMouseDown = function (event) {
        var canvasCoords = { x: event.offsetX, y: event.offsetY };
        var viewportCoords = this.drawer.getViewport().fromCanvas(canvasCoords);
        var realCoords = this.drawer.getViewport().toReal(viewportCoords);
        switch (event.button) {
            case 0:
                // левая кнопка мыши
                console.log('not implemented');
                break;
            case 1:
                // средняя кнопка мыши
                break;
            case 2:
                // правая кнопка мыши
                var newPos = { x: realCoords.x, y: realCoords.y };
                this.drawer.getViewport().setPos(newPos);
                break;
            default:
                console.warn("unexpected button " + event.button);
                break;
        }
        return false;
    };
    CanvasInputController.prototype.onMouseMove = function (event) {
        var canvasCoords = { x: event.offsetX, y: event.offsetY };
        var viewportCoords = this.drawer.getViewport().fromCanvas(canvasCoords);
        var realCoords = this.drawer.getViewport().toReal(viewportCoords);
    };
    CanvasInputController.prototype.onMouseEnter = function (event) {
    };
    CanvasInputController.prototype.onMouseOut = function (event) {
    };
    CanvasInputController.prototype.onMouseWheel = function (event) {
        if (typeof event.originalEvent !== 'undefined' && event.originalEvent instanceof MouseEvent) {
            var params = normalizeWheel(event.originalEvent);
            if (params.spinY > 0) {
                // на себя
                this.drawer.getViewport().scaleBy(1.0 / SCALE_STEP);
            }
            else {
                // от себя
                this.drawer.getViewport().scaleBy(SCALE_STEP);
            }
        }
        // capture all scrolling over map
        return false;
    };
    CanvasInputController.prototype.onContextMenu = function (event) {
        event.preventDefault();
        return false;
    };
    return CanvasInputController;
}());
//# sourceMappingURL=canvas_input_controller.js.map