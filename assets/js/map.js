"use strict";

const
    SIMULATION_MODE_STEP_BY_STEP = true,
    SIMULATION_MODE_CONTINIOUS = false;

var SimulationMode = {
    STEP_BY_STEP: true,
    CONTINUOUS: false
};

const ObjectType = {
    Obstacle: 1,
    NPC: 10,
    Player: 100,
};

/**
 * Class for holding and NOT drawing list of map objects (including players)
 * @param {HTMLCanvasElement} elem
 * @constructor
 */
function Map(elem) {
    /**
     * List of map objects
     * @type {DrawableObject[]}
     */

    this.objectsById = {};
    this.myObject = null; // объект текущего игрока
    this.lastCursorPositionReal = null;
    this.simulationMode = SimulationMode.CONTINUOUS; // @type {boolean}
    this.simulationTime = 0; // игровое время
    this.stateSyncTime = 0; // время, в которое мы получили состояние серера
   /**
     * @type {Player[]}
     */
    this.players = {};
    this.latency = 0; // пинг до сервера
    this.timeCorrection = 0; // сколько нужно прибавить к текущему времени, чтобы получить серверное
    /**
     * Flag, is auto drawing enabled
     * @type {boolean}
     */
    this.isAnimating = false;
    /**
     * List of animations time (for diagram)
     * @type {Array}
     */
    this.animations = [];
    this.prevAnimationTime = null;

    this.initHandlers(this.elem);

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
}

Map.prototype.mainLoop = function() {
    if (this.isAnimating) {
        requestAnimationFrame(this.mainLoop().bind(this));
    }
};

Map.prototype.initHandlers = function(elem) {
    $(elem).on('mousedown', function(event) {
        var canvasCoords = new Point2D(event.offsetX, event.offsetY);
        var viewportCoords = this.viewport.fromCanvas(canvasCoords);
        var realCoords = this.viewport.toReal(viewportCoords);
        switch (event.button) {
            case 0:
                // левая кнопка мыши
                $(this).trigger('game:click', realCoords);
                break;
            case 1:
                // средняя кнопка мыши
                break;
            case 2:
                // правая кнопка мыши
                if (this.isAnimating) {
                    this.viewportAdjustPoint = realCoords;
                    this.viewportAdjustPointVP = viewportCoords;
                } else {
                    this.viewport.x = realCoords.x;
                    this.viewport.y = realCoords.y;
                    this.draw_();
                }
                break;
            default:
                console.warn("unexpected button " + event.button);
                break;
        }
        return false;
    }.bind(this));

    $(elem).on('mousemove', function(event) {
        var canvasCoords = new Point2D(event.offsetX, event.offsetY);
        var viewportCoords = this.viewport.fromCanvas(canvasCoords);
        var realCoords = this.viewport.toReal(viewportCoords);
        this.viewportAdjustPointVP = canvasCoords;
        this.lastCursorPositionReal = realCoords;
    }.bind(this));

    $(elem).on('mouseenter', function() {}.bind(this));

    $(elem).on('mouseout', function() {
        this.lastCursorPositionReal = null;
    }.bind(this));

    $(elem).on('mousewheel DOMMouseScroll', function(ev) {
        var params = normalizeWheel(ev.originalEvent);
        if (params.spinY > 0) {
            // на себя
            this.viewport.scaleBy(0.97);
        } else {
            // от себя
            this.viewport.scaleBy(1.05);
        }

        this.lastCursorPositionReal = this.viewport.toReal(new Point2D(ev.offsetX, ev.offsetY));

        // capture all scrolling over map
        return false;
    }.bind(this));
    $(elem).on('contextmenu', function (event) {
        event.preventDefault();
        return false;
    }.bind(this));
};

// возвращает текущее время, оттуалкиваясь от которого мы рисуем объекты
Map.prototype.getCurrentSimulationTime = function () {
    if (this.simulationMode === SIMULATION_MODE_STEP_BY_STEP) {
        // отдаём то, что нам сказали в последний раз, так как итерируем по шагам
        return this.simulationTime;
    } else {
        var timeSinceSync = Date.now() - this.stateSyncTime;
        // игровое время + плюс сколько прошло с момента синхронизации + коррекция
        return this.simulationTime + timeSinceSync + this.timeCorrection;
    }
};

Map.prototype.addObject = function (id, objectType, coords) {
    if (objectType === undefined) {
        objectType = ObjectType.Obstacle;
    }

    if (coords === undefined || coords === null) {
        coords = {x: 0, y: 0};
    }
    var obj = new DrawableObject(id, objectType, coords);

    this.objects.push(obj);
    this.objectsById[obj.id] = obj;

    return obj;
};
Map.prototype.removeObject = function (obj) {
    var index = this.objects.indexOf(obj);

    this.objects.splice(index, 1);
    delete this.objectsById[obj.id];
};

Map.prototype.removeAllObjects = function () {
    this.objectsById = {};
    this.objects = [];
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
        if (this.objects[i].hasOwnProperty('player') && this.objects[i].player != null && this.objects[i].player.id === playerId) {
            this.objects.splice(i, 1);
            break;
        }
    }
};

Map.prototype.removeAllPlayers = function () {
    this.players = {};
};

Map.prototype.clear = function () {
    this.removeAllObjects();
    this.removeAllPlayers();
};

Map.prototype.toggleAutoDrawing = function () {
    this.isAnimating = !this.isAnimating;

    if (this.isAnimating) {
        this.draw();
    }
};

Map.prototype.adjustViewport = function () {
    // disable viewport adjustment
    // return;

    // disable adjusting to out object
    // if (this.myObject != null) {
    //     // следование за "своим" объектом
    //     var realWidth = this.elem.width * 0.3 / this.viewport.scale;
    //     var realHeight = this.elem.height * 0.3 / this.viewport.scale;
    //
    //     var serverTime = Date.now() + this.timeCorrection;
    //     var objPosReal = this.myObject.getApproximatedPosition(serverTime);
    //
    //     if (objPosReal.x < this.viewport.x - realWidth / 2) { this.viewport.x = objPosReal.x + realWidth / 2; }
    //     if (objPosReal.x > this.viewport.x + realWidth / 2) { this.viewport.x = objPosReal.x - realWidth / 2; }
    //     if (objPosReal.y < this.viewport.y - realHeight / 2) { this.viewport.y = objPosReal.y + realHeight / 2; }
    //     if (objPosReal.y > this.viewport.y + realHeight / 2) { this.viewport.y = objPosReal.y - realHeight / 2; }
    //
    //     // пока что жёстко отключим перемещение вьюпорта.
    //     // в дальнейшем надо бы его переделать.
    //     // иногда случается, что пользователь начал движение в другую сторону и вьюпорт дёргается немного при развороте
    //     // сильнее всего чувствуется на мобильнике
    //     return;
    // }

    if (this.viewportAdjustPoint == null) {
        return;
    }

    var dx = (this.viewportAdjustPoint.x - this.viewport.x) * 0.1;
    var dy = (this.viewportAdjustPoint.y - this.viewport.y) * 0.1;

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

    this.lastCursorPositionReal = this.viewportToReal(this.viewportAdjustPointVP);
};

/**
 *
 * @param {Point2D} a
 * @param {Point2D} b
 * @returns {number}
 */
Map.prototype.distance = function (a, b) {
    return Math.sqrt(Math.pow(a.x - b.x, 2) + Math.pow(a.y - b.y, 2));
};

/**
 *
 * @param {bool} mode
 */
Map.prototype.setSimulationMode = function (mode) {
    if (this.simulationMode !== mode) {
        this.simulationMode = mode;
        var modeName = (this.simulationMode === SimulationMode.STEP_BY_STEP ? "StepByStep" : "Continious");
        $('.simulationMode .value').html(modeName);
    } else {

    }
};

Map.prototype.updateServerState = function (state) {
    this.setSimulationMode(state.simulation_by_step);
    this.simulationTime = state.simulation_time;
    this.stateSyncTime = Date.now();
    console.log(state);
};