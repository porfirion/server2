"use strict";

/**
 * @constructor
 * @param {DrawableObject} obj
 */
abstract class DrawableObjectLayer {
    protected obj: DrawableObject;

    constructor(obj: DrawableObject) {
        this.obj = obj;
        if (typeof this.obj == "undefined" || this.obj == null) {
            console.warn("empty object!");
        }
    }

    abstract draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void;

    public setObject(obj: DrawableObject): void {
        this.obj = obj;
    }

    static drawTextCentered(ctx: CanvasRenderingContext2D, text: string, x: number, y: number) {
        let measure = ctx.measureText(text);
        ctx.fillText(text, x - measure.width / 2, y);
    };
}


class IdLayer extends DrawableObjectLayer {
    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean) {
        ctx.save();
        try {
            let size = Math.round(useScale ? 12 * viewport.getScale() : 12);
            ctx.font = size + "px serif";
            ctx.fillStyle = "black";
            DrawableObjectLayer.drawTextCentered(ctx, this.obj.getId().toString(), 0, size / 3);
        } finally {
            ctx.restore();
        }
    }

}

class CoordsLayer extends DrawableObjectLayer {
    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getPosition().toString(), 0, 0);
    }
}

class CircleLayer extends DrawableObjectLayer {
    private fillColor: string | null = null;
    private strokeColor: string | null = null;

    constructor(obj: DrawableObject, fillColor: string | null = null, strokeColor: string | null | undefined) {
        super(obj);
        this.fillColor = fillColor;
        if (typeof strokeColor != "undefined") {
            this.strokeColor = strokeColor;
        }
    }

    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean) {
        ctx.lineWidth = 1;

        ctx.beginPath();
        let radius = useScale ? this.obj.getBoundingCircle() * viewport.getScale() : this.obj.getBoundingCircle();
        ctx.arc(0, 0, radius, 0, Math.PI * 2);
        ctx.closePath();

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

class RectLayer extends DrawableObjectLayer {
    private fillColor: string | null = null;
    private borderColor: string | null = null;

    constructor(obj: DrawableObject, fillColor: string | null, borderColor: string | null | undefined) {
        super(obj);
        this.fillColor = fillColor;

        if (typeof borderColor != "undefined") {
            this.borderColor = borderColor;
        }
    }

    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void {
        let size: number = this.obj.getBoundingCircle();
        if (useScale) {
            size *= viewport.getScale();
        }

        if (this.fillColor != null) {
            ctx.fillStyle = this.fillColor;
            ctx.fillRect(-size, -size, size * 2, size * 2);
        }

        if (this.borderColor != null) {
            ctx.strokeStyle = this.borderColor;
            ctx.rect(-size, -size, size * 2, size * 2);
        }
    }
}

class ImageLayer extends DrawableObjectLayer {
    private image: HTMLImageElement;
    private btm: ImageBitmap | null = null;
    constructor(obj: DrawableObject, image: HTMLImageElement) {
        super(obj);
        this.image = image;
        createImageBitmap(this.image).then((btm) => {
            this.btm = btm;
        });
    }

    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void {

        if (this.btm != null) {
            let width = this.btm.width;
            let height = this.btm.height;

            let max = Math.max(width, height);

            let coeff = this.obj.getBoundingCircle() * 2 / max;

            if (useScale) {
                coeff *= viewport.getScale();
            }

            width *= coeff;
            height *= coeff;

            ctx.drawImage(this.btm, -width / 2, -height / 2, width, height);
        }
    }
}

