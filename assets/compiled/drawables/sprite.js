"use strict";
var Sprite = /** @class */ (function () {
    function Sprite(name) {
        var _this = this;
        this.data = null;
        this.name = name;
        this.img = new Image();
        this.img.onload = function () {
            createImageBitmap(_this.img).then(function (value) {
                _this.data = value;
            });
        };
        this.img.src = Sprite.getSrcByName(name);
    }
    Sprite.prototype.draw = function (ctx, viewport, useScale) {
        console.error('not implemented');
    };
    Sprite.getSrcByName = function (name) {
        return name;
    };
    return Sprite;
}());
//# sourceMappingURL=sprite.js.map