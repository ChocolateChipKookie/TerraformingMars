import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainer, SubContainerElement } from '../components/Container';
import LinkButton from '../components/LinkButton';
import { ROUTES } from '../constants/routes';
import styles from '../styles/PlayersPage.module.css';

function PlayersPage() {
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchPlayers();
  }, []);

  const fetchPlayers = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/players');
      if (!response.ok) {
        throw new Error('Failed to fetch players');
      }
      const data = await response.json();
      setPlayers(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout>
        <Container title="Players" titleStyle="banner">
          <div className={styles.loading}>Loading players...</div>
        </Container>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <Container title="Players" titleStyle="banner">
          <div className={styles.error}>Error: {error}</div>
        </Container>
      </Layout>
    );
  }

  const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <Layout>
      <Container title="Players">
        <SubContainer title="All Players">
          <SubContainerElement>
            <div className={styles.playersList}>
              {players.length > 0 ? (
                players.map(player => (
                  <a
                    key={player.id}
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      navigate(ROUTES.playerDetails(player.id));
                    }}
                    className={styles.playerCard}
                  >
                    <div className={styles.playerName}>{player.name}</div>
                    <div className={styles.playerInfo}>
                      <span>Created: {formatDate(player.created_at)}</span>
                    </div>
                  </a>
                ))
              ) : (
                <div className={styles.noPlayers}>No players found</div>
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
        padding: '0 10px'
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

export default PlayersPage;
