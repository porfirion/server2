class GameObject {
}

// describes visible game region and whole game state
class GameState {
    visibleObjects: GameObject[] = [];

    constructor() {}

    processMessage(msg: ServerMessage): void {
        // add/remove/update visible objects (work with drawer)
        // adjust whole game state (day/night, victory, ...)
    }
}