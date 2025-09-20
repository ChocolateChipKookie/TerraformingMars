import React from 'react';
import styles from '../styles/GamePage.module.css';
import { InputField } from './GameContainers/FormFields';

// Numeric input component with increment/decrement buttons
export const NumericInputWithButtons = React.memo(({
  value,
  onChange,
  onDecrement,
  onIncrement,
  readOnly = false,
}) => {
  if (readOnly) {
    return (
      <div className={styles.numericInputContainer}>
        <InputField
          type="text"
          className={styles.numericInput}
          value={value}
          readOnly={true}
          style={{ textAlign: 'center', width: '100%' }}
        />
      </div>
    );
  }
  
  return (
    <div className={styles.numericInputContainer}>
      <button className={styles.numericInputButton} onClick={onDecrement}>
        −
      </button>
      <InputField
        type="text"
        inputMode="numeric"
        pattern="[0-9]*"
        className={styles.numericInput}
        value={value}
        onChange={onChange}
      />
      <button className={styles.numericInputButton} onClick={onIncrement}>
        +
      </button>
    </div>
  );
});

// Player names header component for grid layouts
export const PlayerNamesHeader = React.memo(({ players }) => {
  return (
    <div className={styles.playerNamesHeader}>
      <div className={styles.playerNamesEmpty}>
        {/* Empty space for alignment */}
      </div>
      <div className={styles.playerFieldsContainer}>
        {players.map((player, playerIndex) => (
          <div key={playerIndex} className={styles.playerNameCell}>
            {player.name || `P${playerIndex + 1}`}
          </div>
        ))}
      </div>
    </div>
  );
});