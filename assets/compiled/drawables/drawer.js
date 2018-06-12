"use strict";
/**
 * Class for holding and drawing list of map objects (including players)
 * Can be a wrapper to some framework
 */
var Drawer = /** @class */ (function () {
    function Drawer(ctx, width, height) {
        this.timeDrawer = new TimeDrawer();
        this.prevAnimationTime = null;
        this.useCanvasScale = true;
        this.gridSize = 50;
        this.bgColor = "#3fba54";
        this.mainAxisColor = 'rgba(0, 0, 0, 0.055)';
        this.secondaryAxisColor = 'rgba(0, 0, 0, 0.055)';
        this.ctx = ctx;
        this.objects = [];
        this.objectsById = new Map();
        this.nextObjectId = 1;
        this.viewport = new Viewport(0, 0, 1, width, height);
        this.canvasSize = { width: width, height: height };
    }
    Drawer.prototype.getViewport = function () {
        return this.viewport;
    };
    /**
     * Создаёт новый объект и возвращает его
     * @return {DrawableObject}
     */
    Drawer.prototype.createObject = function () {
        var id = this.nextObjectId++;
        var obj = new DrawableObject(id);
        this.objects.push(obj);
        this.objectsById.set(id, obj);
        return obj;
    };
    Drawer.prototype.removeObject = function (objectId) {
    };
    Drawer.prototype.setCanvasSize = function (width, height) {
        // this.elem.width = this.elem.clientWidth;
        // this.elem.height = this.elem.clientHeight;
        // this.viewport.setCanvasSize(this.elem.width, this.elem.height);
        this.viewport.setCanvasSize(width, height);
        this.canvasSize = { width: width, height: height };
    };
    Drawer.prototype.draw = function () {
        if (this.ctx == null)
            return;
        this.ctx.clearRect(0, 0, this.canvasSize.width, this.canvasSize.height);
        this.ctx.save();
        if (typeof this.bgColor != "undefined") {
            this.ctx.save();
            this.ctx.fillStyle = this.bgColor;
            this.ctx.fillRect(0, 0, this.canvasSize.width, this.canvasSize.height);
            this.ctx.restore();
        }
        this.ctx.save();
        this.drawGrid();
        this.ctx.restore();
        this.ctx.save();
        this.drawObjects();
        this.ctx.restore();
        // this.drawAnchors();
        this.ctx.save();
        var now = Date.now();
        if (this.prevAnimationTime !== null) {
            this.timeDrawer.addAnimationTime(now - this.prevAnimationTime);
        }
        this.prevAnimationTime = now;
        this.ctx.globalAlpha = 0.6;
        this.ctx.translate(0, 0);
        this.ctx.scale(1, 1);
        this.ctx.drawImage(this.timeDrawer.getTimeCanvas(), 0, 0, 300, 100);
        this.ctx.restore();
        this.ctx.restore();
    };
    Drawer.prototype.drawObjects = function () {
        if (this.ctx == null)
            return;
        var ctx = this.ctx;
        var realViewport = this.viewport.getRealDimensions(); // real position and size of viewport
        ctx.save();
        ctx.setTransform(1, 0, 0, 1, 0, 0);
        if (this.useCanvasScale) {
            // console.log("using canvas scale");
            // применяем скейл ко всему канвасу, чтобы работал аппаратный зум
            ctx.scale(this.viewport.getScale(), this.viewport.getScale());
        }
        for (var i = 0; i < this.objects.length; i++) {
            var obj = this.objects[i];
            // будем рисовать только те объекты, которые попадают во вьюпорт
            if (Drawer.rectContainsPoint(realViewport, obj.getPosition(), obj.getBoundingCircle())) {
                // рисуем текущее положение объекта по серверу
                ctx.save();
                // если мы до этого уже применили скейл ко всему канвасу,
                // то здесь его применять уже не нужны и наоборот
                var canvasPos = this.viewport.fromRealToCanvas(obj.getPosition(), !this.useCanvasScale);
                ctx.translate(canvasPos.x, canvasPos.y);
                ctx.rotate(obj.getRotation());
                obj.draw(ctx, this.viewport, !this.useCanvasScale);
                ctx.restore();
            }
        }
        ctx.restore();
    };
    Drawer.prototype.drawGrid = function () {
        if (this.ctx == null)
            return;
        var realViewport = this.viewport.getRealDimensions();
        var viewportRealWidth = realViewport.width;
        var viewportRealHeight = realViewport.height;
        var ctx = this.ctx;
        var leftColReal = Math.ceil((realViewport.left) / this.gridSize) * this.gridSize;
        var colCount = Math.max(Math.ceil(realViewport.width / this.gridSize), 1);
        var topRowReal = Math.floor((realViewport.top) / this.gridSize) * this.gridSize;
        var rowCount = Math.max(Math.ceil(realViewport.height / this.gridSize), 1);
        ctx.save();
        if (this.useCanvasScale) {
            ctx.scale(this.viewport.getScale(), this.viewport.getScale());
        }
        ctx.strokeStyle = '#ccc';
        // рисуем вертикали
        var height = this.useCanvasScale ? viewportRealHeight : viewportRealHeight * this.viewport.getScale();
        for (var i = 0; i < colCount; i++) {
            var rx = leftColReal + i * this.gridSize, x = this.viewport.fromRealToCanvasX(rx, !this.useCanvasScale);
            if (rx === 0) {
                ctx.strokeStyle = this.mainAxisColor;
            }
            else {
                ctx.strokeStyle = this.secondaryAxisColor;
            }
            ctx.beginPath();
            ctx.moveTo(x, 0);
            ctx.lineTo(x, height);
            ctx.stroke();
        }
        // рисуем горизонтали
        var width = this.useCanvasScale ? viewportRealWidth : viewportRealWidth * this.viewport.getScale();
        for (var j = 0; j < rowCount; j++) {
            var ry = topRowReal - j * this.gridSize, y = this.viewport.fromRealToCanvasY(ry, !this.useCanvasScale);
            if (ry === 0) {
                ctx.strokeStyle = this.mainAxisColor;
            }
            else {
                ctx.strokeStyle = this.secondaryAxisColor;
            }
            ctx.beginPath();
            ctx.moveTo(0, y);
            ctx.lineTo(width, y);
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
        var centerCanvasPos = this.viewport.fromRealToCanvas(this.viewport.getRealDimensions(), !this.useCanvasScale);
        ctx.beginPath();
        ctx.moveTo(centerCanvasPos.x - 15, centerCanvasPos.y);
        ctx.lineTo(centerCanvasPos.x + 15, centerCanvasPos.y);
        ctx.moveTo(centerCanvasPos.x, centerCanvasPos.y - 15);
        ctx.lineTo(centerCanvasPos.x, centerCanvasPos.y + 15);
        ctx.stroke();
        // рисуем границы области
        ctx.lineWidth = this.useCanvasScale ? 10 : 10 * this.viewport.getScale();
        ctx.beginPath();
        var vlt = this.viewport.fromRealToCanvas({ x: -5000, y: -5000 }, !this.useCanvasScale);
        var vrb = this.viewport.fromRealToCanvas({ x: 5000, y: 5000 }, !this.useCanvasScale);
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
        ctx.fillText(t, this.canvasSize.width / 2 - ctx.measureText(t).width / 2, 10);
        ctx.fillText(b, this.canvasSize.width / 2 - ctx.measureText(b).width / 2, this.canvasSize.height);
        ctx.fillText(l, 0, this.canvasSize.height / 2 + 3);
        ctx.fillText(r, this.canvasSize.width - ctx.measureText(r).width, this.canvasSize.height / 2 + 3);
    };
    Drawer.rectContainsPoint = function (rect, point, radius) {
        return rect.left <= (point.x + radius) && point.x - radius <= rect.right &&
            rect.bottom <= (point.y + radius) && (point.y - radius) <= rect.top;
    };
    return Drawer;
}());
//# sourceMappingURL=drawer.js.map