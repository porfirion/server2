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
	var _this = this;

	$(elem).on('click', function(event) {
		_this.points.push({x: event.offsetX, y: event.offsetY, color: randomColor()});
	})
}

Map.prototype.draw = function() {
	var rectSize = 50;
	var ctx = this.ctx;
	var elem = this.elem;

	elem.width = elem.clientWidth;
	elem.height = elem.clientHeight;

	// у контекста нет ширины и высоты - они есть только у элемента canvas
	// ctx.width = this.elem.clientWidth;
	// ctx.height = this.elem.clientHeight;

	// console.log(ctx, this.elem.width, this.elem.height);
	// console.log(ctx, this.elem.clientWidth, this.elem.clientHeight);

	ctx.clearRect(0, 0, ctx.width, ctx.height);
	ctx.font = "16px serif";

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
	ctx.ellipse(300, 150, 10, 10, 0, 0, Math.PI * 2);
	ctx.stroke();
	ctx.beginPath();
	ctx.ellipse(150, 300, 10, 10, 0, 0, Math.PI * 2);
	ctx.stroke();
	ctx.beginPath();
	ctx.ellipse(300, 300, 10, 10, 0, 0, Math.PI * 2);
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

	ctx.ellipse(elem.width, elem.height, 50, 50, 50, 50, 0);
	ctx.stroke();

	ctx.strokeText('width*height', elem.width, elem.height);

	ctx.beginPath();
	ctx.ellipse(elem.width, elem.height, 20, 20, 0, 0, Math.PI * 2);
	ctx.stroke();
	
	ctx.beginPath();
	ctx.strokeStyle = 'black';
	ctx.rect(1, 1, 300, 300);
	ctx.stroke();

	// console.log(this.points.length);

	for (var i = 0; i < this.points.length; i++) {
		ctx.beginPath();
		ctx.strokeStyle = this.points[i].color;
		ctx.ellipse(this.points[i].x, this.points[i].y, 5, 5, 0, 0, Math.PI * 2);
		// ctx.closePath();
		ctx.stroke();
	}


	window.requestAnimationFrame(this.draw.bind(this));
}

function randomColor() {
	return '#'
		+ (Math.round(Math.random() * 255)).toString(16)
		+ (Math.round(Math.random() * 255)).toString(16)
		+ (Math.round(Math.random() * 255)).toString(16);
}