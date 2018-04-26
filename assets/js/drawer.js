"use strict";
var DrawMode;
(function (DrawMode) {
    DrawMode[DrawMode["ONLY_SERVER"] = 1] = "ONLY_SERVER";
    DrawMode[DrawMode["ONLY_REAL"] = 2] = "ONLY_REAL";
    DrawMode[DrawMode["BOTH"] = 3] = "BOTH";
})(DrawMode || (DrawMode = {}));
var MAIN_AXIS_COLOR = '#333', SECONDARY_AXIS_COLOR = '#ccc';
/**
 * Class for holding and drawing list of map objects (including players)
 * Can be a wrapper to some framework
 * @param {HTMLCanvasElement} elem
 * @constructor
 */
var Drawer = /** @class */ (function () {
    function Drawer(elem) {
        this.elem = elem;
        this.ctx = this.elem.getContext("2d");
        this.viewport = new Viewport();
        this.objects = [];
        this.objectsById = {};
        this.gridSize = 100;
        this.timeCanvas = document.createElement("canvas");
        this.timeCanvas.width = 300;
        this.timeCanvas.height = 100;
        this.nextObjectId = 1;
    }
    Drawer.prototype.draw = function () {
        this.elem.width = this.elem.clientWidth;
        this.elem.height = this.elem.clientHeight;
        this.viewport.setCanvasSize(this.elem.width, this.elem.height);
        this.ctx.font = "16px serif";
        // this.adjustViewport();
        this.drawGrid();
        this.drawObjects();
        // this.drawAnchors();
        // var now = Date.now();
        // if (this.prevAnimationTime !== null) {
        //     this.animations.push(now - this.prevAnimationTime);
        //
        //     if (this.animations.length > 100) {
        //         this.animations.shift();
        //     }
        // }
        //
        // this.prevAnimationTime = now;
        // this.drawTime();
    };
    Drawer.prototype.drawObjects = function () {
        var ctx = this.ctx;
        var realViewport = this.viewport.getRealDimensions(); // real position and size of viewport
        ctx.save();
        ctx.scale(this.viewport.scale, this.viewport.scale);
        ctx.lineWidth = 1;
        for (var i = 0; i < this.objects.length; i++) {
            var obj = this.objects[i];
            var objSizeReal = obj.size; // half of object size
            if (this.drawMode === DrawMode.ONLY_SERVER || this.drawMode === DrawMode.BOTH) {
                var objPos = obj.getPosition();
                if (this.rectContainsPoint(realViewport, objPos, objSizeReal)) {
                    // рисуем текущее положение объекта по серверу
                    ctx.save();
                    var serverPos = this.viewport.fromRealToCanvas(objPos, false);
                    ctx.translate(serverPos.x, serverPos.y);
                    ctx.lineWidth = 1;
                    ctx.setLineDash([4, 2]);
                    ctx.beginPath();
                    ctx.arc(0, 0, 10, 0, Math.PI * 2);
                    ctx.strokeStyle = '#aaaaaa';
                    ctx.closePath();
                    ctx.stroke();
                    ctx.restore();
                }
            }
        }
        ctx.restore();
    };
    Drawer.prototype.drawGrid = function () {
        var realViewport = this.viewport.getRealDimensions();
        var viewportRealWidth = realViewport.width;
        var viewportRealHeight = realViewport.height;
        var ctx = this.ctx;
        var leftColReal = Math.ceil((realViewport.left) / this.gridSize) * this.gridSize;
        var colCount = Math.max(Math.ceil(realViewport.width / this.gridSize), 1);
        var topRowReal = Math.floor((realViewport.top) / this.gridSize) * this.gridSize;
        var rowCount = Math.max(Math.ceil(realViewport.height / this.gridSize), 1);
        ctx.save();
        ctx.scale(this.viewport.scale, this.viewport.scale);
        ctx.strokeStyle = '#ccc';
        // рисуем вертикали
        for (var i = 0; i < colCount; i++) {
            var rx = leftColReal + i * this.gridSize, x = this.viewport.realXToCanvasWithScale(rx);
            if (rx === 0) {
                ctx.strokeStyle = MAIN_AXIS_COLOR;
            }
            else {
                ctx.strokeStyle = SECONDARY_AXIS_COLOR;
            }
            ctx.beginPath();
            ctx.moveTo(x, 0);
            ctx.lineTo(x, viewportRealHeight);
            ctx.stroke();
        }
        // рисуем горизонтали
        for (var j = 0; j < rowCount; j++) {
            var ry = topRowReal - j * this.gridSize, y = this.viewport.realYToCanvasWithScale(ry);
            if (ry === 0) {
                ctx.strokeStyle = MAIN_AXIS_COLOR;
            }
            else {
                ctx.strokeStyle = SECONDARY_AXIS_COLOR;
            }
            ctx.beginPath();
            ctx.moveTo(0, y);
            ctx.lineTo(viewportRealWidth, y);
            ctx.stroke();
        }
        // // рисуем вращающийся курсор только для непрерывной анимации
        // if (this.lastCursorPositionReal) {
        //     ctx.save();
        //
        //     ctx.strokeStyle = 'magenta';
        //     ctx.lineWidth = 2;
        //     ctx.setLineDash([12, 6]);
        //     this.prevOffset = (this.prevOffset + 0.5) % 18;
        //     ctx.lineDashOffset = this.prevOffset;
        //
        //     let cursorViewport = this.viewport.fromReal(this.lastCursorPositionReal);
        //
        //     ctx.beginPath();
        //     ctx.arc(cursorViewport.x, cursorViewport.y, 20, 0, Math.PI * 2);
        //     ctx.stroke();
        //
        //     ctx.restore();
        // }
        // рисуем центр
        ctx.strokeStyle = 'lime';
        ctx.beginPath();
        // ctx.ellipse(realWidth / 2, realHeight / 2, 10, 10, 0, 0, Math.PI * 2);
        ctx.moveTo(viewportRealWidth / 2 - 15, viewportRealHeight / 2);
        ctx.lineTo(viewportRealWidth / 2 + 15, viewportRealHeight / 2);
        ctx.moveTo(viewportRealWidth / 2, viewportRealHeight / 2 - 15);
        ctx.lineTo(viewportRealWidth / 2, viewportRealHeight / 2 + 15);
        // ctx.arc(realWidth / 2, realHeight / 2, 10, 0, Math.PI * 2);
        ctx.stroke();
        // рисуем границы области
        ctx.lineWidth = 10;
        ctx.beginPath();
        var vlt = this.viewport.fromRealToCanvas({ x: -5000, y: -5000 }, false);
        var vrb = this.viewport.fromRealToCanvas({ x: 5000, y: 5000 }, false);
        ctx.rect(vlt.x, vlt.y, vrb.x - vlt.x, vrb.y - vlt.y);
        ctx.stroke();
        // ctx.globalAlpha = 0.6;
        ctx.restore();
        // Выводим размеры вьюпорта
        var l = Math.round(realViewport.left);
        var t = Math.round(realViewport.top);
        var r = Math.round(realViewport.right);
        var b = Math.round(realViewport.bottom);
        ctx.fillStyle = 'black';
        ctx.font = '14px serif';
        ctx.fillText(t, this.elem.width / 2 - ctx.measureText(t).width / 2, 10);
        ctx.fillText(b, this.elem.width / 2 - ctx.measureText(b).width / 2, this.elem.height);
        ctx.fillText(l, 0, this.elem.height / 2 + 3);
        ctx.fillText(r, this.elem.width - ctx.measureText(r).width, this.elem.height / 2 + 3);
    };
    /**
     * Создаёт новый объект и возвращает его
     * @return {DrawableObject}
     */
    Drawer.prototype.createObject = function () {
        var id = this.nextObjectId++;
        var obj = new DrawableObject(id);
        this.objects.push(obj);
        this.objectsById[id] = obj;
        return obj;
    };
    Drawer.prototype.removeObject = function (objectId) {
    };
    Drawer.prototype.getObject = function () {
        if (this.objects) {
        }
    };
    Drawer.rectContainsPoint = function (rect, point, radius) {
        return rect.left <= (point.x + radius) && point.x - radius <= rect.right &&
            rect.top <= (point.y + radius) && (point.y - radius) <= rect.bottom;
    };
    return Drawer;
}());
