import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

export default function Leaderboard() {
  const location = useLocation();
  const navigate = useNavigate();
  // Access data passed from the navigate state
  const results = location.state?.leaderboardData || [];

  return (
    <div className="min-h-screen bg-slate-950 p-8 flex flex-col items-center">
      <h1 className="text-4xl font-black text-white mb-8 italic uppercase italic tracking-tighter">Your Game History</h1>
      
      <div className="w-full max-w-2xl bg-slate-900 rounded-3xl border border-slate-800 overflow-hidden shadow-2xl">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-800 text-slate-400 uppercase text-xs font-bold tracking-widest">
              <th className="p-5">Outcome</th>
              <th className="p-5">Played At</th>
            </tr>
          </thead>
          <tbody className="text-white">
            {results.length > 0 ? results.map((res, i) => (
              <tr key={i} className="border-t border-slate-800 hover:bg-slate-800/50 transition-colors">
                <td className="p-5">
                  <span className={`px-3 py-1 rounded-full text-xs font-black ${
                    res.outcome === 'WON' ? 'bg-green-500/20 text-green-400' : 
                    res.outcome === 'LOST' ? 'bg-red-500/20 text-red-400' : 'bg-slate-700 text-slate-300'
                  }`}>
                    {res.outcome}
                  </span>
                </td>
                <td className="p-5 text-slate-400 font-mono text-sm">
                  {new Date(res.played_at).toLocaleString()}
                </td>
              </tr>
            )) : (
              <tr>
                <td colSpan="2" className="p-10 text-center text-slate-500">No games found for this player.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <button 
        onClick={() => navigate('/')}
        className="mt-8 text-blue-400 font-bold hover:text-blue-300 transition-colors uppercase tracking-widest text-sm"
      >
        ← Back to Lobby
      </button>
    </div>
  );
}