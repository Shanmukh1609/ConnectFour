import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useGame } from '../context/GameContext';

export default function WaitingRoom() {
  const { gameId } = useParams();
  const navigate = useNavigate();
  const { connect, gameState, assignedId } = useGame();
  const playerName = localStorage.getItem('playerName');

  useEffect(() => {
    // 1. Establish connection immediately
    if (playerName) {
      connect(playerName, gameId);
    }
  }, [playerName, gameId, connect]);

  useEffect(() => {
  if (gameState) {
    console.log("Backend Status Received:", gameState.status); // Should log "IN_PROGRESS"

    // Match the backend string constant exactly
    if (gameState.status === "IN_PROGRESS") {
      console.log("Match started! Redirecting...");
      
      // Ensure we have a valid ID for the URL
      const targetId = assignedId || gameId;
      if (targetId && targetId !== 'new') {
        navigate(`/game/${targetId}`);
      }
    }
  }
}, [gameState, assignedId, gameId, navigate]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-slate-950 p-6">
      <div className="bg-slate-900 border border-slate-800 p-10 rounded-3xl shadow-2xl text-center max-w-md w-full">
        <h2 className="text-3xl font-bold text-white mb-2">Matchmaking</h2>
        <p className="text-slate-400 mb-8">Waiting for an opponent to join...</p>

        {/* Display logic for Game ID */}
        <div className="bg-slate-950 border border-blue-500/30 px-6 py-4 rounded-2xl mb-10">
          <p className="text-xs uppercase tracking-widest text-slate-500 mb-2 font-bold">Share your Game ID</p>
          <p className="text-2xl font-mono text-blue-400 font-bold tracking-wider">
            {gameId !== 'new' ? gameId : (assignedId || "FETCHING...")}
          </p>
        </div>

        {/* Loading Animation */}
        <div className="flex flex-col items-center gap-4">
          <div className="relative">
            <div className="w-16 h-16 border-4 border-blue-500/20 rounded-full"></div>
            <div className="absolute top-0 w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
          </div>
          <p className="text-blue-500 font-medium animate-pulse">Searching for players...</p>
        </div>
      </div>
    </div>
  );
}