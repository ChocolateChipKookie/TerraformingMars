import React, { useCallback, useMemo } from 'react';
import Container from '../Container';
import { SubContainer } from '../Container';
import { PlayerNamesHeader } from '../Common';
import styles from '../../styles/GamePage.module.css';
import { gameData, GAME_CONSTANTS } from '../../data/gameData';
import { SelectField } from './FormFields';

// Award Button component
function AwardButton({ awardConfig, playerIndex, handlers, readOnly = false }) {
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
      disabled={isDisabled || readOnly}
    />
  );
}


// Helper component for award row
const AwardRow = React.memo(({ config, gameState, handlers, readOnly = false }) => {
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
      {isCustomizable && !readOnly ? (
        <SelectField
          value={award}
          onChange={(e) => handleUpdate(e.target.value)}
          options={getDropdownOptions()}
          placeholder="Select Award"
          className={styles.objectiveDropdown}
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
            readOnly={readOnly}
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
  getFundedAwardsCount,
  readOnly = false
}) {
  return (
    <Container title="Awards" titleStyle="banner">
      <SubContainer>
        <PlayerNamesHeader players={playerManager.players} />

        {awards.selected.map((award, index) => {
          const isVenusLocked = gameConfig.expansions["Venus Next"]
            && (gameData.awards.Venus || []).includes(award);
          return (
          <AwardRow
            key={index}
            config={{
              award,
              index,
              isCustomizable: gameConfig.expansions["Milestones & Awards"] && !isVenusLocked,
            }}
            gameState={{
              awards,
              players: playerManager.players,
            }}
            handlers={{ cyclePlacement, isAwardFunded, getFundedAwardsCount }}
            readOnly={readOnly}
          />
          );
        })}
      </SubContainer>
    </Container>
  );
}

export default AwardsContainer;