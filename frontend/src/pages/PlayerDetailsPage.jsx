import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainer, SubContainerElement } from '../components/Container';
import LinkButton from '../components/LinkButton';
import { InputField } from '../components/GameContainers/FormFields';
import { ROUTES } from '../constants/routes';
import { ratingApi } from '../services/api';
import styles from '../styles/GamePage.module.css';
import listStyles from '../styles/GamesPage.module.css';

function PlayerDetailsPage() {
  const { playerId } = useParams();
  const navigate = useNavigate();
  const [playerInfo, setPlayerInfo] = useState(null);
  const [ratingHistory, setRatingHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showProvisionalInfo, setShowProvisionalInfo] = useState(false);

  // Must mirror backend rating constants (OpenSkill default mu/sigma, EstablishedSigmaThreshold).
  const INITIAL_SIGMA = 8.333;
  const ESTABLISHED_SIGMA = 6.0;

  const provisionalProgress = (sigma) => {
    const pct = ((INITIAL_SIGMA - sigma) / (INITIAL_SIGMA - ESTABLISHED_SIGMA)) * 100;
    return Math.max(0, Math.min(99, pct));
  };

  const provisionalInfo = (sigma, games) =>
    `Provisional rating — based on ${games} game${games === 1 ? '' : 's'}, ${provisionalProgress(sigma).toFixed(0)}% to established.`;

  useEffect(() => {
    fetchAll();
  }, [playerId]);

  const fetchAll = async () => {
    try {
      setLoading(true);
      const [infoResp, history] = await Promise.all([
        fetch(`/api/players/${playerId}/extended`).then(r => {
          if (!r.ok) throw new Error('Failed to fetch player information');
          return r.json();
        }),
        ratingApi.getHistory(playerId),
      ]);
      setPlayerInfo(infoResp);
      // Most recent first.
      setRatingHistory([...history].reverse());
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

        {ratingHistory.length > 0 && (
          <>
            <SubContainerElement>
              <label>Rating:</label>
              <InputField
                className={styles.optionInput}
                value={
                  <span
                    style={{
                      color: ratingHistory[0].established ? 'inherit' : '#aaa',
                      cursor: ratingHistory[0].established ? 'inherit' : 'help',
                    }}
                    title={ratingHistory[0].established ? undefined : provisionalInfo(ratingHistory[0].sigma, ratingHistory.length)}
                    onClick={ratingHistory[0].established ? undefined : () => setShowProvisionalInfo(v => !v)}
                  >
                    {ratingHistory[0].ordinal.toFixed(1)}
                    {!ratingHistory[0].established && '*'}
                  </span>
                }
                readOnly={true}
              />
            </SubContainerElement>
            {!ratingHistory[0].established && showProvisionalInfo && (
              <div style={{ padding: '0 1rem 0.5rem', fontSize: '0.9rem', color: '#666', fontStyle: 'italic' }}>
                {provisionalInfo(ratingHistory[0].sigma, ratingHistory.length)}
              </div>
            )}
          </>
        )}
      </Container>

      <Container title="Rating History">
        <SubContainer>
          <SubContainerElement>
            <div className={listStyles.gamesList}>
              {ratingHistory.length === 0 ? (
                <div className={listStyles.noGames}>No rated games yet.</div>
              ) : (
                ratingHistory.map(e => (
                  <a
                    key={e.game_id}
                    href="#"
                    onClick={(ev) => {
                      ev.preventDefault();
                      navigate(ROUTES.gameDetails(e.game_id), {
                        state: { from: 'player', playerId },
                      });
                    }}
                    className={listStyles.gameLink}
                  >
                    <div className={listStyles.gameName}>{e.game_name}</div>
                    <div className={listStyles.gameInfo}>
                      <span>{formatDate(e.date)}</span>
                      <span>Rating: {e.ordinal.toFixed(1)}</span>
                      <span style={{
                        color: e.delta_ordinal > 0 ? '#2e7d32' : e.delta_ordinal < 0 ? '#c62828' : 'inherit',
                      }}>
                        Δ {e.delta_ordinal > 0 ? '+' : ''}{e.delta_ordinal.toFixed(1)}
                      </span>
                    </div>
                  </a>
                ))
              )}
            </div>
          </SubContainerElement>
        </SubContainer>
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
