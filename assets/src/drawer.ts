"use strict";

const
    MAIN_AXIS_COLOR = 'rgba(0, 0, 0, 0.055)',
    SECONDARY_AXIS_COLOR = 'rgba(0, 0, 0, 0.055)',
    BG_COLOR = "#3fba54",
    GRID_SIZE = 50,
    USE_CANVAS_SCALE = true;

interface Rectangle {
    left: number
    top: number
    right: number
    bottom: number
}

interface Point2D {
    x: number
    y: number
}

/**
 * Class for holding and drawing list of map objects (including players)
 * Can be a wrapper to some framework
 * @param {HTMLCanvasElement} elem
 * @constructor
 */
class Drawer {
    private ctx: CanvasRenderingContext2D;
    private viewport: Viewport;

    private objects: DrawableObject[];
    private objectsById: Map<number, DrawableObject>;

    private gridSize: number;
    private nextObjectId: number;
    private timeDrawer: TimeDrawer = new TimeDrawer();
    private prevAnimationTime: number | null = null;
    private canvasSize: { width: number, height: number };

    constructor(ctx: CanvasRenderingContext2D, width: number, height: number) {
        this.ctx = ctx;

        this.objects = [];
        this.objectsById = new Map<number, DrawableObject>();
        this.gridSize = GRID_SIZE;
        this.nextObjectId = 1;

        this.viewport = new Viewport(0, 0, 1, width, height);
        this.canvasSize = {width: width, height: height};
    }

    /**
     * Создаёт новый объект и возвращает его
     * @return {DrawableObject}
     */
    public createObject(): DrawableObject {
        let id = this.nextObjectId++;
        let obj = new DrawableObject(id);
        this.objects.push(obj);
        this.objectsById.set(id, obj);

        return obj;
    }

    public removeObject(objectId: number): void {

    }

    public setCanvasSize(width: number, height: number) {
        // this.elem.width = this.elem.clientWidth;
        // this.elem.height = this.elem.clientHeight;
        // this.viewport.setCanvasSize(this.elem.width, this.elem.height);
        this.viewport.setCanvasSize(width, height);
        this.canvasSize = {width: width, height: height};
    }

    public draw() {
        this.ctx.clearRect(0, 0, this.canvasSize.width, this.canvasSize.height);
        this.ctx.save();

        if (typeof BG_COLOR != "undefined") {
            this.ctx.save();
            this.ctx.fillStyle = BG_COLOR;
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
        let now = Date.now();
        if (this.prevAnimationTime !== null) {
            this.timeDrawer.addAnimationTime(now - this.prevAnimationTime);
        }
        this.prevAnimationTime = now;

        this.drawTime();
        this.ctx.restore();

        this.ctx.restore()
    }

    private drawObjects() {
        const ctx = this.ctx;
        const realViewport = this.viewport.getRealDimensions(); // real position and size of viewport

        ctx.save();

        if (USE_CANVAS_SCALE) {
            // console.log("using canvas scale");
            // применяем скейл ко всему канвасу, чтобы работал аппаратный зум
            ctx.scale(this.viewport.getScale(), this.viewport.getScale());
        }

        for (let i = 0; i < this.objects.length; i++) {
            let obj = this.objects[i];
            // будем рисовать только те объекты, которые попадают во вьюпорт
            if (Drawer.rectContainsPoint(realViewport, obj.getPosition(), obj.getBoundingCircle())) {
                // рисуем текущее положение объекта по серверу
                ctx.save();

                // если мы до этого уже применили скейл ко всему канвасу,
                // то здесь его применять уже не нужны и наоборот
                let canvasPos = this.viewport.fromRealToCanvas(obj.getPosition(), !USE_CANVAS_SCALE);

                ctx.translate(canvasPos.x, canvasPos.y);
                obj.draw(ctx, this.viewport, !USE_CANVAS_SCALE);

                ctx.restore();
            }
        }

        ctx.restore();
    }

    private drawGrid() {
        const realViewport = this.viewport.getRealDimensions();

        const viewportRealWidth = realViewport.width;
        const viewportRealHeight = realViewport.height;

        const ctx = this.ctx;

        const leftColReal = Math.ceil((realViewport.left) / this.gridSize) * this.gridSize;
        const colCount = Math.max(Math.ceil(realViewport.width / this.gridSize), 1);
        const topRowReal = Math.floor((realViewport.top) / this.gridSize) * this.gridSize;
        const rowCount = Math.max(Math.ceil(realViewport.height / this.gridSize), 1);

        ctx.save();
        if (USE_CANVAS_SCALE) {
            ctx.scale(this.viewport.getScale(), this.viewport.getScale());
        }

        ctx.strokeStyle = '#ccc';

        // рисуем вертикали
        let height = USE_CANVAS_SCALE ? viewportRealHeight : viewportRealHeight * this.viewport.getScale();
        for (let i = 0; i < colCount; i++) {
            let rx = leftColReal + i * this.gridSize,
                x = this.viewport.fromRealToCanvasX(rx, !USE_CANVAS_SCALE);

            if (rx === 0) {
                ctx.strokeStyle = MAIN_AXIS_COLOR;
            } else {
                ctx.strokeStyle = SECONDARY_AXIS_COLOR;
            }
            ctx.beginPath();
            ctx.moveTo(x, 0);
            ctx.lineTo(x, height);
            ctx.stroke();
        }

        // рисуем горизонтали
        let width = USE_CANVAS_SCALE ? viewportRealWidth : viewportRealWidth * this.viewport.getScale();
        for (let j = 0; j < rowCount; j++) {
            let ry = topRowReal - j * this.gridSize,
                y = this.viewport.fromRealToCanvasY(ry, !USE_CANVAS_SCALE);

            if (ry === 0) {
                ctx.strokeStyle = MAIN_AXIS_COLOR;
            } else {
                ctx.strokeStyle = SECONDARY_AXIS_COLOR;
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
        let centerCanvasPos = this.viewport.fromRealToCanvas(this.viewport.getRealDimensions(), !USE_CANVAS_SCALE);
        ctx.beginPath();
        ctx.moveTo(centerCanvasPos.x - 15, centerCanvasPos.y);
        ctx.lineTo(centerCanvasPos.x + 15, centerCanvasPos.y);
        ctx.moveTo(centerCanvasPos.x, centerCanvasPos.y - 15);
        ctx.lineTo(centerCanvasPos.x, centerCanvasPos.y + 15);
        ctx.stroke();

        // рисуем границы области
        ctx.lineWidth = USE_CANVAS_SCALE ? 10 : 10 * this.viewport.getScale();
        ctx.beginPath();
        let vlt = this.viewport.fromRealToCanvas({x: -5000, y: -5000}, !USE_CANVAS_SCALE);
        let vrb = this.viewport.fromRealToCanvas({x: 5000, y: 5000}, !USE_CANVAS_SCALE);
        ctx.rect(vlt.x, vlt.y, vrb.x - vlt.x, vrb.y - vlt.y);
        ctx.stroke();

        // ctx.globalAlpha = 0.6;

        ctx.restore();

        // Выводим размеры вьюпорта
        let l = Math.round(realViewport.left).toString();
        let t = Math.round(realViewport.top).toString();
        let r = Math.round(realViewport.right).toString();
        let b = Math.round(realViewport.bottom).toString();
        ctx.fillStyle = 'black';
        ctx.font = '14px serif';

        ctx.fillText(t, this.canvasSize.width / 2 - ctx.measureText(t).width / 2, 10);
        ctx.fillText(b, this.canvasSize.width / 2 - ctx.measureText(b).width / 2, this.canvasSize.height);
        ctx.fillText(l, 0, this.canvasSize.height / 2 + 3);
        ctx.fillText(r, this.canvasSize.width - ctx.measureText(r).width, this.canvasSize.height / 2 + 3);
    }

    private drawTime(): void {
        this.ctx.save();
        this.ctx.globalAlpha = 0.7;
        this.ctx.drawImage(this.timeDrawer.getTimeCanvas(), 0, 0, 300, 100);
        this.ctx.restore();
    }

    private static rectContainsPoint(rect: Rectangle, point: Point2D, radius: number): boolean {
        return rect.left <= (point.x + radius) && point.x - radius <= rect.right &&
            rect.bottom <= (point.y + radius) && (point.y - radius) <= rect.top;
    }
}