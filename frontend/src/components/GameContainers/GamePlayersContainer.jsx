import React from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import styles from '../../styles/GamePage.module.css';
import { SelectField } from './FormFields';

// PlayerInput component
const PlayerInput = React.memo(
  ({ index, player, corporations, onUpdate, selectedCorporations, availablePlayers, selectedPlayers, readOnly = false, isLegacyMode = false }) => {
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
      <div className={isLegacyMode ? styles.playerInputDivCentered : styles.playerInputDiv}>
        <SelectField
          className={styles.containerInput}
          value={player.name}
          onChange={(e) => onUpdate(index, "name", e.target.value)}
          readOnly={readOnly}
          placeholder="Select Player"
        >
          {availablePlayersForSelection.map((p) => (
            <option key={p.id} value={p.name}>
              {p.name}
            </option>
          ))}
        </SelectField>
        {!isLegacyMode && (
          <SelectField
            className={styles.containerInput}
            value={player.corporation}
            onChange={(e) => onUpdate(index, "corporation", e.target.value)}
            readOnly={readOnly}
            placeholder="Select Corporation"
          >
            {availableCorporations.map((corp) => (
              <option key={corp} value={corp}>
                {corp}
              </option>
            ))}
          </SelectField>
        )}
      </div>
    );
  },
);

function GamePlayersContainer({
  playerManager,
  availablePlayers,
  selectedCorporations,
  getAvailableCorporations,
  readOnly = false,
  isLegacyMode = false
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
            readOnly={readOnly}
            isLegacyMode={isLegacyMode}
          />
        ))}
      </SubContainer>
    </Container>
  );
}

export default GamePlayersContainer;