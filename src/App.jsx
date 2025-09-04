import React from 'react';
import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import AddGamePage from './pages/AddGamePage';
import AddPlayerPage from './pages/AddPlayerPage';
import PlayerStatsPage from './pages/PlayerStatsPage';
import GamesPage from './pages/GamesPage';
import { ROUTES } from './constants/routes';
import './App.css';

function App() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<HomePage />} />
      <Route path={ROUTES.ADD_GAME} element={<AddGamePage />} />
      <Route path={ROUTES.ADD_PLAYER} element={<AddPlayerPage />} />
      <Route path={ROUTES.PLAYER_STATS} element={<PlayerStatsPage />} />
      <Route path={ROUTES.PLAYED_GAMES} element={<GamesPage />} />
      <Route path="/games/:id" element={<div>Game Detail Page (TODO)</div>} />
    </Routes>
  );
}

export default App;
