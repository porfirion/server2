"use strict";

function onLoad() {
    let elem = window.document.getElementById('map');
    let drawer = new Drawer(elem);
    let clickPos = null;
    let clickViewportPos = null;

    drawer.draw();
    updateViewportSize();

    function updateViewportSize() {
        document.getElementById('viewportSize').innerHTML = JSON.stringify(drawer.viewport.getRealDimensions(), numberPrecisionLimiter, " ");
    }

    window.addEventListener("resize", () => {
        updateViewportSize();
        drawer.draw();
    });

    elem.addEventListener('mousewheel', function (ev) {
        let params = normalizeWheel(ev);
        if (params.spinY > 0) {
            // на себя
            drawer.viewport.scaleBy(0.97);
        } else {
            // от себя
            drawer.viewport.scaleBy(1.05);

        }
        drawer.draw();
        updateViewportSize();

        // capture all scrolling over map
        return false;
    });

    elem.addEventListener('contextmenu', function (event) {
        event.preventDefault();
        return false;
    });

    elem.addEventListener('mousedown', function (event) {
        let canvasCoords = new Point2D(event.offsetX, event.offsetY);
        let viewportCoords = drawer.viewport.fromCanvas(canvasCoords);
        let realCoords = drawer.viewport.toReal(viewportCoords);

        switch (event.button) {
            case 0:
                clickPos = canvasCoords;
                clickViewportPos = drawer.viewport.getRealDimensions();
                // левая кнопка мыши
                break;
            case 1:
                // средняя кнопка мыши
                break;
            case 2:
                // правая кнопка мыши
                drawer.viewport.setPos(realCoords);
                drawer.draw();
                updateViewportSize();
                break;
            default:
                console.warn("unexpected button " + event.button);
                break;
        }
        return false;
    });

    elem.addEventListener('mouseup', function () {
        clickPos = null;
        clickViewportPos = null;
    });
    elem.addEventListener('mouseout', function () {
        clickPos = null;
        clickViewportPos = null;
    });

    elem.addEventListener('mousemove', function (event) {
        let canvasCoords = new Point2D(event.offsetX, event.offsetY),
            viewportCoords = drawer.viewport.fromCanvas(canvasCoords),
            realCoords = drawer.viewport.toReal(viewportCoords);

        if (clickPos != null) {
            let
                dx = (canvasCoords.x - clickPos.x) / drawer.viewport.scale,
                dy = (canvasCoords.y - clickPos.y) / drawer.viewport.scale
            drawer.viewport.setPos(new Point2D(
                clickViewportPos.x - dx,
                clickViewportPos.y + dy
            ));
            drawer.draw();
            updateViewportSize();

            // мы передвинули вьюпорт, так что нужно по новой рассчитать положение указателя
            viewportCoords = drawer.viewport.fromCanvas(canvasCoords);
            realCoords = drawer.viewport.toReal(viewportCoords);
        }

        let realOptCoords = drawer.viewport.fromCanvasToReal(canvasCoords);

        document.getElementById('mouseCoords').innerHTML = ''
            + 'canvas  : ' + JSON.stringify(canvasCoords, numberPrecisionLimiter) + "\n"
            + 'viewport: ' + JSON.stringify(viewportCoords, numberPrecisionLimiter) + "\n"
            + 'real    : ' + JSON.stringify(realCoords, numberPrecisionLimiter) + "\n"
            + 'realOpt : ' + JSON.stringify(realOptCoords, numberPrecisionLimiter) + "\n"
            + 'dx: ' + (realCoords.x - realOptCoords.x).toFixed(8) + ' dy: ' + (realCoords.y - realOptCoords.y).toFixed(8)
        ;
    });
}
