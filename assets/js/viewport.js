/**
 * Point2D
 * @param {Number} x
 * @param {Number} y
 * @constructor
 */
function Point2D(x, y) {
    this.x = x;
    this.y = y;
}

/**
 * @class
 * Viewport - assumed as a window from screen to real world;
 * @param {Number} x
 * @param {Number} y
 * @param {Number} scale
 * @param {Number} width
 * @param {Number} height
 * @constructor
 */
function Viewport(x, y, scale, width, height) {
    this.maxScale = 0.005;
    this.minScale = 50;

    this.realX = typeof x !== 'undefined' ? Number(x) : 0;
    this.realY = typeof y !== 'undefined' ? Number(y) : 0;
    this.scale = typeof scale !== 'undefined' ? Number(scale) : 1.0;

    this.canvasWidth = typeof width !== 'undefined' ? Number(width) : 1000;
    this.canvasHeight = typeof height !== 'undefined' ? Number(height) : 1000;

    this.canvasWidth_half = this.canvasWidth / 2;
    this.canvasHeight_half = this.canvasHeight / 2;

    this.updateRealSize();
}

/**
 * Move viewport center to (x; y)
 * @param {Point2D} p
 */
Viewport.prototype.setPos = function (p) {
    this.realX = p.x;
    this.realY = p.y;

    this.updateRealSize();
};

/**
 * Multiply current scale by coeff
 * @param coeff
 * @public
 */
Viewport.prototype.scaleBy = function (coeff) {
    this.setScale(this.scale * coeff);
};

/**
 * Set exactly this scale
 * @param value
 * @public
 */
Viewport.prototype.setScale = function (value) {
    this.scale = Math.min(Math.max(value, this.maxScale), this.minScale);

    this.updateRealSize();
};

Viewport.prototype.setCanvasSize = function (width, height) {
    // canvas size
    this.canvasWidth = width;
    this.canvasHeight = height;

    this.canvasWidth_half = this.canvasWidth / 2;
    this.canvasHeight_half = this.canvasHeight / 2;

    this.updateRealSize();
};

/**
 * @private
 */
Viewport.prototype.updateRealSize = function () {
    this.realXScaled = this.realX * this.scale;
    this.realYScaled = this.realY * this.scale;

    this.realWidth = this.canvasWidth / this.scale;
    this.realHeight = this.canvasHeight / this.scale;

    this.realWidth_half = this.realWidth / 2;
    this.realHeight_half = this.realHeight / 2;

    this.realLeft = this.realX - this.realWidth_half;
    this.realTop = this.realY + this.realHeight_half;
    this.realRight = this.realX + this.realWidth_half;
    this.realBottom = this.realY - this.realHeight_half;
};

Viewport.prototype.getScale = function() {
    return this.scale;
}

/**
 * Returns real world coords and size of viewport
 * @returns {{x: number, y: number, width: number, height: number, left: number, top: number, right: number, bottom: number}}
 */
Viewport.prototype.getRealDimensions = function () {
    return {
        // реальное положение и размер вьюпорта
        x: this.realX,
        y: this.realY,
        width: this.realWidth,
        height: this.realHeight,
        scale: this.scale,

        // для удобства отдаём ещё и реальные границы
        left: this.realLeft,
        top: this.realTop,
        right: this.realRight,
        bottom: this.realBottom
    };
};

/**
 * Рассчитывает X на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param {number} rx реальная координата по X
 * @returns {number} координата X на канве
 */
Viewport.prototype.realXToCanvasWithScale = function (rx) {
    /**
     * Normally canvas x is:
     * cX = this.canvasWidth_half + (rX - this.realX) * this.scale,
     *
     * but!
     * this.realWidth = this.canvasWidth / this.scale;
     * this.canvasWidth = this.realWidth * this.scale;
     *
     * this.canvasWidth_half = this.realWidth * this.scale / 2
     *                       = this.realWidth / 2 * this.scale
     *
     * =>
     * cX = this.realWidth / 2 * this.scale + (rX - this.realX) * this.scale
     *    = ( this.realWidth / 2 + (rX - this.realX) ) * this.scale
     *
     * So, if we apply scale on canvas itself, then we do not need to multiply on this.scale
     *
     * cX = this.realWidth / 2 + (rX - this.realX)
     *    = this.realWidth_half + (rX - this.realX)
     *    = this.realWidth_half - this.realX + rX
     *    = (-1) * (this.realX - this.realWidth_half) + rx
     *    = (-1) * this.realLeft + rx
     *    = rx - this.realLeft
     */

    return rx - this.realLeft;
};

/**
 * Рассчитывает Y на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 *
 * Получено по аналогии с realXToCanvasWithScale
 *
 * @param ry реальная координата Y
 * @returns {number} координата Y на канве
 */
Viewport.prototype.realYToCanvasWithScale = function (ry) {
    return this.realTop - ry;
};

/**
 * Translate real world coords into viewport coords
 * @param realPos
 * @returns {Point2D} viewport coords
 *
 *
 *
 * vX = (realPos.x - this.realX) * this.scale,
 * vY = (realPos.y - this.realY) * this.scale
 *
 * NO OPTIMIZATION FOR CLARITY
 *
 * open brackets
 * vX = realPos.x * this.scale - this.realX * this.scale,
 * vY = realPos.y * this.scale - this.realY * this.scale
 *
 * and apply optimization
 * vX = realPos.x * this.scale - this.realXScaled,
 * vY = realPos.y * this.scale - this.realYScaled
 */
Viewport.prototype.fromReal = function (realPos) {
    let
        vX = (realPos.x - this.realX) * this.scale,
        vY = (realPos.y - this.realY) * this.scale;

    return new Point2D(vX, vY);
};

/**
 * Translate viewport coords into canvas coords
 * @param viewportPos
 * @returns {Point2D} canvas coords
 */
Viewport.prototype.toCanvas = function (viewportPos) {
    let
        cX = this.canvasWidth_half + viewportPos.x,
        cY = this.canvasHeight_half - viewportPos.y;

    return new Point2D(cX, cY);
};

/**
 * scale is applied by default
 * @param {Point2D} realPos real world coords
 * @param {Boolean} [applyScale = true]
 * @returns {Point2D} canvas coords
 */
Viewport.prototype.fromRealToCanvas = function (realPos, applyScale) {
    if (typeof applyScale === 'undefined' || Boolean(applyScale) === true) {
        //
        // return toCanvas(fromReal(realPos));
        //
        // let
        //     rX = realPos.x,
        //     rY = realPos.y,
        //
        //     // content of fromReal
        //     vX = (rX - this.realX) * this.scale,
        //     vY = (rY - this.realY) * this.scale,
        //
        //     // content of toCanvas
        //     cX = this.canvasWidth / 2 + vX,
        //     cY = this.canvasHeight / 2 - vY;
        //
        //     // replace vX and vY with their values
        //     cX = this.canvasWidth  / 2 + (rX - this.realX) * this.scale;
        //     cY = this.canvasHeight / 2 - (rY - this.realY) * this.scale;
        //
        //     // replace with optimized halfs
        //     cX = this.canvasWidth_half + (rX - this.realX) * this.scale;
        //     cY = this.canvasHeight_half - (rY - this.realY) * this.scale;
        //
        //     cX = this.canvasWidth_half + rX * this.scale - this.realX * this.scale;
        //     cY = this.canvasHeight_half - rY * this.scale + this.realY * this.scale;
        //
        //     cX = this.canvasWidth / 2 - this.realX * this.scale + rX * this.scale;
        //     cY = this.canvasHeight / 2 + this.realY * this.scale - rY * this.scale;
        //
        //     cX = this.realWidth * this.scale / 2 - this.realX * this.scale + rX * this.scale;
        //     cY = this.realHeight * this.scale / 2 + this.realY * this.scale - rY * this.scale;
        //
        //     cX = this.realWidth / 2 * this.scale - this.realX * this.scale + rX * this.scale;
        //     cY = this.realHeight / 2 * this.scale + this.realY * this.scale - rY * this.scale;
        //
        //     cX = this.realWidth_half * this.scale - this.realX * this.scale + rX * this.scale;
        //     cY = this.realHeight_half * this.scale + this.realY * this.scale - rY * this.scale;
        //
        //     cX = - this.realX * this.scale + this.realWidth_half * this.scale + rX * this.scale;
        //     cY = + this.realY * this.scale + this.realHeight_half * this.scale - rY * this.scale;
        //
        //     cX = - (this.realX * this.scale - this.realWidth_half * this.scale) + rX * this.scale;
        //     cY = + this.realY * this.scale + this.realHeight_half * this.scale - rY * this.scale;
        //
        //     cX = - (this.realX - this.realWidth_half ) * this.scale + rX * this.scale;
        //     cY = + (this.realY + this.realHeight_half) * this.scale - rY * this.scale;
        //
        //     cX = - (this.realLeft ) * this.scale + rX * this.scale;
        //     cY = + (this.realTop) * this.scale - rY * this.scale;
        //
        //     cX = (-this.realLeft  + rX ) * this.scale;
        //     cY = (this.realTop - rY) * this.scale;
        //
        //     cX = -(this.realLeft - rX) * this.scale;
        //     cY = (this.realTop - rY) * this.scale;
        //
        //     cX = (rX - this.realLeft) * this.scale;
        //     cY = (this.realTop - rY) * this.scale;
        //
        // return new Point2D(cX, cY);
        //
        // i.e.
        return new Point2D(
            (realPos.x - this.realLeft) * this.scale,
            (this.realTop - realPos.y) * this.scale
        );
    } else {
        // for explanation see realXToCanvasWithScale and realYToCanvasWithScale

        return new Point2D(
            realPos.x - this.realLeft,
            this.realTop - realPos.y
        );
    }
};

/**
 * Translate viewport coords into real world coords
 * @param {Point2D} viewportPos
 * @returns {Point2D} real world coords
 */
Viewport.prototype.toReal = function (viewportPos) {
    let
        rX = this.realX + (viewportPos.x / this.scale),
        rY = this.realY + (viewportPos.y / this.scale);

    return new Point2D(rX, rY);
};

/**
 * Translate canvas coords to viewport coords
 * @param canvasPos
 * @returns {Point2D}
 */
Viewport.prototype.fromCanvas = function (canvasPos) {
    let
        vX = canvasPos.x - this.canvasWidth_half,
        vY = -(canvasPos.y - this.canvasHeight_half);

    return new Point2D(vX, vY);
};

/**
 * Translate canvas coords to real world coords
 * @param {Point2D} canvasPos
 * @returns {Point2D} real world coords
 *
 *
 * return this.toReal(this.fromCanvas(canvasPos));
 *
 * let
 *     vX = canvasPos.x - this.canvasWidth_half,
 *     vY = -(canvasPos.y - this.canvasHeight_half),
 *
 *     rX = this.realX + (vX / this.scale),
 *     rY = this.realY + (vY / this.scale);
 *
 * rX = this.realX + ((canvasPos.x - this.canvasWidth_half) / this.scale);
 * rY = this.realY + (-(canvasPos.y - this.canvasHeight_half) / this.scale);
 *
 * rX = this.realX + (canvasPos.x - this.canvasWidth_half) / this.scale;
 * rY = this.realY + (this.canvasHeight_half - canvasPos.y) / this.scale;
 *
 * rX = this.realX + canvasPos.x / this.scale - this.canvasWidth_half / this.scale;
 * rY = this.realY + this.canvasHeight_half / this.scale - canvasPos.y / this.scale;
 *
 * rX = this.realX + canvasPos.x / this.scale - this.realWidth_half;
 * rY = this.realY + this.realHeight_half - canvasPos.y / this.scale;
 *
 * rX = this.realX - this.realWidth_half  + canvasPos.x / this.scale;
 * rY = this.realY + this.realHeight_half - canvasPos.y / this.scale;
 *
 * this.realLeft = this.realX - this.realWidth_half;
 * this.realTop = this.realY + this.realHeight_half;
 *
 * rX = this.realLeft  + canvasPos.x / this.scale;
 * rY = this.realTop - canvasPos.y / this.scale;
 */
Viewport.prototype.fromCanvasToReal = function (canvasPos) {
    return new Point2D(
        this.realLeft + canvasPos.x / this.scale,
        this.realTop - canvasPos.y / this.scale
    );
};