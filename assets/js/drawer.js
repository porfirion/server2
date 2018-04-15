const
    DRAW_MODE_ONLY_SERVER = 1,
    DRAW_MODE_ONLY_REAL = 2,
    DRAW_MODE_BOTH = 3;

/**
 * Viewport - is a window from sceen to real world;
 * @param x
 * @param y
 * @param scale
 * @constructor
 */
function Viewport(x, y, scale, width, height) {
    this.realX = typeof x !== 'undefined' ? Number(x) : 0;
    this.realY = typeof y !== 'undefined' ? Number(y) : 0;
    this.scale = typeof scale !== 'undefined' ? Number(scale) : 1.0;

    this.realWidth = 1000;
    this.realHeight = 1000;

    if (typeof width !== 'undefined' && typeof height !== 'undefined') {
        this.updateSize(width, height);
    }
}

Viewport.prototype.updateSize = function(width, height) {
    this.canvasWidth = width;
    this.canvasHeight = height;

    // размер вьюпорта в реальных пикселях
    this.realWidth = width / this.scale;
    this.realHeight = height / this.scale;
    
    this.realWidth_half = this.realWidth / 2;
    this.realHeight_half = this.realHeight / 2;
};

/**
 * Returns real world coords and size of viewport
 * @returns {{x: number, y: number, width: number, height: number, left: number, top: number, right: number, bottom: number}}
 */
Viewport.prototype.getRealDimensions = function() {
    return {
        // реальное положение и размер вьюпорта
        x: this.realX,
        y: this.realY,
        width: this.realWidth,
        height: this.realHeight,

        // для удобства отдаём ещё и реальные границы
        left: this.realX - this.realWidth_half,
        top: this.realY - this.realHeight_half,
        right: this.realX + this.realWidth_half,
        bottom: this.realY + this.realHeight_half
    };
};

/**
 * Рассчитывает X на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param xr реальная координата по X
 * @returns {number} координата X на канве
 */
Viewport.prototype.xToCanvas = function(xr) {
    return this.realWidth_half + (xr - this.x);
};

/**
 * Рассчитывает Y на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param yr реальная координата Y
 * @returns {number} координата Y на канве
 */
Viewport.prototype.yToCanvas = function(yr) {
    return this.realHeight_half - (yr - this.y);
};

/**
 * Рассчитывает положение объекта на канве без учёта скейла,
 * так как скейл применён на самом вьюпорте
 * @param realPos реальное положение объекта
 * @returns {{x: number, y: number}} координаты объекта на канве
 */
Viewport.prototype.fromReal = function(realPos) {
    return {
        x: this.realWidth_half + (realPos.x - this.x),
        y: this.realHeight_half - (realPos.y - this.y)
    };
};

/**
 * Translate viewport coords into real world coords
 * @param {Point} viewportPos положение курсора на канве
 * @returns {Point} реальное положение курсора
 */
Viewport.prototype.toReal = function(viewportPos) {
    return {
        x: this.x + (viewportPos.x / this.scale - this.realWidth_half),
        y: this.y - (viewportPos.y / this.scale - this.realHeight_half)
    };
};

Viewport.prototype.fromCanvas = function(canvasPos) {
    return {
        x: canvasPos.x - this.realWidth_half,
        y: -canvasPos.y + this.realHeight_half
    };
};


/**
 * Class for holding and drawing list of map objects (including players)
 * Can be a wrapper to some framework
 * @param {HTMLCanvasElement} elem
 * @constructor
 */
function Drawer(elem) {
    this.elem = elem;
    /**
     * Canvas for drawing
     * @type {CanvasRenderingContext2D}
     */
    this.ctx = this.elem.getContext("2d");

    /**
     * Position (real) and scale of viewport
     * @type {{x: number, y: number, scale: number}}
     */
    this.viewport = new Viewport();

    // запомненное состояние вращающегося курсора
    this.prevOffset = 0;
    /**
     * List of all objects in the entire map
     * @type {Array}
     */
    this.objects = [];
    this.objectsById = {};
    this.gridSize = 100;

    this.timeCanvas = document.createElement("canvas");
    this.timeCanvas.width = 300;
    this.timeCanvas.height = 100;

    // this.drawMode = DRAW_MODE_ONLY_SERVER;
    this.drawMode = DRAW_MODE_BOTH;
}

/**
 * Wrapper for automatic drawing
 * @private
 */
Drawer.prototype.draw = function() {
    if (this.isAnimating) {
        this.draw_();

        requestAnimationFrame(this.draw.bind(this));
    }
};
Drawer.prototype.forceDraw = function() {
    this.draw_();
};

/**
 * Exokicit drawing int canvas
 * @private
 */
Drawer.prototype.draw_ = function() {
    var ctx = this.ctx;
    var elem = this.elem;

    elem.width = elem.clientWidth;
    elem.height = elem.clientHeight - 1;

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
};

Map.prototype.drawTime = function() {
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

    var fillTextRight = function(text, right, top) {
        ctx.fillText(text, right - ctx.measureText(text).width, top);
    };

    fillTextRight('FPS: ' + Math.round(1000 / average), 290, 15);
    fillTextRight('min: ' + min, 290, 30);
    fillTextRight('average: ' + Math.round(average), 290, 45);
    fillTextRight('max: ' + max, 290, 60);

    ctx.fillText('viewport: (x: ' + Math.round(this.viewport.x * 100) / 100 + '; y: ' + Math.round(this.viewport.y * 100) / 100 + ')', 15, 15);
    ctx.fillText('scale: ' + Math.round(this.viewport.scale * 100) / 100, 15, 30);
    ctx.fillText('latency: ' + this.latency.toFixed(0) + ' ms', 15, 45);
    ctx.fillText('time correction: ' + this.timeCorrection.toFixed(1) + ' ms', 15, 60);

    this.ctx.save();
    this.ctx.globalAlpha = 0.7;
    this.ctx.drawImage(this.timeCanvas, 0, 0, 300, 100);
    this.ctx.restore();
};

Drawer.prototype.drawObjects = function() {
    var ctx = this.ctx;
    var viewportReal = this.getRealViewport(); // real position and size of viewport

    var serverTime = this.getCurrentSimulationTime();

    ctx.save();
    ctx.scale(this.viewport.scale, this.viewport.scale);
    ctx.lineWidth = 1;

    for (var i = 0; i < this.objects.length; i++) {
        var obj = this.objects[i];
        var objSizeReal = obj.size; // half of object size

        if (this.drawMode === DRAW_MODE_ONLY_SERVER || this.drawMode === DRAW_MODE_BOTH) {
            var objPosServer = obj.getLastServerPosition();
            if (this.rectContainsPoint(viewportReal, objPosServer, objSizeReal)) {
                // рисуем текущее положение объекта по серверу
                ctx.save();
                var serverPos = this.realToViewport(objPosServer);
                ctx.translate(serverPos.x, serverPos.y);
                ctx.lineWidth = 1;
                ctx.setLineDash([4, 2]);
                ctx.beginPath();
                ctx.arc(0, 0, 10, 0, Math.PI * 2);
                ctx.strokeStyle = '#aaaaaa';
                ctx.closePath();
                ctx.stroke();
                ctx.restore();
            }
        }
        if (this.drawMode === DRAW_MODE_ONLY_REAL || this.drawMode === DRAW_MODE_BOTH) {
            var objPosReal = obj.getApproximatedPosition(serverTime);
            if (this.rectContainsPoint(viewportReal, objPosReal, objSizeReal) || true) {
                ctx.save();
                var objPosViewport = this.realToViewport(objPosReal); // viewport position of object
                ctx.translate(objPosViewport.x, objPosViewport.y);
                obj.draw(ctx);
                ctx.restore();
            }
        }

        // var isVisible = viewportReal.left <= objPosReal.x + objHalfSizeReal && objPosReal.x - objHalfSizeReal <= viewportReal.right &&
        // 	viewportReal.top <= (objPosReal.y + objHalfSizeReal) && (objPosReal.y - objHalfSizeReal) <= viewportReal.bottom;
    }
    ctx.restore();
};

Drawer.prototype.drawGrid = function() {
    var viewportW = this.elem.width / this.viewport.scale;
    var viewportH = this.elem.height / this.viewport.scale;

    var realViewport = this.viewport.getRealDimensions();

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
        var x = this.xToCanvas(leftCol + i * this.gridSize);
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
        var y = this.yToCanvas(topRow + j * this.gridSize);
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

    if (this.isAnimating) {
        // рисуем вращающийся курсор только для непрерывной анимации
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
    }

    // рисуем центр
    ctx.strokeStyle = 'lime';
    ctx.beginPath();
    // ctx.ellipse(realWidth / 2, realHeight / 2, 10, 10, 0, 0, Math.PI * 2);
    ctx.moveTo(viewportW / 2 - 15, viewportH / 2);
    ctx.lineTo(viewportW / 2 + 15, viewportH / 2);
    ctx.moveTo(viewportW / 2, viewportH / 2 - 15);
    ctx.lineTo(viewportW / 2, viewportH / 2 + 15);
    // ctx.arc(realWidth / 2, realHeight / 2, 10, 0, Math.PI * 2);
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
Drawer.prototype.drawAnchors = function() {
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

    ctx.restore();
};