import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Layout from "../components/Layout";
import Container from "../components/Container";
import { SubContainer, SubContainerElement } from "../components/SubContainer";
import LinkButton from "../components/LinkButton";
import { formStyles } from "../styles/formStyles";

// Local components for this page only
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

function Checkbox({ checked, disabled, onChange, children }) {
  const checkmarkStyle = {
    ...formStyles.checkmark,
    backgroundColor: disabled ? "#777" : "#333",
  };

  const rectangleStyle = {
    position: "absolute",
    display: checked ? "block" : "none",
    left: "50%",
    top: "50%",
    transform: "translate(-50%, -50%)",
    width: "12px",
    height: "12px",
    backgroundColor: "white",
  };

  return (
    <label style={formStyles.checkboxLabel}>
      <input
        type="checkbox"
        checked={checked}
        disabled={disabled}
        onChange={onChange}
        style={formStyles.hiddenCheckbox}
      />
      {children}
      <span style={checkmarkStyle}>
        {checked && <span style={rectangleStyle} />}
      </span>
    </label>
  );
}

function PlayerInput({ index, player, corporations, onUpdate }) {
  return (
    <div style={formStyles.playerInputDiv}>
      <input
        type="text"
        style={formStyles.playerInput}
        placeholder={`Player ${index + 1} name`}
        value={player.name}
        onChange={(e) => onUpdate(index, "name", e.target.value)}
      />
      <select
        style={formStyles.playerInput}
        value={player.corporation}
        onChange={(e) => onUpdate(index, "corporation", e.target.value)}
      >
        <option value="">Select Corporation</option>
        {corporations.map((corp) => (
          <option key={corp} value={corp}>
            {corp}
          </option>
        ))}
      </select>
    </div>
  );
}

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
  colonies: [
    "Callisto",
    "Ceres",
    "Europa",
    "Enceladus",
    "Ganymede",
    "Io",
    "Luna",
    "Miranda",
    "Pluto",
    "Titan",
    "Triton",
  ],
};

function AddGamePage() {
  const navigate = useNavigate();

  // Basic game options
  const [name, setName] = useState("");
  const [date, setDate] = useState("");
  const [map, setMap] = useState("Tharsis");
  const [generations, setGenerations] = useState(0);
  const [playerNumber, setPlayerNumber] = useState(4);
  const [players, setPlayers] = useState([]);

  // Expansions
  const [expansions, setExpansions] = useState({
    "Base Game": true,
    "Draft": true,
    "Corporate Era": true,
    "Prelude": true,
    "Prelude 2": false,
    "Venus Next": false,
    "Colonies": false,
    "Turmoil": false,
    "Milestones & Awards": false,
    "Promo": false,
  });
  const [showExpansions, setShowExpansions] = useState(true);

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

  const updatePlayerData = (index, field, value) => {
    const newPlayers = [...players];
    newPlayers[index][field] = value;
    setPlayers(newPlayers);
  };

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
          <SubContainerElement
            label="Name:"
            input={{
              type: "text",
              className: "option-input",
              value: name,
              onChange: (e) => setName(e.target.value),
            }}
          />

          <SubContainerElement
            label="Date:"
            input={{
              type: "date",
              className: "option-input",
              value: date,
              onChange: (e) => setDate(e.target.value),
              required: true,
            }}
          />

          <SubContainerElement label="Map:">
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
            <div>
              <label>Select expansion:</label>
              <input
                type="button"
                value={showExpansions ? "Hide" : "Show"}
                style={{
                  ...formStyles.optionInput,
                  cursor: "pointer",
                }}
                onClick={() => setShowExpansions(!showExpansions)}
              />
            </div>
            {showExpansions && (
              <div style={formStyles.subcontainerBox}>
                {Object.entries(expansions).map(([key, value]) => (
                  <Checkbox
                    key={key}
                    checked={value}
                    disabled={key === "Base Game"}
                    onChange={() => handleExpansionChange(key)}
                  >
                    {key}
                  </Checkbox>
                ))}
              </div>
            )}
          </SubContainerElement>

          <SubContainerElement
            label="Number of generations"
            input={{
              type: "number",
              className: "option-input",
              min: 1,
              max: 20,
              value: generations,
              onChange: (e) => setGenerations(parseInt(e.target.value) || 0),
            }}
          />

          <SubContainerElement
            label="Number of players"
            input={{
              type: "number",
              min: 1,
              max: 6,
              value: playerNumber,
              onChange: (e) =>
                setPlayerNumber(
                  Math.min(5, Math.max(1, parseInt(e.target.value) || 1)),
                ),
            }}
          >
            <div style={formStyles.subcontainerBox}>
              {players.map((player, index) => (
                <PlayerInput
                  key={index}
                  index={index}
                  player={player}
                  corporations={getAvailableCorporations()}
                  onUpdate={updatePlayerData}
                />
              ))}
            </div>
          </SubContainerElement>
        </SubContainer>
      </Container>

      <Container title="Milestones" titleStyle="banner">
        <SubContainer>
          <SubContainerElement>
            <div style={formStyles.subcontainerBox}>
              {(expansions["Milestones & Awards"] ? selectedMilestones : getAvailableMilestones()).map((milestone, index) => (
                <div
                  key={index}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    margin: "10px",
                  }}
                >
                  {expansions["Milestones & Awards"] ? (
                    <select
                      style={{
                        textAlign: "center",
                        fontFamily: "inherit",
                        fontSize: "15px",
                        background: "inherit",
                        height: "35px",
                        width: "225px",
                        boxSizing: "border-box",
                        margin: "0 5px",
                        cursor: "pointer",
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
                        textAlign: "center",
                        fontFamily: "inherit",
                        fontSize: "15px",
                        height: "35px",
                        width: "225px",
                        boxSizing: "border-box",
                        margin: "0 5px",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        fontWeight: "bold",
                      }}
                    >
                      {milestone}
                    </div>
                  )}

                  <select
                    style={{
                      textAlign: "center",
                      fontFamily: "Courier New, Courier, monospace",
                      fontWeight: "bolder",
                      fontSize: "15px",
                      background: "inherit",
                      height: "35px",
                      width: "225px",
                      boxSizing: "border-box",
                      margin: "0 5px",
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
            </div>
          </SubContainerElement>
        </SubContainer>
      </Container>

      <Container title="Awards" titleStyle="banner">
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
                        padding: "10px",
                        fontSize: "18px",
                        fontWeight: "bold",
                        borderBottom: "2px solid #999",
                        textAlign: "left",
                        width: "200px",
                      }}
                    >
                      Award
                    </th>
                    {players.map((player, index) => (
                      <th
                        key={index}
                        style={{
                          padding: "10px",
                          fontSize: "16px",
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
                  {(expansions["Milestones & Awards"] ? selectedAwards : getAvailableAwards()).map((award, awardIndex) => (
                    <tr key={awardIndex}>
                      <td style={{ padding: "10px", fontWeight: "bold" }}>
                        {expansions["Milestones & Awards"] ? (
                          <select
                            style={{
                              textAlign: "center",
                              fontFamily: "inherit",
                              fontSize: "15px",
                              background: "inherit",
                              height: "35px",
                              width: "100%",
                              boxSizing: "border-box",
                              borderColor: "#999",
                              borderRadius: "3px",
                            }}
                            value={award || ""}
                            onChange={(e) => updateSelectedAward(awardIndex, e.target.value)}
                          >
                            {!award && <option value="">Select Award</option>}
                            {getAvailableAwardsForDropdown(award).map(a => (
                              <option key={a} value={a}>{a}</option>
                            ))}
                          </select>
                        ) : (
                          award
                        )}
                      </td>
                      {players.map((player, playerIndex) => (
                        <td key={playerIndex} style={{ padding: "10px", textAlign: "center" }}>
                          <button
                            style={{
                              width: "30px",
                              height: "30px",
                              border: "2px solid #666",
                              cursor: (!award || (!isAwardFunded(award) && getFundedAwardsCount() >= 3)) ? "default" : "pointer",
                              backgroundColor:
                                (!award || (!isAwardFunded(award) && getFundedAwardsCount() >= 3)) ? "#666666" : // Much lighter gray when disabled
                                  awardPlacements[award]?.[playerIndex] === 1 ? "#FFD700" : // Gold
                                    awardPlacements[award]?.[playerIndex] === 2 ? "#CCC" : // Silver 
                                      "#444", // None/Black
                              userSelect: "none",
                              WebkitUserSelect: "none",
                              MozUserSelect: "none",
                              msUserSelect: "none",
                              WebkitTouchCallout: "none",
                              WebkitTapHighlightColor: "transparent",
                            }}
                            onClick={() => cyclePlacement(award, playerIndex)}
                            disabled={
                              !award ||
                              (!isAwardFunded(award) && getFundedAwardsCount() >= 3)
                            }
                          >
                          </button>
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </SubContainerElement>
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
