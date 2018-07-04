jQuery(document).ready(function () {
    myName = randomName();
    $('.playerName').html(myName);
    client = new WsClient("ws://" + window.location.host + "/ws", myName);
    client.on('message', onmessage);
    client.on('close', onclose);
    client.on('open', function () {
        syncTimeTimer = setInterval(function () {
            //client.sendMessage(MessageType.SYNC_TIME, {time: 0});
            //console.log('sent');
        }, 10000);
    });
    client.on(WsClient.TimeSynced, function (latency, timeCorrection) {
        $('.latency .value').html(client.latencies[client.latencies.length - 1]);
        // Коррекцию выбираем как среднее из последних полученных
        var currentCorrection = client.timeCorrections.reduce(function (sum, a) {
            return sum + a
        }, 0) / (client.timeCorrections.length || 1);
        $('.timeCorrection .value').html(currentCorrection.toFixed(3));
        map.latency = latency;
        map.timeCorrection = currentCorrection;
    });
    client.connect();

    $('#chat_form').submit(function (event) {
        event.preventDefault();

        var inp = $('.chat_input');
        var text = inp.val();
        try {
            client.sendMessage(MessageType.TEXT, {Text: text});
        } catch (err) {
            showMessage("Unable to send " + text, "text-danger");
            console.error(err);
        }

        inp.val('');

        return false;
    });

    var elem = document.getElementById("map");
    // var wrapper = document.getElementById('map-wrapper');
    // console.log(wrapper);
    //elem.width = wrapper.clientWidth;
    //elem.height = wrapper.clientHeight;

    map = new Map(elem);

    // запуск анимации, если она ещё не была начата
    $(document.body).on('click', '.drawButton', function () {
        map.forceDraw();
        return false;
    });

    // запуск анимации, если она ещё не была начата
    $(document.body).on('click', '.autoDrawButton', function () {
        map.toggleAutoDrawing();
        return false;
    });

    // центрирование вьюпорта (0:0)
    $(document.body).on('click', '.centerButton', function () {
        map.viewport.x = 0;
        map.viewport.y = 0;
        map.viewportAdjustPoint = null;
        map.forceDraw();
        return false;
    });

    $(document.body).on('click', '.simulateButton', function () {
        client.sendMessage(MessageType.SIMULATE_MESSAGE, {steps: 1});
    });

    // перемещение вьюпорта при помощи кнопок навигации
    $(document).on('click', '.floatingButton', function () {
        var x = $(this).data('x');
        var y = $(this).data('y');

        map.viewport.x += Number(x) * map.viewport.scale * 20;
        map.viewport.y += Number(y) * map.viewport.scale * 20;
        map.forceDraw();
    });

    $(document).on('click', '.zoomIn', function() {
        map.viewport.scale *= 1.1;
        map.forceDraw();
        return false;
    });
    $(document).on('click', '.zoomOut', function() {
        map.viewport.scale /= 1.1;
        map.forceDraw();
        return false;
    });

    $(map).on('game:click', function (event, data) {
        console.log('clicked at ', data);

        client.sendMessage(MessageType.ACTION_MESSAGE, {
            action_type: 'move',
            action_data: data
        });
    });

    $(document).on('click', '.chatToggler', function () {
        $('#chat').toggle();
    });

    // setInterval(function() {
    // 	console.log('send last pos ', map.lastCursorPosition, map.lastCursorPositionReal);
    // 	if (map.lastCursorPositionReal != null) {
    // 		client.sendMessage(MessageType.ACTION_MESSAGE, {
    // 			action_type: 'accelerate',
    // 			action_data: map.lastCursorPositionReal,
    // 		});
    // 	}
    // }, 1000);

    $(document).on('keydown', function (e) {
        switch (e.keyCode) {
            case 32:
                // space
                client.sendMessage(MessageType.SIMULATE_MESSAGE, {steps: 1});
                break;
            default:
                console.log('Key pressed', e.keyCode);
                break;
        }
    });

    $(document).on('click', '.reloadButton', function(event) {
        window.location.reload(true);

        event.stopPropagation();
        event.preventDefault();
        return false;
    });

    $(document).on("click", ".changeSimulationMode", function() {
        client.sendMessage(MessageType.CHANGE_SIMULATION_MODE, {step_by_step: !map.simulationMode});
    });

    // map.toggleAutoDrawing();
});