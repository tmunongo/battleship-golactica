import { Target } from "lucide-react";
import React, { useEffect, useState } from "react";
import "./App.css";
import { Button } from "./components/ui/button";
import { Card } from "./components/ui/card";

function App() {
  const GRID_SIZE = 10;
  const SHIP_SIZES = [5, 4, 3, 2, 1];

  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [gameState, setGameState] = useState({
    playerBoard: Array(GRID_SIZE)
      .fill()
      .map(() => Array(GRID_SIZE).fill({ hasShip: false, hit: false })),
    opponentBoard: Array(GRID_SIZE)
      .fill()
      .map(() => Array(GRID_SIZE).fill({ hit: false })),
    currentTurn: 0,
    playerIndex: null,
    status: "connecting",
    score: { hits: 0, misses: 0 },
    gameOver: false,
  });

  const [playerBoard, setPlayerBoard] = useState(
    Array(GRID_SIZE)
      .fill()
      .map(() => Array(GRID_SIZE).fill("empty"))
  );
  const [selectedCell, setSelectedCell] = useState<{
    row: number;
    col: number;
  } | null>(null);

  const getCellColor = (status) => {
    switch (status) {
      case "empty":
        return "bg-blue-100";
      case "ship":
        return "bg-gray-500";
      case "hit":
        return "bg-red-400";
      case "miss":
        return "bg-blue-300";
      default:
        return "bg-blue-100";
    }
  };

  const handleCellClick = (rowIndex, colIndex) => {
    setSelectedCell({ row: rowIndex, col: colIndex });
  };

  const handleFire = () => {
    // do nothing if no selected cell
    if (!selectedCell) {
      return;
    }

    const newBoard = [...playerBoard];
    const { row, col } = selectedCell;

    // server handles hit or miss logic

    setPlayerBoard(newBoard);
    setSelectedCell(null);
  };

  // create websocket connection
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      setSocket(ws);
      setGameState((prevState) => ({
        ...prevState,
        status: "connected",
      }));
    };
  }, []);

  return (
    <div className="p-4 w-full max-w-4xl mx-auto">
      <Card className="p-6">
        <h2 className="text-2xl font-bold mb-4">Battleship Golactica</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Player board */}
          <div>
            <h3 className="text-xl mb-2">Your board</h3>
            <div className="grid grid-cols-10 gap-1">
              {playerBoard.map((row, rowIndex) =>
                row.map((cell, colIndex) => (
                  <div
                    key={`${rowIndex}-${colIndex}`}
                    className={`
                          w-8 h-8 border cursor-pointer
                          ${getCellColor(cell)}
                          ${
                            selectedCell?.row === rowIndex &&
                            selectedCell?.col === colIndex
                              ? "ring-2 ring-yellow-400"
                              : ""
                          }
                        `}
                    onClick={() => handleCellClick(rowIndex, colIndex)}
                  >
                    {cell === "hit" && (
                      <Target className="w-6 h-6 text-white m-1" />
                    )}
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Gamepad */}
          <div className="flex flex-col gap-4">
            <div className="p-4 bg-gray-100 rounded">
              <h3 className="font-bold mb-2">Selected target</h3>
              {selectedCell ? (
                <p>
                  Row: {selectedCell.row + 1} x Col: {selectedCell.col + 1}
                </p>
              ) : (
                <p>Select a target cell</p>
              )}
            </div>

            <Button
              onClick={handleFire}
              disabled={!selectedCell}
              className="w-full"
            >
              Fire torpedo
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}

export default App;
