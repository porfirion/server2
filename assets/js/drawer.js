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
        var context = this.elem.getContext("2d");
        if (context != null) {
            this.ctx = context;
        }
        else {
            console.error("no context");
        }
        this.viewport = new Viewport(0, 0, 1, 1000, 1000);
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
        ctx.scale(this.viewport.getScale(), this.viewport.getScale());
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
        ctx.scale(this.viewport.getScale(), this.viewport.getScale());
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
        var l = Math.round(realViewport.left).toString();
        var t = Math.round(realViewport.top).toString();
        var r = Math.round(realViewport.right).toString();
        var b = Math.round(realViewport.bottom).toString();
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
    Drawer.prototype.drawTime = function () {
        var ctx = this.timeCanvas.getContext('2d');
        ctx.clearRect(0, 0, this.timeCanvas.width, this.timeCanvas.height);
        ctx.fillStyle = 'white';
        ctx.fillRect(0, 0, this.timeCanvas.width, this.timeCanvas.height);
        var min = Infinity;
        var max = -Infinity;
        var average = 0;
        ctx.strokeStyle = '1px black';
        for (var i = 0; i < this.animations.length; i++) {
            average += this.animations[i];
            if (this.animations[i] > max) {
                max = this.animations[i];
            }
            if (this.animations[i] < min) {
                min = this.animations[i];
            }
            ctx.beginPath();
            ctx.moveTo(i * 3, 100);
            ctx.lineTo(i * 3, 100 - this.animations[i]);
            ctx.stroke();
        }
        average = average / this.animations.length;
        ctx.fillStyle = 'black';
        var fillTextRight = function (text, right, top) {
            ctx.fillText(text, right - ctx.measureText(text).width, top);
        };
        fillTextRight('FPS: ' + Math.round(1000 / average), 290, 15);
        fillTextRight('min: ' + min, 290, 30);
        fillTextRight('average: ' + Math.round(average), 290, 45);
        fillTextRight('max: ' + max, 290, 60);
        ctx.fillText('viewport: (x: ' + Math.round(this.viewport.x * 100) / 100 + '; y: ' + Math.round(this.viewport.y * 100) / 100 + ')', 15, 15);
        ctx.fillText('scale: ' + Math.round(this.viewport.scale * 100) / 100, 15, 30);
        ctx.fillText('latency: ' + this.latency.toFixed(0) + ' ms', 15, 45);
        ctx.fillText('time correction: ' + this.timeCorrection.toFixed(1) + ' ms', 15, 60);
        this.ctx.save();
        this.ctx.globalAlpha = 0.7;
        this.ctx.drawImage(this.timeCanvas, 0, 0, 300, 100);
        this.ctx.restore();
    };
    Drawer.rectContainsPoint = function (rect, point, radius) {
        return rect.left <= (point.x + radius) && point.x - radius <= rect.right &&
            rect.top <= (point.y + radius) && (point.y - radius) <= rect.bottom;
    };
    return Drawer;
}());
module.exports = Drawer;
