import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainer, SubContainerElement } from '../components/Container';
import LinkButton from '../components/LinkButton';
import styles from '../styles/GamesPage.module.css';

function GamesPage() {
  const navigate = useNavigate();
  const [games, setGames] = useState([]);
  const [searchValue, setSearchValue] = useState('');

  useEffect(() => {
    fetch('http://localhost:8080/api/games')
      .then(response => response.json())
      .then(data => {
        if (data && Array.isArray(data)) {
          setGames(data);
        } else {
          setGames([]);
        }
      })
      .catch(error => {
        console.error('Error loading games:', error);
        setGames([]);
      });
  }, []);

  const redirectToGame = () => {
    const game = games.find(g => g.name === searchValue);
    if (game) {
      navigate(`/game/${game.id}`);
    }
  };

  const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    });
  };

  const formatExpansions = (expansions) => {
    if (!expansions) return '';
    return Object.entries(expansions)
      .filter(([_, enabled]) => enabled)
      .map(([name]) => name)
      .join(', ');
  };

  return (
    <Layout>
      <Container title="Played Games">
        <SubContainer title="Search Games">
          <SubContainerElement>
            <div className={styles.searchContainer}>
              <label className={styles.searchLabel}>Search:</label>
              <input 
                type="text" 
                placeholder="Enter game name" 
                value={searchValue}
                onChange={(e) => setSearchValue(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && redirectToGame()}
                list="datalist-games"
                className={styles.searchInput}
              />
              <datalist id="datalist-games">
                {games.map(game => (
                  <option key={game.id} value={game.name} />
                ))}
              </datalist>
              <button onClick={redirectToGame} className={styles.searchButton}>
                GO
              </button>
            </div>
          </SubContainerElement>
        </SubContainer>

        <SubContainer title="All Games">
          <SubContainerElement>
            <div className={styles.gamesList}>
              {games.length > 0 ? (
                games.map(game => (
                  <a 
                    key={game.id} 
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      navigate(`/game/${game.id}`);
                    }}
                    className={styles.gameLink}
                  >
                    <div className={styles.gameName}>{game.name}</div>
                    <div className={styles.gameInfo}>
                      <span>Date: {formatDate(game.date)}</span>
                      <span>Map: {game.map}</span>
                      <span>Gen: {game.generations}</span>
                      {formatExpansions(game.expansions) && (
                        <span>Exp: {formatExpansions(game.expansions)}</span>
                      )}
                    </div>
                  </a>
                ))
              ) : (
                <div className={styles.noGames}>No games found</div>
              )}
            </div>
          </SubContainerElement>
        </SubContainer>
      </Container>

      <div style={{ 
        display: 'flex', 
        justifyContent: 'center',
        margin: '1% auto',
        maxWidth: '900px',
        padding: '0'
      }}>
        <LinkButton 
          onClick={() => navigate('/')}
          style={{ width: '200px' }}
        >
          Main Page
        </LinkButton>
      </div>
    </Layout>
  );
}

export default GamesPage;