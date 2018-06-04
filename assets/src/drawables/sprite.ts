class Sprite implements Drawable {
    private img: HTMLImageElement;
    private name: string;

    public data: ImageBitmap | null = null;

    constructor(name: string) {
        this.name = name;
        this.img = new Image();
        this.img.onload = () => {
            createImageBitmap(this.img).then((value: ImageBitmap) => {
                this.data = value;
            });
        };
        this.img.src = Sprite.getSrcByName(name);
    }

    draw(ctx: CanvasRenderingContext2D, viewport: Viewport, useScale: boolean): void {
        console.error('not implemented');
    }

    static getSrcByName(name: string): string {
        return name;
    }
}