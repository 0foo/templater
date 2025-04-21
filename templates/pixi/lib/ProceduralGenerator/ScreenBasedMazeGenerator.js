import { ProceduralGeneratorCore } from './ProceduralGeneratorCore.js';

export class ScreenBasedMazeGenerator extends ProceduralGeneratorCore {
  constructor(screenWidth, screenHeight, tileSize = 64, wallIndex = 0, pathIndex = 1) {
    const width = Math.floor(screenWidth / tileSize / 2);
    const height = Math.floor(screenHeight / tileSize / 2);
    super(width, height, wallIndex, pathIndex);
  }
}
