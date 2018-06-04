"use strict";
var SpriteFactory = /** @class */ (function () {
    function SpriteFactory() {
        this.sprites = new Map();
    }
    SpriteFactory.prototype.getSprite = function (name) {
        var existing = this.sprites.get(name);
        if (typeof existing != 'undefined') {
            return existing;
        }
        else {
            existing = new Sprite(name);
            this.sprites.set(name, existing);
            return existing;
        }
    };
    return SpriteFactory;
}());
//# sourceMappingURL=sprite_factory.js.map