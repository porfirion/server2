"use strict";
var Sprite = /** @class */ (function () {
    function Sprite(name) {
        var _this = this;
        this.data = null;
        this.img = new Image();
        this.img.onload = function () {
            createImageBitmap(_this.img).then(function (value) {
                _this.data = value;
            });
        };
    }
    Sprite.prototype.draw = function (ctx, viewport, useScale) {
    };
    return Sprite;
}());
//# sourceMappingURL=sprite.js.map