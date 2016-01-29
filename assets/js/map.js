"use strict";

var canvas;

function Map(elem) {
	canvas = elem;
	this.elem = elem;
	this.ctx = this.elem.getContext("2d");
	var _this = this;

	this.viewport = {
		x: 0,
		y: 0,
		scale: 3,
	};
	this.points = [];
	this.objects = [];

	this.gridSize = 100;

	this.prevOffset = 0;

	this.lastCursorPosition = null;

	this.players = {};

	$(elem).on('click', function(event) {
		this.points.push({x: event.offsetX, y: event.offsetY, color: randomColor()});

		var real = this.viewportToReal({x: event.offsetX, y: event.offsetY});

		// for (var i = 0; i < this.objects.length; i++) {
		// 	if (this.distance(real, this.objects[i]) < this.objects[i].size) {
		// 		this.objects[i].active = true;
		// 	}
		// }
		
		$(_this).trigger('game:click', real);
	}.bind(this));

	$(elem).on('mousemove', function(ev) {
		this.lastCursorPosition = {x: ev.offsetX, y: ev.offsetY};
	}.bind(this));

	$(elem).on('mouseenter', function(ev) {
		
	}.bind(this));

	$(elem).on('mouseout', function(ev) {
		// this.lastCursorPosition = {x: this.elem.width / 2, y: this.elem.height / 2};
		this.lastCursorPosition = null;
	}.bind(this));

	$(elem).on('mousewheel DOMMouseScroll', function(ev) {
		var params = normalizeWheel(ev.originalEvent);
		if (params.spinY > 0) {
			// на себя
			this.viewport.scale *= 1.1;
			this.viewport.scale = Math.min(this.viewport.scale, 50);
		} else {
			// от себя
			this.viewport.scale *= 0.95;
			this.viewport.scale = Math.max(this.viewport.scale, 0.005);
		}
		return false;
	}.bind(this));

	this.isAnimating = false;
	this.animations = [];	
	this.prevAnimationTime = null;

	this.fillObjects();
}

Map.prototype.fillObjects = function() {
	// static objects
	var i;
	for (i = 0; i < 1000; i++) {
		this.addObject(ObjectType.Obstacle, {
					x: Math.random() * 10000 - 5000,
					y: Math.random() * 10000 - 5000
			}
		).setSize(Math.random() * 200).setColor(randomColor());
	}

	var chaos = function(_this) {
		_this.setDirection({
			x: Math.random() * 100 - 50,
			y: Math.random() * 100 - 50,
		});
		_this.setSpeed(Math.random() * 80 + 20);

		setTimeout(function() { chaos(_this); }, 20000 * Math.random());
	};

	for (i = 0; i < 10; i++) {
		var obj = this.addObject(ObjectType.NPC);
		obj.setSize(50);
		chaos(obj);
	}
}

Map.prototype.addObject = function(objectType, coords) {
	if (typeof objectType == 'undefined') {
		objectType = ObjectType.Obstacle;
	}

	if (typeof coords == 'undefined' || coords == null) {
		coords = { x: 0, y: 0};		
	}
	var obj = new MapObject(objectType, coords);

	this.objects.push(obj);

	return obj;
}
Map.prototype.removeObject = function(obj) {
	var index = this.objects.indexOf(obj);

	this.objects.splice(index, 1);
}

Map.prototype.addPlayer = function(player) {
	var obj = this.addObject(ObjectType.Player, player.state.position);

	obj.setPlayer(player);
	obj.setSize(100);

	this.players[player.id] = obj;
	return obj;
}

Map.prototype.removePlayer = function(playerId) {
	this.removeObject(this.players[playerId]);
}

Map.prototype.draw = function() {
	if (this.isAnimating) {
		return;
	}
	
	this.isAnimating = true;
	this._draw();
}

Map.prototype._draw = function() {
	var ctx = this.ctx;
	var elem = this.elem;

	elem.width = elem.clientWidth;
	elem.height = elem.clientHeight;

	// у контекста нет ширины и высоты - они есть только у элемента canvas
	// ctx.width = this.elem.clientWidth;
	// ctx.height = this.elem.clientHeight;

	ctx.font = "16px serif";

	this.adjustViewport();

	this.drawGrid();

	this.drawObjects();

	// this.drawAnchors();

	var now = Date.now();
	if (this.prevAnimationTime != null) {
		this.animations.push(now - this.prevAnimationTime);

		if (this.animations.length > 100) {
			this.animations.shift();
		}
	}

	this.prevAnimationTime = now;

	this.drawTime();

	window.requestAnimationFrame(this._draw.bind(this));
}

Map.prototype.drawTime = function() {
	var ctx = this.ctx;
	ctx.save();

	ctx.globalAlpha = 0.7;
	ctx.fillStyle = 'white';
	ctx.fillRect(this.elem.width - 300, 0, 300, 100);
	// ctx.clearRect(this.elem.width - 300, 0, 300, 100);

	var min = Infinity;
	var max = -Infinity;
	var average = 0;

	ctx.strokeStyle = '1px black';
	

	for (var i = 0; i < this.animations.length; i++) {
		average += this.animations[i];
		if (this.animations[i] > max)
			max = this.animations[i];
		if (this.animations[i] < min)
			min = this.animations[i];

		ctx.beginPath();
		ctx.moveTo(this.elem.width - 300 + i * 3, 100);
		ctx.lineTo(this.elem.width - 300 + i * 3, 100 - this.animations[i]);
		ctx.stroke();
	}
	average = average / this.animations.length;

	ctx.fillStyle = 'black';

	ctx.fillText('fps: ' + Math.round(1000 / average, 0), this.elem.width - ctx.measureText('fps: ' + Math.round(1000 / average, 0)).width - 10, 15);
	ctx.fillText('min: ' + min, this.elem.width - ctx.measureText('min: ' + min).width - 10, 30);
	ctx.fillText('average: ' + Math.round(average, 2), this.elem.width - ctx.measureText('average: ' + Math.round(average, 2)).width - 10, 45);
	ctx.fillText('max: ' + max, this.elem.width - ctx.measureText('max: ' + max).width - 10, 60);

	ctx.fillText('viewport: (x: ' + Math.round(this.viewport.x * 100) / 100 + '; y: ' + Math.round(this.viewport.y * 100) / 100 + ')', this.elem.width - 285, 15);
	ctx.fillText('scale: ' + Math.round(this.viewport.scale * 100) / 100, this.elem.width - 285, 30);

	ctx.restore();
}

Map.prototype.adjustViewport = function() {
	if (this.lastCursorPosition == null) 
		return;

	this.viewport.x += (this.lastCursorPosition.x - this.elem.width / 2) * this.viewport.scale * 0.02;
	this.viewport.y -= (this.lastCursorPosition.y - this.elem.height / 2) * this.viewport.scale * 0.02;

	this.viewport.x = Math.min(5000, Math.max(-5000, this.viewport.x));
	this.viewport.y = Math.min(5000, Math.max(-5000, this.viewport.y));
}


Map.prototype.drawObjects = function() {
	var ctx = this.ctx;

	var real = this.getRealViewport();

	ctx.save();
	ctx.lineWidth = 1;
	for (var i = 0; i < this.objects.length; i++) {
		var obj = this.objects[i];
		var pos = obj.getPos();
		var s2 = obj.size / 2;
		if (pos.x + s2 > real.x - real.w / 2 && real.x + real.w / 2 > pos.x - s2
			&& pos.y + s2 > real.y - real.h / 2 && real.y + real.h / 2 > pos.y - s2) {

			var vp = this.realToViewport(pos);
			var os = (obj.size) / this.viewport.scale;
			
			if (obj.type == ObjectType.NPC || obj.type == ObjectType.Player) {
				ctx.lineWidth = 2;
				ctx.fillStyle = obj.color;
				ctx.beginPath();
				ctx.arc(vp.x, vp.y, os, 0, Math.PI * 2);
				// ctx.rect(vp.x - os / 2, vp.y - os / 2, os, os);
				ctx.fill();
				ctx.strokeStyle = 'yellow';
				ctx.stroke();

				// ctx.strokeStyle = 'yellow';
				// ctx.beginPath();
				// ctx.rect(vp.x - os / 2, vp.y - os / 2, os, os);
				// ctx.stroke();

				ctx.lineWidth = 1;
				ctx.strokeStyle = 'blue';
				ctx.beginPath();
				ctx.moveTo(vp.x, vp.y);
				var len = obj.speed / this.viewport.scale * 10;
				ctx.lineTo(vp.x + obj.direction.x * len, vp.y - obj.direction.y * len);
				ctx.stroke();

				if (obj.type == ObjectType.Player) {
					var measure = ctx.measureText(obj.player.name);
					ctx.strokeText(obj.player.name, vp.x - measure.width / 2, vp.y - os / 2);
				}

			} else {
				ctx.lineWidth = 1;
				ctx.strokeStyle = obj.color;
				ctx.beginPath();
				ctx.rect(vp.x - os / 2, vp.y - os / 2, os, os);
				ctx.stroke();
			}
			
			if (obj.active) {
				ctx.save();
				//ctx.globalAlpha = 0.3;
				ctx.fillStyle = obj.color;
				ctx.fill();
				ctx.restore();
			}
		}
	}
	ctx.restore();
}

Map.prototype.drawGrid = function() {
	var viewportW = this.elem.width;
	var viewportH = this.elem.height;

	var real = this.getRealViewport();
	var viewGridSize = this.gridSize / this.scale;

	var ctx = this.ctx;

	var leftCol = Math.ceil((real.x - real.w / 2) / this.gridSize) * this.gridSize;
	var topRow = Math.ceil((real.y - real.h / 2) / this.gridSize) * this.gridSize;

	ctx.save();

	// рисуем вертикали
	ctx.strokeStyle = '#ccc';
	for (var i = 0; i < (real.w / this.gridSize); i++) {
		var x = this.xRealToView(leftCol + i * this.gridSize);
		if (leftCol + i * this.gridSize == 0) {
			ctx.strokeStyle = '#888';
		} else {
			ctx.strokeStyle = '#ccc';
		}
		ctx.beginPath();
		ctx.moveTo(x, 0);
		ctx.lineTo(x, viewportH);
		ctx.stroke();
	}

	// рисуем горизонтали
	for (var j = 0; j < (real.h / this.gridSize); j++) {
		var y = this.yRealToView(topRow + j * this.gridSize);
		if (topRow + j * this.gridSize == 0) {
			ctx.strokeStyle = '#888';
		} else {
			ctx.strokeStyle = '#ccc';
		}
		ctx.beginPath();
		ctx.moveTo(0, y);
		ctx.lineTo(viewportW, y);
		ctx.stroke();	
	}

	// рисуем курсор
	if (this.lastCursorPosition) {
		ctx.save();
		ctx.strokeStyle = 'magenta';
		ctx.lineWidth = 2;
		ctx.setLineDash([10, 5]);
		this.prevOffset = (this.prevOffset + 0.5) % 15;
		ctx.lineDashOffset = this.prevOffset;

		ctx.beginPath();
		// ctx.ellipse(this.lastCursorPosition.x, this.lastCursorPosition.y, 20, 20, 0, 0, Math.PI * 2);
		ctx.arc(this.lastCursorPosition.x, this.lastCursorPosition.y, 20, 0, Math.PI * 2);
		ctx.stroke();
		ctx.restore();
	}
	
	

	// рисуем центр
	ctx.strokeStyle = 'lime';
	ctx.beginPath();
	// ctx.ellipse(viewportW / 2, viewportH / 2, 10, 10, 0, 0, Math.PI * 2);
	ctx.arc(viewportW / 2, viewportH / 2, 10, 0, Math.PI * 2);
	ctx.stroke();

	// рисуем границы области
	ctx.beginPath();
	var vlt = this.realToViewport({x: -5000, y: -5000});
	var vrb = this.realToViewport({x: 5000, y: 5000});
	ctx.rect(vlt.x, vlt.y, vrb.x - vlt.x, vrb.y - vlt.y);
	ctx.stroke();

	// ctx.globalAlpha = 0.6;
	ctx.fillStyle = 'black';
	ctx.font = '14px serif';
	var real = this.getRealViewport();
	var l = Math.round(real.x - real.w / 2);
	var t = Math.round(real.y + real.h / 2);
	var r = Math.round(real.x + real.w / 2);
	var b = Math.round(real.y - real.h / 2);
	
	ctx.fillText(t, this.elem.width / 2 - ctx.measureText(t).width / 2, 10);
	ctx.fillText(b, this.elem.width / 2 - ctx.measureText(b).width / 2, this.elem.height);

	
	ctx.fillText(l, 0, this.elem.height / 2 + 3);
	ctx.fillText(r, this.elem.width - ctx.measureText(r).width, this.elem.height / 2 + 3);

	ctx.restore();

}

Map.prototype.drawAnchors = function() {
	var ctx = this.ctx;
	var elem = this.elem;

	ctx.save();

	ctx.globalAlpha = 0.2;

	ctx.beginPath();
	ctx.rect(10, 10, elem.width - 20, elem.height - 20);
	ctx.strokeStyle = 'green';
	ctx.closePath();
	ctx.stroke();

	ctx.beginPath();
	ctx.rect(0, 0, elem.width, elem.height);
	ctx.closePath();
	ctx.stroke();

	ctx.beginPath();
	ctx.strokeStyle = 'blue';

	ctx.moveTo(100,0);
	ctx.lineTo(300,150);
	ctx.moveTo(100,0);
	ctx.lineTo(150, 300);
	ctx.moveTo(100,0);
	ctx.lineTo(300,300);
	ctx.closePath();
	ctx.stroke();

	ctx.fillText('300*150', 300, 175);
	ctx.fillText('150*300', 150, 325);
	ctx.fillText('300*300', 300, 325);

	ctx.beginPath();
	// ctx.ellipse(300, 150, 10, 10, 0, 0, Math.PI * 2);
	ctx.arc(300, 150, 10, 0, Math.PI * 2);
	ctx.stroke();
	ctx.beginPath();
	// ctx.ellipse(150, 300, 10, 10, 0, 0, Math.PI * 2);
	ctx.arc(150, 300, 10, 0, Math.PI * 2);
	ctx.stroke();
	ctx.beginPath();
	// ctx.ellipse(300, 300, 10, 10, 0, 0, Math.PI * 2);
	ctx.arc(300, 300, 10, 0, Math.PI * 2);
	ctx.stroke();

	ctx.beginPath();
	ctx.strokeStyle = 'red';
	var gradient = ctx.createLinearGradient(0,0,elem.width,elem.height);
	gradient.addColorStop(0,"yellow");
	gradient.addColorStop(0.3,"blue");
	gradient.addColorStop(0.7,"red");
	gradient.addColorStop(1,"purple");
	ctx.strokeStyle = gradient;
	ctx.moveTo(0, 0);
	ctx.lineTo(elem.width, elem.height);

	// ctx.ellipse(elem.width, elem.height, 50, 50, 50, 50, 0);
	ctx.arc(elem.width, elem.height, 50, 50, 0);
	ctx.stroke();

	ctx.strokeText('width*height', elem.width, elem.height);

	ctx.beginPath();
	// ctx.ellipse(elem.width, elem.height, 20, 20, 0, 0, Math.PI * 2);
	ctx.arc(elem.width, elem.height, 20, 0, Math.PI * 2);
	ctx.stroke();
	
	ctx.beginPath();
	ctx.strokeStyle = 'black';
	ctx.rect(1, 1, 300, 300);
	ctx.stroke();

	// console.log(this.points.length);

	for (var i = 0; i < this.points.length; i++) {
		ctx.beginPath();
		ctx.strokeStyle = this.points[i].color;
		// ctx.ellipse(this.points[i].x, this.points[i].y, 5, 5, 0, 0, Math.PI * 2);
		ctx.arc(this.points[i].x, this.points[i].y, 5, 0, Math.PI * 2);
		// ctx.closePath();
		ctx.stroke();
	}

	ctx.restore();
}

Map.prototype.xRealToView = function(xr) {
	var viewportW = this.elem.width;
	return viewportW / 2 + (xr - this.viewport.x) * (1 / this.viewport.scale);
}

Map.prototype.yRealToView = function(yr) {
	var viewportH = this.elem.height;
	return viewportH / 2 - (yr - this.viewport.y) * (1 / this.viewport.scale)
}

Map.prototype.viewportToReal = function(pos) {
	var viewportW = this.elem.width;
	var viewportH = this.elem.height;
	var realPos = {
		x: this.viewport.x + (pos.x - viewportW / 2) * this.viewport.scale,
		y: this.viewport.y - (pos.y - viewportH / 2) * this.viewport.scale
	}
	return realPos;
}

Map.prototype.realToViewport = function(pos) {
	var viewportW = this.elem.width;
	var viewportH = this.elem.height;
	var viewPos = {
		x: viewportW / 2 + (pos.x - this.viewport.x) * (1 / this.viewport.scale),
		y: viewportH / 2 - (pos.y - this.viewport.y) * (1 / this.viewport.scale)
	}

	return viewPos;
}

Map.prototype.getRealViewport = function() {
	var viewportW = this.elem.width;
	var viewportH = this.elem.height;

	return {
		x: this.viewport.x, 
		y: this.viewport.y, 
		w: viewportW * this.viewport.scale,
		h: viewportH * this.viewport.scale,
	}
}

Map.prototype.distance = function(a, b) {
	return Math.sqrt(Math.pow(a.x - b.x, 2) + Math.pow(a.y - b.y, 2));
}

function randomComponent() {
	var comp = (Math.round(Math.random() * 255)).toString(16);
	if (comp.length < 2) {
		comp = '0' + comp;
	}
	return comp;
}

function randomColor() {
	return '#'
		+ randomComponent()
		+ randomComponent()
		+ randomComponent()
}

// Reasonable defaults
var PIXEL_STEP  = 10;
var LINE_HEIGHT = 40;
var PAGE_HEIGHT = 800;

function normalizeWheel(/*object*/ event) /*object*/ {
  var sX = 0, sY = 0,       // spinX, spinY
      pX = 0, pY = 0;       // pixelX, pixelY

  // Legacy
  if ('detail'      in event) { sY = event.detail; }
  if ('wheelDelta'  in event) { sY = -event.wheelDelta / 120; }
  if ('wheelDeltaY' in event) { sY = -event.wheelDeltaY / 120; }
  if ('wheelDeltaX' in event) { sX = -event.wheelDeltaX / 120; }

  // side scrolling on FF with DOMMouseScroll
  if ( 'axis' in event && event.axis === event.HORIZONTAL_AXIS ) {
    sX = sY;
    sY = 0;
  }

  pX = sX * PIXEL_STEP;
  pY = sY * PIXEL_STEP;

  if ('deltaY' in event) { pY = event.deltaY; }
  if ('deltaX' in event) { pX = event.deltaX; }

  if ((pX || pY) && event.deltaMode) {
    if (event.deltaMode == 1) {          // delta in LINE units
      pX *= LINE_HEIGHT;
      pY *= LINE_HEIGHT;
    } else {                             // delta in PAGE units
      pX *= PAGE_HEIGHT;
      pY *= PAGE_HEIGHT;
    }
  }

  // Fall-back if spin cannot be determined
  if (pX && !sX) { sX = (pX < 1) ? -1 : 1; }
  if (pY && !sY) { sY = (pY < 1) ? -1 : 1; }

  return { spinX  : sX,
           spinY  : sY,
           pixelX : pX,
           pixelY : pY };
}
