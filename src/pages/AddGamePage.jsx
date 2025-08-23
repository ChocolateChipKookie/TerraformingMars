import React, { useState, useEffect, useCallback, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import Layout from "../components/Layout";
import Container from "../components/Container";
import { SubContainer, SubContainerElement } from "../components/SubContainer";
import LinkButton from "../components/LinkButton";
import { formStyles } from "../styles/formStyles";
import { gameData, GAME_CONSTANTS } from "../data/gameData";

// Local components for this page only
function AwardButton({
  award,
  playerIndex,
  awardPlacements,
  onCyclePlacement,
  isAwardFunded,
  getFundedAwardsCount,
}) {
  const isDisabled =
    !award ||
    (!isAwardFunded(award) &&
      getFundedAwardsCount() >= GAME_CONSTANTS.MAX_AWARDS_FUNDED);
  const placement = awardPlacements[award]?.[playerIndex] || 0;

  const getBackgroundColor = () => {
    if (isDisabled) return "#666666";
    if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.GOLD) return "#FFD700"; // Gold
    if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.SILVER) return "#CCC"; // Silver
    return "#444"; // None/Black
  };

  return (
    <button
      style={{
        width: "3rem",
        height: "3rem",
        border: "2px solid #666",
        cursor: isDisabled ? "default" : "pointer",
        backgroundColor: getBackgroundColor(),
        userSelect: "none",
        WebkitUserSelect: "none",
        MozUserSelect: "none",
        msUserSelect: "none",
        WebkitTouchCallout: "none",
        WebkitTapHighlightColor: "transparent",
      }}
      onClick={() => onCyclePlacement(award, playerIndex)}
      disabled={isDisabled}
    />
  );
}

function NumericInputWithButtons({
  value,
  onChange,
  onDecrement,
  onIncrement,
}) {
  const buttonStyle = {
    width: "3rem",
    minWidth: "3rem",
    height: "3rem",
    fontSize: "2.5rem",
    fontWeight: "bold",
    border: "1px solid #777",
    borderRadius: "3px",
    background: "inherit",
    cursor: "pointer",
    flexShrink: 0,
  };

  return (
    <div
      style={{
        float: "right",
        display: "flex",
        gap: "5px",
        alignItems: "center",
        width: "50%",
      }}
    >
      <button style={buttonStyle} onClick={onDecrement}>
        −
      </button>
      <input
        type="text"
        inputMode="numeric"
        pattern="[0-9]*"
        style={{
          textAlign: "center",
          fontFamily: "inherit",
          fontSize: "2rem",
          fontWeight: "bold",
          background: "inherit",
          height: "3rem",
          flex: 1,
          minWidth: 0,
          boxSizing: "border-box",
          border: "1px solid #777",
          borderRadius: "3px",
        }}
        value={value}
        onChange={onChange}
      />
      <button style={buttonStyle} onClick={onIncrement}>
        +
      </button>
    </div>
  );
}

function PointInput({
  label,
  players,
  playerScores,
  field,
  onChange,
  readOnly = false,
  placeholder = "0",
}) {
  return (
    <div
      style={{
        ...formStyles.playerInputDiv,
        marginBottom: "10px",
      }}
    >
      <div
        style={{
          ...formStyles.milestoneLabel,
          width: "28%",
        }}
      >
        {label}
      </div>
      <div
        style={{
          display: "flex",
          gap: "5px",
          width: "68%",
          justifyContent: "space-around",
        }}
      >
        {players.map((_, playerIndex) =>
          readOnly ? (
            <div
              key={playerIndex}
              style={{
                textAlign: "center",
                fontFamily: "inherit",
                fontSize: "1.5rem",
                height: "3rem",
                lineHeight: "3rem",
                fontWeight: "bold",
                flex: 1,
                border: "1px solid #999",
                borderRadius: "3px",
                backgroundColor:
                  field === "totalPoints" ? "#f0f0f0" : "inherit",
              }}
            >
              {playerScores[playerIndex]?.[field] || 0}
            </div>
          ) : (
            <input
              key={playerIndex}
              type="text"
              inputMode="numeric"
              pattern="[0-9]*"
              style={{
                textAlign: "center",
                fontFamily: "inherit",
                fontSize: "1.5rem",
                background: "inherit",
                height: "3rem",
                flex: 1,
                minWidth: 0,
                boxSizing: "border-box",
                borderWidth: "1px",
                borderStyle: "solid",
                borderColor: "#999",
                borderRadius: "3px",
              }}
              value={playerScores[playerIndex]?.[field] || ""}
              onChange={(e) => {
                const value = e.target.value;
                if (value === "" || /^\d+$/.test(value)) {
                  onChange(playerIndex, field, value);
                }
              }}
              placeholder={placeholder}
            />
          ),
        )}
      </div>
    </div>
  );
}

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
      style={{
        ...(!showText && {
          width: "auto",
          height: "auto",
          margin: "5px",
          display: "inline-block",
          lineHeight: "normal",
        }),
        ...(showText && {
          ...formStyles.checkboxLabel,
          display: "flex",
          gap: "8px",
          justifyContent: "flex-start",
          marginLeft: "0",
        }),
        cursor: disabled ? "not-allowed" : "pointer",
      }}
    >
      <input
        type="checkbox"
        checked={checked}
        disabled={disabled}
        onChange={onChange}
        style={formStyles.hiddenCheckbox}
      />
      <div
        style={{
          ...formStyles.expansionIconStyle,
          backgroundColor: config.backgroundColor,
          color: config.symbolColor,
          opacity: isShaded ? 0.4 : 1,
          filter: isShaded ? "saturate(0.3)" : "none",
          cursor: disabled ? "not-allowed" : "pointer",
        }}
        title={expansion}
      >
        {config.svg ? (
          <svg
            width="28"
            height="28"
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

// Custom hook for managing game objectives (milestones/awards)
function useGameObjectives(type, map, expansions, playerNumber) {
  const [selected, setSelected] = useState([]);
  const [data, setData] = useState({});

  const isAward = type === "award";
  const dataKey = isAward ? "awards" : "milestones";
  const additionalDataKey = isAward
    ? "additionalAwards"
    : "additionalMilestones";
  const venusSlots = isAward
    ? GAME_CONSTANTS.VENUS_AWARD_SLOTS
    : GAME_CONSTANTS.VENUS_MILESTONE_SLOTS;
  const defaultSlots = isAward
    ? GAME_CONSTANTS.DEFAULT_AWARD_SLOTS
    : GAME_CONSTANTS.DEFAULT_MILESTONE_SLOTS;

  // Get available objectives based on expansions
  const getAvailable = useCallback(() => {
    let available = [...(gameData[dataKey][map] || [])];
    if (expansions["Venus Next"]) {
      available = [...available, ...gameData[dataKey].Venus];
    }
    if (expansions["Milestones & Awards"]) {
      available = [...available, ...gameData[additionalDataKey]];
    }
    return [...new Set(available)];
  }, [map, expansions, dataKey, additionalDataKey]);

  // Get available objectives for dropdown (excluding already selected)
  const getAvailableForDropdown = useCallback(
    (current) => {
      const all = getAvailable();
      return all.filter((item) => item === current || !selected.includes(item));
    },
    [getAvailable, selected],
  );

  // Initialize objectives when map or expansions change
  useEffect(() => {
    if (expansions["Milestones & Awards"]) {
      const mapObjectives = gameData[dataKey][map] || [];
      const numSlots = expansions["Venus Next"] ? venusSlots : defaultSlots;

      const defaultObjectives = [...mapObjectives];
      if (expansions["Venus Next"]) {
        defaultObjectives.push(...gameData[dataKey].Venus);
      }

      setSelected(defaultObjectives.slice(0, numSlots));
    } else {
      const available = getAvailable();
      setSelected(available);

      // Initialize data structure
      if (isAward) {
        const newData = {};
        available.forEach((item) => {
          newData[item] = data[item] || {};
          for (let i = 0; i < playerNumber; i++) {
            newData[item][i] = newData[item][i] || 0;
          }
        });
        setData(newData);
      } else {
        const newData = {};
        available.forEach((item) => {
          newData[item] = data[item] ?? -1;
        });
        setData(newData);
      }
    }
  }, [
    map,
    expansions["Milestones & Awards"],
    expansions["Venus Next"],
    playerNumber,
  ]);

  // Update selected objective at index
  const updateSelected = useCallback(
    (index, newValue) => {
      const newSelected = [...selected];
      const oldValue = newSelected[index];
      newSelected[index] = newValue;
      setSelected(newSelected);

      // Update data structure
      const newData = { ...data };
      if (oldValue) {
        delete newData[oldValue];
      }
      if (newValue) {
        if (isAward) {
          newData[newValue] = {};
          for (let i = 0; i < playerNumber; i++) {
            newData[newValue][i] = 0;
          }
        } else {
          newData[newValue] = -1;
        }
      }
      setData(newData);
    },
    [selected, data, playerNumber, isAward],
  );

  return {
    selected,
    data,
    setData,
    getAvailable,
    getAvailableForDropdown,
    updateSelected,
  };
}

const PlayerInput = React.memo(
  ({ index, player, corporations, onUpdate, selectedCorporations }) => {
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

    return (
      <div style={formStyles.playerInputDiv}>
        <input
          type="text"
          style={formStyles.containerInput}
          placeholder={`Player ${index + 1} name`}
          value={player.name}
          onChange={(e) => onUpdate(index, "name", e.target.value)}
        />
        <select
          style={formStyles.containerInput}
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

function AddGamePage() {
  const navigate = useNavigate();

  // Basic game options
  const [name, setName] = useState("");
  const [date, setDate] = useState("");
  const [map, setMap] = useState("Tharsis");
  const [generations, setGenerations] = useState(
    GAME_CONSTANTS.DEFAULT_GENERATIONS,
  );
  const [playerNumber, setPlayerNumber] = useState(
    GAME_CONSTANTS.DEFAULT_PLAYER_COUNT,
  );
  const [players, setPlayers] = useState([]);

  // Expansions
  const [expansions, setExpansions] = useState({
    "Base Game": true,
    Draft: true,
    "Corporate Era": true,
    Prelude: true,
    "Prelude 2": true,
    "Venus Next": false,
    Colonies: false,
    Turmoil: false,
    "Milestones & Awards": false,
    Promo: false,
  });
  const [expandedExpansions, setExpandedExpansions] = useState(false);

  // Use custom hooks for milestones and awards
  const milestones = useGameObjectives(
    "milestone",
    map,
    expansions,
    playerNumber,
  );
  const awards = useGameObjectives("award", map, expansions, playerNumber);

  // Points
  const [playerScores, setPlayerScores] = useState([]);

  // Set default date on mount
  useEffect(() => {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, "0");
    const day = String(today.getDate()).padStart(2, "0");
    setDate(`${year}-${month}-${day}`);
  }, []);

  // Update players when number changes
  useEffect(() => {
    const newPlayers = [];
    for (let i = 0; i < playerNumber; i++) {
      newPlayers.push({
        name: players[i]?.name || "",
        corporation: players[i]?.corporation || "",
      });
    }
    setPlayers(newPlayers);
  }, [playerNumber]);

  // Get max players from map data
  const maxPlayers =
    gameData.maps[map]?.maxPlayers || GAME_CONSTANTS.DEFAULT_MAX_PLAYERS;

  // Adjust player count when map changes
  useEffect(() => {
    if (playerNumber > maxPlayers) {
      setPlayerNumber(maxPlayers);
    }
  }, [map, maxPlayers]);

  // Update player scores when number changes
  useEffect(() => {
    const newScores = [];
    for (let i = 0; i < playerNumber; i++) {
      newScores.push({
        terraformingRating: playerScores[i]?.terraformingRating || "",
        cities: playerScores[i]?.cities || "",
        greeneries: playerScores[i]?.greeneries || "",
        cards: playerScores[i]?.cards || "",
        turmoilPoints: playerScores[i]?.turmoilPoints || "",
        milestonePoints: 0,
        awardPoints: 0,
        totalPoints: 0,
      });
    }
    setPlayerScores(newScores);
  }, [playerNumber]);

  const updatePlayerData = React.useCallback((index, field, value) => {
    setPlayers((prevPlayers) => {
      const newPlayers = [...prevPlayers];
      newPlayers[index][field] = value;
      return newPlayers;
    });
  }, []);

  const handleExpansionChange = (expansion) => {
    setExpansions({ ...expansions, [expansion]: !expansions[expansion] });
  };

  const getAvailableCorporations = () => {
    let availableCorporations = [];

    // Add corporations from each selected expansion
    Object.entries(expansions).forEach(([expansion, isSelected]) => {
      if (isSelected && gameData.corporationsByExpansion[expansion]) {
        availableCorporations = [
          ...availableCorporations,
          ...gameData.corporationsByExpansion[expansion],
        ];
      }
    });

    // Remove duplicates and sort
    return [...new Set(availableCorporations)].sort();
  };

  const updateMilestoneWinner = (milestone, playerIndex) => {
    // Count how many milestones are currently selected
    const currentlySelected = Object.values(milestones.data).filter(
      (idx) => idx !== -1,
    ).length;

    // If trying to select a new milestone and already have max claimed, don't allow
    if (
      playerIndex !== -1 &&
      milestones.data[milestone] === -1 &&
      currentlySelected >= GAME_CONSTANTS.MAX_MILESTONES_CLAIMED
    ) {
      return;
    }

    milestones.setData({
      ...milestones.data,
      [milestone]: playerIndex,
    });
  };

  const getSelectedMilestonesCount = () => {
    return Object.values(milestones.data).filter((idx) => idx !== -1).length;
  };

  const cyclePlacement = (award, playerIndex) => {
    // Check if this award has any placements, if not and we already have max funded awards, don't allow
    const currentAwardHasPlacements =
      awards.data[award] &&
      Object.values(awards.data[award]).some((placement) => placement > 0);
    const fundedAwardsCount = getFundedAwardsCount();

    if (
      !currentAwardHasPlacements &&
      fundedAwardsCount >= GAME_CONSTANTS.MAX_AWARDS_FUNDED
    ) {
      return; // Can't add new award if already have max funded
    }

    const newPlacements = { ...awards.data };
    if (!newPlacements[award]) {
      newPlacements[award] = {};
    }

    const currentPlacement =
      newPlacements[award][playerIndex] || GAME_CONSTANTS.AWARD_PLACEMENT.NONE;
    // Cycle through placements: none -> gold -> silver -> none
    const placementCount = Object.keys(GAME_CONSTANTS.AWARD_PLACEMENT).length;
    newPlacements[award][playerIndex] = (currentPlacement + 1) % placementCount;

    awards.setData(newPlacements);
  };

  const getFundedAwardsCount = () => {
    return Object.keys(awards.data).filter(
      (award) =>
        awards.data[award] &&
        Object.values(awards.data[award]).some((placement) => placement > 0),
    ).length;
  };

  const isAwardFunded = (award) => {
    return (
      awards.data[award] &&
      Object.values(awards.data[award]).some((placement) => placement > 0)
    );
  };

  const updatePlayerScore = (playerIndex, field, value) => {
    setPlayerScores((currentScores) => {
      const newScores = [...currentScores];
      newScores[playerIndex][field] = value;

      // Recalculate total for this player
      const score = newScores[playerIndex];
      const terraformingRating = parseInt(score.terraformingRating || 0);
      const cities = parseInt(score.cities || 0);
      const greeneries = parseInt(score.greeneries || 0);
      const cards = parseInt(score.cards || 0);
      const turmoilPoints = parseInt(score.turmoilPoints || 0);

      score.totalPoints =
        terraformingRating +
        cities +
        greeneries +
        cards +
        turmoilPoints +
        (score.milestonePoints || 0) +
        (score.awardPoints || 0);

      return newScores;
    });
  };

  // Calculate points whenever milestones or awards change
  useEffect(() => {
    if (playerScores.length === 0) return;

    const newScores = playerScores.map((score, playerIndex) => {
      // Calculate milestone points
      let milestonePoints = 0;
      Object.values(milestones.data).forEach((winnerIndex) => {
        if (winnerIndex === playerIndex) {
          milestonePoints += GAME_CONSTANTS.MILESTONE_POINTS;
        }
      });

      // Calculate award points based on placement
      let awardPoints = 0;
      Object.keys(awards.data).forEach((award) => {
        const placement = awards.data[award]?.[playerIndex];
        if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.GOLD) {
          awardPoints += GAME_CONSTANTS.AWARD_POINTS.GOLD;
        } else if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.SILVER) {
          awardPoints += GAME_CONSTANTS.AWARD_POINTS.SILVER;
        }
      });

      // Calculate total points
      const terraformingRating = parseInt(score.terraformingRating || 0);
      const cities = parseInt(score.cities || 0);
      const greeneries = parseInt(score.greeneries || 0);
      const cards = parseInt(score.cards || 0);
      const turmoilPoints = parseInt(score.turmoilPoints || 0);

      const totalPoints =
        terraformingRating +
        cities +
        greeneries +
        cards +
        turmoilPoints +
        milestonePoints +
        awardPoints;

      return {
        ...score,
        milestonePoints,
        awardPoints,
        totalPoints,
      };
    });

    setPlayerScores(newScores);
  }, [milestones.data, awards.data]);

  return (
    <Layout>
      <Container title="Add game" />

      <Container title="Options" titleStyle="banner">
        <SubContainer>
          <SubContainerElement>
            <label>Name:</label>
            <input
              type="text"
              style={formStyles.optionInput}
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </SubContainerElement>

          <SubContainerElement>
            <label>Date:</label>
            <input
              type="date"
              style={formStyles.optionInput}
              value={date}
              onChange={(e) => setDate(e.target.value)}
              required={true}
            />
          </SubContainerElement>

          <SubContainerElement>
            <label>Map:</label>
            <select
              style={formStyles.optionInput}
              value={map}
              onChange={(e) => setMap(e.target.value)}
            >
              {Object.keys(gameData.maps).map((m) => (
                <option key={m} value={m}>
                  {m}
                </option>
              ))}
            </select>
          </SubContainerElement>

          <SubContainerElement>
            <label>Generations:</label>
            <NumericInputWithButtons
              value={generations}
              onChange={(e) => {
                const val = parseInt(e.target.value) || 0;
                setGenerations(
                  Math.min(
                    GAME_CONSTANTS.MAX_GENERATIONS,
                    Math.max(GAME_CONSTANTS.MIN_GENERATIONS, val),
                  ),
                );
              }}
              onDecrement={() =>
                setGenerations(
                  Math.max(GAME_CONSTANTS.MIN_GENERATIONS, generations - 1),
                )
              }
              onIncrement={() =>
                setGenerations(
                  Math.min(GAME_CONSTANTS.MAX_GENERATIONS, generations + 1),
                )
              }
            />
          </SubContainerElement>

          <SubContainerElement>
            <label>Players:</label>
            <NumericInputWithButtons
              value={playerNumber}
              onChange={(e) =>
                setPlayerNumber(
                  Math.min(
                    maxPlayers,
                    Math.max(
                      GAME_CONSTANTS.MIN_PLAYERS,
                      parseInt(e.target.value) || GAME_CONSTANTS.MIN_PLAYERS,
                    ),
                  ),
                )
              }
              onDecrement={() =>
                setPlayerNumber(
                  Math.max(GAME_CONSTANTS.MIN_PLAYERS, playerNumber - 1),
                )
              }
              onIncrement={() =>
                setPlayerNumber(Math.min(maxPlayers, playerNumber + 1))
              }
            />
          </SubContainerElement>

          <SubContainerElement>
            <div style={formStyles.subcontainerBox}>
              <div
                style={{
                  display: "flex",
                  alignItems: "flex-start",
                  padding: "15px",
                  gap: "15px",
                }}
              >
                {/* Content container - shared by both bar and list */}
                <div style={{ flex: 1 }}>
                  {!expandedExpansions && (
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "space-between",
                      }}
                    >
                      {Object.entries(expansions).map(([key, value]) => (
                        <ExpansionIcon
                          key={key}
                          expansion={key}
                          checked={value}
                          disabled={key === "Base Game"}
                          onChange={() => handleExpansionChange(key)}
                          showText={false}
                        >
                          {key}
                        </ExpansionIcon>
                      ))}
                    </div>
                  )}

                  {expandedExpansions && (
                    <div
                      style={{
                        display: "flex",
                        flexDirection: "column",
                        alignItems: "flex-start",
                        margin: "5px",
                      }}
                    >
                      {Object.entries(expansions).map(([key, value]) => (
                        <ExpansionIcon
                          key={key + "_expanded"}
                          expansion={key}
                          checked={value}
                          disabled={key === "Base Game"}
                          onChange={() => handleExpansionChange(key)}
                          showText={true}
                        >
                          {key}
                        </ExpansionIcon>
                      ))}
                    </div>
                  )}
                </div>

                {/* Expand/Collapse button container */}
                <div
                  style={{
                    display: "flex",
                    alignItems: "center",
                    margin: "5px",
                  }}
                >
                  <div
                    style={{
                      ...formStyles.expansionIconStyle,
                      backgroundColor: "#333",
                      color: "#FFF",
                      userSelect: "none",
                      cursor: "pointer",
                    }}
                    onClick={() => setExpandedExpansions(!expandedExpansions)}
                    title={expandedExpansions ? "Collapse" : "Expand"}
                  >
                    ☰
                  </div>
                </div>
              </div>
            </div>
          </SubContainerElement>
        </SubContainer>
      </Container>

      <Container title="Players" titleStyle="banner">
        <SubContainer>
          {players.map((player, index) => (
            <PlayerInput
              key={index}
              index={index}
              player={player}
              corporations={getAvailableCorporations()}
              onUpdate={updatePlayerData}
              selectedCorporations={players
                .map((p) => p.corporation)
                .filter((corp) => corp !== "")}
            />
          ))}
        </SubContainer>
      </Container>

      <Container title="Milestones" titleStyle="banner">
        <SubContainer>
          {milestones.selected.map((milestone, index) => (
            <div
              key={index}
              style={{
                ...formStyles.playerInputDiv,
              }}
            >
              {expansions["Milestones & Awards"] ? (
                <select
                  style={{
                    ...formStyles.containerInput,
                  }}
                  value={milestone || ""}
                  onChange={(e) =>
                    milestones.updateSelected(index, e.target.value)
                  }
                >
                  {!milestone && <option value="">Select Milestone</option>}
                  {milestones.getAvailableForDropdown(milestone).map((m) => (
                    <option key={m} value={m}>
                      {m}
                    </option>
                  ))}
                </select>
              ) : (
                <div
                  style={{
                    ...formStyles.milestoneLabel,
                  }}
                >
                  {milestone}
                </div>
              )}

              <select
                style={{
                  ...formStyles.containerInput,
                }}
                value={milestones.data[milestone] ?? -1}
                onChange={(e) =>
                  updateMilestoneWinner(milestone, parseInt(e.target.value))
                }
                disabled={
                  !milestone ||
                  (milestones.data[milestone] === -1 &&
                    getSelectedMilestonesCount() >=
                      GAME_CONSTANTS.MAX_MILESTONES_CLAIMED)
                }
              >
                <option value={-1}>Not achieved</option>
                {players.map((p, i) => (
                  <option key={i} value={i}>
                    {p.name || `Player ${i + 1}`}
                  </option>
                ))}
              </select>
            </div>
          ))}
        </SubContainer>
      </Container>

      <Container title="Awards" titleStyle="banner">
        <SubContainer>
          {/* Player names header */}
          <div
            style={{
              ...formStyles.playerInputDiv,
              marginBottom: "15px",
              borderBottom: "2px solid #999",
              paddingBottom: "10px",
            }}
          >
            <div style={{ width: "28%" }}>
              {/* Empty space for alignment */}
            </div>
            <div
              style={{
                display: "flex",
                gap: "5px",
                width: "68%",
                justifyContent: "space-around",
              }}
            >
              {players.map((player, playerIndex) => (
                <div
                  key={playerIndex}
                  style={{
                    fontSize: "1.2rem",
                    fontWeight: "bold",
                    textAlign: "center",
                    flex: 1,
                  }}
                >
                  {player.name || `P${playerIndex + 1}`}
                </div>
              ))}
            </div>
          </div>

          {awards.selected.map((award, index) => (
            <div
              key={index}
              style={{
                ...formStyles.playerInputDiv,
                marginBottom: "10px",
              }}
            >
              {expansions["Milestones & Awards"] ? (
                <select
                  style={{
                    ...formStyles.containerInput,
                    width: "28%",
                  }}
                  value={award || ""}
                  onChange={(e) => awards.updateSelected(index, e.target.value)}
                >
                  {!award && <option value="">Select Award</option>}
                  {awards.getAvailableForDropdown(award).map((a) => (
                    <option key={a} value={a}>
                      {a}
                    </option>
                  ))}
                </select>
              ) : (
                <div
                  style={{
                    ...formStyles.milestoneLabel,
                    width: "28%",
                  }}
                >
                  {award}
                </div>
              )}

              <div
                style={{
                  display: "flex",
                  gap: "5px",
                  width: "68%",
                  justifyContent: "space-around",
                }}
              >
                {players.map((player, playerIndex) => (
                  <AwardButton
                    key={playerIndex}
                    award={award}
                    playerIndex={playerIndex}
                    awardPlacements={awards.data}
                    onCyclePlacement={cyclePlacement}
                    isAwardFunded={isAwardFunded}
                    getFundedAwardsCount={getFundedAwardsCount}
                  />
                ))}
              </div>
            </div>
          ))}
        </SubContainer>
      </Container>

      <Container title="Points" titleStyle="banner">
        <SubContainer>
          {/* Player names header */}
          <div
            style={{
              ...formStyles.playerInputDiv,
              marginBottom: "15px",
              borderBottom: "2px solid #999",
              paddingBottom: "10px",
            }}
          >
            <div style={{ width: "28%" }}>
              {/* Empty space for alignment */}
            </div>
            <div
              style={{
                display: "flex",
                gap: "5px",
                width: "68%",
                justifyContent: "space-around",
              }}
            >
              {players.map((player, playerIndex) => (
                <div
                  key={playerIndex}
                  style={{
                    fontSize: "1.2rem",
                    fontWeight: "bold",
                    textAlign: "center",
                    flex: 1,
                  }}
                >
                  {player.name || `P${playerIndex + 1}`}
                </div>
              ))}
            </div>
          </div>

          <PointInput
            label="TR"
            players={players}
            playerScores={playerScores}
            field="terraformingRating"
            onChange={updatePlayerScore}
            placeholder={GAME_CONSTANTS.DEFAULT_TR.toString()}
          />
          <PointInput
            label="Cities"
            players={players}
            playerScores={playerScores}
            field="cities"
            onChange={updatePlayerScore}
          />
          <PointInput
            label="Greeneries"
            players={players}
            playerScores={playerScores}
            field="greeneries"
            onChange={updatePlayerScore}
          />
          <PointInput
            label="Cards"
            players={players}
            playerScores={playerScores}
            field="cards"
            onChange={updatePlayerScore}
          />
          {expansions.Turmoil && (
            <PointInput
              label="Turmoil"
              players={players}
              playerScores={playerScores}
              field="turmoilPoints"
              onChange={updatePlayerScore}
            />
          )}

          <div
            style={{
              borderTop: "2px solid #999",
              margin: "15px 0",
              paddingTop: "15px",
            }}
          >
            <PointInput
              label="Milestones"
              players={players}
              playerScores={playerScores}
              field="milestonePoints"
              readOnly={true}
            />
            <PointInput
              label="Awards"
              players={players}
              playerScores={playerScores}
              field="awardPoints"
              readOnly={true}
            />
          </div>

          <div
            style={{
              borderTop: "3px solid #000",
              marginTop: "15px",
              paddingTop: "15px",
            }}
          >
            <PointInput
              label="Total"
              players={players}
              playerScores={playerScores}
              field="totalPoints"
              readOnly={true}
            />
          </div>
        </SubContainer>
      </Container>

      <LinkButton onClick={() => navigate("/")}>Main page</LinkButton>
    </Layout>
  );
}

export default AddGamePage;
