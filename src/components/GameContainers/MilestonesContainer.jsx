import React, { useCallback, useMemo } from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import styles from '../../styles/GamePage.module.css';
import { GAME_CONSTANTS } from '../../data/gameData';
import { SelectField } from './FormFields';


// Helper component for milestone/award label display
const ObjectiveLabel = React.memo(({ value }) => {
  return <div className={styles.milestoneLabel}>{value}</div>;
});

// Helper component for milestone row
const MilestoneRow = React.memo(({ config, gameState, handlers, readOnly = false }) => {
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
        <SelectField
          value={milestone}
          onChange={(e) => handleUpdate(e.target.value)}
          options={getDropdownOptions()}
          placeholder="Select Milestone"
          readOnly={readOnly}
          className={styles.milestoneDropdown}
        />
      ) : (
        <ObjectiveLabel value={milestone} />
      )}

      <SelectField
        className={styles.containerInput}
        value={milestone ? (milestones.data[milestone] ?? -1) : -1}
        onChange={handleWinnerChange}
        disabled={isMilestoneDisabled}
        readOnly={readOnly}
        displayFormatter={(val) => {
          if (val === -1 || val == null) return "Not achieved";
          return players[val]?.name || `Player ${val + 1}`;
        }}
      >
        <option value={-1}>Not achieved</option>
        {players.map((p, i) => (
          <option key={i} value={i}>
            {p.name || `Player ${i + 1}`}
          </option>
        ))}
      </SelectField>
    </div>
  );
});

function MilestonesContainer({ 
  milestones,
  gameConfig,
  playerManager,
  getSelectedMilestonesCount,
  updateMilestoneWinner,
  readOnly = false
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
            readOnly={readOnly}
          />
        ))}
      </SubContainer>
    </Container>
  );
}

export default MilestonesContainer;