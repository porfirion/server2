"use strict";
/**
 * Representation of any visible object
 */
var DrawableObject = /** @class */ (function () {
    function DrawableObject(id) {
        /**
         * Radius of bounding circle
         */
        this.size = 10;
        this.position = { x: 0, y: 0 };
        /**
         * список слоёв для отрисовки
         */
        this.layers = [];
        this.layersByName = new Map();
        this.id = id;
    }
    DrawableObject.prototype.getId = function () {
        return this.id;
    };
    DrawableObject.prototype.setPosition = function (position) {
        this.position = position;
    };
    DrawableObject.prototype.getPosition = function () {
        return this.position;
    };
    DrawableObject.prototype.setSize = function (size) {
        this.size = size;
        return this;
    };
    /**
     * @param {string} name
     * @param {DrawableObjectLayer} layer
     */
    DrawableObject.prototype.addLayer = function (name, layer) {
        layer.setObject(this);
        this.layers.push(layer);
        this.layersByName.set(name, layer);
    };
    DrawableObject.prototype.removeLayer = function (name) {
        if (this.layersByName.has(name)) {
            var layer = this.layersByName.get(name);
            this.layersByName.delete(name);
            for (var i = 0; i < this.layers.length; i++) {
                if (this.layers[i] === layer) {
                    this.layers.splice(i, 1);
                    return;
                }
            }
            console.warn("can't find layer in list");
        }
        else {
            console.warn("no layer with name " + name);
        }
    };
    DrawableObject.prototype.draw = function (ctx, viewport, useScale) {
        for (var i = 0; i < this.layers.length; i++) {
            this.layers[i].draw(ctx, viewport, useScale);
        }
    };
    DrawableObject.prototype.getBoundingCircle = function () {
        return this.size;
    };
    return DrawableObject;
}());
//# sourceMappingURL=drawable_object.js.map