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
      const response = await fetch(`http://localhost:8080/api/players/${playerId}/extended`);
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
      <Container title="Player Information">
        <SubContainerElement>
          <label>Name:</label>
          <InputField
            className={styles.optionInput}
            value={playerInfo.player.name}
            readOnly={true}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Role:</label>
          <InputField
            className={styles.optionInput}
            value={playerInfo.player.role}
            readOnly={true}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Created At:</label>
          <InputField
            className={styles.optionInput}
            value={formatDate(playerInfo.player.created_at)}
            readOnly={true}
          />
        </SubContainerElement>
      </Container>

      <Container title="Statistics">
        <SubContainerElement>
          <label>Total Games Played:</label>
          <InputField
            className={styles.optionInput}
            value={playerInfo.total_games_played}
            readOnly={true}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Total Games Won:</label>
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
        margin: '1% auto',
        maxWidth: '900px',
        padding: '0'
      }}>
        <LinkButton
          onClick={() => navigate(ROUTES.PLAYERS)}
          style={{ width: '200px' }}
        >
          Back to Players
        </LinkButton>
      </div>
    </Layout>
  );
}

export default PlayerDetailsPage;
