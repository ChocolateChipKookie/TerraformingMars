import React from 'react';
import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import AddGamePage from './pages/AddGamePage';
import AddPlayerPage from './pages/AddPlayerPage';
import PlayerStatsPage from './pages/PlayerStatsPage';
import GamesPage from './pages/GamesPage';
import './App.css';

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/add-game" element={<AddGamePage />} />
      <Route path="/add-player" element={<AddPlayerPage />} />
      <Route path="/player-stats" element={<PlayerStatsPage />} />
      <Route path="/games" element={<GamesPage />} />
      <Route path="/games/:id" element={<div>Game Detail Page (TODO)</div>} />
    </Routes>
  );
}

export default App;