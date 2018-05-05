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

    abstract draw(ctx: CanvasRenderingContext2D);

    static drawTextCentered(ctx, text, x, y) {
        let measure = ctx.measureText(text);
        ctx.fillText(text, x - measure.width / 2, y);
    };
}


class IdLayer extends DrawableObjectLayer {
    draw(ctx) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getId().toString(), 0, 0);
    }

}
class CoordsLayer extends DrawableObjectLayer {
    draw(ctx) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getPosition().toString(), 0, 0);
    }
}

