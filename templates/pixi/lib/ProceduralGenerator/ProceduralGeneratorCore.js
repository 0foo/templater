export class ProceduralGeneratorCore {
    constructor(width, height, wallIndex = 0, pathIndex = 1) {
      this.width = width;
      this.height = height;
      this.wallTiles = Array.isArray(wallIndex) ? wallIndex : [wallIndex];
      this.pathTiles = Array.isArray(pathIndex) ? pathIndex : [pathIndex];
    }
  
    // Helper: pick a random item from an array
    randomTile(tiles) {
      return tiles[Math.floor(Math.random() * tiles.length)];
    }
  
    generateMaze() {
      const rows = this.height * 2 + 1;
      const cols = this.width * 2 + 1;
  
      // Fill the map with random wall tiles
      const maze = Array.from({ length: rows }, () =>
        Array.from({ length: cols }, () => this.randomTile(this.wallTiles))
      );
  
      const carve = (x, y) => {
        maze[y][x] = this.randomTile(this.pathTiles);
  
        const dirs = [
          [0, -2], [2, 0], [0, 2], [-2, 0]
        ].sort(() => Math.random() - 0.5);
  
        for (const [dx, dy] of dirs) {
          const nx = x + dx;
          const ny = y + dy;
  
          if (
            ny > 0 && ny < rows &&
            nx > 0 && nx < cols &&
            this.wallTiles.includes(maze[ny][nx])
          ) {
            maze[y + dy / 2][x + dx / 2] = this.randomTile(this.pathTiles);
            carve(nx, ny);
          }
        }
      };
  
      carve(1, 1);
      return maze;
    }
  }
  