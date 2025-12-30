import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';

export default function LandingPage() {
  const [showP1Modal, setShowP1Modal] = useState(false);
  const [showP2Modal, setShowP2Modal] = useState(false);
  const [username, setUsername] = useState('');
  const [gameIdInput, setGameIdInput] = useState('');
  const navigate = useNavigate();

  const handleStartGame = (mode) => {
    if (!username) return toast.error("Username is required!");
    
    localStorage.setItem('playerName', username);
    
    if (mode === 'create') {
      // Player 1: Start fresh matchmaking/waiting
      navigate(`/waiting/new`);
    } else {
      // Player 2: Join specific game
      if (!gameIdInput) return toast.error("Enter a Game ID!");
      navigate(`/waiting/${gameIdInput}`);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-slate-900 via-blue-900 to-slate-900">
      <h1 className="text-6xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-red-500 mb-12 drop-shadow-lg">
        CONNECT FOUR
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 w-full max-w-2xl px-4">
        {/* Create Game Card */}
        <button 
          onClick={() => setShowP1Modal(true)}
          className="group p-8 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl hover:bg-blue-600/20 transition-all duration-300 transform hover:-translate-y-2"
        >
          <div className="text-4xl mb-4">🎮</div>
          <h2 className="text-2xl font-bold text-blue-400">Create Game</h2>
          <p className="text-slate-400 text-sm mt-2">Start a new room and wait for opponents</p>
        </button>

        {/* Join Game Card */}
        <button 
          onClick={() => setShowP2Modal(true)}
          className="group p-8 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl hover:bg-red-600/20 transition-all duration-300 transform hover:-translate-y-2"
        >
          <div className="text-4xl mb-4">🔑</div>
          <h2 className="text-2xl font-bold text-red-400">Join by ID</h2>
          <p className="text-slate-400 text-sm mt-2">Enter a code shared by your friend</p>
        </button>
      </div>

      {/* Modals */}
      {(showP1Modal || showP2Modal) && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-slate-800 border border-slate-700 p-8 rounded-3xl w-full max-w-md shadow-2xl">
            <h3 className="text-2xl font-bold mb-6 text-white">
              {showP1Modal ? "Start New Game" : "Join Existing Game"}
            </h3>
            
            <input
              type="text"
              placeholder="Your Username"
              className="w-full bg-slate-900 border border-slate-700 p-4 rounded-xl text-white mb-4 outline-none focus:ring-2 ring-blue-500"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />

            {showP2Modal && (
              <input
                type="text"
                placeholder="Game ID (e.g. 2024...)"
                className="w-full bg-slate-900 border border-slate-700 p-4 rounded-xl text-white mb-6 outline-none focus:ring-2 ring-red-500"
                value={gameIdInput}
                onChange={(e) => setGameIdInput(e.target.value)}
              />
            )}

            <div className="flex gap-4">
              <button 
                onClick={() => { setShowP1Modal(false); setShowP2Modal(false); }}
                className="flex-1 px-6 py-3 rounded-xl font-semibold text-slate-400 hover:bg-slate-700 transition"
              >
                Cancel
              </button>
              <button 
                onClick={() => handleStartGame(showP1Modal ? 'create' : 'join')}
                className="flex-1 bg-gradient-to-r from-blue-600 to-blue-500 px-6 py-3 rounded-xl font-bold text-white hover:scale-105 transition"
              >
                Let's Play
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}