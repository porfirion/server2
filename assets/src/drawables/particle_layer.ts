interface ParticleEmitter {
    emit(): Drawable;
}

class ParticleLayer extends DrawableObjectLayer {
    private startTime: number;
    private emitter: ParticleEmitter;

    private particles: { drawable: Drawable, pos: { x: number, y: number } }[] = [];

    constructor(obj: DrawableObject, emitter: ParticleEmitter) {
        super(obj);
        this.startTime = Date.now();
        this.emitter = emitter;

        for (let i = 0; i < 10; i++) {
            this.particles.push({
                drawable: emitter.emit(),
                pos: {x: Math.random() * 40 - 20, y: Math.random() * 40 - 20}
            });
        }
    };

    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void {
        let passedTime = Date.now() - this.startTime;

        for (let i = 0; i < this.particles.length; i++) {
            ctx.save();
            ctx.translate(this.particles[i].pos.x, this.particles[i].pos.y);
            this.particles[i].drawable.draw(ctx, viewport, useScale);
            ctx.restore();
        }
    }
}