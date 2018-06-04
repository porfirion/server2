"use strict";
var Particle = /** @class */ (function () {
    function Particle(color, size) {
        this.gradient = null;
        this.color = color;
        this.size = size;
    }
    Particle.prototype.draw = function (ctx, viewport, useScale) {
        if (this.gradient == null) {
            this.gradient = ctx.createRadialGradient(0, 0, 0, 0, 0, 5);
            this.gradient.addColorStop(0.0, "rgba(255, 255, 255, 0.9)");
            this.gradient.addColorStop(0.6, this.color + "70");
            this.gradient.addColorStop(1, this.color + "00");
        }
        ctx.fillStyle = this.gradient;
        ctx.beginPath();
        ctx.arc(0, 0, 10, 0, Math.PI * 2);
        ctx.closePath();
        ctx.fill();
    };
    return Particle;
}());
//# sourceMappingURL=particle.js.map