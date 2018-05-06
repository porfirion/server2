"use strict";

/**
 * @constructor
 * @param {DrawableObject} obj
 */
abstract class DrawableObjectLayer {
    protected obj: DrawableObject;

    constructor(obj: DrawableObject) {
        this.obj = obj
    }

    abstract draw(ctx: CanvasRenderingContext2D): void;

    static drawTextCentered(ctx: CanvasRenderingContext2D, text: string, x: number, y: number) {
        let measure = ctx.measureText(text);
        ctx.fillText(text, x - measure.width / 2, y);
    };
}


class IdLayer extends DrawableObjectLayer {
    draw(ctx: CanvasRenderingContext2D) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getId().toString(), 0, 0);
    }

}

class CoordsLayer extends DrawableObjectLayer {
    draw(ctx: CanvasRenderingContext2D) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getPosition().toString(), 0, 0);
    }
}

class CircleLayer extends DrawableObjectLayer {
    private strokeColor: string | null = null;
    private fillColor: string | null = null;

    draw(ctx: CanvasRenderingContext2D) {
        ctx.lineWidth = 1;

        ctx.beginPath();
        ctx.arc(0, 0, this.obj.getBoundingCircle(), 0, Math.PI * 2);
        if (this.fillColor != null) {

            ctx.fillStyle = this.fillColor;
            ctx.fill();
        }
        if (this.strokeColor != null) {
            ctx.strokeStyle = this.strokeColor;
            ctx.stroke();
        }
    }
}

