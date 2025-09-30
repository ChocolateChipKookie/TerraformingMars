import React, { useCallback } from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import { PlayerNamesHeader } from '../Common';
import styles from '../../styles/GamePage.module.css';

// Helper component for read-only score display
const ReadOnlyScoreCell = React.memo(({ playerIndex, playerScores, field }) => {
  const isTotal = field === "totalPoints";

  return (
    <div
      key={playerIndex}
      className={isTotal ? styles.scoreCellTotal : styles.scoreCell}
    >
      {playerScores[playerIndex]?.[field] || 0}
    </div>
  );
});

// Helper component for editable score input
const EditableScoreInput = React.memo(
  ({ playerIndex, playerScores, field, onChange, placeholder }) => {
    const handleChange = useCallback(
      (e) => {
        const value = e.target.value;
        if (value === "" || /^\d+$/.test(value)) {
          onChange(playerIndex, field, value);
        }
      },
      [onChange, playerIndex, field],
    );

    return (
      <input
        key={playerIndex}
        type="text"
        inputMode="numeric"
        pattern="[0-9]*"
        className={styles.scoreInput}
        value={playerScores[playerIndex]?.[field] || ""}
        onChange={handleChange}
        placeholder={placeholder}
      />
    );
  },
);

function PointInput({ config, gameState, options = {} }) {
  const { label, field } = config;
  const { players, playerScores } = gameState;
  const { onChange, readOnly = false, placeholder = "0" } = options;
  
  return (
    <div className={styles.pointInputContainer}>
      <div className={styles.pointInputLabel}>{label}</div>
      <div className={styles.playerFieldsContainer}>
        {players.map((_, playerIndex) =>
          readOnly ? (
            <ReadOnlyScoreCell
              key={playerIndex}
              playerIndex={playerIndex}
              playerScores={playerScores}
              field={field}
            />
          ) : (
            <EditableScoreInput
              key={playerIndex}
              playerIndex={playerIndex}
              playerScores={playerScores}
              field={field}
              onChange={onChange}
              placeholder={placeholder}
            />
          ),
        )}
      </div>
    </div>
  );
}

function PointsContainer({
  playerManager,
  gameConfig,
  readOnly = false,
  isLegacyMode = false
}) {
  const sharedGameState = {
    players: playerManager.players,
    playerScores: playerManager.playerScores,
  };
  const editableOptions = {
    onChange: playerManager.updatePlayerScore,
    readOnly: readOnly,
  };
  const readOnlyOptions = { readOnly: true };

  return (
    <Container title="Points" titleStyle="banner">
      <SubContainer>
        <PlayerNamesHeader players={playerManager.players} />

        <>
          <PointInput
            config={{ label: "TR", field: "terraformingRating" }}
            gameState={sharedGameState}
            options={editableOptions}
          />
          <PointInput
            config={{ label: "Cities", field: "cities" }}
            gameState={sharedGameState}
            options={editableOptions}
          />
          <PointInput
            config={{ label: "Greeneries", field: "greeneries" }}
            gameState={sharedGameState}
            options={editableOptions}
          />
          <PointInput
            config={{ label: "Cards", field: "cards" }}
            gameState={sharedGameState}
            options={editableOptions}
          />
          {gameConfig.expansions.Turmoil && (
            <PointInput
              config={{ label: "Turmoil", field: "turmoilPoints" }}
              gameState={sharedGameState}
              options={editableOptions}
            />
          )}

          <div className={styles.scoreSeparator}>
            <PointInput
              config={{ label: "Milestones", field: "milestonePoints" }}
              gameState={sharedGameState}
              options={isLegacyMode ? editableOptions : readOnlyOptions}
            />
            <PointInput
              config={{ label: "Awards", field: "awardPoints" }}
              gameState={sharedGameState}
              options={isLegacyMode ? editableOptions : readOnlyOptions}
            />
          </div>

          <div className={styles.scoreSeparatorFinal}>
            <PointInput
              config={{ label: "Total", field: "totalPoints" }}
              gameState={sharedGameState}
              options={readOnlyOptions}
            />
          </div>
        </>
      </SubContainer>
    </Container>
  );
}

export default PointsContainer;