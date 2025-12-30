import { useEffect, useRef, useState } from 'react';
import { toast } from 'react-hot-toast';

export default function useWebSocket(playerName, initialGameId) {
  const [gameState, setGameState] = useState(null);
  const [assignedId, setAssignedId] = useState(initialGameId === 'new' ? "" : initialGameId);
  const socketRef = useRef(null);

  useEffect(() => {
    // Only connect if we have a name
    if (!playerName) return;

    const socket = new WebSocket(`ws://${process.env.REACT_APP_API_URL}/enterRoomWS`);
    socketRef.current = socket;

    socket.onopen = () => {
      // Send the initial JSON required by hub.go
      const initData = { playerName };
      if (initialGameId !== 'new') {
        initData.gameId = initialGameId;
      }
      socket.send(JSON.stringify(initData));
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      // Update local game ID if backend shares it
      if (data.gameId) setAssignedId(data.gameId);
      
      // Update game state for the board
      if (data.board) {
        setGameState(data);
      } else if (data.message) {
        toast(data.message);
      }
    };

    return () => socket.close();
  }, [playerName, initialGameId]);

  const sendMove = (col) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ column: col }));
    }
  };

  return { gameState, assignedId, sendMove };
}