import React, { createContext, useContext, useRef, useState, useCallback } from 'react';
import { toast } from 'react-hot-toast';

const GameContext = createContext();

export const GameProvider = ({ children }) => {

  const [gameState, setGameState] = useState(null);
  const [assignedId, setAssignedId] = useState("");
  const socketRef = useRef(null);

  const apiBase = process.env.REACT_APP_API_URL;

// Check if we are running locally or on the web
const isLocal = apiBase.includes('localhost');

const httpUrl = isLocal ? `http://${apiBase}` : `https://${apiBase}`;
const wsUrl = isLocal ? `ws://${apiBase}` : `wss://${apiBase}`;



  const connect = useCallback(async (playerName, initialGameId) => {
    // 1. Prevent duplicate connections
    if (socketRef.current && (socketRef.current.readyState === WebSocket.OPEN || 
        socketRef.current.readyState === WebSocket.CONNECTING)) {
      return;
    }

    try {
      // 2. Call SetCookie with playerName as a query parameter
      const resp = await fetch(`${httpUrl}/setCookie?playerName=${encodeURIComponent(playerName)}`, {
        method: 'GET',
        credentials: 'include' 
      });
      
      if (!resp.ok) throw new Error("Failed to initialize session");
      
      const authData = await resp.json();
      console.log("Session Initialized. UserID:", authData.userId);

      // 3. Initiate WebSocket
      const socket = new WebSocket(`${wsUrl}/enterRoomWS`);
      socketRef.current = socket;

      socket.onopen = () => {
        console.log("WebSocket Connected for:", playerName);
        const gid = initialGameId === 'new' ? null : initialGameId;
        
        // CRITICAL: You MUST send the userId back so the backend can link the session
        socket.send(JSON.stringify({ 
          playerName, 
          gameId: gid,
          userId: authData.userId // Send this so JoinOrReconnect works
        }));
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log("Socket Data Received:", data);

          if (data.gameId) setAssignedId(data.gameId);
          
          if (data.board) {
            setGameState(data);
          } else if (data.message) {
            toast(data.message);
          }
        } catch (err) {
          console.error("Error parsing socket data:", err);
        }
      };

      socket.onclose = (e) => {
        console.log("WebSocket Connection Closed", e.reason);
        socketRef.current = null; // Clear ref on close
      };

      socket.onerror = (err) => {
        console.error("WebSocket Error:", err);
        toast.error("Connection error occurred.");
      };

    } catch (err) {
      console.error("Failed to initialize session:", err);
      toast.error("Session initialization failed.");
    }
  }, []); // Wrapped in useCallback to prevent unnecessary re-renders

  const sendMove = (col) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ column: col }));
    } else {
      toast.error("Not connected to game server.");
    }
  };

  const resetGame = () => {
    setGameState(null);
    setAssignedId("");
    if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
    }
  };

  return (
    <GameContext.Provider value={{ gameState, assignedId, connect, sendMove, resetGame }}>
      {children}
    </GameContext.Provider>
  );
};

export const useGame = () => useContext(GameContext);