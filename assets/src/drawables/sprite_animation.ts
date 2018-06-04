class AnimationStep {
    duration: number;
    sprite: Sprite;

    constructor(duration: number, sprite: Sprite) {
        this.duration = duration;
        this.sprite = sprite;
    }
}

class SpriteAnimation implements Drawable {
    private startTime: number;
    private steps: AnimationStep[];
    private perpetual: boolean;
    private totalDuration: number;

    public onfinish: (() => void) | null = null;

    constructor(steps: AnimationStep[], perpetual: boolean = true) {
        this.startTime = Date.now();
        this.steps = steps;
        this.perpetual = perpetual;

        let sum = 0;
        for (let i = 0; i < steps.length; i++) {
            sum += this.steps[i].duration;
        }
        this.totalDuration = sum;
    }

    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void {
        let passedTime = Date.now() - this.startTime;
        if (!this.perpetual && passedTime > this.totalDuration) {
            if (this.onfinish != null) {
                this.onfinish();
            }
            return;
        }

        passedTime %= this.totalDuration;

        let i = 0;
        while (i < this.steps.length && this.steps[i].duration > passedTime) {
            passedTime -= this.steps[i].duration;
        }

        if (i < this.steps.length) {
            this.steps[i].sprite.draw(ctx, viewport, useScale);
        } else {
            console.warn("out of sprites range")
        }
    }
}