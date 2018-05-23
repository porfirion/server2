"use strict";
var TimeDrawer = /** @class */ (function () {
    function TimeDrawer() {
        /**
         * List of animations time (for diagram)
         */
        this.animations = [];
        this.timeCanvas = document.createElement("CANVAS");
        this.ctx = this.timeCanvas.getContext("2d");
        this.timeCanvas.width = 300;
        this.timeCanvas.height = 100;
    }
    TimeDrawer.prototype.getTimeCanvas = function () {
        return this.timeCanvas;
    };
    TimeDrawer.prototype.addAnimationTime = function (time) {
        this.animations.push(time);
        if (this.animations.length > 100) {
            this.animations.shift();
        }
        if (this.ctx != null) {
            this.drawTime(this.ctx);
        }
        else {
            console.log("CTX IS NULL");
        }
    };
    TimeDrawer.prototype.drawTime = function (ctx) {
        ctx.clearRect(0, 0, this.timeCanvas.width, this.timeCanvas.height);
        ctx.fillStyle = 'white';
        ctx.fillRect(0, 0, this.timeCanvas.width, this.timeCanvas.height);
        var min = Infinity;
        var max = -Infinity;
        var average = 0;
        ctx.strokeStyle = '1px black';
        for (var i = 0; i < this.animations.length; i++) {
            average += this.animations[i];
            if (this.animations[i] > max) {
                max = this.animations[i];
            }
            if (this.animations[i] < min) {
                min = this.animations[i];
            }
            ctx.beginPath();
            ctx.moveTo(i * 3, 100);
            ctx.lineTo(i * 3, 100 - this.animations[i]);
            ctx.stroke();
        }
        average = average / this.animations.length;
        ctx.fillStyle = 'black';
        var fillTextRight = function (text, right, top) {
            ctx.fillText(text, right - ctx.measureText(text).width, top);
        };
        fillTextRight('FPS: ' + Math.round(1000 / average), 290, 15);
        fillTextRight('min: ' + min, 290, 30);
        fillTextRight('average: ' + Math.round(average), 290, 45);
        fillTextRight('max: ' + max, 290, 60);
        // ctx.fillText('viewport: (x: ' + Math.round(this.viewport.x * 100) / 100 + '; y: ' + Math.round(this.viewport.y * 100) / 100 + ')', 15, 15);
        // ctx.fillText('scale: ' + Math.round(this.viewport.scale * 100) / 100, 15, 30);
        // ctx.fillText('latency: ' + this.latency.toFixed(0) + ' ms', 15, 45);
        // ctx.fillText('time correction: ' + this.timeCorrection.toFixed(1) + ' ms', 15, 60);
    };
    return TimeDrawer;
}());
//# sourceMappingURL=time_drawer.js.map