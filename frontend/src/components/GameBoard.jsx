import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useGame } from '../context/GameContext';
import { toast } from 'react-hot-toast';

export default function GameBoard() {

const apiBase = process.env.REACT_APP_API_URL;

// Check if we are running locally or on the web
const isLocal = apiBase.includes('localhost');

const httpUrl = isLocal ? `http://${apiBase}` : `https://${apiBase}`;
const wsUrl = isLocal ? `ws://${apiBase}` : `wss://${apiBase}`;

// Use these for your connections
console.log("Connecting to:", wsUrl + "/enterRoomWS");

  const { gameId } = useParams();
  const navigate = useNavigate();
  const { gameState, sendMove, resetGame } = useGame();
  const [turnTimer, setTurnTimer] = useState(45);
  const playerName = localStorage.getItem('playerName');
  const [reconnectTimer, setReconnectTimer] = useState(30);
  
  // Controls when the "Won/Lost" popup actually appears
  const [showModal, setShowModal] = useState(false);

  const handleBackToLobby = () => {
    resetGame();
    navigate('/');
  }

  // Turn Timer logic
  useEffect(() => {
    if (gameState?.status === "IN_PROGRESS") {
      setTurnTimer(45);
      const interval = setInterval(() => setTurnTimer(t => (t > 0 ? t - 1 : 0)), 1000);
      return () => clearInterval(interval);
    }
  }, [gameState?.current_turn, gameState?.status]);

  // Add this to GameBoard.jsx
// Move this below the line: const players = gameState.players || [];
// OR define players inside the effect as shown below:

useEffect(() => {
  // 1. Define players inside the effect so it's accessible
  const currentPlayers = gameState?.players || [];
  
  // 2. Look for the offline player
  const offlinePlayer = currentPlayers.find(p => !p.is_online && !p.is_bot);
  
  if (offlinePlayer && gameState?.status === "IN_PROGRESS") {
    const interval = setInterval(() => {
      setReconnectTimer(t => (t <= 1 ? 0 : t - 1));
    }, 1000);
    return () => clearInterval(interval);
  } else {
    setReconnectTimer(30); 
  }
}, [gameState?.players, gameState?.status]); // Add players to dependency array

  // Handle the delay for the Winning Modal
  useEffect(() => {
    if (gameState?.status === "FINISHED") {
      // 3-second delay so players can see the winning disc land
      const timer = setTimeout(() => setShowModal(true), 3000);
      return () => clearTimeout(timer);
    } else {
      setShowModal(false);
    }
  }, [gameState?.status]);

  if (!gameState) return (
    <div className="h-screen flex items-center justify-center bg-slate-950 text-white font-mono animate-pulse">
      Synchronizing Game State...
    </div>
  );

  const players = gameState.players || [];
  const localPlayerIndex = players.findIndex(p => p.username === playerName);
  const localPlayerNumber = localPlayerIndex + 1;

  const getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
};

// Inside your component:
const userIdFromCookie = getCookie('connectFourUserId');
console.log("Current User ID from Cookie:", userIdFromCookie);

  return (
    <div className="min-h-screen bg-slate-950 p-4 md:p-8 flex flex-col items-center">
      
      {/* WINNER TOP BANNER: Visible immediately when game ends, before modal pops up */}
      {gameState.status === "FINISHED" && (
        <div className="mb-6 py-2 px-6 bg-blue-600/20 border border-blue-500/50 rounded-full animate-fade-in">
          <span className="text-blue-400 font-black uppercase tracking-widest text-sm">
            {/* Change this line */}
            {gameState.winner === 0 ? "Game Tied!" : `${players[gameState.winner-1]?.username} has won!`}
          </span>
        </div>
      )}

      {/* RECONNECTION WARNING */}
{players.some(p => !p.is_online && !p.is_bot) && gameState.status === "IN_PROGRESS" && (
  <div className="fixed top-20 z-50 bg-red-600 text-white px-8 py-4 rounded-2xl shadow-2xl border-4 border-red-400 animate-bounce">
    <div className="flex flex-col items-center">
      <span className="font-black uppercase tracking-tighter">Opponent Disconnected!</span>
      <span className="text-sm font-bold">Winning by forfeit in: {reconnectTimer}s</span>
    </div>
  </div>
)}

      {/* PLAYER TABS */}
      <div className="grid grid-cols-2 gap-6 w-full max-w-4xl mb-10">
        {players.map((p, i) => {
          const isActive = gameState.current_turn - 1 === i && gameState.status === "IN_PROGRESS";
          console.log(isActive,gameState.current_turn-1,"jfrn",i)
          const discColorClass = i === 0 
            ? "from-red-600 to-red-400 border-red-800" 
            : "from-yellow-500 to-yellow-300 border-yellow-700";

          return (
            <div 
              key={i} 
              className={`p-5 rounded-3xl border-2 transition-all duration-500 shadow-lg ${
                isActive 
                ? 'bg-slate-800 border-blue-500 shadow-blue-500/20 scale-105 z-10' 
                : 'bg-slate-900/40 border-slate-800 opacity-40'
              }`}
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <div className={`w-3.5 h-3.5 rounded-full ${p.is_online || p.is_bot ? 'bg-green-500 shadow-[0_0_10px_#22c55e]' : 'bg-red-500 shadow-[0_0_10px_#ef4444]'}`} />
                  <div className={`w-5 h-5 rounded-full bg-gradient-to-tr border-b-2 ${discColorClass}`} />
                  <span className="text-xl font-black text-white truncate">
                    {p.username} {p.is_bot ? '🤖' : '👤'}
                  </span>
                </div>
                {isActive && (
                  <span className="text-blue-400 font-mono font-bold text-lg">{turnTimer}s</span>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* THE BOARD */}
      <div className="relative bg-blue-700 p-4 md:p-6 rounded-[2.5rem] shadow-2xl border-b-[12px] border-blue-800">
        <div className="flex gap-3 md:gap-4">
          {[...Array(7)].map((_, colIdx) => (
            <div 
              key={colIdx} 
              className="flex flex-col gap-3 md:gap-4 group cursor-pointer p-1 rounded-2xl hover:bg-blue-600/40 transition-colors"
              onClick={() => {
                if (gameState.status !== "IN_PROGRESS") return;
                if (gameState.current_turn - 1 === localPlayerIndex) {
                    console.log("Clicked Successfully")
                    sendMove(colIdx + 1);
                } else {
                    toast.error("It's not your turn!");
                }
              }}
            >
              {[...Array(6)].map((_, rowIdx) => {
                const cell = gameState.board[rowIdx][colIdx];
                return (
                  <div key={rowIdx} className="w-12 h-12 md:w-20 md:h-20 rounded-full bg-slate-950 shadow-[inset_0_4px_10px_rgba(0,0,0,0.8)] flex items-center justify-center overflow-hidden">
                    {cell === 1 && (
                      <div className="w-[88%] h-[88%] bg-gradient-to-tr from-red-700 to-red-500 rounded-full shadow-xl animate-drop border-b-4 border-red-900" />
                    )}
                    {cell === 2 && (
                      <div className="w-[88%] h-[88%] bg-gradient-to-tr from-yellow-500 to-yellow-300 rounded-full shadow-xl animate-drop border-b-4 border-yellow-700" />
                    )}
                  </div>
                );
              })}
            </div>
          ))}
        </div>
      </div>

      {/* GAME FINISHED MODAL */}
     {/* GAME FINISHED MODAL */}
{showModal && (() => {
  let resultText = gameState.winner === 0 ? "DRAW" : (gameState.winner === localPlayerNumber ? "WON" : "LOST");
  
  const handleViewLeaderboard = async () => {
    try {
      // 1. Call the Go API using the current player's name
      const response = await fetch(`${httpUrl}/leaderBoard?playerName=${encodeURIComponent(playerName)}`);
      if (!response.ok) throw new Error("Failed to fetch leaderboard");
      
      const data = await response.json();
      
      // 2. Navigate to the frontend route and pass the data via state
      navigate('/leaderboard', { state: { leaderboardData: data } });
    } catch (err) {
      console.error("Leaderboard Error:", err);
      toast.error("Could not load leaderboard data.");
    }
  };

  return (
    <div className="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center z-50 p-6 animate-fade-in">
      <div className="bg-slate-900 border border-slate-800 p-12 rounded-[3rem] text-center max-w-lg w-full shadow-2xl">
        <h2 className="text-6xl font-black mb-4 italic tracking-tighter uppercase text-white">
          {resultText}
        </h2>
        <p className="text-slate-400 mb-8 font-medium">
          {gameState.winner === 0 
            ? "No more valid moves remaining." 
            : `${players[gameState.winner - 1]?.username} connected four discs!`}
        </p>
        
        <div className="flex flex-col gap-4">
          {/* Main Action Button */}
          <button 
            onClick={handleBackToLobby} 
            className="w-full bg-blue-600 hover:bg-blue-500 text-white font-black py-5 rounded-2xl transition-all uppercase tracking-widest shadow-lg shadow-blue-600/20"
          >
            Back to Lobby
          </button>

          {/* Side-by-Side Secondary Buttons */}
          <div className="flex gap-4">
            <button 
              onClick={handleViewLeaderboard}
              className="flex-1 bg-slate-800 hover:bg-slate-700 text-white font-bold py-4 rounded-xl transition-all border border-slate-700 uppercase text-xs tracking-wider"
            >
              View Leaderboard
            </button>
            <button 
              onClick={() => setShowModal(false)}
              className="flex-1 bg-slate-800 hover:bg-slate-700 text-white font-bold py-4 rounded-xl transition-all border border-slate-700 uppercase text-xs tracking-wider"
            >
              Review Board
            </button>
          </div>
        </div>
      </div>
    </div>
  );
})()}
    </div>
  );
}