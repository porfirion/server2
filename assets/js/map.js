function Map(elem) {
	this.elem = elem;
	this.ctx = this.elem.getContext("2d");

	this.viewport = {
		x: 0,
		y: 0,
		scale: 1,
	};
}

Map.prototype.draw = function() {
	var rectSize = 50;
	var c = this.ctx;

	c.clearRect(0,0,c.width, c.height);

	c.beginPath();

	
	c.moveTo(0,0);
	c.lineTo(300,300);
	c.stroke();

}