import React from 'react';
import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import GamePage from './pages/GamePage';
import AddPlayerPage from './pages/AddPlayerPage';
import PlayerStatsPage from './pages/PlayerStatsPage';
import GamesPage from './pages/GamesPage';
import { ROUTES } from './constants/routes';
import './App.css';

function App() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<HomePage />} />
      <Route path={ROUTES.ADD_GAME} element={<GamePage />} />
      <Route path={ROUTES.ADD_PLAYER} element={<AddPlayerPage />} />
      <Route path={ROUTES.PLAYER_STATS} element={<PlayerStatsPage />} />
      <Route path={ROUTES.PLAYED_GAMES} element={<GamesPage />} />
      <Route path={ROUTES.GAME_DETAILS} element={<GamePage />} />
    </Routes>
  );
}

export default App;
