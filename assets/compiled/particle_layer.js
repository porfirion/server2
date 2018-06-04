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
var ParticleLayer = /** @class */ (function (_super) {
    __extends(ParticleLayer, _super);
    function ParticleLayer(obj, emitter) {
        var _this = _super.call(this, obj) || this;
        _this.particles = [];
        _this.startTime = Date.now();
        _this.emitter = emitter;
        for (var i = 0; i < 10; i++) {
            _this.particles.push({
                drawable: emitter.emit(),
                pos: { x: Math.random() * 40 - 20, y: Math.random() * 40 - 20 }
            });
        }
        return _this;
    }
    ;
    ParticleLayer.prototype.draw = function (ctx, viewport, useScale) {
        var passedTime = Date.now() - this.startTime;
        for (var i = 0; i < this.particles.length; i++) {
            ctx.save();
            ctx.translate(this.particles[i].pos.x, this.particles[i].pos.y);
            this.particles[i].drawable.draw(ctx, viewport, useScale);
            ctx.restore();
        }
    };
    return ParticleLayer;
}(DrawableObjectLayer));
//# sourceMappingURL=particle_layer.js.map