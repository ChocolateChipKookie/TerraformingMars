import React from 'react';
import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import GamePage from './pages/GamePage';
import PlayerFormPage from './pages/PlayerFormPage';
import PlayersPage from './pages/PlayersPage';
import PlayerDetailsPage from './pages/PlayerDetailsPage';
import GamesPage from './pages/GamesPage';
import { ROUTES } from './constants/routes';
import './App.css';

function App() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<HomePage />} />
      <Route path={ROUTES.ADD_GAME} element={<GamePage />} />
      <Route path={ROUTES.ADD_PLAYER} element={<PlayerFormPage />} />
      <Route path={ROUTES.EDIT_PLAYER} element={<PlayerFormPage />} />
      <Route path={ROUTES.PLAYERS} element={<PlayersPage />} />
      <Route path={ROUTES.PLAYER_DETAILS} element={<PlayerDetailsPage />} />
      <Route path={ROUTES.PLAYED_GAMES} element={<GamesPage />} />
      <Route path={ROUTES.GAME_DETAILS} element={<GamePage />} />
    </Routes>
  );
}

export default App;
