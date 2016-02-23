function MapObject(id, type, pos, size, color) {
	this.id = id;
	this.type = type;

	this.pos = pos || {x: 0, y: 0};
	this.size = 10;

	if (type == ObjectType.Obstacle) {
		this.color = 'lightblue';
	} else if (type == ObjectType.Player) {
		this.color = 'lime';
	} else if (type == ObjectType.NPC) { 
		this.color = '#cccccc';
		// this.color = 'red';
	} else {
		this.color = color || randomColor();
	}

	this.speed = 0;
	this.posTime = Date.now();

	this.direction = {x: 0, y: 1};

	this.destination = {x: 0, y: 0};
	this.destinationTime = Date.now();

	this.isAnimating = false;

	this.player = null;

	return this;
}

MapObject.prototype.setSize = function(size) {
	this.size = size;

	return this;
}

MapObject.prototype.setColor = function(color) {
	this.color = color;

	return this;
}

MapObject.prototype.getPos = function() {
	if (!this.isMoving) {
		return this.pos;
	} else {
		var now = Date.now();
		var passed = this.speed * (now - this.posTime) / 1000;
		return {
			x: this.pos.x + this.direction.x * passed,
			y: this.pos.y + this.direction.y * passed,
		}
	}
}

MapObject.prototype.setDirection = function(newDirection) {
	var _this = this;

	this.pos = this.getPos();
	this.posTime = Date.now();

	var sum = Math.abs(newDirection.x) + Math.abs(newDirection.y);

	this.direction = {
		x: newDirection.x / sum,
		y: newDirection.y / sum,
	};
}

MapObject.prototype.stop = function() {
	this.pos = this.getPos();
	this.posTime = Date.now();

	this.speed = 0;
	this.isMoving = false;
}

MapObject.prototype.setSpeed = function(speed) {
	this.speed = speed;
	this.isMoving = speed != 0;
}

// MapObject.prototype.adjustState = function(pos, direction, speed, posTime) {
// 	this.pos = pos;
// 	if (typeof direction != 'undefined') 
// 		this.direction = direction;

// 	if (typeof speed != 'undefined') 
// 		this.setSpeed(speed);

// 	if (typeof posTime != 'undefined') 
// 		this.posTime = posTime;
// 	else
// 		// TODO вот это надо бы убрать. синхронизация всегда происходит по времени и мы не должны получать статус без времени
// 		this.posTime = Date.now();
// }


MapObject.prototype.adjustState = function(obj) {
	this.pos = obj.startPosition;
	this.posTime = obj.startTime;

	this.destination = obj.destinationPosition;
	this.destinationTime = obj.destinationTime;

	this.direction = obj.direction;

	this.setSpeed(obj.speed);

	this.serverPosition = obj.position;
}

MapObject.prototype.getViewPos = function(viewport) {

}

MapObject.prototype.setPlayer = function(player) {
	this.player = player;
	if (player.isMe) {
		this.setColor('#FE2D77');
	}

	// $(player).on('change.position', (function() {
	// 	console.log('adjusting player ' + this.player.name + ' position', this);
	// 	this.adjustState(this.player.state.position);
	// 	console.log(this);

	// 	// console.log('player changed position - need to update map object');
	// }).bind(this));
}

var ObjectType = {
	Obstacle: 1,
	NPC: 10,
	Player: 100,
}
