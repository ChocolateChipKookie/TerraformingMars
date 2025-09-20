// API Service Layer
// Centralized location for all API calls

const API_BASE_URL = 'http://localhost:8080/api';

// Helper function for handling responses
const handleResponse = async (response) => {
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }
  return response.json();
};

// Game API calls
export const gameApi = {
  // Get all games
  getAll: async () => {
    const response = await fetch(`${API_BASE_URL}/games`);
    return handleResponse(response);
  },

  // Get single game by ID
  getById: async (gameId) => {
    const response = await fetch(`${API_BASE_URL}/games/${gameId}`);
    return handleResponse(response);
  },

  // Create new game
  create: async (gameData) => {
    const response = await fetch(`${API_BASE_URL}/games`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(gameData),
    });
    return handleResponse(response);
  },

  // Update existing game
  update: async (gameId, gameData) => {
    const response = await fetch(`${API_BASE_URL}/games/${gameId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(gameData),
    });
    return handleResponse(response);
  },
};

// Player API calls
export const playerApi = {
  // Get all players
  getAll: async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/players`);
      if (response.ok) {
        return await response.json() || [];
      }
      return [];
    } catch (err) {
      console.error('Failed to fetch players:', err);
      return [];
    }
  },

  // Create new player
  create: async (playerData) => {
    const response = await fetch(`${API_BASE_URL}/players`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(playerData),
    });
    return handleResponse(response);
  },
};

// Player Stats API calls
export const playerStatsApi = {
  // Get all player stats
  getAll: async () => {
    const response = await fetch(`${API_BASE_URL}/player-stats`);
    return handleResponse(response);
  },
};