import React, { useMemo } from 'react';
import Container from '../Container';
import { SubContainer, SubContainerElement } from '../Container';
import styles from '../../styles/GamePage.module.css';
import { gameData, GAME_CONSTANTS } from '../../data/gameData';
import { SelectField, InputField } from './FormFields';

// ExpansionIcon component
function ExpansionIcon({
  expansion,
  checked,
  disabled,
  onChange,
  children,
  showText = true,
}) {
  const getExpansionConfig = (expansion) => {
    switch (expansion) {
      case "Base Game":
        return {
          backgroundColor: "transparent",
          symbolColor: "#000",
          svg: null,
        };
      case "Corporate Era":
        return {
          backgroundColor: "#FF0000",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "0deg",
        };
      case "Venus Next":
        return {
          backgroundColor: "#87CEEB",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "180deg",
        };
      case "Prelude":
        return {
          backgroundColor: "#FFC0CB",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "270deg",
        };
      case "Prelude 2":
        return {
          backgroundColor: "#FF69B4",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "270deg",
        };
      case "Colonies":
        return {
          backgroundColor: "#808080",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "0deg",
        };
      case "Turmoil":
        return {
          backgroundColor: "#FFA500",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "180deg",
        };
      case "Milestones & Awards":
        return {
          backgroundColor: "#FFD700",
          symbolColor: "#000",
          svg: <rect x="6" y="6" width="12" height="12" fill="currentColor" />,
        };
      case "Promo":
        return {
          backgroundColor: "#2F2F2F",
          symbolColor: "#FFF",
          svg: <circle cx="12" cy="12" r="8" fill="currentColor" />,
        };
      case "Draft":
        return {
          backgroundColor: "#4B0082",
          symbolColor: "#FFF",
          svg: null,
          fallbackSymbol: "D",
        };
      default:
        return {
          backgroundColor: "#666",
          symbolColor: "#FFF",
          svg: null,
          fallbackSymbol: "?",
        };
    }
  };

  const config = getExpansionConfig(expansion);
  const isShaded = !checked;

  return (
    <label
      className={
        showText
          ? disabled
            ? styles.expansionLabelExpandedDisabled
            : styles.expansionLabelExpanded
          : disabled
            ? styles.expansionLabelCompactDisabled
            : styles.expansionLabelCompact
      }
    >
      <input
        type="checkbox"
        checked={checked}
        disabled={disabled}
        onChange={onChange}
        className={styles.hiddenCheckbox}
      />
      <div
        className={
          isShaded
            ? styles.expansionIconShaded
            : disabled
              ? styles.expansionIconDisabled
              : styles.expansionIcon
        }
        style={{
          backgroundColor: config.backgroundColor,
          color: config.symbolColor,
        }}
        title={expansion}
      >
        {config.svg ? (
          <svg
            className={styles.expansionIconSvg}
            viewBox="0 0 24 24"
            style={{
              transform: config.rotation
                ? `rotate(${config.rotation})`
                : undefined,
            }}
          >
            {config.svg}
          </svg>
        ) : (
          config.fallbackSymbol || ""
        )}
      </div>
      {showText && <span>{children}</span>}
    </label>
  );
}

// Local components for this container only
function NumericInputWithButtons({
  value,
  onChange,
  onDecrement,
  onIncrement,
  readOnly = false,
}) {
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
}

function GameOptionsContainer({ 
  gameConfig, 
  dispatch,
  playerManager,
  maxPlayers,
  readOnly = false
}) {
  // Get map options memoized
  const mapOptions = useMemo(() => {
    return Object.keys(gameData.maps).map((m) => (
      <option key={m} value={m}>
        {m}
      </option>
    ));
  }, []);

  const expansionEntries = useMemo(() => {
    return Object.entries(gameConfig.expansions);
  }, [gameConfig.expansions]);

  return (
    <Container title="Options" titleStyle="banner">
      <SubContainer>
        <SubContainerElement>
          <label>Name:</label>
          <InputField
            type="text"
            className={styles.optionInput}
            value={gameConfig.name}
            onChange={(e) =>
              dispatch({
                type: "SET_FIELD",
                field: "name",
                value: e.target.value,
              })
            }
            readOnly={readOnly}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Date:</label>
          <InputField
            type="date"
            className={styles.optionInput}
            value={gameConfig.date}
            onChange={(e) =>
              dispatch({ type: "SET_DATE", value: e.target.value })
            }
            required={true}
            readOnly={readOnly}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Map:</label>
          <SelectField
            className={styles.optionInput}
            value={gameConfig.map}
            onChange={(e) =>
              dispatch({ type: "SET_MAP", value: e.target.value })
            }
            readOnly={readOnly}
          >
            {mapOptions}
          </SelectField>
        </SubContainerElement>

        <SubContainerElement>
          <label>Generations:</label>
          <NumericInputWithButtons
            value={gameConfig.generations}
            onChange={(e) => {
              const val = parseInt(e.target.value) || 0;
              dispatch({ type: "SET_GENERATIONS", value: val });
            }}
            onDecrement={() =>
              dispatch({
                type: "SET_GENERATIONS",
                value: gameConfig.generations - 1,
              })
            }
            onIncrement={() =>
              dispatch({
                type: "SET_GENERATIONS",
                value: gameConfig.generations + 1,
              })
            }
            readOnly={readOnly}
          />
        </SubContainerElement>

        <SubContainerElement>
          <label>Players:</label>
          <NumericInputWithButtons
            value={playerManager.playerNumber}
            onChange={(e) => {
              const val =
                parseInt(e.target.value) || GAME_CONSTANTS.MIN_PLAYERS;
              playerManager.setPlayerNumber(Math.min(maxPlayers, val));
            }}
            onDecrement={() =>
              playerManager.setPlayerNumber(
                Math.max(
                  GAME_CONSTANTS.MIN_PLAYERS,
                  playerManager.playerNumber - 1,
                ),
              )
            }
            onIncrement={() =>
              playerManager.setPlayerNumber(
                Math.min(maxPlayers, playerManager.playerNumber + 1),
              )
            }
            readOnly={readOnly}
          />
        </SubContainerElement>

        <SubContainerElement>
          <div className={styles.subcontainerBox}>
            <div className={styles.expansionsContainer}>
              {/* Content container - shared by both bar and list */}
              <div className={styles.expansionsContent}>
                {!gameConfig.expandedExpansions && (
                  <div className={styles.expansionsCompactView}>
                    {expansionEntries.map(([key, value]) => (
                      <ExpansionIcon
                        key={key}
                        expansion={key}
                        checked={value}
                        disabled={key === "Base Game" || readOnly}
                        onChange={() =>
                          readOnly ? null : dispatch({
                            type: "TOGGLE_EXPANSION",
                            expansion: key,
                          })
                        }
                        showText={false}
                      >
                        {key}
                      </ExpansionIcon>
                    ))}
                  </div>
                )}

                {gameConfig.expandedExpansions && (
                  <div className={styles.expansionsExpandedView}>
                    {expansionEntries.map(([key, value]) => (
                      <ExpansionIcon
                        key={key + "_expanded"}
                        expansion={key}
                        checked={value}
                        disabled={key === "Base Game" || readOnly}
                        onChange={() =>
                          readOnly ? null : dispatch({
                            type: "TOGGLE_EXPANSION",
                            expansion: key,
                          })
                        }
                        showText={true}
                      >
                        {key}
                      </ExpansionIcon>
                    ))}
                  </div>
                )}
              </div>

              {/* Expand/Collapse button container */}
              <div className={styles.expansionsToggleContainer}>
                <div
                  className={styles.expansionsToggleButton}
                  onClick={() => dispatch({ type: "TOGGLE_EXPANDED_VIEW" })}
                  title={
                    gameConfig.expandedExpansions ? "Collapse" : "Expand"
                  }
                >
                  ☰
                </div>
              </div>
            </div>
          </div>
        </SubContainerElement>
      </SubContainer>
    </Container>
  );
}

export default GameOptionsContainer;