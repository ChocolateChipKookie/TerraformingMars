import React from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import LinkButton from '../components/LinkButton';
import styles from '../styles/HomePage.module.css';
import { ROUTES } from '../constants/routes';

// Navigation button configuration
const NAVIGATION_BUTTONS = [
  {
    id: 'add-game',
    label: 'Add game',
    route: ROUTES.ADD_GAME,
    enabled: true,
    description: 'Track a new Terraforming Mars game',
  },
  {
    id: 'add-player',
    label: 'Add player',
    route: ROUTES.ADD_PLAYER,
    enabled: true,
    description: 'Add a new player to the system',
  },
  {
    id: 'player-stats',
    label: 'Player stats',
    route: ROUTES.PLAYER_STATS,
    enabled: false,
    description: 'View player statistics and performance',
  },
  {
    id: 'played-games',
    label: 'Played games',
    route: ROUTES.PLAYED_GAMES,
    enabled: true,
    description: 'Browse game history and results',
  },
  {
    id: 'game-timer',
    label: 'Game timer',
    route: ROUTES.GAME_TIMER,
    enabled: false,
    description: 'Time your game sessions',
  },
];

function HomePage() {
  const navigate = useNavigate();

  const handleNavigation = (route) => {
    navigate(route);
  };

  return (
    <Layout>
      <div className={styles.container}>
        <div className={styles.buttonGrid}>
          {NAVIGATION_BUTTONS.map((button) => (
            <LinkButton
              key={button.id}
              onClick={button.enabled ? () => handleNavigation(button.route) : undefined}
              disabled={!button.enabled}
              title={button.description}
            >
              {button.label}
            </LinkButton>
          ))}
        </div>
      </div>
    </Layout>
  );
}

export default HomePage;
