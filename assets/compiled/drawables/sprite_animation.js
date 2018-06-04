"use strict";
var AnimationStep = /** @class */ (function () {
    function AnimationStep(duration, sprite) {
        this.duration = duration;
        this.sprite = sprite;
    }
    return AnimationStep;
}());
var SpriteAnimation = /** @class */ (function () {
    function SpriteAnimation(steps, perpetual) {
        if (perpetual === void 0) { perpetual = true; }
        this.onfinish = null;
        this.startTime = Date.now();
        this.steps = steps;
        this.perpetual = perpetual;
        var sum = 0;
        for (var i = 0; i < steps.length; i++) {
            sum += this.steps[i].duration;
        }
        this.totalDuration = sum;
    }
    SpriteAnimation.prototype.draw = function (ctx, viewport, useScale) {
        var passedTime = Date.now() - this.startTime;
        if (!this.perpetual && passedTime > this.totalDuration) {
            if (this.onfinish != null) {
                this.onfinish();
            }
            return;
        }
        passedTime %= this.totalDuration;
        var i = 0;
        while (i < this.steps.length && this.steps[i].duration > passedTime) {
            passedTime -= this.steps[i].duration;
        }
        if (i < this.steps.length) {
            this.steps[i].sprite.draw(ctx, viewport, useScale);
        }
        else {
            console.warn("out of sprites range");
        }
    };
    return SpriteAnimation;
}());
//# sourceMappingURL=sprite_animation.js.map