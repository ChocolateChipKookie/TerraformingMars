import React, {
  useState,
  useEffect,
  useCallback,
  useMemo,
  useReducer,
} from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { ROUTES } from "../constants/routes";
import Layout from "../components/Layout";
import Container from "../components/Container";
import AuthenticationContainer from "../components/AuthenticationContainer";
import GameOptionsContainer from "../components/GameContainers/GameOptionsContainer";
import GamePlayersContainer from "../components/GameContainers/GamePlayersContainer";
import MilestonesContainer from "../components/GameContainers/MilestonesContainer";
import AwardsContainer from "../components/GameContainers/AwardsContainer";
import PointsContainer from "../components/GameContainers/PointsContainer";
import { SubContainer, SubContainerElement } from "../components/Container";
import LinkButton from "../components/LinkButton";
import styles from "../styles/GamePage.module.css";
import { gameData, GAME_CONSTANTS } from "../data/gameData";

// Local components for this page only

function NumericInputWithButtons({
  value,
  onChange,
  onDecrement,
  onIncrement,
}) {
  return (
    <div className={styles.numericInputContainer}>
      <button className={styles.numericInputButton} onClick={onDecrement}>
        −
      </button>
      <input
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



// Game configuration reducer
const gameConfigInitialState = {
  name: "",
  date: "",
  map: "Tharsis",
  generations: GAME_CONSTANTS.DEFAULT_GENERATIONS,
  expansions: {
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
  },
  expandedExpansions: false,
};

const gameConfigReducer = (state, action) => {
  switch (action.type) {
    case "SET_FIELD":
      return { ...state, [action.field]: action.value };
    case "SET_MAP":
      return { ...state, map: action.value };
    case "TOGGLE_EXPANSION":
      return {
        ...state,
        expansions: {
          ...state.expansions,
          [action.expansion]: !state.expansions[action.expansion],
        },
      };
    case "SET_EXPANSION":
      return {
        ...state,
        expansions: {
          ...state.expansions,
          [action.expansion]: action.value,
        },
      };
    case "TOGGLE_EXPANDED_VIEW":
      return { ...state, expandedExpansions: !state.expandedExpansions };
    case "SET_DATE":
      return { ...state, date: action.value };
    case "SET_GENERATIONS":
      const generations = Math.min(
        GAME_CONSTANTS.MAX_GENERATIONS,
        Math.max(GAME_CONSTANTS.MIN_GENERATIONS, action.value),
      );
      return { ...state, generations };
    default:
      return state;
  }
};

// Custom hook for player management
function usePlayerManagement(
  initialPlayerCount = GAME_CONSTANTS.DEFAULT_PLAYER_COUNT,
  maxPlayers = GAME_CONSTANTS.DEFAULT_MAX_PLAYERS,
) {
  const [playerNumber, setPlayerNumber] = useState(initialPlayerCount);
  const [players, setPlayers] = useState([]);
  const [playerScores, setPlayerScores] = useState([]);

  // Initialize/update players when count changes
  useEffect(() => {
    setPlayers((prevPlayers) => {
      const newPlayers = [];
      for (let i = 0; i < playerNumber; i++) {
        newPlayers.push({
          name: prevPlayers[i]?.name || "",
          corporation: prevPlayers[i]?.corporation || "",
        });
      }
      return newPlayers;
    });

    setPlayerScores((prevScores) => {
      const newScores = [];
      for (let i = 0; i < playerNumber; i++) {
        newScores.push({
          terraformingRating: prevScores[i]?.terraformingRating || "",
          cities: prevScores[i]?.cities || "",
          greeneries: prevScores[i]?.greeneries || "",
          cards: prevScores[i]?.cards || "",
          turmoilPoints: prevScores[i]?.turmoilPoints || "",
          milestonePoints: 0,
          awardPoints: 0,
          totalPoints: 0,
        });
      }
      return newScores;
    });
  }, [playerNumber]);

  const updatePlayerData = useCallback((index, field, value) => {
    setPlayers((prevPlayers) => {
      const newPlayers = [...prevPlayers];
      newPlayers[index][field] = value;
      return newPlayers;
    });
  }, []);

  const updatePlayerScore = useCallback((playerIndex, field, value) => {
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
  }, []);

  const setPlayerCount = useCallback((count) => {
    const validCount = Math.max(
      GAME_CONSTANTS.MIN_PLAYERS,
      Math.min(count, maxPlayers),
    );
    setPlayerNumber(validCount);
  }, [maxPlayers]);

  return {
    playerNumber,
    players,
    playerScores,
    setPlayerNumber: setPlayerCount,
    setPlayerScores,
    updatePlayerData,
    updatePlayerScore,
  };
}



// Custom hook for managing game objectives (milestones/awards)
function useGameObjectives(type, map, expansions, playerNumber) {
  const [selected, setSelected] = useState([]);
  const [data, setData] = useState({});
  const [isInitialized, setIsInitialized] = useState(false);

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
    // Skip initialization if already initialized (loading existing game)
    if (isInitialized) return;
    
    if (expansions["Milestones & Awards"]) {
      const mapObjectives = gameData[dataKey][map] || [];
      const numSlots = expansions["Venus Next"] ? venusSlots : defaultSlots;

      const defaultObjectives = [...mapObjectives];
      if (expansions["Venus Next"]) {
        defaultObjectives.push(...gameData[dataKey].Venus);
      }

      setSelected(defaultObjectives.slice(0, numSlots));
    } else {
      // Calculate available items inline to avoid dependency issues
      let available = [...(gameData[dataKey][map] || [])];
      if (expansions["Venus Next"]) {
        available = [...available, ...gameData[dataKey].Venus];
      }
      if (expansions["Milestones & Awards"]) {
        available = [...available, ...gameData[additionalDataKey]];
      }
      available = [...new Set(available)];

      setSelected(available);

      // Initialize data structure
      if (isAward) {
        const newData = {};
        available.forEach((item) => {
          newData[item] = {};
          for (let i = 0; i < playerNumber; i++) {
            newData[item][i] = 0;
          }
        });
        setData(newData);
      } else {
        const newData = {};
        available.forEach((item) => {
          newData[item] = -1;
        });
        setData(newData);
      }
    }
  }, [
    map,
    expansions,
    playerNumber,
    dataKey,
    additionalDataKey,
    venusSlots,
    defaultSlots,
    isAward,
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

  // Function to load existing game data
  const loadGameData = useCallback((selectedItems, itemData) => {
    setSelected(selectedItems);
    setData(itemData);
    setIsInitialized(true);
  }, []);

  return {
    selected,
    setSelected,
    data,
    setData,
    getAvailable,
    getAvailableForDropdown,
    updateSelected,
    loadGameData,
  };
}


function GamePage() {
  const navigate = useNavigate();
  const { gameId } = useParams();
  const [searchParams] = useSearchParams();
  
  // Determine mode: 'view', 'edit', or 'add'
  const [mode, setMode] = useState(() => {
    if (!gameId) return 'add'; // No gameId means adding new game
    return searchParams.get('edit') === 'true' ? 'edit' : 'view';
  });
  
  // Loading state for fetching game data
  const [isLoading, setIsLoading] = useState(false);
  const [loadError, setLoadError] = useState(null);

  // Use reducer for game configuration
  const [gameConfig, dispatch] = useReducer(
    gameConfigReducer,
    gameConfigInitialState,
  );

  // Get max players from map data (memoized to prevent re-renders)
  const maxPlayers = useMemo(
    () =>
      gameData.maps[gameConfig.map]?.maxPlayers ||
      GAME_CONSTANTS.DEFAULT_MAX_PLAYERS,
    [gameConfig.map],
  );

  // Use custom hook for player management
  const playerManager = usePlayerManagement(GAME_CONSTANTS.DEFAULT_PLAYER_COUNT, maxPlayers);

  // Authentication state
  const [actorName, setActorName] = useState('');
  const [actorPassword, setActorPassword] = useState('');
  
  // Available players from backend
  const [availablePlayers, setAvailablePlayers] = useState([]);
  const [playersLoading, setPlayersLoading] = useState(true);
  
  // Trigger for refetching data
  const [shouldRefetch, setShouldRefetch] = useState(false);

  // Use custom hooks for milestones and awards
  const milestones = useGameObjectives(
    "milestone",
    gameConfig.map,
    gameConfig.expansions,
    playerManager.playerNumber,
  );
  const awards = useGameObjectives(
    "award",
    gameConfig.map,
    gameConfig.expansions,
    playerManager.playerNumber,
  );

  // Fetch game data if in view/edit mode
  useEffect(() => {
    const fetchGameData = async () => {
      if (!gameId) return;
      
      setIsLoading(true);
      setLoadError(null);
      
      try {
        const response = await fetch(`http://localhost:8080/api/games/${gameId}`);
        if (!response.ok) {
          throw new Error('Failed to fetch game');
        }
        
        const gameData = await response.json();
        
        console.log('Fetched game data:', gameData);
        
        // Parse and populate the game data into state
        // Game basic info
        dispatch({ type: "SET_FIELD", field: "name", value: gameData.game.name });
        dispatch({ type: "SET_FIELD", field: "date", value: gameData.game.date });
        dispatch({ type: "SET_MAP", value: gameData.game.map });
        dispatch({ type: "SET_GENERATIONS", value: gameData.game.generations });
        
        // Set expansions one by one to maintain order
        if (gameData.game.expansions) {
          Object.keys(gameConfigInitialState.expansions).forEach(exp => {
            const value = gameData.game.expansions[exp] || false;
            dispatch({ type: "SET_EXPANSION", expansion: exp, value: value });
          });
        }
        
        // Set players
        if (gameData.game_players && gameData.players) {
          playerManager.setPlayerNumber(gameData.game_players.length);
          
          // Map game_players to our player format
          const mappedPlayers = gameData.game_players.map((gp, index) => {
            const player = gameData.players.find(p => p.id === gp.player_id);
            return {
              name: player ? player.name : '',
              corporation: gp.corporation
            };
          });
          
          // Update player data
          mappedPlayers.forEach((player, index) => {
            playerManager.updatePlayerData(index, 'name', player.name);
            playerManager.updatePlayerData(index, 'corporation', player.corporation);
          });
          
          // Update player scores
          gameData.game_players.forEach((gp, index) => {
            playerManager.updatePlayerScore(index, 'terraformingRating', String(gp.terraforming_rating));
            playerManager.updatePlayerScore(index, 'cities', String(gp.cities));
            playerManager.updatePlayerScore(index, 'greeneries', String(gp.greeneries));
            playerManager.updatePlayerScore(index, 'cards', String(gp.cards));
            playerManager.updatePlayerScore(index, 'turmoilPoints', String(gp.turmoil_points));
          });
        }
        
        // Set milestones
        if (gameData.milestones) {
          const milestoneData = {};
          const milestoneNames = [];
          gameData.milestones.forEach(m => {
            // Find the player index from game_player_id
            const winnerIndex = m.winner_game_player_id 
              ? gameData.game_players.findIndex(gp => gp.id === m.winner_game_player_id)
              : -1;
            milestoneData[m.name] = winnerIndex;
            milestoneNames.push(m.name);
          });
          milestones.loadGameData(milestoneNames, milestoneData);
        }
        
        // Set awards
        if (gameData.awards && gameData.award_placements) {
          const awardData = {};
          const awardNames = [];
          gameData.awards.forEach(a => {
            awardData[a.name] = {};
            awardNames.push(a.name);
            // Find placements for this award
            const placements = gameData.award_placements.filter(p => p.award_id === a.id);
            placements.forEach(p => {
              const playerIndex = gameData.game_players.findIndex(gp => gp.id === p.game_player_id);
              if (playerIndex !== -1) {
                awardData[a.name][playerIndex] = p.placement;
              }
            });
            // Initialize empty placements for players without awards
            for (let i = 0; i < gameData.game_players.length; i++) {
              if (awardData[a.name][i] === undefined) {
                awardData[a.name][i] = 0;
              }
            }
          });
          awards.loadGameData(awardNames, awardData);
        }
        
      } catch (err) {
        console.error('Failed to fetch game:', err);
        setLoadError(err.message);
      } finally {
        setIsLoading(false);
      }
    };
    
    if (shouldRefetch || gameId) {
      fetchGameData();
      setShouldRefetch(false);
    }
  }, [gameId, shouldRefetch]); // Fetch when gameId changes or shouldRefetch is triggered

  // Set default date for new games and fetch players
  useEffect(() => {
    // Only set default date for new games
    if (mode === 'add') {
      const today = new Date();
      const year = today.getFullYear();
      const month = String(today.getMonth() + 1).padStart(2, "0");
      const day = String(today.getDate()).padStart(2, "0");
      dispatch({ type: "SET_DATE", value: `${year}-${month}-${day}` });
    }
    
    // Always fetch available players
    const fetchPlayers = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/players');
        if (response.ok) {
          const data = await response.json();
          // Sort all players alphabetically by name
          const sortedPlayers = (data || []).sort((a, b) => a.name.localeCompare(b.name));
          setAvailablePlayers(sortedPlayers);
        }
      } catch (err) {
        console.error('Failed to fetch players:', err);
      } finally {
        setPlayersLoading(false);
      }
    };
    
    fetchPlayers();
  }, [mode]);

  // Adjust player count when map changes
  useEffect(() => {
    if (playerManager.playerNumber > maxPlayers) {
      playerManager.setPlayerNumber(maxPlayers);
    }
  }, [maxPlayers, playerManager.playerNumber, playerManager.setPlayerNumber]);

  const getAvailableCorporations = useCallback(() => {
    let availableCorporations = [];

    // Add corporations from each selected expansion
    Object.entries(gameConfig.expansions).forEach(([expansion, isSelected]) => {
      if (isSelected && gameData.corporationsByExpansion[expansion]) {
        availableCorporations = [
          ...availableCorporations,
          ...gameData.corporationsByExpansion[expansion],
        ];
      }
    });

    // Remove duplicates and sort
    return [...new Set(availableCorporations)].sort();
  }, [gameConfig.expansions]);

  const getSelectedMilestonesCount = useMemo(() => {
    return Object.values(milestones.data).filter((idx) => idx !== -1).length;
  }, [milestones.data]);

  const updateMilestoneWinner = useCallback(
    (milestone, playerIndex) => {
      // If trying to select a new milestone and already have max claimed, don't allow
      if (
        playerIndex !== -1 &&
        milestones.data[milestone] === -1 &&
        getSelectedMilestonesCount >= GAME_CONSTANTS.MAX_MILESTONES_CLAIMED
      ) {
        return;
      }

      milestones.setData({
        ...milestones.data,
        [milestone]: playerIndex,
      });
    },
    [milestones.data, milestones.setData, getSelectedMilestonesCount],
  );

  // Memoized calculations for performance
  const mapOptions = useMemo(() => {
    return Object.keys(gameData.maps).map((m) => (
      <option key={m} value={m}>
        {m}
      </option>
    ));
  }, []);

  const selectedCorporations = useMemo(() => {
    return playerManager.players
      .map((p) => p.corporation)
      .filter((corp) => corp !== "");
  }, [playerManager.players]);

  const expansionEntries = useMemo(() => {
    return Object.entries(gameConfig.expansions);
  }, [gameConfig.expansions]);

  const getFundedAwardsCount = useCallback(() => {
    return Object.keys(awards.data).filter(
      (award) =>
        awards.data[award] &&
        Object.values(awards.data[award]).some((placement) => placement > 0),
    ).length;
  }, [awards.data]);

  const isAwardFunded = useCallback(
    (award) => {
      return (
        awards.data[award] &&
        Object.values(awards.data[award]).some((placement) => placement > 0)
      );
    },
    [awards.data],
  );

  const cyclePlacement = useCallback(
    (award, playerIndex) => {
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
        newPlacements[award][playerIndex] ||
        GAME_CONSTANTS.AWARD_PLACEMENT.NONE;
      // Cycle through placements: none -> gold -> silver -> none
      const placementCount = Object.keys(GAME_CONSTANTS.AWARD_PLACEMENT).length;
      newPlacements[award][playerIndex] =
        (currentPlacement + 1) % placementCount;

      awards.setData(newPlacements);
    },
    [awards.data, awards.setData, getFundedAwardsCount],
  );

  // Calculate points whenever milestones or awards change
  useEffect(() => {
    if (playerManager.playerScores.length === 0) return;

    const newScores = playerManager.playerScores.map((score, playerIndex) => {
      // Calculate milestone points - memoized for performance
      const milestonePoints = Object.values(milestones.data).reduce(
        (points, winnerIndex) => {
          return winnerIndex === playerIndex
            ? points + GAME_CONSTANTS.MILESTONE_POINTS
            : points;
        },
        0,
      );

      // Calculate award points based on placement - memoized for performance
      const awardPoints = Object.entries(awards.data).reduce(
        (points, [award, placements]) => {
          const placement = placements?.[playerIndex];
          if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.GOLD) {
            return points + GAME_CONSTANTS.AWARD_POINTS.GOLD;
          } else if (placement === GAME_CONSTANTS.AWARD_PLACEMENT.SILVER) {
            return points + GAME_CONSTANTS.AWARD_POINTS.SILVER;
          }
          return points;
        },
        0,
      );

      // Only update if milestone or award points changed
      if (
        score.milestonePoints === milestonePoints &&
        score.awardPoints === awardPoints
      ) {
        return score;
      }

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

    // Only update if scores actually changed
    const hasChanged = newScores.some((newScore, index) => {
      const oldScore = playerManager.playerScores[index];
      return (
        newScore.milestonePoints !== oldScore.milestonePoints ||
        newScore.awardPoints !== oldScore.awardPoints ||
        newScore.totalPoints !== oldScore.totalPoints
      );
    });

    if (hasChanged) {
      playerManager.setPlayerScores(newScores);
    }
  }, [milestones.data, awards.data]);

  const handleSubmitGame = async () => {
    // Validate authentication
    if (!actorName || !actorPassword) {
      alert('Please enter your username and password to submit the game.');
      return;
    }

    // Format the game data for submission
    const gameData = {
      name: gameConfig.name,
      date: gameConfig.date,
      map: gameConfig.map,
      generations: gameConfig.generations,
      expansions: gameConfig.expansions,
      note: null,

      // Format players with scores
      players: playerManager.players.map((player, index) => ({
        name: player.name,
        corporation: player.corporation,
        terraforming_rating: parseInt(playerManager.playerScores[index].terraformingRating) || 0,
        cities: parseInt(playerManager.playerScores[index].cities) || 0,
        greeneries: parseInt(playerManager.playerScores[index].greeneries) || 0,
        cards: parseInt(playerManager.playerScores[index].cards) || 0,
        turmoil_points: parseInt(playerManager.playerScores[index].turmoilPoints) || 0,
      })),

      // Format milestones
      milestones: Object.entries(milestones.data).map(([name, winnerIndex]) => ({
        name: name,
        winner_game_player_index: winnerIndex === -1 ? null : winnerIndex,
      })),

      // Format awards with placements
      awards: Object.entries(awards.data).map(([name, placements]) => ({
        name: name,
        placements: Object.entries(placements)
          .filter(([_, placement]) => placement > 0)
          .map(([playerIndex, placement]) => ({
            player_index: parseInt(playerIndex),
            placement: placement,
          })),
      })),

      images: [],

      // Authentication from form inputs
      actor_name: actorName,
      actor_password: actorPassword,
    };

    console.log('Sending game data:', gameData);
    console.log('Milestones selected:', milestones.selected);
    console.log('Milestones data:', milestones.data);
    console.log('Awards selected:', awards.selected);
    console.log('Awards data:', awards.data);

    try {
      const isEditing = mode === 'edit' && gameId;
      const url = isEditing 
        ? `http://localhost:8080/api/games/${gameId}`
        : 'http://localhost:8080/api/games';
      
      const response = await fetch(url, {
        method: isEditing ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(gameData),
      });

      if (response.ok) {
        const result = await response.json();
        const message = isEditing 
          ? `Game updated successfully!`
          : `Game created successfully! Game ID: ${result.game.id}`;
        alert(message);
        navigate(ROUTES.PLAYED_GAMES);
      } else {
        const error = await response.json();
        alert(`Failed to ${isEditing ? 'update' : 'create'} game: ${error.error || 'Unknown error'}`);
      }
    } catch (error) {
      console.error('Error submitting game:', error);
      alert('Failed to submit game. Please check if the backend is running.');
    }
  };

  // Handle cancel edit - reload original data
  const handleCancelEdit = () => {
    setMode('view');
    // Trigger re-fetch of game data to discard changes
    setShouldRefetch(true);
  };

  // Get page title based on mode
  const getPageTitle = () => {
    if (mode === 'add') return 'Add Game';
    if (mode === 'edit') return 'Edit Game';
    return 'Game Details';
  };

  // Handle loading and error states
  if (isLoading) {
    return (
      <Layout>
        <Container title="Loading..." />
      </Layout>
    );
  }

  if (loadError) {
    return (
      <Layout>
        <Container title="Error" />
        <Container>
          <p>Failed to load game: {loadError}</p>
          <LinkButton onClick={() => navigate(ROUTES.PLAYED_GAMES)}>
            Back to Games
          </LinkButton>
        </Container>
      </Layout>
    );
  }

  return (
    <Layout>
      <Container title={getPageTitle()} />

      <GameOptionsContainer 
        gameConfig={gameConfig}
        dispatch={dispatch}
        playerManager={playerManager}
        maxPlayers={maxPlayers}
        readOnly={mode === 'view'}
      />

      <GamePlayersContainer
        playerManager={playerManager}
        availablePlayers={availablePlayers}
        selectedCorporations={selectedCorporations}
        getAvailableCorporations={getAvailableCorporations}
        readOnly={mode === 'view'}
      />

      <MilestonesContainer
        milestones={milestones}
        gameConfig={gameConfig}
        playerManager={playerManager}
        getSelectedMilestonesCount={getSelectedMilestonesCount}
        updateMilestoneWinner={updateMilestoneWinner}
        readOnly={mode === 'view'}
      />

      <AwardsContainer
        awards={awards}
        gameConfig={gameConfig}
        playerManager={playerManager}
        cyclePlacement={cyclePlacement}
        isAwardFunded={isAwardFunded}
        getFundedAwardsCount={getFundedAwardsCount}
        readOnly={mode === 'view'}
      />

      <PointsContainer
        playerManager={playerManager}
        gameConfig={gameConfig}
        readOnly={mode === 'view'}
      />

      {/* Only show authentication in add/edit modes */}
      {(mode === 'add' || mode === 'edit') && (
        <AuthenticationContainer
          actorName={actorName}
          setActorName={setActorName}
          actorPassword={actorPassword}
          setActorPassword={setActorPassword}
          players={availablePlayers}
          playersLoading={playersLoading}
          title="Authentication"
        />
      )}

      {/* Dynamic buttons based on mode */}
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between',
        gap: '2rem', 
        margin: '1% auto',
        maxWidth: '900px',
        padding: '0'
      }}>
        {mode === 'view' ? (
          <>
            <LinkButton 
              onClick={() => navigate(ROUTES.PLAYED_GAMES)}
              style={{ width: 'calc(50% - 1rem)' }}
            >
              Back to Games
            </LinkButton>
            <LinkButton 
              onClick={() => setMode('edit')} 
              style={{ backgroundColor: '#2196F3', width: 'calc(50% - 1rem)' }}
            >
              Edit Game
            </LinkButton>
          </>
        ) : mode === 'edit' ? (
          <>
            <LinkButton 
              onClick={handleCancelEdit}
              style={{ width: 'calc(50% - 1rem)' }}
            >
              Cancel
            </LinkButton>
            <LinkButton 
              onClick={handleSubmitGame} 
              style={{ backgroundColor: '#4CAF50', width: 'calc(50% - 1rem)' }}
            >
              Save Changes
            </LinkButton>
          </>
        ) : (
          <>
            <LinkButton 
              onClick={() => navigate(ROUTES.HOME)}
              style={{ width: 'calc(50% - 1rem)' }}
            >
              Main page
            </LinkButton>
            <LinkButton 
              onClick={handleSubmitGame} 
              style={{ backgroundColor: '#4CAF50', width: 'calc(50% - 1rem)' }}
            >
              Submit Game
            </LinkButton>
          </>
        )}
      </div>
    </Layout>
  );
}

export default GamePage;
