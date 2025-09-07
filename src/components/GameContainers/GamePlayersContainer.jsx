import React from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import styles from '../../styles/AddGamePage.module.css';

// PlayerInput component
const PlayerInput = React.memo(
  ({ index, player, corporations, onUpdate, selectedCorporations, availablePlayers, selectedPlayers }) => {
    const availableCorporations = React.useMemo(
      () =>
        corporations.filter(
          (corp) =>
            corp === "Beginner" ||
            !selectedCorporations.includes(corp) ||
            corp === player.corporation,
        ),
      [corporations, selectedCorporations, player.corporation],
    );

    // Filter out already selected players
    const availablePlayersForSelection = React.useMemo(
      () =>
        availablePlayers.filter(
          (p) =>
            !selectedPlayers.includes(p.name) ||
            p.name === player.name
        ),
      [availablePlayers, selectedPlayers, player.name],
    );

    return (
      <div className={styles.playerInputDiv}>
        <select
          className={styles.containerInput}
          value={player.name}
          onChange={(e) => onUpdate(index, "name", e.target.value)}
        >
          <option value="">Select Player</option>
          {availablePlayersForSelection.map((p) => (
            <option key={p.id} value={p.name}>
              {p.name}
            </option>
          ))}
        </select>
        <select
          className={styles.containerInput}
          value={player.corporation}
          onChange={(e) => onUpdate(index, "corporation", e.target.value)}
        >
          <option value="">Select Corporation</option>
          {availableCorporations.map((corp) => (
            <option key={corp} value={corp}>
              {corp}
            </option>
          ))}
        </select>
      </div>
    );
  },
);

function GamePlayersContainer({ 
  playerManager,
  availablePlayers,
  selectedCorporations,
  getAvailableCorporations
}) {
  return (
    <Container title="Players" titleStyle="banner">
      <SubContainer>
        {playerManager.players.map((player, index) => (
          <PlayerInput
            key={index}
            index={index}
            player={player}
            corporations={getAvailableCorporations()}
            onUpdate={playerManager.updatePlayerData}
            selectedCorporations={selectedCorporations}
            availablePlayers={availablePlayers}
            selectedPlayers={playerManager.players.map(p => p.name).filter(n => n)}
          />
        ))}
      </SubContainer>
    </Container>
  );
}

export default GamePlayersContainer;