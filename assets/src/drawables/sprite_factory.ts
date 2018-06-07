class SpriteFactory {
    private sprites: Map<String, Sprite> = new Map<String, Sprite>();

    constructor() {}

    public getSprite(name: string): Sprite {
        let existing = this.sprites.get(name);
        if (typeof existing != 'undefined') {
            return existing;
        } else {
            existing = new Sprite(name);
            this.sprites.set(name, existing);
            return existing;
        }
    }

}