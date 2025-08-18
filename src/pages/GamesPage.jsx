import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function GamesPage() {
  const navigate = useNavigate();
  const [games, setGames] = useState([]);
  const [searchValue, setSearchValue] = useState('');
  const [showGameList, setShowGameList] = useState(false);

  useEffect(() => {
    fetch('/data/games.json')
      .then(response => response.json())
      .then(data => setGames(data))
      .catch(error => console.error('Error loading games:', error));
  }, []);

  const redirectToGame = () => {
    const game = games.find(g => g.name === searchValue);
    if (game) {
      window.location.href = `/games/${game.id}.html`;
    }
  };

  const hideContainer = () => {
    setShowGameList(!showGameList);
  };

  return (
    <div>
      <div>Played games</div>
      
      <div>
        <div>Links</div>
        <div>
          <div>
            <div>
              <label>Search:</label>
              <button onClick={redirectToGame}>GO</button>
              <input 
                type="text" 
                placeholder="Game" 
                value={searchValue}
                onChange={(e) => setSearchValue(e.target.value)}
                list="datalist-games"
              />
              <datalist id="datalist-games">
                {games.map(game => (
                  <option key={game.id} value={game.name} />
                ))}
              </datalist>
            </div>
          </div>
          <div>
            <div>
              <label>Select game:</label>
              <button onClick={hideContainer}>
                {showGameList ? "Hide" : "Show"}
              </button>
            </div>
            <div style={{ display: showGameList ? 'block' : 'none' }}>
              {games.map(game => (
                <a key={game.id} href={`/games/${game.id}.html`}>
                  {game.name}
                </a>
              ))}
            </div>
          </div>
        </div>
      </div>
      <div>
        <button onClick={() => navigate('/')}>Main page</button>
      </div>
    </div>
  );
}

export default GamesPage;