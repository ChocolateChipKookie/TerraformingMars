import React from 'react';
import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import GamePage from './pages/GamePage';
import AddPlayerPage from './pages/AddPlayerPage';
import PlayersPage from './pages/PlayersPage';
import GamesPage from './pages/GamesPage';
import { ROUTES } from './constants/routes';
import './App.css';

function App() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<HomePage />} />
      <Route path={ROUTES.ADD_GAME} element={<GamePage />} />
      <Route path={ROUTES.ADD_PLAYER} element={<AddPlayerPage />} />
      <Route path={ROUTES.PLAYERS} element={<PlayersPage />} />
      <Route path={ROUTES.PLAYED_GAMES} element={<GamesPage />} />
      <Route path={ROUTES.GAME_DETAILS} element={<GamePage />} />
    </Routes>
  );
}

export default App;
