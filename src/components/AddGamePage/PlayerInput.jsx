import React from "react";
import styles from "../../styles/AddGamePage.module.css";

// Helper component for player names header
const PlayerNamesHeader = React.memo(({ players }) => {
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

export { PlayerInput, PlayerNamesHeader };