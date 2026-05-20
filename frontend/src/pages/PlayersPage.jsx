import React, { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainer, SubContainerElement } from '../components/Container';
import LinkButton from '../components/LinkButton';
import { ROUTES } from '../constants/routes';
import { playerApi } from '../services/api';
import styles from '../styles/PlayersPage.module.css';

const PROVISIONAL_INFO = 'Provisional rating — based on too few games for the model to be confident yet.';

function PlayersPage() {
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [sortBy, setSortBy] = useState('name'); // 'name' | 'rating'
  const [openProvisional, setOpenProvisional] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    playerApi.getAll()
      .then(data => setPlayers(Array.isArray(data) ? data : []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const sortedPlayers = useMemo(() => {
    const arr = [...players];
    if (sortBy === 'rating') {
      // Established first (by ordinal desc), then unestablished (by ordinal desc), then unrated last.
      arr.sort((a, b) => {
        const tier = (p) => p.rating ? (p.rating.established ? 0 : 1) : 2;
        const ta = tier(a), tb = tier(b);
        if (ta !== tb) return ta - tb;
        if (ta === 2) return a.name.localeCompare(b.name);
        return b.rating.ordinal - a.rating.ordinal;
      });
    } else {
      arr.sort((a, b) => a.name.localeCompare(b.name));
    }
    return arr;
  }, [players, sortBy]);

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

  const renderRating = (player) => {
    if (!player.rating) return null;
    const { ordinal, established } = player.rating;
    if (established) return <>Rating: {ordinal.toFixed(1)}</>;
    const handleClick = (e) => {
      e.preventDefault();
      e.stopPropagation();
      setOpenProvisional(openProvisional === player.id ? null : player.id);
    };
    return (
      <span
        title={PROVISIONAL_INFO}
        onClick={handleClick}
        style={{ color: '#aaa', cursor: 'help' }}
      >
        Rating: {ordinal.toFixed(1)}*
      </span>
    );
  };

  return (
    <Layout>
      <Container title="Players">
        <SubContainer title="All Players">
          <SubContainerElement>
            <div className={styles.sortContainer}>
              <label htmlFor="sort-select">Sort by:</label>
              <select
                id="sort-select"
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className={styles.sortSelect}
              >
                <option value="name">Name</option>
                <option value="rating">Rating</option>
              </select>
            </div>
            <div className={styles.playersList}>
              {sortedPlayers.length > 0 ? (
                sortedPlayers.map(player => (
                  <React.Fragment key={player.id}>
                    <a
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
                        {renderRating(player)}
                      </div>
                    </a>
                    {openProvisional === player.id && (
                      <div style={{ padding: '0 1.5rem', fontSize: '0.9rem', color: '#666', fontStyle: 'italic' }}>
                        {PROVISIONAL_INFO}
                      </div>
                    )}
                  </React.Fragment>
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
