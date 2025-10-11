import sharedGameData from '../../backend/shared/game-data.json';

// Game rule constants - all from shared JSON
export const GAME_CONSTANTS = {
  MAX_MILESTONES_CLAIMED: sharedGameData.constants.maxMilestonesClaimed,
  MAX_AWARDS_FUNDED: sharedGameData.constants.maxAwardsFunded,
  MILESTONE_POINTS: sharedGameData.constants.milestonePoints,
  AWARD_POINTS: {
    GOLD: sharedGameData.constants.awardPointsGold,
    SILVER: sharedGameData.constants.awardPointsSilver,
  },
  AWARD_PLACEMENT: {
    NONE: sharedGameData.constants.awardPlacementNone,
    GOLD: sharedGameData.constants.awardPlacementGold,
    SILVER: sharedGameData.constants.awardPlacementSilver,
  },
  DEFAULT_MILESTONE_SLOTS: sharedGameData.constants.defaultMilestoneSlots,
  VENUS_MILESTONE_SLOTS: sharedGameData.constants.venusMilestoneSlots,
  DEFAULT_AWARD_SLOTS: sharedGameData.constants.defaultAwardSlots,
  VENUS_AWARD_SLOTS: sharedGameData.constants.venusAwardSlots,
  MIN_PLAYERS: sharedGameData.constants.minPlayers,
  DEFAULT_MAX_PLAYERS: sharedGameData.constants.defaultMaxPlayers,
  MIN_GENERATIONS: sharedGameData.constants.minGenerations,
  MAX_GENERATIONS: sharedGameData.constants.maxGenerations,
  DEFAULT_GENERATIONS: sharedGameData.constants.defaultGenerations,
  DEFAULT_PLAYER_COUNT: sharedGameData.constants.defaultPlayerCount,
  DEFAULT_TR: sharedGameData.constants.defaultTR,
};

// Transform shared data to match expected format
const maps = {};
Object.entries(sharedGameData.maps).forEach(([mapName, mapData]) => {
  // Skip Venus as it's not a real map
  if (mapName === 'Venus') return;

  maps[mapName] = {
    name: mapName,
    maxPlayers: mapData.maxPlayers,
  };
});

// Extract milestones from maps
const milestones = {};
Object.entries(sharedGameData.maps).forEach(([mapName, mapData]) => {
  milestones[mapName] = mapData.milestones;
});

// Extract awards from maps
const awards = {};
Object.entries(sharedGameData.maps).forEach(([mapName, mapData]) => {
  awards[mapName] = mapData.awards;
});

export const gameData = {
  maps,
  corporationsByExpansion: sharedGameData.corporations,
  milestones,
  additionalMilestones: sharedGameData.allMilestones,
  awards,
  additionalAwards: sharedGameData.allAwards,
};

