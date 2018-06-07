// Reasonable defaults
var PIXEL_STEP = 10;
var LINE_HEIGHT = 40;
var PAGE_HEIGHT = 800;

function generateUUID() {
    var d = new Date().getTime();
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = (d + Math.random() * 16) % 16 | 0;
        d = Math.floor(d / 16);
        return (c === 'x' ? r : (r & 0x7 | 0x8)).toString(16);
    });
}

/**
 *
 * @param event
 * @returns {{spinX: number, spinY: number, pixelX: number, pixelY: number}}
 */
function normalizeWheel(event) {
    var sX = 0, sY = 0,       // spinX, spinY
        pX = 0, pY = 0;       // pixelX, pixelY

    // Legacy
    if ('detail' in event) {
        sY = event.detail;
    }
    if ('wheelDelta' in event) {
        sY = -event.wheelDelta / 120;
    }
    if ('wheelDeltaY' in event) {
        sY = -event.wheelDeltaY / 120;
    }
    if ('wheelDeltaX' in event) {
        sX = -event.wheelDeltaX / 120;
    }

    // side scrolling on FF with DOMMouseScroll
    if ('axis' in event && event.axis === event.HORIZONTAL_AXIS) {
        sX = sY;
        sY = 0;
    }

    pX = sX * PIXEL_STEP;
    pY = sY * PIXEL_STEP;

    if ('deltaY' in event) { pY = event.deltaY; }
    if ('deltaX' in event) { pX = event.deltaX; }

    if ((pX || pY) && event.deltaMode) {
        if (event.deltaMode == 1) {          // delta in LINE units
            pX *= LINE_HEIGHT;
            pY *= LINE_HEIGHT;
        } else {                             // delta in PAGE units
            pX *= PAGE_HEIGHT;
            pY *= PAGE_HEIGHT;
        }
    }

    // Fall-back if spin cannot be determined
    if (pX && !sX) { sX = (pX < 1) ? -1 : 1; }
    if (pY && !sY) { sY = (pY < 1) ? -1 : 1; }

    return {
        spinX: sX,
        spinY: sY,
        pixelX: pX,
        pixelY: pY
    };
}

/**
 * 0 - 255 in hex
 * @returns {int}
 */
function randomColorComponent() {
    var comp = (Math.round(Math.random() * 255)).toString(16);
    if (comp.length < 2) {
        comp = '0' + comp;
    }
    return comp;
}

/**
 * Random color hex representation
 * @returns {string}
 */
function randomColor() {
    return '#' + randomColorComponent() + randomColorComponent() + randomColorComponent();
}

function randomName() {
    var names = [
        "Ivan", "Mikhail", "Ilya", "Sergey", "Alexander", "Egor", "Diman", "Alexey"
    ];
    var suffix = [
        "Eagle eye", "Morning star", "Black hoof", "Three finger", "Yellow wolf", "Wasp nest", "Short bull", "Hard needle", "Fetid datura", "Curved horn"
    ];

    var rndNameId = Math.round(Math.random() * (names.length - 1));
    var rndSuffixId = Math.round(Math.random() * (suffix.length - 1));
    // var rndNum = Math.round(Math.random() * 1000000);

    // console.log(rndNameId, rndNum, names[rndNameId]);

    return names[rndNameId] + " " + suffix[rndSuffixId];
}

function numberPrecisionLimiter(key, value) {
    if (typeof value === 'number') {
        return parseFloat(value.toFixed(3));
    }
    return value;
}

function Extend(Child, Parent) {
    var F = function () { };
    F.prototype = Parent.prototype;
    Child.prototype = new F();
    Child.prototype.constructor = Child;
    Child.superclass = Parent.prototype;
}

function drawTextCentered(ctx, text, x, y) {
    var measure = ctx.measureText(text);
    ctx.fillText(text, x - measure.width / 2, y);
}