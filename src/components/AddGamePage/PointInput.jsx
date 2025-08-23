import React, { useCallback } from "react";
import styles from "../../styles/AddGamePage.module.css";

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

export default PointInput;