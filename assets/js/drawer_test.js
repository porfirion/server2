"use strict";

function onLoad() {
    let canvas = window.document.getElementById('map');
    let drawer = new Drawer(canvas);
    drawer.draw();

    window.addEventListener("resize", () => {
        drawer.draw();
    })
}
