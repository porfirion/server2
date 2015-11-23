var canvas;

function Map(elem) {
	canvas = elem;
	this.elem = elem;
	this.ctx = this.elem.getContext("2d");

	this.viewport = {
		x: 0,
		y: 0,
		scale: 1,
	};
	this.points = [];
	this.objects = [];

	this.gridSize = 100;

	this.prevOffset = 0;

	this.lastCursorPosition = null;

	$(elem).on('click', function(event) {
		console.log(event);
		this.points.push({x: event.offsetX, y: event.offsetY, color: randomColor()});
	}.bind(this));

	$(elem).on('mousemove', function(ev) {
		this.lastCursorPosition = {x: ev.offsetX, y: ev.offsetY};
	}.bind(this));

	$(elem).on('mouseenter', function(ev) {
		
	}.bind(this));

	$(elem).on('mouseout', function(ev) {
		this.lastCursorPosition = {x: this.elem.width / 2, y: this.elem.height / 2};
	}.bind(this));

	$(elem).on('mousewheel DOMMouseScroll', function(ev) {
		var params = normalizeWheel(ev.originalEvent);
		if (params.spinY > 0) {
			// на себя
			this.viewport.scale *= 1.1;
		} else {
			// от себя
			this.viewport.scale *= 0.95;
		}
		console.log(params);
		// console.log(ev.wheelDelta);
		// console.log(ev.detail);
		// console.log(ev.originalEvent.wheelDelta);
		// console.log(ev.originalEvent.detail);
		// console.log(ev);
		return false;
	}.bind(this));

	this.isAnimating = false;

	this.objects.push({
		x: 0, 
		y: 0,
		color: 'red',
		size: 20,
	});
	this.objects.push({
		x: 100, 
		y: 0,
		color: 'green',
		size: 50,
	});
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

	if (this.lastCursorPosition == null) {
		this.lastCursorPosition = {x: this.elem.width / 2, y: this.elem.height / 2};
	}

	elem.width = elem.clientWidth;
	elem.height = elem.clientHeight;

	// у контекста нет ширины и высоты - они есть только у элемента canvas
	// ctx.width = this.elem.clientWidth;
	// ctx.height = this.elem.clientHeight;

	// console.log(ctx, this.elem.width, this.elem.height);
	// console.log(ctx, this.elem.clientWidth, this.elem.clientHeight);

	ctx.clearRect(0, 0, ctx.width, ctx.height);
	ctx.font = "16px serif";

	this.adjustViewport();

	this.drawGrid();

	this.drawObjects();

	// this.drawAnchors();

	window.requestAnimationFrame(this._draw.bind(this));
}

Map.prototype.adjustViewport = function() {
	this.viewport.x += (this.lastCursorPosition.x - this.elem.width / 2) * this.viewport.scale * 0.01;
	this.viewport.y -= (this.lastCursorPosition.y - this.elem.height / 2) * this.viewport.scale * 0.01;
}


Map.prototype.drawObjects = function() {
	var ctx = this.ctx;
	ctx.save();
	for (var i = 0; i < this.objects.length; i++) {
		var obj = this.objects[i];
		ctx.beginPath();
		ctx.strokeStyle = obj.color;
		var vp = this.realToViewport(obj);
		var os = (obj.size) / this.viewport.scale;
		ctx.rect(vp.x - os / 2, vp.y - os / 2, os, os);
		ctx.stroke();
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

	ctx.strokeStyle = '#ccc';
	for (var i = 0; i < (real.w / this.gridSize); i++) {
		var x = this.xRealToView(leftCol + i * this.gridSize);
		ctx.beginPath();
		ctx.moveTo(x, 0);
		ctx.lineTo(x, viewportH);
		ctx.stroke();
	}

	for (var j = 0; j < (real.h / this.gridSize); j++) {
		var y = this.yRealToView(topRow + j * this.gridSize);
		ctx.beginPath();
		ctx.moveTo(0, y);
		ctx.lineTo(viewportW, y);
		ctx.stroke();	
	}

	ctx.save();
	ctx.strokeStyle = 'magenta';
	ctx.lineWidth = 2;
	ctx.setLineDash([10, 5]);
	this.prevOffset = (this.prevOffset + 1.0) % 14;
	ctx.lineDashOffset = this.prevOffset;
	ctx.beginPath();
	// ctx.ellipse(this.lastCursorPosition.x, this.lastCursorPosition.y, 20, 20, 0, 0, Math.PI * 2);
	ctx.arc(this.lastCursorPosition.x, this.lastCursorPosition.y, 20, 0, Math.PI * 2);
	ctx.stroke();
	ctx.restore();

	ctx.strokeStyle = 'lime';
	ctx.beginPath();
	// ctx.ellipse(viewportW / 2, viewportH / 2, 10, 10, 0, 0, Math.PI * 2);
	ctx.arc(viewportW / 2, viewportH / 2, 10, 0, Math.PI * 2);
	ctx.stroke();

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

function randomColor() {
	return '#'
		+ (Math.round(Math.random() * 255)).toString(16)
		+ (Math.round(Math.random() * 255)).toString(16)
		+ (Math.round(Math.random() * 255)).toString(16);
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