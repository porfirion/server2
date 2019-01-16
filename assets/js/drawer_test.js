"use strict";

const OBJECTS_DISTANCE = 0;
const OBJECTS_COUNT = 1000;
const BOUND = 5000;

const objects = [];

function onLoad() {
    let spriteFactory = new SpriteFactory();

    window.elem = window.document.getElementById('map');
    let drawer = new Drawer(elem.getContext("2d"), elem.clientWidth, elem.clientHeight);
    let clickPos = null;
    let clickViewportPos = null;

    updateViewportSize();
    updateViewportInfo();

    function initObjects() {
        // очищаем объекты
        objects.forEach(function(obj) {
            drawer.removeObject(obj.obj.id);
        });
        objects.splice(0, objects.length);

        document.getElementById("info").innerHTML = "Objects count: " + OBJECTS_COUNT;

        for (let i = 0; i < OBJECTS_COUNT; i++) {
            let obj = drawer.createObject();
            obj.setSize(Math.random() * 5 + 10);
            obj.setPosition({
                x: Math.random() * OBJECTS_DISTANCE * 2 - OBJECTS_DISTANCE,
                y: Math.random() * OBJECTS_DISTANCE * 2 - OBJECTS_DISTANCE
            });
            // obj.addLayer("rect", new RectLayer(obj, randomColor()/*, randomColor()*/));
            // obj.addLayer("circle", new CircleLayer(obj, randomColor()/*, randomColor()*/));
            // obj.addLayer("id", new IdLayer(obj));
            obj.addLayer("image", new ImageLayer(obj, spriteFactory.getSprite("/assets/img/pig.png")));

            let x = Math.random() - 0.5;
            let y = Math.random() - 0.5;

            // square
            // let abs = Math.abs(x) + Math.abs(y);
            // let speed = {x: x / abs, y: y / abs};

            // circle
            // let abs = Math.sqrt(x * x + y * y);
            // let speed = {x: x / abs, y: y / abs};

            let speed = {x: x, y: y};

            if (i < 10) {
                speed = {x: 0, y: 0};
                obj.setSize(20);
                obj.addLayer("id", new IdLayer(obj));
                obj.setPosition({x: i * 45 - 205, y: 0});
            }

            if (i === 0) {
                obj.addLayer("particle", new ParticleLayer(obj, {
                    emit: () => {
                        return new Particle(randomColor(), 10);
                    }
                }));
            } else if (i === 1) {

            }

            objects.push({
                obj: obj,
                start: obj.getPosition(),
                speed: speed,
                startTime: Date.now(),
            });
        }
    }
    initObjects();

    // drawer.viewport.setScale(0.095); // for global view
    drawer.viewport.setScale(3.500);// for detailed view

    drawer.draw();

    function move() {
        let now = Date.now();
        for (let i = 10; i < OBJECTS_COUNT; i++) {
            let obj = objects[i];

            let delta = (now - obj.startTime);

            let current = {
                x: objects[i].start.x + objects[i].speed.x * delta,
                y: objects[i].start.y + objects[i].speed.y * delta,
            };

            if (Math.abs(current.x) > BOUND) {
                obj.start = current;
                obj.startTime = now;

                current.x = Math.min(Math.max(-BOUND, current.x), BOUND);
                obj.speed.x = -obj.speed.x;
            }
            if (Math.abs(current.y) > BOUND) {
                obj.start = current;
                obj.startTime = now;

                current.y = Math.min(Math.max(-BOUND, current.y), BOUND);
                obj.speed.y = -obj.speed.y;
            }
            objects[i].obj.setPosition(current);
            let x = objects[i].speed.x;
            let y = objects[i].speed.y;
            objects[i].obj.setRotation(Math.PI + Math.atan2(x, y));
        }

        drawer.draw();

        requestAnimationFrame(move);
    }

    requestAnimationFrame(move);

    function updateViewportSize() {
        elem.width = elem.clientWidth;
        elem.height = elem.clientHeight;
        drawer.setCanvasSize(elem.clientWidth, elem.clientHeight);
    }

    function updateViewportInfo() {
        document.getElementById('viewportSize').innerHTML = JSON.stringify(drawer.viewport.getRealDimensions(), numberPrecisionLimiter, " ");
    }

    window.addEventListener("resize", () => {
        updateViewportSize();
        updateViewportInfo();
        drawer.draw();
    });

    elem.addEventListener('mousewheel', function(ev) {
        let params = normalizeWheel(ev);
        if (params.spinY > 0) {
            // на себя
            drawer.viewport.scaleBy(0.96);
        } else {
            // от себя
            drawer.viewport.scaleBy(1.05);
        }

        requestAnimationFrame(drawer.draw.bind(drawer));
        updateViewportInfo();

        // capture all scrolling over map
        return false;
    });

    elem.addEventListener('contextmenu', function(event) {
        event.preventDefault();
        return false;
    });

    elem.addEventListener('mousedown', function(event) {
        let canvasCoords = new Point2D(event.offsetX, event.offsetY);
        let viewportCoords = drawer.viewport.fromCanvas(canvasCoords);
        let realCoords = drawer.viewport.toReal(viewportCoords);

        switch (event.button) {
            case 0:
                clickPos = canvasCoords;
                clickViewportPos = drawer.viewport.getRealDimensions();
                // левая кнопка мыши
                break;
            case 1:
                // средняя кнопка мыши
                break;
            case 2:
                // правая кнопка мыши
                drawer.viewport.setPos(realCoords);
                // updateViewportSize();
                drawer.draw();
                break;
            default:
                console.warn("unexpected button " + event.button);
                break;
        }
        return false;
    });

    elem.addEventListener('mouseup', function() {
        clickPos = null;
        clickViewportPos = null;
    });
    elem.addEventListener('mouseout', function() {
        clickPos = null;
        clickViewportPos = null;
    });

    elem.addEventListener('mousemove', function(event) {
        let canvasCoords = new Point2D(event.offsetX, event.offsetY),
            viewportCoords = drawer.viewport.fromCanvas(canvasCoords),
            realCoords = drawer.viewport.toReal(viewportCoords);

        if (clickPos != null) {
            let
                dx = (canvasCoords.x - clickPos.x) / drawer.viewport.getScale(),
                dy = (canvasCoords.y - clickPos.y) / drawer.viewport.getScale();

            drawer.viewport.setPos(new Point2D(
                clickViewportPos.x - dx,
                clickViewportPos.y + dy
            ));
            updateViewportInfo();
            drawer.draw();

            // мы передвинули вьюпорт, так что нужно по новой рассчитать положение указателя
            viewportCoords = drawer.viewport.fromCanvas(canvasCoords);
            realCoords = drawer.viewport.toReal(viewportCoords);
        }

        let realOptCoords = drawer.viewport.fromCanvasToReal(canvasCoords);

        document.getElementById('mouseCoords').innerHTML = ''
            + 'canvas  : ' + JSON.stringify(canvasCoords, numberPrecisionLimiter) + "\n"
            + 'viewport: ' + JSON.stringify(viewportCoords, numberPrecisionLimiter) + "\n"
            + 'real    : ' + JSON.stringify(realCoords, numberPrecisionLimiter) + "\n"
            + 'realOpt : ' + JSON.stringify(realOptCoords, numberPrecisionLimiter) + "\n"
            + 'dx: ' + (realCoords.x - realOptCoords.x).toFixed(8) + ' dy: ' + (realCoords.y - realOptCoords.y).toFixed(8)
        ;
    });

    document.getElementById("restartButton").addEventListener("click", function() {
        initObjects();
    })
}
