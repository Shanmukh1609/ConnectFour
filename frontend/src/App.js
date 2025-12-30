import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { GameProvider } from './context/GameContext'; // Import the provider
import LandingPage from './components/LandingPage';
import WaitingRoom from './components/WaitingRoom';
import GameBoard from './components/GameBoard';
import Leaderboard from './components/LeaderBoard';

function App() {
  return (
    // GameProvider must wrap the Router to provide context to all routes
    <GameProvider>
      <Router>
        <div className="min-h-screen bg-slate-950 text-white">
          <Toaster 
            position="top-center" 
            toastOptions={{
              className: 'bg-slate-800 text-white border border-slate-700',
            }} 
          />
          <Routes>
            <Route path="/" element={<LandingPage />} />
            <Route path="/waiting/:gameId" element={<WaitingRoom />} />
            <Route path="/game/:gameId" element={<GameBoard />} />
            <Route path="leaderboard" element={<Leaderboard/>} />
          </Routes>
        </div>
      </Router>
    </GameProvider>
  );
}

export default App;