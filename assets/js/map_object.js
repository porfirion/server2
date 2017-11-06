"use strict";

/**
 * Base class for all map objects
 * @param {Number} id
 * @param {ObjectType} type
 * @param {{x: Number, y: Number}} pos
 * @param {Number} size
 * @param {String} color
 * @returns {MapObject}
 * @constructor
 */
function MapObject(id, type, pos, size, color) {
	this.id = id;
    /**
     * @type {ObjectType}
     */
	this.type = type;

	this.size = 10;

	if (type === ObjectType.Obstacle) {
		this.color = 'lightblue';
	} else if (type === ObjectType.Player) {
		this.color = 'lightblue';
	} else if (type === ObjectType.NPC) {
		this.color = '#cccccc';
		// this.color = 'red';
	} else {
		this.color = color || randomColor();
	}

	this.speed = {x: 0, y: 0};
	this.isAnimating = false;

    /**
	 * Player associated with this object
     * @type {null}
     */
	this.player = null;

    /**
	 * Last actual server position
     * @type {{x: Number, y: Number}}
     */
    this.serverPosition = null;
    /**
	 * Time (server) when object position was actual
     * @type {int}
     */
    this.adjustServerTime = null;

	return this;
}

MapObject.prototype.setSize = function(size) {
	this.size = size;

	return this;
};

MapObject.prototype.setColor = function(color) {
	this.color = color;

	return this;
};

MapObject.prototype.getApproximatedPosition = function(currentTime) {
	var passedTime = (currentTime - this.adjustServerTime) / 1000;
	var newPos = {x: this.serverPosition.x, y: this.serverPosition.y};
	newPos.x += this.speed.x * passedTime;
	newPos.y += this.speed.y * passedTime;
	return newPos;
};

MapObject.prototype.getLastServerPosition = function() {
	return this.serverPosition;
}

MapObject.prototype.setSpeed = function(speed) {
	this.speed = speed;
	this.isMoving = speed.x !== 0 || speed.y !== 0 ? true : false;
};

MapObject.prototype.adjustState = function(obj, time) {
    this.adjustServerTime = time;
	this.serverPosition = obj.position;
	this.setSpeed(obj.speed);
	this.setSize(obj.size);

	// this.posTime = obj.startTime;
	// this.destination = obj.destinationPosition;
	// this.destinationTime = obj.destinationTime;
	// this.direction = obj.direction;
};

MapObject.prototype.setPlayer = function(player) {
	this.player = player;
	if (player.isMe) {
		this.setColor('#FE2D77');
	}
};

/**
 *
 * @param {CanvasRenderingContext2D} ctx
 */
MapObject.prototype.draw = function(ctx) {
	var drawTextCentered = function (ctx, text, x, y) {
		var measure = ctx.measureText(text);
		ctx.fillText(text, x - measure.width / 2, y);
	};

	if (this.type === ObjectType.NPC || this.type === ObjectType.Player) {
		// рисуем закрашенный кружок
		ctx.lineWidth = 1;
		ctx.font = '8px serif';
		ctx.fillStyle = this.color;
		ctx.beginPath();
		ctx.arc(0, 0, this.size / 2, 0, Math.PI * 2);
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
		} else {
			ctx.fillStyle = 'grey';
			drawTextCentered(ctx, "NPC" + this.id, 0, this.size);
		}
	} else {
		// рисуем просто квадрат
		ctx.lineWidth = 1;
		ctx.strokeStyle = this.color;
		ctx.beginPath();
		ctx.arc(0, 0, this.size / 2, 0, Math.PI * 2);
		// ctx.rect(0 - this.size / 2, 0 - this.size / 2, this.size, this.size);
		ctx.stroke();
	}

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
}

const ObjectType = {
	Obstacle: 1,
	NPC: 10,
	Player: 100,
};
