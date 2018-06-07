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
        if (typeof this.obj == "undefined" || this.obj == null) {
            console.warn("empty object!");
        }
    }
    DrawableObjectLayer.prototype.setObject = function (obj) {
        this.obj = obj;
    };
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
    IdLayer.prototype.draw = function (ctx, viewport, useScale) {
        ctx.save();
        try {
            var size = Math.round(useScale ? 12 * viewport.getScale() : 12);
            ctx.font = size + "px serif";
            ctx.fillStyle = "black";
            DrawableObjectLayer.drawTextCentered(ctx, this.obj.getId().toString(), 0, size / 3);
        }
        finally {
            ctx.restore();
        }
    };
    return IdLayer;
}(DrawableObjectLayer));
var CoordsLayer = /** @class */ (function (_super) {
    __extends(CoordsLayer, _super);
    function CoordsLayer() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    CoordsLayer.prototype.draw = function (ctx, viewport, useScale) {
        DrawableObjectLayer.drawTextCentered(ctx, this.obj.getPosition().toString(), 0, 0);
    };
    return CoordsLayer;
}(DrawableObjectLayer));
var CircleLayer = /** @class */ (function (_super) {
    __extends(CircleLayer, _super);
    function CircleLayer(obj, fillColor, strokeColor) {
        if (fillColor === void 0) { fillColor = null; }
        var _this = _super.call(this, obj) || this;
        _this.fillColor = null;
        _this.strokeColor = null;
        _this.fillColor = fillColor;
        if (typeof strokeColor != "undefined") {
            _this.strokeColor = strokeColor;
        }
        return _this;
    }
    CircleLayer.prototype.draw = function (ctx, viewport, useScale) {
        ctx.lineWidth = 1;
        ctx.beginPath();
        var radius = useScale ? this.obj.getBoundingCircle() * viewport.getScale() : this.obj.getBoundingCircle();
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
    };
    return CircleLayer;
}(DrawableObjectLayer));
var RectLayer = /** @class */ (function (_super) {
    __extends(RectLayer, _super);
    function RectLayer(obj, fillColor, borderColor) {
        var _this = _super.call(this, obj) || this;
        _this.fillColor = null;
        _this.borderColor = null;
        _this.fillColor = fillColor;
        if (typeof borderColor != "undefined") {
            _this.borderColor = borderColor;
        }
        return _this;
    }
    RectLayer.prototype.draw = function (ctx, viewport, useScale) {
        var size = this.obj.getBoundingCircle();
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
    };
    return RectLayer;
}(DrawableObjectLayer));
var ImageLayer = /** @class */ (function (_super) {
    __extends(ImageLayer, _super);
    function ImageLayer(obj, image) {
        var _this = _super.call(this, obj) || this;
        _this.image = image;
        return _this;
    }
    ImageLayer.prototype.draw = function (ctx, viewport, useScale) {
        if (this.image.data != null) {
            var width = this.image.data.width;
            var height = this.image.data.height;
            var max = Math.max(width, height);
            var coeff = this.obj.getBoundingCircle() * 2 / max;
            if (useScale) {
                coeff *= viewport.getScale();
            }
            width *= coeff;
            height *= coeff;
            ctx.drawImage(this.image.data, -width / 2, -height / 2, width, height);
        }
    };
    return ImageLayer;
}(DrawableObjectLayer));
//# sourceMappingURL=drawable_object_layer.js.map