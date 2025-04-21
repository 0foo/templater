export class MySpriteSheet {
  constructor(sheet) {
    this.sheet = sheet;

    if (!sheet.framesArray) {
      this.indexSpritesheetTextures(this.sheet);
    }
  }

  indexSpritesheetTextures(sheet) {
    const entries = Object.entries(sheet.textures)
      .sort(([a], [b]) => a.localeCompare(b));

    const framesArray = [];

    entries.forEach(([_, texture], i) => {
      texture.index = i;
      framesArray[i] = texture;
    });

    sheet.framesArray = framesArray;
  }

  drawTile(x, y, tileIndex, app) {
    const texture = this.sheet.framesArray?.[tileIndex];

    if (!texture) {
      console.warn(`Tile index ${tileIndex} not found.`);
      return;
    }

    const sprite = new PIXI.Sprite(texture);
    sprite.x = x;
    sprite.y = y;
    app.stage.addChild(sprite);
    return sprite;
  }
  getTextureIndexes() {
    return this.sheet.framesArray.map((_, i) => i);
  }
  
  drawGrid(tileMap, app, tileSize = 64) {
    for (let row = 0; row < tileMap.length; row++) {
      for (let col = 0; col < tileMap[row].length; col++) {
        const index = tileMap[row][col];
        const x = col * tileSize;
        const y = row * tileSize;
        this.drawTile(x, y, index, app);
      }
    }
  }
  
  static fromPixiAlias(alias) {
    const sheet = PIXI.Assets.get(alias);

    if (!sheet) {
      console.warn(`No spritesheet found for alias "${alias}".`);
      return null;
    }

    return new MySpriteSheet(sheet);
  }
}
