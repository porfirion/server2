declare interface Point2D {
    x: number
    y: number;
}

declare class Viewport {
    constructor(x: number, y: number, scale: number, width: number, height: number);

    setCanvasSize(width: number, height: number): void;
    setPos(pos: Point2D): void;
    setScale(scale: number): void;
    scaleBy(coeff: number): void;

    getScale(): number;
    getRealDimensions(): {x: number, y: number, width: number, height: number, left: number, top: number, right: number, bottom: number};

    fromCanvas(canvasCoords: Point2D): Point2D;
    toReal(viewportCoords: Point2D): Point2D;

    fromCanvasToReal(canvasCoords: Point2D): Point2D;

    realXToCanvasWithScale(x: number): number;
    realYToCanvasWithScale(y: number): number;

    fromRealToCanvas(position: Point2D, applyScale: boolean): Point2D;
}

