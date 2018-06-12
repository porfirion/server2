"use strict";
var GameObject = /** @class */ (function () {
    function GameObject() {
    }
    return GameObject;
}());
// describes visible game region and whole game state
var GameState = /** @class */ (function () {
    function GameState() {
        this.visibleObjects = [];
    }
    GameState.prototype.processMessage = function (msg) {
        // add/remove/update visible objects (work with drawer)
        // adjust whole game state (day/night, victory, ...)
    };
    return GameState;
}());
//# sourceMappingURL=game_state.js.map