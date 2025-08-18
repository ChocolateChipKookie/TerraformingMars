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
  corporations: [
    "Beginner",
    "Credicor",
    "Ecoline",
    "Helion",
    "Interplanetary Cinematics",
    "Inventrix",
    "Mining Guild",
    "Phobolog",
    "Tharsis Republic",
    "Thorgate",
    "UNMI",
    "Saturn Systems",
    "Teractor",
    "Aridor",
    "Arklight",
    "Polyphemos",
    "Poseidon",
    "Stormcraft Incorporated",
    "Cheung Shing Mars",
    "Point Luna",
    "Robinson Industries",
    "Valley Trust",
    "Vitor",
    "Aphrodite",
    "Celestic",
    "Manutech",
    "Morning Star INC",
    "Viron",
    "Arcadian Communities",
    "Recyclon",
    "Splice",
    "Lakefront Resorts",
    "Pristar",
    "Septem Tribus",
    "Terralabs Research",
    "Utopia Invest",
    "Factorum",
    "Mons Insurance",
    "Philares",
    "Ecotec",
    "Tycho Magnetics",
    "Kuiper Cooperative",
    "Spire",
    "Sagitta",
    "Palladin Shipping",
    "Nirgal Enterprises",
    "Astrodrill",
    "Pharmacy Union",
  ].sort(),
  milestones: {
    Tharsis: ["Builder", "Gardener", "Mayor", "Planner", "Terraformer"],
    Hellas: [
      "Diversifier",
      "Tactician",
      "Polar Explorer",
      "Energizer",
      "Rim Settler",
    ],
    Elysium: ["Generalist", "Specialist", "Ecologist", "Tycoon", "Legend"],
    Venus: ["Hoverlord"],
  },
  awards: {
    Tharsis: ["Landlord", "Scientist", "Banker", "Thermalist", "Miner"],
    Hellas: ["Cultivator", "Magnate", "Space Baron", "Excentric", "Contractor"],
    Elysium: [
      "Celebrity",
      "Industrialist",
      "Desert Settler",
      "Estate Dealer",
      "Benefactor",
    ],
    Venus: ["Venuphile"],
  },
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
    Base: true,
    Draft: true,
    "Corporate Era": true,
    Prelude: true,
    Venus: false,
    Colonies: false,
    Turmoil: false,
  });
  const [showExpansions, setShowExpansions] = useState(true);

  // Milestones
  const [milestoneNumber, setMilestoneNumber] = useState(3);
  const [milestones, setMilestones] = useState([]);

  // Awards
  const [awardNumber, setAwardNumber] = useState(3);
  const [awards, setAwards] = useState([]);

  // Colonies
  const [selectedColonies, setSelectedColonies] = useState({});
  const [showColonies, setShowColonies] = useState(true);

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

  // Update milestones when number changes
  useEffect(() => {
    const newMilestones = [];
    for (let i = 0; i < milestoneNumber; i++) {
      newMilestones.push({
        milestone: milestones[i]?.milestone || "",
        playerIndex: milestones[i]?.playerIndex ?? -1, // Store index instead of name
      });
    }
    setMilestones(newMilestones);
  }, [milestoneNumber]);

  // Update awards when number changes
  useEffect(() => {
    const newAwards = [];
    for (let i = 0; i < awardNumber; i++) {
      const playerPoints = {};
      players.forEach((_, playerIndex) => {
        playerPoints[playerIndex] =
          awards[i]?.playerPoints?.[playerIndex] || "";
      });
      newAwards.push({
        award: awards[i]?.award || "",
        playerPoints: playerPoints,
      });
    }
    setAwards(newAwards);
  }, [awardNumber, players.length]);

  const updatePlayerData = (index, field, value) => {
    const newPlayers = [...players];
    newPlayers[index][field] = value;
    setPlayers(newPlayers);
  };

  const handleExpansionChange = (expansion) => {
    setExpansions({ ...expansions, [expansion]: !expansions[expansion] });
  };

  const getAvailableMilestones = () => {
    let availableMilestones = [...(gameData.milestones[map] || [])];
    if (expansions.Venus) {
      availableMilestones = [
        ...availableMilestones,
        ...gameData.milestones.Venus,
      ];
    }
    return availableMilestones.sort();
  };

  const updateMilestone = (index, field, value) => {
    const newMilestones = [...milestones];
    newMilestones[index][field] = value;
    setMilestones(newMilestones);
  };

  const getAvailableAwards = () => {
    let availableAwards = [...(gameData.awards[map] || [])];
    if (expansions.Venus) {
      availableAwards = [...availableAwards, ...gameData.awards.Venus];
    }
    return availableAwards.sort();
  };

  const updateAward = (index, field, value) => {
    const newAwards = [...awards];
    newAwards[index][field] = value;
    setAwards(newAwards);
  };

  const updateAwardPlayerPoints = (awardIndex, playerIndex, points) => {
    const newAwards = [...awards];
    newAwards[awardIndex].playerPoints[playerIndex] = points;
    setAwards(newAwards);
  };

  const handleColonyChange = (colonyName) => {
    setSelectedColonies((prev) => ({
      ...prev,
      [colonyName]: prev[colonyName] !== undefined ? undefined : 0,
    }));
  };

  const increaseColonyNumber = (colonyName) => {
    setSelectedColonies((prev) => ({
      ...prev,
      [colonyName]: (prev[colonyName] + 1) % 4,
    }));
  };

  const getExpectedColonies = () => {
    let expected = playerNumber + 2;
    if (players.some((player) => player.corporation === "Aridor")) {
      expected += 1;
    }
    return expected;
  };

  const getCurrentColonies = () => {
    return Object.keys(selectedColonies).filter(
      (colony) => selectedColonies[colony] !== undefined,
    ).length;
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

    // Calculate award points for all players
    const playerAwardPoints = new Array(players.length).fill(0);

    awards.forEach((award) => {
      if (!award.award) return; // Skip if no award selected

      // Get all player counts for this award (including 0 counts)
      const playerCounts = players.map((_, playerIndex) => ({
        playerIndex,
        count: parseInt(award.playerPoints?.[playerIndex] || 0),
      }));

      // If all values are 0, skip this award (not entered)
      const totalCount = playerCounts.reduce((sum, p) => sum + p.count, 0);
      if (totalCount === 0) return;

      // Sort by count (descending), then by player index for tie-breaking consistency
      playerCounts.sort((a, b) => {
        if (b.count !== a.count) return b.count - a.count;
        return a.playerIndex - b.playerIndex;
      });

      // Find first place (highest count)
      const firstPlaceCount = playerCounts[0].count;
      const firstPlacePlayers = playerCounts.filter(
        (p) => p.count === firstPlaceCount,
      );

      // Award first place points (5 each)
      firstPlacePlayers.forEach((p) => {
        playerAwardPoints[p.playerIndex] += 5;
      });

      // If there's a tie for first place, no second place points
      if (firstPlacePlayers.length === 1) {
        // Find second place (next highest count)
        const remainingPlayers = playerCounts.filter(
          (p) => p.count < firstPlaceCount,
        );
        if (remainingPlayers.length > 0) {
          const secondPlaceCount = remainingPlayers[0].count;
          const secondPlacePlayers = remainingPlayers.filter(
            (p) => p.count === secondPlaceCount,
          );

          // Award second place points (2 each)
          secondPlacePlayers.forEach((p) => {
            playerAwardPoints[p.playerIndex] += 2;
          });
        }
      }
    });

    const newScores = playerScores.map((score, playerIndex) => {
      // Calculate milestone points
      let milestonePoints = 0;
      milestones.forEach((milestone) => {
        if (milestone.playerIndex === playerIndex) {
          milestonePoints += 5;
        }
      });

      // Get calculated award points
      const awardPoints = playerAwardPoints[playerIndex];

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
  }, [milestones, awards]);

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
                    disabled={key === "Base"}
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
                  corporations={gameData.corporations}
                  onUpdate={updatePlayerData}
                />
              ))}
            </div>
          </SubContainerElement>
        </SubContainer>
      </Container>

      <Container title="Milestones" titleStyle="banner">
        <SubContainer>
          <SubContainerElement
            label="Number of milestones"
            input={{
              type: "number",
              min: 0,
              max: 3,
              value: milestoneNumber,
              onChange: (e) =>
                setMilestoneNumber(
                  Math.min(3, Math.max(0, parseInt(e.target.value) || 0)),
                ),
            }}
          >
            {milestoneNumber > 0 && (
              <div style={formStyles.subcontainerBox}>
                {milestones.map((milestone, index) => (
                  <div
                    key={index}
                    style={{
                      height: "auto",
                      textAlign: "center",
                      margin: "10px",
                    }}
                  >
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
                      }}
                      value={milestone.milestone}
                      onChange={(e) =>
                        updateMilestone(index, "milestone", e.target.value)
                      }
                    >
                      <option value="">Select Milestone</option>
                      {getAvailableMilestones().map((m) => (
                        <option key={m} value={m}>
                          {m}
                        </option>
                      ))}
                    </select>

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
                      value={milestone.playerIndex}
                      onChange={(e) =>
                        updateMilestone(
                          index,
                          "playerIndex",
                          parseInt(e.target.value),
                        )
                      }
                    >
                      <option value={-1}>Select Player</option>
                      {players.map((p, i) => (
                        <option key={i} value={i}>
                          {p.name || `Player ${i + 1}`}
                        </option>
                      ))}
                    </select>
                  </div>
                ))}
              </div>
            )}
          </SubContainerElement>
        </SubContainer>
      </Container>

      <Container title="Awards" titleStyle="banner">
        <SubContainer>
          <SubContainerElement
            label="Number of awards"
            input={{
              type: "number",
              min: 0,
              max: 3,
              value: awardNumber,
              onChange: (e) =>
                setAwardNumber(
                  Math.min(3, Math.max(0, parseInt(e.target.value) || 0)),
                ),
            }}
          >
            {awardNumber > 0 && (
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
                            fontFamily: "Courier New, Courier, monospace",
                          }}
                        >
                          {player.name || `Player ${index + 1}`}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {awards.map((award, awardIndex) => (
                      <tr key={awardIndex}>
                        <td style={{ padding: "10px" }}>
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
                            value={award.award}
                            onChange={(e) =>
                              updateAward(awardIndex, "award", e.target.value)
                            }
                          >
                            <option value="">Select Award</option>
                            {getAvailableAwards().map((a) => (
                              <option key={a} value={a}>
                                {a}
                              </option>
                            ))}
                          </select>
                        </td>
                        {players.map((player, playerIndex) => (
                          <td key={playerIndex} style={{ padding: "10px" }}>
                            <input
                              type="text"
                              style={{
                                textAlign: "center",
                                fontFamily: "inherit",
                                fontSize: "15px",
                                background: "inherit",
                                height: "35px",
                                width: "100%",
                                boxSizing: "border-box",
                                borderWidth: "1px",
                                borderStyle: "solid",
                                borderColor: "#999",
                                borderRadius: "3px",
                              }}
                              value={award.playerPoints?.[playerIndex] || ""}
                              onChange={(e) => {
                                const value = e.target.value;
                                // Only allow empty string or valid integers
                                if (value === "" || /^-?\d+$/.test(value)) {
                                  updateAwardPlayerPoints(
                                    awardIndex,
                                    playerIndex,
                                    value,
                                  );
                                }
                              }}
                              placeholder="0"
                            />
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </SubContainerElement>
        </SubContainer>
      </Container>

      {expansions.Colonies && (
        <Container title="Colonies" titleStyle="banner">
          <SubContainer>
            <SubContainerElement>
              <div>
                <label>Select colonies:</label>
                <input
                  type="button"
                  value={showColonies ? "Hide" : "Show"}
                  style={{
                    ...formStyles.optionInput,
                    cursor: "pointer",
                  }}
                  onClick={() => setShowColonies(!showColonies)}
                />
              </div>
              {showColonies && (
                <div style={formStyles.subcontainerBox}>
                  {gameData.colonies.map((colony) => (
                    <Checkbox
                      key={colony}
                      checked={selectedColonies[colony] !== undefined}
                      onChange={() => handleColonyChange(colony)}
                    >
                      {colony}
                    </Checkbox>
                  ))}
                </div>
              )}
            </SubContainerElement>

            {Object.keys(selectedColonies).length > 0 && (
              <SubContainerElement>
                {(() => {
                  const expected = getExpectedColonies();
                  const current = getCurrentColonies();
                  return (
                    expected !== current && (
                      <div
                        style={{
                          backgroundColor: "#ffcccc",
                          border: "2px solid #ff0000",
                          borderRadius: "5px",
                          padding: "10px",
                          margin: "10px 0",
                          textAlign: "center",
                          fontWeight: "bold",
                          color: "#cc0000",
                        }}
                      >
                        WARNING! Wrong number of colonies, {expected} expected!
                      </div>
                    )
                  );
                })()}
                <div
                  style={{
                    display: "flex",
                    flexWrap: "wrap",
                    justifyContent: "space-between",
                    gap: "10px",
                    marginTop: "20px",
                  }}
                >
                  {Object.entries(selectedColonies).map(
                    ([colony, count]) =>
                      count !== undefined && (
                        <div
                          key={colony}
                          style={{
                            display: "flex",
                            flexDirection: "column",
                            alignItems: "center",
                            cursor: "pointer",
                            padding: "10px",
                            border: "2px solid #999",
                            borderRadius: "10px",
                            backgroundColor: "#f0f0f0",
                            flex: "1",
                            userSelect: "none",
                            WebkitUserSelect: "none",
                            MozUserSelect: "none",
                            msUserSelect: "none",
                          }}
                          onClick={() => increaseColonyNumber(colony)}
                        >
                          <img
                            src={`/resources/icons/colonies/${colony.toLowerCase()}.png`}
                            alt={colony}
                            style={{
                              width: "50px",
                              height: "50px",
                              marginBottom: "5px",
                            }}
                            onError={(e) => {
                              e.target.style.display = "none";
                            }}
                          />
                          <div
                            style={{
                              textAlign: "center",
                              fontSize: "14px",
                              fontWeight: "bold",
                            }}
                          >
                            {colony}
                            <br />({count})
                          </div>
                        </div>
                      ),
                  )}
                </div>
              </SubContainerElement>
            )}
          </SubContainer>
        </Container>
      )}

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
