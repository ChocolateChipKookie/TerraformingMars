import React, { useCallback, useMemo } from "react";
import styles from "../../styles/AddGamePage.module.css";
import { GAME_CONSTANTS } from "../../data/gameData";

// Award Button component
function AwardButton({ awardConfig, playerIndex, handlers }) {
  const { award, awardPlacements } = awardConfig;
  const { onCyclePlacement, isAwardFunded, getFundedAwardsCount } = handlers;
  
  const isDisabled =
    !award ||
    (!isAwardFunded(award) &&
      getFundedAwardsCount() >= GAME_CONSTANTS.MAX_AWARDS_FUNDED);
  const placement = awardPlacements[award]?.[playerIndex] || 0;

  const getButtonClass = () => {
    if (isDisabled) return styles.awardButtonDisabled;
    if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.GOLD)
      return styles.awardButtonGold;
    if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.SILVER)
      return styles.awardButtonSilver;
    return styles.awardButtonNone;
  };

  return (
    <button
      className={getButtonClass()}
      onClick={() => onCyclePlacement(award, playerIndex)}
      disabled={isDisabled}
    />
  );
}

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

// Helper component for award row
const AwardRow = React.memo(({ config, gameState, handlers }) => {
  const { award, index, isCustomizable } = config;
  const { awards, players } = gameState;
  const { cyclePlacement, isAwardFunded, getFundedAwardsCount } = handlers;
  
  const handleUpdate = useCallback(
    (newValue) => {
      awards.updateSelected(index, newValue);
    },
    [awards.updateSelected, index],
  );

  const getDropdownOptions = useCallback(() => {
    return awards.getAvailableForDropdown(award);
  }, [awards.getAvailableForDropdown, award]);

  return (
    <div className={styles.pointInputContainer}>
      {isCustomizable ? (
        <ObjectiveSelector
          value={award}
          onUpdate={handleUpdate}
          getAvailableOptions={getDropdownOptions}
          placeholder="Select Award"
        />
      ) : (
        <div className={styles.pointInputLabel}>{award}</div>
      )}

      <div className={styles.playerFieldsContainer}>
        {players.map((player, playerIndex) => (
          <AwardButton
            key={playerIndex}
            awardConfig={{ award, awardPlacements: awards.data }}
            playerIndex={playerIndex}
            handlers={{
              onCyclePlacement: cyclePlacement,
              isAwardFunded,
              getFundedAwardsCount,
            }}
          />
        ))}
      </div>
    </div>
  );
});

export { MilestoneRow, AwardRow };