import { ProceduralGeneratorCore } from './ProceduralGeneratorCore.js';

export class TileBasedMazeGenerator extends ProceduralGeneratorCore {
  constructor(width, height, wallIndex = 0, pathIndex = 1) {
    super(width, height, wallIndex, pathIndex);
  }
}
