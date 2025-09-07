import React, { useCallback, useMemo } from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import styles from '../../styles/AddGamePage.module.css';
import { GAME_CONSTANTS } from '../../data/gameData';

// Helper component for milestone/award selection dropdown
const ObjectiveSelector = React.memo(
  ({ value, onUpdate, getAvailableOptions, placeholder = "Select..." }) => {
    const availableOptions = useMemo(
      () => getAvailableOptions(),
      [getAvailableOptions],
    );

    const handleChange = useCallback(
      (e) => {
        onUpdate(e.target.value);
      },
      [onUpdate],
    );

    return (
      <select
        className={styles.containerInput}
        value={value || ""}
        onChange={handleChange}
      >
        {!value && <option value="">{placeholder}</option>}
        {availableOptions.map((option) => (
          <option key={option} value={option}>
            {option}
          </option>
        ))}
      </select>
    );
  },
);

// Helper component for milestone/award label display
const ObjectiveLabel = React.memo(({ value }) => {
  return <div className={styles.milestoneLabel}>{value}</div>;
});

// Helper component for milestone row
const MilestoneRow = React.memo(({ config, gameState, handlers }) => {
  const { milestone, index, isCustomizable } = config;
  const { milestones, players, getSelectedMilestonesCount } = gameState;
  const { updateMilestoneWinner } = handlers;
  
  const isMilestoneDisabled = useMemo(
    () =>
      !milestone ||
      (milestones.data[milestone] === -1 &&
        getSelectedMilestonesCount >= GAME_CONSTANTS.MAX_MILESTONES_CLAIMED),
    [milestone, milestones.data, getSelectedMilestonesCount],
  );

  const handleUpdate = useCallback(
    (newValue) => {
      milestones.updateSelected(index, newValue);
    },
    [milestones.updateSelected, index],
  );

  const handleWinnerChange = useCallback(
    (e) => {
      updateMilestoneWinner(milestone, parseInt(e.target.value));
    },
    [updateMilestoneWinner, milestone],
  );

  const getDropdownOptions = useCallback(() => {
    return milestones.getAvailableForDropdown(milestone);
  }, [milestones.getAvailableForDropdown, milestone]);

  return (
    <div className={styles.playerInputDiv}>
      {isCustomizable ? (
        <ObjectiveSelector
          value={milestone}
          onUpdate={handleUpdate}
          getAvailableOptions={getDropdownOptions}
          placeholder="Select Milestone"
        />
      ) : (
        <ObjectiveLabel value={milestone} />
      )}

      <select
        className={styles.containerInput}
        value={milestones.data[milestone] ?? -1}
        onChange={handleWinnerChange}
        disabled={isMilestoneDisabled}
      >
        <option value={-1}>Not achieved</option>
        {players.map((p, i) => (
          <option key={i} value={i}>
            {p.name || `Player ${i + 1}`}
          </option>
        ))}
      </select>
    </div>
  );
});

function MilestonesContainer({ 
  milestones,
  gameConfig,
  playerManager,
  getSelectedMilestonesCount,
  updateMilestoneWinner
}) {
  return (
    <Container title="Milestones" titleStyle="banner">
      <SubContainer>
        {milestones.selected.map((milestone, index) => (
          <MilestoneRow
            key={index}
            config={{
              milestone,
              index,
              isCustomizable: gameConfig.expansions["Milestones & Awards"],
            }}
            gameState={{
              milestones,
              players: playerManager.players,
              getSelectedMilestonesCount,
            }}
            handlers={{ updateMilestoneWinner }}
          />
        ))}
      </SubContainer>
    </Container>
  );
}

export default MilestonesContainer;