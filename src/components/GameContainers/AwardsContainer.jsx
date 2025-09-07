import React, { useCallback, useMemo } from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import styles from '../../styles/AddGamePage.module.css';
import { GAME_CONSTANTS } from '../../data/gameData';

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

function AwardsContainer({ 
  awards,
  gameConfig,
  playerManager,
  cyclePlacement,
  isAwardFunded,
  getFundedAwardsCount
}) {
  return (
    <Container title="Awards" titleStyle="banner">
      <SubContainer>
        <PlayerNamesHeader players={playerManager.players} />

        {awards.selected.map((award, index) => (
          <AwardRow
            key={index}
            config={{
              award,
              index,
              isCustomizable: gameConfig.expansions["Milestones & Awards"],
            }}
            gameState={{
              awards,
              players: playerManager.players,
            }}
            handlers={{ cyclePlacement, isAwardFunded, getFundedAwardsCount }}
          />
        ))}
      </SubContainer>
    </Container>
  );
}

export default AwardsContainer;