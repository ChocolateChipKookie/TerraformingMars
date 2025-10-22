import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainerElement } from '../components/Container';
import LinkButton from '../components/LinkButton';
import { InputField } from '../components/GameContainers/FormFields';
import { ROUTES } from '../constants/routes';
import styles from '../styles/GamePage.module.css';

function PlayerDetailsPage() {
  const { playerId } = useParams();
  const navigate = useNavigate();
  const [playerInfo, setPlayerInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchPlayerInfo();
  }, [playerId]);

  const fetchPlayerInfo = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/players/${playerId}/extended`);
      if (!response.ok) {
        throw new Error('Failed to fetch player information');
      }
      const data = await response.json();
      setPlayerInfo(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
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

  if (loading) {
    return (
      <Layout>
        <Container title="Player Details">
          <div className={styles.loading}>Loading player information...</div>
        </Container>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <Container title="Player Details">
          <div className={styles.error}>Error: {error}</div>
        </Container>
      </Layout>
    );
  }

  return (
    <Layout>
      <Container title={playerInfo.player.name}>
        <SubContainerElement>
          <label>Games Played:</label>
          <InputField
            className={styles.optionInput}
            value={playerInfo.total_games_played}
            readOnly={true}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Games Won:</label>
          <InputField
            className={styles.optionInput}
            value={playerInfo.total_games_won}
            readOnly={true}
          />
        </SubContainerElement>
      </Container>

      <div style={{
        display: 'flex',
        justifyContent: 'center',
        gap: '1rem',
        margin: '1% auto',
        maxWidth: '900px',
        padding: '0 10px'
      }}>
        <LinkButton
          onClick={() => navigate(ROUTES.PLAYERS)}
          style={{ width: 'calc(50% - 0.5rem)' }}
        >
          Back to Players
        </LinkButton>
        <LinkButton
          onClick={() => navigate(ROUTES.editPlayer(playerId))}
          style={{ width: 'calc(50% - 0.5rem)' }}
        >
          Edit
        </LinkButton>
      </div>
    </Layout>
  );
}

export default PlayerDetailsPage;
