"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = Object.setPrototypeOf ||
        ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
        function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
/**
 * @constructor
 * @param {DrawableObject} obj
 */
var DrawableObjectLayer = /** @class */ (function () {
    function DrawableObjectLayer(obj) {
        this.obj = obj;
    }
    DrawableObjectLayer.drawTextCentered = function (ctx, text, x, y) {
        var measure = ctx.measureText(text);
        ctx.fillText(text, x - measure.width / 2, y);
    };
    ;
    return DrawableObjectLayer;
}());
var IdLayer = /** @class */ (function (_super) {
    __extends(IdLayer, _super);
    function IdLayer() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    IdLayer.prototype.draw = function (ctx) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getId().toString(), 0, 0);
    };
    return IdLayer;
}(DrawableObjectLayer));
var CoordsLayer = /** @class */ (function (_super) {
    __extends(CoordsLayer, _super);
    function CoordsLayer() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    CoordsLayer.prototype.draw = function (ctx) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getPosition().toString(), 0, 0);
    };
    return CoordsLayer;
}(DrawableObjectLayer));
var CircleLayer = /** @class */ (function (_super) {
    __extends(CircleLayer, _super);
    function CircleLayer() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this.strokeColor = null;
        _this.fillColor = null;
        return _this;
    }
    CircleLayer.prototype.draw = function (ctx) {
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.arc(0, 0, this.obj.getBoundingCircle(), 0, Math.PI * 2);
        if (this.fillColor != null) {
            ctx.fillStyle = this.fillColor;
            ctx.fill();
        }
        if (this.strokeColor != null) {
            ctx.strokeStyle = this.strokeColor;
            ctx.stroke();
        }
    };
    return CircleLayer;
}(DrawableObjectLayer));
//# sourceMappingURL=drawable_object_layer.js.map