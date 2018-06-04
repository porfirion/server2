interface Drawable {
    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void;
}