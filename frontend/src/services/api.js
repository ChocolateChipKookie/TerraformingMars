// API Service Layer
// Centralized location for all API calls

const API_BASE_URL = '/api';

// Helper function for handling responses
const handleResponse = async (response) => {
  if (!response.ok) {
    const error = await response.json().catch((parseError) => {
      console.error(`Failed to parse error response from ${response.url}:`, parseError);
      console.error(`Response status: ${response.status} ${response.statusText}`);
      return { error: 'Unknown error' };
    });
    console.error(`API Error [${response.status}] ${response.url}:`, error.error || error);
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

  // Delete game
  delete: async (gameId, actorName, actorPassword) => {
    const response = await fetch(`${API_BASE_URL}/games/${gameId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        actor_name: actorName,
        actor_password: actorPassword,
      }),
    });
    if (!response.ok) {
      const error = await response.json().catch((parseError) => {
        console.error(`Failed to parse error response from ${response.url}:`, parseError);
        console.error(`Response status: ${response.status} ${response.statusText}`);
        return { error: 'Unknown error' };
      });
      console.error(`API Error [${response.status}] ${response.url}:`, error.error || error);
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }
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

  // Get single player by ID
  getById: async (playerId) => {
    const response = await fetch(`${API_BASE_URL}/players/${playerId}`);
    return handleResponse(response);
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

  // Update existing player
  update: async (playerId, playerData) => {
    const response = await fetch(`${API_BASE_URL}/players/${playerId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(playerData),
    });
    return handleResponse(response);
  },
};

// Rating API calls (derived from full game history, computed backend-side)
export const ratingApi = {
  // Get one player's per-game rating timeline.
  getHistory: async (playerId) => {
    const response = await fetch(`${API_BASE_URL}/players/${playerId}/rating-history`);
    return handleResponse(response);
  },
};