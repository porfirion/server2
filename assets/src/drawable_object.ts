"use strict";

/**
 * Representation of any visible object
 */
class DrawableObject {
    protected id: number;

	/**
     * Radius of bounding circle
     */
	protected size: number = 10;
	protected position: Point2D = {x: 0, y: 0};

    /**
     * список слоёв для отрисовки
     */
    private layers: Array<DrawableObjectLayer> = [];
    private layersByName: Map<string, DrawableObjectLayer> = new Map<string, DrawableObjectLayer>();

    constructor(id: number) {
    	this.id = id;
	}

	getId(): number {
        return this.id;
    }

    getPosition(): Point2D {
        return this.position;
    }

    setSize(size: number): DrawableObject {
        this.size = size;
        return this;
    }

    /**
     * @param {string} name
     * @param {DrawableObjectLayer} layer
     */
    addLayer(name: string, layer: DrawableObjectLayer): void {
        this.layers.push(layer);
        this.layersByName.set(name, layer);
    }

    removeLayer(name: string) {
        if (this.layersByName.has(name)) {
            let layer = this.layersByName.get(name);
            this.layersByName.delete(name);

            for (let i = 0; i < this.layers.length; i++) {
                if (this.layers[i] === layer) {
                    this.layers.splice(i, 1);
                    return;
                }
            }

            console.warn("can't find layer in list");
        } else {
            console.warn("no layer with name " + name);
        }
    }

    draw(ctx: CanvasRenderingContext2D) {
        for (let i = 0; i < this.layers.length; i++) {
            this.layers[i].draw(ctx);
        }
    }

    getBoundingCircle() {
        return this.size;
    }
}
