declare interface Point2D {
    x: number
    y: number;
}

declare class Viewport {
    constructor(x: number, y: number, scale: number, width: number, height: number)

    setCanvasSize(width: number, height: number): void;
    getRealDimensions(): {x: number, y: number, width: number, height: number, left: number, top: number, right: number, bottom: number};
    getScale(): number;

    realXToCanvasWithScale(x: number): number;
    realYToCanvasWithScale(y: number): number;

    fromRealToCanvas(position: Point2D, applyScale: boolean): Point2D;
}

