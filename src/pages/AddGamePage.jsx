import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Layout from "../components/Layout";
import Container from "../components/Container";
import { SubContainer, SubContainerElement } from "../components/SubContainer";
import LinkButton from "../components/LinkButton";
import { formStyles } from "../styles/formStyles";

// Local components for this page only
function AwardButton({ award, playerIndex, awardPlacements, onCyclePlacement, isAwardFunded, getFundedAwardsCount }) {
  const isDisabled = !award || (!isAwardFunded(award) && getFundedAwardsCount() >= 3);
  const placement = awardPlacements[award]?.[playerIndex] || 0;

  const getBackgroundColor = () => {
    if (isDisabled) return "#666666";
    if (placement === 1) return "#FFD700"; // Gold
    if (placement === 2) return "#CCC";    // Silver
    return "#444";                         // None/Black
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

function NumericInputWithButtons({ value, onChange, onDecrement, onIncrement }) {
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
    <div style={{ float: "right", display: "flex", gap: "5px", alignItems: "center", width: "50%" }}>
      <button style={buttonStyle} onClick={onDecrement}>
        −
      </button>
      <input
        type="text"
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

function PointsRow({
  label,
  players,
  playerScores,
  field,
  onChange,
  readOnly = false,
  placeholder = "0",
}) {
  return (
    <tr>
      <td style={{ padding: "5px", fontWeight: "bold" }}>{label}</td>
      {players.map((_, playerIndex) => (
        <td key={playerIndex} style={{ padding: "5px" }}>
          {readOnly ? (
            <div
              style={{
                textAlign: "center",
                fontFamily: "inherit",
                fontSize: "13px",
                height: "30px",
                lineHeight: "30px",
                fontWeight: "bold",
              }}
            >
              {playerScores[playerIndex]?.[field] || 0}
            </div>
          ) : (
            <input
              type="text"
              style={{
                textAlign: "center",
                fontFamily: "inherit",
                fontSize: "13px",
                background: "inherit",
                height: "30px",
                width: "100%",
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
          )}
        </td>
      ))}
    </tr>
  );
}

function ExpansionIcon({ expansion, checked, disabled, onChange, children, showText = true }) {
  const getExpansionConfig = (expansion) => {
    switch (expansion) {
      case "Base Game":
        return {
          backgroundColor: "transparent",
          symbolColor: "#000",
          svg: null
        };
      case "Corporate Era":
        return {
          backgroundColor: "#FF0000",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "0deg"
        };
      case "Venus Next":
        return {
          backgroundColor: "#87CEEB",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "180deg"
        };
      case "Prelude":
        return {
          backgroundColor: "#FFC0CB",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "270deg"
        };
      case "Prelude 2":
        return {
          backgroundColor: "#FF69B4",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "270deg"
        };
      case "Colonies":
        return {
          backgroundColor: "#808080",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "0deg"
        };
      case "Turmoil":
        return {
          backgroundColor: "#FFA500",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "180deg"
        };
      case "Milestones & Awards":
        return {
          backgroundColor: "#FFD700",
          symbolColor: "#000",
          svg: <rect x="6" y="6" width="12" height="12" fill="currentColor" />
        };
      case "Promo":
        return {
          backgroundColor: "#2F2F2F",
          symbolColor: "#FFF",
          svg: <circle cx="12" cy="12" r="8" fill="currentColor" />
        };
      default:
        return {
          backgroundColor: "#666",
          symbolColor: "#FFF",
          svg: null,
          fallbackSymbol: "?"
        };
    }
  };

  const config = getExpansionConfig(expansion);
  const isShaded = !checked;

  return (
    <label style={{
      ...(!showText && {
        width: "auto",
        height: "auto",
        margin: "5px",
        display: "inline-block",
        lineHeight: "normal"
      }),
      ...(showText && {
        ...formStyles.checkboxLabel,
        display: "flex",
        gap: "8px",
        justifyContent: "flex-start",
        marginLeft: "0",
      }),
      cursor: disabled ? "not-allowed" : "pointer",
    }}>
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
            style={{ transform: config.rotation ? `rotate(${config.rotation})` : undefined }}
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

const PlayerInput = React.memo(({ index, player, corporations, onUpdate, selectedCorporations }) => {
  const availableCorporations = React.useMemo(() =>
    corporations.filter(corp =>
      corp === "Beginner" || !selectedCorporations.includes(corp) || corp === player.corporation
    ), [corporations, selectedCorporations, player.corporation]
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
});

const gameData = {
  maps: [
    "Tharsis",
    "Hellas",
    "Elysium",
    "Utopia Planitia",
    "Terra Cimmeria",
    "Vastitas Borealis",
    "Amazonis Planitia",
  ],
  corporationsByExpansion: {
    "Base Game": [
      "Beginner",
      "Credicor",
      "Ecoline",
      "Helion",
      "Mining Guild",
      "Interplanetarty Cinematics",
      "Inventrix",
      "Phobolog",
      "Tharsis Republic",
      "Thorgate",
      "United Nations Mars Initiative",
    ],
    "Corporate Era": [
      "Teractor",
      "Saturn Systems",
    ],
    "Venus Next": [
      "Aphrodite",
      "Celestic",
      "Manutech",
      "Morning Star Inc.",
      "Viron"
    ],
    "Prelude": [
      "Cheung Shing Mars",
      "Point Luna",
      "Robinson Industries",
      "Valley Trust",
      "Vitor",
      "Ecotec",
    ],
    "Prelude 2": [
      "Spire",
      "Sagitta",
      "Palladin Shipping",
      "Nirgal Enterprises"
    ],
    "Colonies": [
      "Aridor",
      "Arklight",
      "Polyphemos",
      "Poseidon",
      "Stormcraft Incorporated"
    ],
    "Turmoil": [
      "Lakefront Resorts",
      "Pristar",
      "Septem Tribus",
      "Terralabs Research",
      "Utopia Invest"
    ],
    "Promo": [
      "Factorum",
      "Mons Insurance",
      "Philares",
      "Arcadian Communities",
      "Recyclon",
      "Splice Tactical Genomics",
      "Astrodrill",
      "Pharmacy Union",
      "Tycho Magnetics",
      "Kuiper Cooperative",
    ],
  },
  milestones: {
    'Amazonis Planitia': ["Terran", "Landshaper", "Merchant", "Sponsor", "Lobbyist"],
    'Elysium': ["Generalist", "Specialist", "Ecologist", "Tycoon", "Legend"],
    'Hellas': [
      "Diversifier",
      "Tactician",
      "Polar Explorer",
      "Energizer",
      "Rim Settler",
    ],
    'Terra Cimmeria': ["Planetologist", "Architect", "Coastguard", "Forester", "Investor"],
    'Tharsis': ["Builder", "Gardener", "Mayor", "Planner", "Terraformer"],
    'Utopia Planitia': ["Specialist", "Pioneer", "Trader", "Metallurgist", "Researcher"],
    'Vastitas Borealis': ["Agronomist", "Spacefarer", "Geologist", "Engineer", "Farmer"],
    'Venus': ["Hoverlord"],
  },
  additionalMilestones: [
    "Briber",
    "Builder",
    "Coastguard",
    "Diversifier",
    "Ecologist",
    "Energizer",
    "Engineer",
    "Farmer",
    "Forester",
    "Fundraiser",
    "Gardener",
    "Generalist",
    "Geologist",
    "Hydrologist",
    "Landshaper",
    "Legend",
    "Lobbyist",
    "Mayor",
    "Merchant",
    "Metallurgist",
    "Philantropist",
    "Pioneer",
    "Planetologist",
    "Planner",
    "Producer",
    "Researcher",
    "Rim Settler",
    "Spacefarer",
    "Sponsor",
    "Tactician",
    "Terraformer",
    "Terran",
    "Thawer",
    "Trader",
    "Tycoon",
  ],
  awards: {
    'Amazonis Planitia': ["Collector", "Innovator", "Constructor", "Manufacturer", "Physicist"],
    'Elysium': [
      "Celebrity",
      "Industrialist",
      "Desert Settler",
      "Estate Dealer",
      "Benefactor",
    ],
    'Hellas': ["Cultivator", "Magnate", "Space Baron", "Excentric", "Contractor"],
    'Terra Cimmeria': ["Electrician", "Founder", "Mogul", "Zoologist", "Forecaster"],
    'Tharsis': ["Landlord", "Scientist", "Banker", "Thermalist", "Miner"],
    'Utopia Planitia': ["Suburbian", "Investor", "Botanist", "Incorporator", "Metropolist"],
    'Vastitas Borealis': ["Traveller", "Landscaper", "Highlander", "Promoter", "Blacksmith"],
    'Venus': ["Venuphile"],
  },
  additionalAwards: [
    "Administrator",
    "Banker",
    "Benefactor",
    "Biologist",
    "Botanist",
    "Celebrity",
    "Collector",
    "Constructor",
    "Contractor",
    "Cultivator",
    "Electrician",
    "Estate Dealer",
    "Excentric",
    "Forecaster",
    "Founder",
    "Highlander",
    "Industrialist",
    "Incorporator",
    "Investor",
    "Landlord",
    "Landscaper",
    "Magnate",
    "Manufacturer",
    "Metropolist",
    "Miner",
    "Mogul",
    "Politician",
    "Promoter",
    "Scientist",
    "Space Baron",
    "Suburbian",
    "Thermalist",
    "Traveller",
    "Visionary",
    "Zoologist",
  ],
};

function AddGamePage() {
  const navigate = useNavigate();

  // Basic game options
  const [name, setName] = useState("");
  const [date, setDate] = useState("");
  const [map, setMap] = useState("Tharsis");
  const [generations, setGenerations] = useState(1);
  const [playerNumber, setPlayerNumber] = useState(4);
  const [players, setPlayers] = useState([]);

  // Expansions
  const [expansions, setExpansions] = useState({
    "Base Game": true,
    "Draft": true,
    "Corporate Era": true,
    "Prelude": true,
    "Prelude 2": true,
    "Venus Next": false,
    "Colonies": false,
    "Turmoil": false,
    "Milestones & Awards": false,
    "Promo": false,
  });
  const [expandedExpansions, setExpandedExpansions] = useState(false);

  // Milestones - object with milestone name as key, player index as value (-1 means not claimed)
  const [milestoneWinners, setMilestoneWinners] = useState({});
  // When Milestones & Awards expansion is enabled, track which milestones are selected
  const [selectedMilestones, setSelectedMilestones] = useState([]);

  // Awards - object with award name as key, and player placements (0=none, 1=gold, 2=silver)
  const [awardPlacements, setAwardPlacements] = useState({});
  // When Milestones & Awards expansion is enabled, track which awards are selected
  const [selectedAwards, setSelectedAwards] = useState([]);


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

  // Calculate max players based on map
  const maxPlayers = (map === "Amazonis Planitia" || map === "Vastitas Borealis") ? 6 : 5;

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

  // Initialize milestones when map or expansions change
  useEffect(() => {
    if (expansions["Milestones & Awards"]) {
      // When Milestones & Awards is enabled, show 5 slots (or 6 with Venus Next)
      const mapMilestones = gameData.milestones[map] || [];
      const numSlots = expansions["Venus Next"] ? 6 : 5;

      // Initialize with the first milestones from the map
      const defaultMilestones = [...mapMilestones];
      if (expansions["Venus Next"]) {
        defaultMilestones.push(...gameData.milestones.Venus);
      }

      setSelectedMilestones(defaultMilestones.slice(0, numSlots));
    } else {
      // When Milestones & Awards is disabled, use standard milestones for the map
      const availableMilestones = getAvailableMilestones();
      const newMilestoneWinners = {};
      availableMilestones.forEach(milestone => {
        newMilestoneWinners[milestone] = milestoneWinners[milestone] ?? -1;
      });
      setMilestoneWinners(newMilestoneWinners);
      setSelectedMilestones(availableMilestones);
    }
  }, [map, expansions["Milestones & Awards"], expansions["Venus Next"]]);

  // Initialize awards when map or expansions change
  useEffect(() => {
    if (expansions["Milestones & Awards"]) {
      // When Milestones & Awards is enabled, show 5 slots (or 6 with Venus Next)
      const mapAwards = gameData.awards[map] || [];
      const numSlots = expansions["Venus Next"] ? 6 : 5;

      // Initialize with the first awards from the map
      const defaultAwards = [...mapAwards];
      if (expansions["Venus Next"]) {
        defaultAwards.push(...gameData.awards.Venus);
      }

      setSelectedAwards(defaultAwards.slice(0, numSlots));
    } else {
      // When Milestones & Awards is disabled, use standard awards for the map
      const availableAwards = getAvailableAwards();
      const newAwardPlacements = {};
      availableAwards.forEach(award => {
        newAwardPlacements[award] = awardPlacements[award] || {};
        // Initialize all players to 0 (no medal) for this award
        for (let i = 0; i < playerNumber; i++) {
          newAwardPlacements[award][i] = newAwardPlacements[award][i] || 0;
        }
      });
      setAwardPlacements(newAwardPlacements);
      setSelectedAwards(availableAwards);
    }
  }, [map, expansions["Milestones & Awards"], expansions["Venus Next"]]);

  const updatePlayerData = React.useCallback((index, field, value) => {
    setPlayers(prevPlayers => {
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
          ...gameData.corporationsByExpansion[expansion]
        ];
      }
    });

    // Remove duplicates and sort
    return [...new Set(availableCorporations)].sort();
  };

  const getAvailableMilestones = () => {
    let availableMilestones = [...(gameData.milestones[map] || [])];
    if (expansions["Venus Next"]) {
      availableMilestones = [
        ...availableMilestones,
        ...gameData.milestones.Venus,
      ];
    }
    if (expansions["Milestones & Awards"]) {
      // When Milestones & Awards expansion is enabled, add additional milestones to choose from
      availableMilestones = [
        ...availableMilestones,
        ...gameData.additionalMilestones,
      ];
    }
    // Remove duplicates
    return [...new Set(availableMilestones)];
  };

  const updateMilestoneWinner = (milestone, playerIndex) => {
    // Count how many milestones are currently selected
    const currentlySelected = Object.values(milestoneWinners).filter(idx => idx !== -1).length;

    // If trying to select a new milestone and already have 3, don't allow
    if (playerIndex !== -1 && milestoneWinners[milestone] === -1 && currentlySelected >= 3) {
      return;
    }

    setMilestoneWinners({
      ...milestoneWinners,
      [milestone]: playerIndex
    });
  };

  const getSelectedMilestonesCount = () => {
    return Object.values(milestoneWinners).filter(idx => idx !== -1).length;
  };

  const updateSelectedMilestone = (index, newMilestone) => {
    // When using Milestones & Awards expansion, allow changing which milestones are used
    const newSelectedMilestones = [...selectedMilestones];
    const oldMilestone = newSelectedMilestones[index];
    newSelectedMilestones[index] = newMilestone;
    setSelectedMilestones(newSelectedMilestones);

    // Clear the winner for this milestone when it changes
    const newMilestoneWinners = { ...milestoneWinners };
    if (oldMilestone) {
      delete newMilestoneWinners[oldMilestone];
    }
    if (newMilestone) {
      newMilestoneWinners[newMilestone] = -1;
    }
    setMilestoneWinners(newMilestoneWinners);
  };

  const getAvailableMilestonesForDropdown = (currentMilestone) => {
    // Get all available milestones
    const allMilestones = getAvailableMilestones();
    // Filter out milestones that are already selected (except the current one)
    return allMilestones.filter(m =>
      m === currentMilestone || !selectedMilestones.includes(m)
    );
  };

  const getAvailableAwards = () => {
    let availableAwards = [...(gameData.awards[map] || [])];
    if (expansions["Venus Next"]) {
      availableAwards = [...availableAwards, ...gameData.awards.Venus];
    }
    if (expansions["Milestones & Awards"]) {
      // When Milestones & Awards expansion is enabled, add additional awards to choose from
      availableAwards = [
        ...availableAwards,
        ...gameData.additionalAwards,
      ];
    }
    // Remove duplicates
    return [...new Set(availableAwards)];
  };

  const cyclePlacement = (award, playerIndex) => {
    // Check if this award has any placements, if not and we already have 3 funded awards, don't allow
    const currentAwardHasPlacements = awardPlacements[award] && Object.values(awardPlacements[award]).some(placement => placement > 0);
    const fundedAwardsCount = getFundedAwardsCount();

    if (!currentAwardHasPlacements && fundedAwardsCount >= 3) {
      return; // Can't add new award if already have 3 funded
    }

    const newPlacements = { ...awardPlacements };
    if (!newPlacements[award]) {
      newPlacements[award] = {};
    }

    const currentPlacement = newPlacements[award][playerIndex] || 0;
    // Cycle: 0 (none) -> 1 (gold) -> 2 (silver) -> 0 (none)
    newPlacements[award][playerIndex] = (currentPlacement + 1) % 3;

    setAwardPlacements(newPlacements);
  };

  const getFundedAwardsCount = () => {
    return Object.keys(awardPlacements).filter(award =>
      awardPlacements[award] && Object.values(awardPlacements[award]).some(placement => placement > 0)
    ).length;
  };

  const isAwardFunded = (award) => {
    return awardPlacements[award] && Object.values(awardPlacements[award]).some(placement => placement > 0);
  };

  const updateSelectedAward = (index, newAward) => {
    // When using Milestones & Awards expansion, allow changing which awards are used
    const newSelectedAwards = [...selectedAwards];
    const oldAward = newSelectedAwards[index];
    newSelectedAwards[index] = newAward;
    setSelectedAwards(newSelectedAwards);

    // Clear the placements for this award when it changes
    const newAwardPlacements = { ...awardPlacements };
    if (oldAward) {
      delete newAwardPlacements[oldAward];
    }
    if (newAward) {
      newAwardPlacements[newAward] = {};
      // Initialize all players to 0 (no medal) for this award
      for (let i = 0; i < playerNumber; i++) {
        newAwardPlacements[newAward][i] = 0;
      }
    }
    setAwardPlacements(newAwardPlacements);
  };

  const getAvailableAwardsForDropdown = (currentAward) => {
    // Get all available awards
    const allAwards = getAvailableAwards();
    // Filter out awards that are already selected (except the current one)
    return allAwards.filter(a =>
      a === currentAward || !selectedAwards.includes(a)
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
      Object.values(milestoneWinners).forEach((winnerIndex) => {
        if (winnerIndex === playerIndex) {
          milestonePoints += 5;
        }
      });

      // Calculate award points - gold = 5 points, silver = 2 points
      let awardPoints = 0;
      Object.keys(awardPlacements).forEach((award) => {
        const placement = awardPlacements[award]?.[playerIndex];
        if (placement === 1) { // Gold
          awardPoints += 5;
        } else if (placement === 2) { // Silver
          awardPoints += 2;
        }
      });

      // Calculate total points
      const terraformingRating = parseInt(score.terraformingRating || 0);
      const cities = parseInt(score.cities || 0) * 2;
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
  }, [milestoneWinners, awardPlacements]);

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
              {gameData.maps.map((m) => (
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
                setGenerations(Math.min(20, Math.max(0, val)));
              }}
              onDecrement={() => setGenerations(Math.max(0, generations - 1))}
              onIncrement={() => setGenerations(Math.min(20, generations + 1))}
            />
          </SubContainerElement>

          <SubContainerElement>
            <label>Players:</label>
            <NumericInputWithButtons
              value={playerNumber}
              onChange={(e) =>
                setPlayerNumber(
                  Math.min(maxPlayers, Math.max(1, parseInt(e.target.value) || 1)),
                )
              }
              onDecrement={() => setPlayerNumber(Math.max(1, playerNumber - 1))}
              onIncrement={() => setPlayerNumber(Math.min(maxPlayers, playerNumber + 1))}
            />
          </SubContainerElement>

          <SubContainerElement>
            <div style={formStyles.subcontainerBox}>
              <div style={{
                display: "flex",
                alignItems: "flex-start",
                padding: "15px",
                gap: "15px"
              }}>
                {/* Content container - shared by both bar and list */}
                <div style={{ flex: 1 }}>
                  {!expandedExpansions && (
                    <div style={{
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "space-between"
                    }}>
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
                    <div style={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "flex-start",
                      margin: "5px"
                    }}>
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
                <div style={{
                  display: "flex",
                  alignItems: "center",
                  margin: "5px"
                }}>
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
              selectedCorporations={players.map(p => p.corporation).filter(corp => corp !== "")}
            />
          ))}
        </SubContainer>
      </Container>

      <Container title="Milestones" titleStyle="banner">
        <SubContainer>
          {(expansions["Milestones & Awards"] ? selectedMilestones : getAvailableMilestones()).map((milestone, index) => (
            <div
              key={index}
              style={{
                ...formStyles.playerInputDiv
              }}
            >
              {expansions["Milestones & Awards"] ? (
                <select
                  style={{
                    ...formStyles.containerInput
                  }}
                  value={milestone || ""}
                  onChange={(e) => updateSelectedMilestone(index, e.target.value)}
                >
                  {!milestone && <option value="">Select Milestone</option>}
                  {getAvailableMilestonesForDropdown(milestone).map(m => (
                    <option key={m} value={m}>{m}</option>
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
                  ...formStyles.containerInput
                }}
                value={milestoneWinners[milestone] ?? -1}
                onChange={(e) =>
                  updateMilestoneWinner(milestone, parseInt(e.target.value))
                }
                disabled={
                  !milestone ||
                  (milestoneWinners[milestone] === -1 &&
                    getSelectedMilestonesCount() >= 3)
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
              paddingBottom: "10px"
            }}
          >
            <div style={{ width: "28%" }}>
              {/* Empty space for alignment */}
            </div>
            <div style={{
              display: "flex",
              gap: "5px",
              width: "68%",
              justifyContent: "space-around"
            }}>
              {players.map((player, playerIndex) => (
                <div
                  key={playerIndex}
                  style={{
                    fontSize: "1.2rem",
                    fontWeight: "bold",
                    textAlign: "center",
                    flex: 1
                  }}
                >
                  {player.name || `P${playerIndex + 1}`}
                </div>
              ))}
            </div>
          </div>

          {(expansions["Milestones & Awards"] ? selectedAwards : getAvailableAwards()).map((award, index) => (
            <div
              key={index}
              style={{
                ...formStyles.playerInputDiv,
                marginBottom: "10px"
              }}
            >
              {expansions["Milestones & Awards"] ? (
                <select
                  style={{
                    ...formStyles.containerInput,
                    width: "28%"
                  }}
                  value={award || ""}
                  onChange={(e) => updateSelectedAward(index, e.target.value)}
                >
                  {!award && <option value="">Select Award</option>}
                  {getAvailableAwardsForDropdown(award).map(a => (
                    <option key={a} value={a}>{a}</option>
                  ))}
                </select>
              ) : (
                <div
                  style={{
                    ...formStyles.milestoneLabel,
                    width: "28%"
                  }}
                >
                  {award}
                </div>
              )}

              <div style={{
                display: "flex",
                gap: "5px",
                width: "68%",
                justifyContent: "space-around"
              }}>
                {players.map((player, playerIndex) => (
                  <AwardButton
                    key={playerIndex}
                    award={award}
                    playerIndex={playerIndex}
                    awardPlacements={awardPlacements}
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
          <SubContainerElement>
            <div style={formStyles.subcontainerBox}>
              <table
                style={{
                  width: "100%",
                  borderCollapse: "collapse",
                  marginTop: "10px",
                  marginBottom: "10px",
                }}
              >
                <thead>
                  <tr>
                    <th
                      style={{
                        padding: "5px",
                        fontSize: "14px",
                        fontWeight: "bold",
                        borderBottom: "2px solid #999",
                        textAlign: "left",
                      }}
                    >
                      Player
                    </th>
                    {players.map((player, index) => (
                      <th
                        key={index}
                        style={{
                          padding: "5px",
                          fontSize: "14px",
                          fontWeight: "bold",
                          borderBottom: "2px solid #999",
                          textAlign: "center",
                          fontFamily: "Courier New, Courier, monospace",
                        }}
                      >
                        {player.name || `Player ${index + 1}`}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  <PointsRow
                    label="TR"
                    players={players}
                    playerScores={playerScores}
                    field="terraformingRating"
                    onChange={updatePlayerScore}
                    placeholder="20"
                  />
                  <PointsRow
                    label="Cities"
                    players={players}
                    playerScores={playerScores}
                    field="cities"
                    onChange={updatePlayerScore}
                  />
                  <PointsRow
                    label="Greeneries"
                    players={players}
                    playerScores={playerScores}
                    field="greeneries"
                    onChange={updatePlayerScore}
                  />
                  <PointsRow
                    label="Cards"
                    players={players}
                    playerScores={playerScores}
                    field="cards"
                    onChange={updatePlayerScore}
                  />
                  {expansions.Turmoil && (
                    <PointsRow
                      label="Turmoil"
                      players={players}
                      playerScores={playerScores}
                      field="turmoilPoints"
                      onChange={updatePlayerScore}
                    />
                  )}
                  <tr style={{ borderTop: "2px solid #999" }}>
                    <td
                      colSpan={players.length + 1}
                      style={{ padding: "0px", height: "5px" }}
                    ></td>
                  </tr>
                  <PointsRow
                    label="Milestones"
                    players={players}
                    playerScores={playerScores}
                    field="milestonePoints"
                    readOnly={true}
                  />
                  <PointsRow
                    label="Awards"
                    players={players}
                    playerScores={playerScores}
                    field="awardPoints"
                    readOnly={true}
                  />
                  <tr style={{ borderTop: "3px solid #000" }}>
                    <td
                      colSpan={players.length + 1}
                      style={{ padding: "0px", height: "5px" }}
                    ></td>
                  </tr>
                  <PointsRow
                    label="Total"
                    players={players}
                    playerScores={playerScores}
                    field="totalPoints"
                    readOnly={true}
                  />
                </tbody>
              </table>
            </div>
          </SubContainerElement>
        </SubContainer>
      </Container>

      <LinkButton onClick={() => navigate("/")}>Main page</LinkButton>
    </Layout>
  );
}

export default AddGamePage;
