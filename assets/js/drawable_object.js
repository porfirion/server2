"use strict";
var ObjectType;
(function (ObjectType) {
    ObjectType[ObjectType["NPC"] = 0] = "NPC";
    ObjectType[ObjectType["Player"] = 1] = "Player";
    ObjectType[ObjectType["Obstacle"] = 2] = "Obstacle";
})(ObjectType || (ObjectType = {}));
/**
 * Representation of any visible object
 * @param {Number} id
 * @param {Number} size
 * @returns {DrawableObject}
 * @constructor
 */
var DrawableObject = /** @class */ (function () {
    function DrawableObject(id) {
        this.position = new Point2D(0, 0);
        /**
         * список слоёв для отрисовки
         */
        this.layers = [];
        this.layersByName = new Map();
        this.removeLayer = function (name) {
            if (this.layersByName.hasOwnProperty(name)) {
                var layer = this.layersByName[name];
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
        this.id = id;
    }
    DrawableObject.prototype.getId = function () {
        return this.id;
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
        this.layers.push(layer);
        this.layersByName.set(name, layer);
    };
    DrawableObject.prototype.draw = function (ctx) {
        for (var i = 0; i < this.layers.length; i++) {
            this.layers[i].draw(ctx);
        }
        if (this.type === ObjectType.NPC || this.type === ObjectType.Player) {
            // рисуем закрашенный кружок
            ctx.lineWidth = 1;
            ctx.font = '8px serif';
            ctx.fillStyle = this.color;
            ctx.beginPath();
            ctx.arc(0, 0, this.size, 0, Math.PI * 2);
            ctx.fill();
            ctx.strokeStyle = '#777';
            ctx.stroke();
            // рисуем вектор движения
            ctx.lineWidth = 1;
            ctx.strokeStyle = 'blue';
            ctx.beginPath();
            ctx.moveTo(0, 0);
            ctx.lineTo(this.speed.x, 0 - this.speed.y);
            ctx.stroke();
            if (this.type === ObjectType.Player) {
                ctx.fillStyle = 'blue';
                drawTextCentered(ctx, this.player.name, 0, this.size);
            }
            else {
                ctx.fillStyle = 'grey';
                drawTextCentered(ctx, "NPC" + this.id, 0, this.size);
            }
        }
        else {
            // рисуем просто круг
            ctx.lineWidth = 1;
            ctx.strokeStyle = this.color;
            ctx.beginPath();
            ctx.arc(0, 0, this.size, 0, Math.PI * 2);
            // ctx.rect(0 - this.size / 2, 0 - this.size / 2, this.size, this.size);
            ctx.stroke();
        }
        // id объекта
        ctx.font = '10px serif';
        ctx.fillStyle = 'blue';
        drawTextCentered(ctx, this.id, 0, 2);
        // координаты объекта
        var str = Math.round(this.serverPosition.x, 0) + ':' + Math.round(this.serverPosition.y, 0);
        ctx.fillStyle = 'black';
        drawTextCentered(ctx, str, 0, 0 - this.size / 1.5);
        // if (this.active) {
        // 	ctx.save();
        // 	//ctx.globalAlpha = 0.3;
        // 	ctx.fillStyle = this.color;
        // 	ctx.fill();
        // 	ctx.restore();
        // }
    };
    return DrawableObject;
}());
