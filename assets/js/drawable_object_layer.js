"use strict";

/**
 * @constructor
 * @param {DrawableObject} obj
 */
function DrawableObjectLayer(obj) {
    this.obj = obj;
}

DrawableObjectLayer.prototype.draw = function() {};
DrawableObjectLayer.prototype.drawTextCentered = function (ctx, text, x, y) {
    var measure = ctx.measureText(text);
    ctx.fillText(text, x - measure.width / 2, y);
};


function IdLayer(obj) {
    DrawableObjectLayer.apply(this, arguments);
}
IdLayer.prototype.draw = function(ctx) {
    this.drawTextCentered(ctx, this.obj.id, 0, 0);
};

/**
 *
 * @param {Point2D} coords
 * @constructor
 */
function CoordsLayer(coords) {
    this.coords = coords
}

CoordsLayer.prototype.draw = function() {

};

