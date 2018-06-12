const SCALE_STEP = 1.05;

class CanvasInputController {
    private canvas: HTMLCanvasElement;
    private drawer: Drawer;
    private gameState: GameState;

    constructor(canvas: HTMLCanvasElement, drawer: Drawer, gameState:GameState) {
        this.canvas = canvas;
        this.drawer = drawer;
        this.gameState = gameState;

        this.initHandlers();
    }

    initHandlers(): void {
        $(this.canvas).on('mousedown', this.onMouseDown.bind(this));
        $(this.canvas).on('mousemove', this.onMouseMove.bind(this));
        $(this.canvas).on('mouseenter', this.onMouseEnter.bind(this));
        $(this.canvas).on('mouseout', this.onMouseOut.bind(this));
        $(this.canvas).on('mousewheel DOMMouseScroll', this.onMouseWheel.bind(this));
        $(this.canvas).on('contextmenu', this.onContextMenu.bind(this));
    }

    onMouseDown(event: MouseEvent): boolean {
        let canvasCoords = {x: event.offsetX, y: event.offsetY};

        let viewportCoords = this.drawer.getViewport().fromCanvas(canvasCoords);
        let realCoords = this.drawer.getViewport().toReal(viewportCoords);
        
        switch (event.button) {
            case 0:
                // левая кнопка мыши
                console.log('not implemented');
                break;
            case 1:
                // средняя кнопка мыши
                break;
            case 2:
                // правая кнопка мыши
                let newPos = {x: realCoords.x, y: realCoords.y};
                this.drawer.getViewport().setPos(newPos);
                break;
            default:
                console.warn("unexpected button " + event.button);
                break;
        }
        return false;
    }

    onMouseMove(event: MouseEvent): void {
        let canvasCoords = {x: event.offsetX, y: event.offsetY};
        let viewportCoords = this.drawer.getViewport().fromCanvas(canvasCoords);
        let realCoords = this.drawer.getViewport().toReal(viewportCoords);
    }

    onMouseEnter(event: MouseEvent): void {
    }

    onMouseOut(event: MouseEvent): void {
    }

    onMouseWheel(event: {originalEvent: WheelEvent}): boolean {
        let params = normalizeWheel(event.originalEvent);
        if (params.spinY > 0) {
            // на себя
            this.drawer.getViewport().scaleBy(1.0/SCALE_STEP);
        } else {
            // от себя
            this.drawer.getViewport().scaleBy(SCALE_STEP);
        }

        // capture all scrolling over map
        return false;
    }

    onContextMenu(event: Event): boolean {
        event.preventDefault();
        return false;
    }
}