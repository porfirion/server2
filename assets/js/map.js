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
	this.myObject = null; // объект текущего игрока
	this.objectsById = {};
	this.gridSize = 100;
	this.prevOffset = 0;
	this.lastCursorPositionReal = null;
	this.players = {};
	this.latency = 0; // пинг до сервера
	this.timeCorrection = 0; // сколько нужно прибавить к текущему времени, чтобы получить серверное
	this.timeCanvas = document.createElement("canvas");
	this.timeCanvas.width = 300;
	this.timeCanvas.height = 100;

	$(elem).on('click', function (event) {
		this.points.push({x: event.offsetX, y: event.offsetY, color: randomColor()});
		var real = this.viewportToReal({x: event.offsetX, y: event.offsetY});

		// for (var i = 0; i < this.objects.length; i++) {
		// 	if (this.distance(real, this.objects[i]) < this.objects[i].size) {
		// 		this.objects[i].active = true;
		// 	}
		// }

		$(_this).trigger('game:click', real);
	}.bind(this));

	$(elem).on('mousemove', function (ev) {
		this.lastCursorPositionReal = this.viewportToReal({x: ev.offsetX, y: ev.offsetY});
	}.bind(this));

	$(elem).on('mouseenter', function () {

	}.bind(this));

	$(elem).on('mouseout', function () {
		this.lastCursorPositionReal = null;
	}.bind(this));

	$(elem).on('mousewheel DOMMouseScroll', function (ev) {
		var params = normalizeWheel(ev.originalEvent);
		if (params.spinY > 0) {
			// на себя
			this.viewport.scale *= 0.97;
			this.viewport.scale = Math.max(this.viewport.scale, 0.005);
		} else {
			// от себя
			this.viewport.scale *= 1.05;
			this.viewport.scale = Math.min(this.viewport.scale, 50);
		}

		this.lastCursorPositionReal = this.viewportToReal({x: ev.offsetX, y: ev.offsetY});
		// capture all scrolling over map
		return false;
	}.bind(this));

	$(elem).on('contextmenu', function (event) {
		this.viewportAdjustPoint = this.viewportToReal({x: event.offsetX, y: event.offsetY});

		return false;
	}.bind(this));

	this.isAnimating = false;
	this.animations = [];
	this.prevAnimationTime = null;

	// this.fillObjects();
}

Map.prototype.fillObjects = function () {
	// static objects
	var i;
	for (i = 0; i < 1000; i++) {
		this.addObject(1000000 + i, ObjectType.Obstacle, {
			x: Math.random() * 10000 - 5000,
			y: Math.random() * 10000 - 5000
		}).setSize(Math.random() * 200).setColor(randomColor());
	}

	var chaos = function (_this) {
		_this.setSpeed({
			x: Math.random() * 100 - 50,
			y: Math.random() * 100 - 50,
		});

		setTimeout(function () {
			chaos(_this);
		}, 20000 * Math.random());
	};

	for (i = 0; i < 10; i++) {
		var obj = this.addObject(1100000 + i, ObjectType.NPC);
		obj.setSize(50);
		chaos(obj);
	}
};

Map.prototype.addObject = function (id, objectType, coords) {
	if (objectType === undefined) {
		objectType = ObjectType.Obstacle;
	}

	if (coords === undefined || coords === null) {
		coords = {x: 0, y: 0};
	}
	var obj = new MapObject(id, objectType, coords);

	this.objects.push(obj);
	this.objectsById[obj.id] = obj;

	return obj;
};
Map.prototype.removeObject = function (obj) {
	var index = this.objects.indexOf(obj);

	this.objects.splice(index, 1);
	delete this.objectsById[obj.id];
};

Map.prototype.updateObjectPosition = function (obj, time) {
	var mapObject = null;

	if (this.objectsById.hasOwnProperty(obj.id)) {
		mapObject = this.objectsById[obj.id];
	} else {
		mapObject = this.addObject(obj.id, obj.objectType);
		if (obj.userId) {
			var player = this.players[obj.userId];
			mapObject.setPlayer(player);

			if (player.isMe) {
				this.myObject = mapObject;
			}
		}
	}

	if (mapObject) {
		mapObject.adjustState(obj, time);
	}
};

Map.prototype.addPlayer = function (player) {
	// добавление объекта происходит при синхронизации объектов
	// var obj = this.addObject(ObjectType.Player, player.state.position);

	// obj.setPlayer(player);
	// obj.setSize(100);

	this.players[player.id] = player;
};

Map.prototype.removePlayer = function (playerId) {
	// this.removeObject(this.players[playerId]);

	delete this.players[playerId];
	for (var i = 0; i < this.objects.length; i++) {
		console.log(this.objects[i]);
		if (this.objects[i].hasOwnProperty('player') && this.objects[i].player != null && this.objects[i].player.id == playerId) {
			this.objects.splice(i, 1);
			break;
		}
	}
};

Map.prototype.draw = function () {
	if (this.isAnimating) {
		return;
	}

	this.isAnimating = true;
	this.draw_();
};

Map.prototype.draw_ = function () {
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
	if (this.prevAnimationTime !== null) {
		this.animations.push(now - this.prevAnimationTime);

		if (this.animations.length > 100) {
			this.animations.shift();
		}
	}

	this.prevAnimationTime = now;

	this.drawTime();

	window.requestAnimationFrame(this.draw_.bind(this));
};

Map.prototype.drawTime = function () {
	var ctx = this.timeCanvas.getContext('2d');
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

	ctx.fillText('FPS: ' + Math.round(1000 / average, 0), 300 - ctx.measureText('FPS: ' + Math.round(1000 / average, 0)).width - 10, 15);
	ctx.fillText('min: ' + min, 300 - ctx.measureText('min: ' + min).width - 10, 30);
	ctx.fillText('average: ' + Math.round(average, 2), 300 - ctx.measureText('average: ' + Math.round(average, 2)).width - 10, 45);
	ctx.fillText('max: ' + max, 300 - ctx.measureText('max: ' + max).width - 10, 60);

	ctx.fillText('viewport: (x: ' + Math.round(this.viewport.x * 100) / 100 + '; y: ' + Math.round(this.viewport.y * 100) / 100 + ')', 15, 15);
	ctx.fillText('scale: ' + Math.round(this.viewport.scale * 100) / 100, 15, 30);
	ctx.fillText('latency: ' + this.latency.toFixed(0) + ' ms', 15, 45);
	ctx.fillText('time correction: ' + this.timeCorrection.toFixed(1) + ' ms', 15, 60);

	this.ctx.save();
	this.ctx.globalAlpha = 0.7;
	this.ctx.drawImage(this.timeCanvas, 0, 0, 300, 100);
	this.ctx.restore();
};

Map.prototype.adjustViewport = function () {
	// disable viewport adjustement
	// return;

	if (this.myObject != null) {

		var viewportW = this.elem.width * 0.3 / this.viewport.scale;
		var viewportH = this.elem.height * 0.3 / this.viewport.scale;

		var serverTime = Date.now() + this.timeCorrection;
		var objPosReal = this.myObject.getApproximatedPosition(serverTime);

		if (objPosReal.x < this.viewport.x - viewportW / 2) { this.viewport.x = objPosReal.x + viewportW / 2; }
		if (objPosReal.x > this.viewport.x + viewportW / 2) { this.viewport.x = objPosReal.x - viewportW / 2; }
		if (objPosReal.y < this.viewport.y - viewportH / 2) { this.viewport.y = objPosReal.y + viewportH / 2; }
		if (objPosReal.y > this.viewport.y + viewportH / 2) { this.viewport.y = objPosReal.y - viewportH / 2; }

		// пока что жёстко отключим перемещение вьюпорта.
		// в дальнейшем надо бы его переделать.
		// иногда случается, что пользователь начал движение в другую сторону и вьюпорт дёргается немного при развороте
		// сильнее всего чувствуется на мобильнике
		return;
	}

	if (this.viewportAdjustPoint == null) {
		return;
	}

	var dx = (this.viewportAdjustPoint.x - this.viewport.x) * 0.02;
	var dy = (this.viewportAdjustPoint.y - this.viewport.y) * 0.02;

	if (Math.abs(dx) + Math.abs(dy) < 0.01 / this.viewport.scale) {
		this.viewport.x = this.viewportAdjustPoint.x;
		this.viewport.y = this.viewportAdjustPoint.y;
		this.viewportAdjustPoint = null;
		return;
	}

	this.viewport.x += dx;
	this.viewport.y += dy;

	this.viewport.x = Math.min(5000, Math.max(-5000, this.viewport.x));
	this.viewport.y = Math.min(5000, Math.max(-5000, this.viewport.y));
};

Map.prototype.rectContainsPoint = function (rect, point, radius) {
	return rect.left <= (point.x + radius) && point.x - radius <= rect.right &&
		rect.top <= (point.y + radius) && (point.y - radius) <= rect.bottom;
};

Map.prototype.drawObjects = function () {
	var ctx = this.ctx;
	var viewportReal = this.getRealViewport(); // real position and size of viewport

	var serverTime = Date.now() + this.timeCorrection;

	ctx.save();
	ctx.scale(this.viewport.scale, this.viewport.scale);
	ctx.lineWidth = 1;

	for (var i = 0; i < this.objects.length; i++) {
		var obj = this.objects[i];

		var objPosReal = obj.getApproximatedPosition(serverTime);
		var objHalfSizeReal = obj.size / 2; // half of object size

		var isVisible = this.rectContainsPoint(viewportReal, objPosReal, objHalfSizeReal)
			|| this.rectContainsPoint(viewportReal, obj.getLastServerPosition(), objHalfSizeReal);
		// var isVisible = viewportReal.left <= objPosReal.x + objHalfSizeReal && objPosReal.x - objHalfSizeReal <= viewportReal.right &&
		// 	viewportReal.top <= (objPosReal.y + objHalfSizeReal) && (objPosReal.y - objHalfSizeReal) <= viewportReal.bottom;

		if (!isVisible) {
			continue;
		}

		var objPosViewport = this.realToViewport(objPosReal); // viewport position of object
		ctx.save();
		ctx.translate(objPosViewport.x, objPosViewport.y);
		obj.draw(ctx);
		ctx.restore();

		// рисуем текущее положение объекта по серверу
		ctx.save();
		var serverPos = this.realToViewport(obj.getLastServerPosition());
		ctx.translate(serverPos.x, serverPos.y);
		ctx.lineWidth = 1;
		ctx.setLineDash([4, 2]);
		ctx.beginPath();
		ctx.arc(0, 0, objHalfSizeReal, 0, Math.PI * 2);
		ctx.strokeStyle = '#aaaaaa';
		ctx.closePath();
		ctx.stroke();
		ctx.restore();
	}
	ctx.restore();
};

Map.prototype.drawGrid = function () {
	var viewportW = this.elem.width / this.viewport.scale;
	var viewportH = this.elem.height / this.viewport.scale;

	var realViewport = this.getRealViewport();

	var ctx = this.ctx;

	var leftCol = Math.ceil((realViewport.left) / this.gridSize) * this.gridSize;
	var colCount = Math.max(Math.ceil(realViewport.width / this.gridSize), 1);
	var topRow = Math.ceil((realViewport.top) / this.gridSize) * this.gridSize;
	var rowCount = Math.max(Math.ceil(realViewport.height / this.gridSize), 1);

	ctx.save();
	ctx.scale(this.viewport.scale, this.viewport.scale);
	ctx.strokeStyle = '#ccc';

	// рисуем вертикали
	for (var i = 0; i < colCount; i++) {
		var x = this.xRealToViewport(leftCol + i * this.gridSize);
		if (leftCol + i * this.gridSize === 0) {
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
	for (var j = 0; j < rowCount; j++) {
		var y = this.yRealToViewport(topRow + j * this.gridSize);
		if (topRow + j * this.gridSize === 0) {
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
	if (this.lastCursorPositionReal) {
		ctx.save();

		ctx.strokeStyle = 'magenta';
		ctx.lineWidth = 2;
		ctx.setLineDash([12, 6]);
		this.prevOffset = (this.prevOffset + 0.5) % 18;
		ctx.lineDashOffset = this.prevOffset;

		var cursorViewport = this.realToViewport(this.lastCursorPositionReal);

		ctx.beginPath();
		ctx.arc(cursorViewport.x, cursorViewport.y, 20, 0, Math.PI * 2);
		ctx.stroke();

		ctx.restore();
	}


	// рисуем центр
	ctx.strokeStyle = 'lime';
	ctx.beginPath();
	// ctx.ellipse(viewportW / 2, viewportH / 2, 10, 10, 0, 0, Math.PI * 2);
	ctx.moveTo(viewportW / 2 - 15, viewportH / 2);
	ctx.lineTo(viewportW / 2 + 15, viewportH / 2);
	ctx.moveTo(viewportW / 2, viewportH / 2 - 15);
	ctx.lineTo(viewportW / 2, viewportH / 2 + 15);
	// ctx.arc(viewportW / 2, viewportH / 2, 10, 0, Math.PI * 2);
	ctx.stroke();

	// рисуем границы области
	ctx.lineWidth = 10;
	ctx.beginPath();
	var vlt = this.realToViewport({x: -5000, y: -5000});
	var vrb = this.realToViewport({x: 5000, y: 5000});
	ctx.rect(vlt.x, vlt.y, vrb.x - vlt.x, vrb.y - vlt.y);
	ctx.stroke();

	// ctx.globalAlpha = 0.6;

	ctx.restore();

	// Выводим размеры вьюпорта
	var l = Math.round(realViewport.left);
	var t = Math.round(realViewport.top);
	var r = Math.round(realViewport.right);
	var b = Math.round(realViewport.bottom);
	ctx.fillStyle = 'black';
	ctx.font = '14px serif';
	ctx.fillText(t, this.elem.width / 2 - ctx.measureText(t).width / 2, 10);
	ctx.fillText(b, this.elem.width / 2 - ctx.measureText(b).width / 2, this.elem.height);
	ctx.fillText(l, 0, this.elem.height / 2 + 3);
	ctx.fillText(r, this.elem.width - ctx.measureText(r).width, this.elem.height / 2 + 3);
};

/**
 * Бесполезная штука, которая рисует линии для отладки
 */
Map.prototype.drawAnchors = function() {
	var ctx = this.ctx;
	var elem = this.elem;

	ctx.save();

	ctx.scale(this.viewport.scale, this.viewport.scale);

	ctx.globalAlpha = 0.7;

	var realViewport = this.getRealViewport();

	// рисуем рамки чуть меньше вьюпорта
	ctx.beginPath();
	ctx.rect(10, 10, realViewport.width - 20, realViewport.height - 20);
	ctx.strokeStyle = 'green';
	ctx.closePath();
	ctx.stroke();

	// рисуем рамку ровно в размер вьюпорта
	ctx.beginPath();
	ctx.rect(0, 0, realViewport.width, realViewport.height);
	ctx.closePath();
	ctx.stroke();

	ctx.beginPath();
	ctx.strokeStyle = 'blue';

	ctx.moveTo(100, 0);
	ctx.lineTo(300, 150);
	ctx.moveTo(100, 0);
	ctx.lineTo(150, 300);
	ctx.moveTo(100, 0);
	ctx.lineTo(300, 300);
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
	var gradient = ctx.createLinearGradient(0, 0, elem.width, elem.height);
	gradient.addColorStop(0, "yellow");
	gradient.addColorStop(0.3, "blue");
	gradient.addColorStop(0.7, "red");
	gradient.addColorStop(1, "purple");
	ctx.strokeStyle = gradient;
	ctx.moveTo(0, 0);
	ctx.lineTo(realViewport.width, realViewport.height);
	ctx.arc(realViewport.width, realViewport.height, 50, 50, 0);
	ctx.stroke();

	ctx.strokeText('width*height', realViewport.width, realViewport.height);

	ctx.beginPath();
	// ctx.ellipse(elem.width, elem.height, 20, 20, 0, 0, Math.PI * 2);
	ctx.arc(elem.width / this.viewport.scale, elem.height / this.viewport.scale, 20, 0, Math.PI * 2);
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
};

/**
 * Рассчитывает X на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param xr реальная координата по X
 * @returns {number} координата X на канве
 */
Map.prototype.xRealToViewport = function (xr) {
	var viewportW = this.elem.width / this.viewport.scale;
	return viewportW / 2 + (xr - this.viewport.x);
};

/**
 * Рассчитывает Y на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param yr реальная координата Y
 * @returns {number} координата Y на канве
 */
Map.prototype.yRealToViewport = function (yr) {
	var viewportH = this.elem.height / this.viewport.scale;
	return viewportH / 2 - (yr - this.viewport.y);
};

/**
 * Рассчитывает положение объекта на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param pos реальное положение объекта
 * @returns {{x: number, y: number}} координаты объекта на канве
 */
Map.prototype.realToViewport = function (pos) {
	var viewportW = this.elem.width / this.viewport.scale;
	var viewportH = this.elem.height / this.viewport.scale;
	return {
		x: viewportW / 2 + (pos.x - this.viewport.x),
		y: viewportH / 2 - (pos.y - this.viewport.y)
	};
};

/**
 * рассчитывает с учётом скейла
 * @param pos положение курсора на канве
 * @returns {{x: number, y: number}} реальное положение курсора
 */
Map.prototype.viewportToReal = function (pos) {
	//реальные размеры вьюпорта в пикселях
	var viewportW = this.elem.width / this.viewport.scale;
	var viewportH = this.elem.height / this.viewport.scale;

	return {
		x: this.viewport.x + (pos.x / this.viewport.scale - viewportW / 2),
		y: this.viewport.y - (pos.y / this.viewport.scale - viewportH / 2)
	};
};


/**
 * Возвращает реальные параметры вьюпорта
 * @returns {{x: number, y: number, width: number, height: number, left: number, top: number, right: number, bottom: number}}
 */
Map.prototype.getRealViewport = function () {
	// размер вьюпорта в пикселях
	var viewportW = this.elem.width / this.viewport.scale;
	var viewportH = this.elem.height / this.viewport.scale;

	return {
		// реальное положение и размер вьюпорта
		x: this.viewport.x,
		y: this.viewport.y,
		width: viewportW,
		height: viewportH,

		// для удобства отдаём ещё и реальные границы
		left: this.viewport.x - viewportW / 2,
		top: this.viewport.y - viewportH / 2,
		right: this.viewport.x + viewportW / 2,
		bottom: this.viewport.y + viewportH / 2
	};
};

Map.prototype.distance = function (a, b) {
	return Math.sqrt(Math.pow(a.x - b.x, 2) + Math.pow(a.y - b.y, 2));
};

Map.prototype.drawTextCentered = function (ctx, text, x, y) {
	var measure = ctx.measureText(text);
	ctx.fillText(text, x - measure.width / 2, y);
};

function randomComponent() {
	var comp = (Math.round(Math.random() * 255)).toString(16);
	if (comp.length < 2) {
		comp = '0' + comp;
	}
	return comp;
}

function randomColor() {
	return '#' + randomComponent() + randomComponent() + randomComponent();
}

// Reasonable defaults
var PIXEL_STEP = 10;
var LINE_HEIGHT = 40;
var PAGE_HEIGHT = 800;

function normalizeWheel(/*object*/ event) /*object*/ {
	var sX = 0, sY = 0,       // spinX, spinY
		pX = 0, pY = 0;       // pixelX, pixelY

	// Legacy
	if ('detail' in event) {
		sY = event.detail;
	}
	if ('wheelDelta' in event) {
		sY = -event.wheelDelta / 120;
	}
	if ('wheelDeltaY' in event) {
		sY = -event.wheelDeltaY / 120;
	}
	if ('wheelDeltaX' in event) {
		sX = -event.wheelDeltaX / 120;
	}

	// side scrolling on FF with DOMMouseScroll
	if ('axis' in event && event.axis === event.HORIZONTAL_AXIS) {
		sX = sY;
		sY = 0;
	}

	pX = sX * PIXEL_STEP;
	pY = sY * PIXEL_STEP;

	if ('deltaY' in event) {
		pY = event.deltaY;
	}
	if ('deltaX' in event) {
		pX = event.deltaX;
	}

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
	if (pX && !sX) {
		sX = (pX < 1) ? -1 : 1;
	}
	if (pY && !sY) {
		sY = (pY < 1) ? -1 : 1;
	}

	return {
		spinX: sX,
		spinY: sY,
		pixelX: pX,
		pixelY: pY
	};
}
