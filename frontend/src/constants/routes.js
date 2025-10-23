export const ROUTES = {
  HOME: '/',
  ADD_GAME: '/add-game',
  ADD_PLAYER: '/add-player',
  EDIT_PLAYER: '/players/:playerId/edit',
  PLAYERS: '/players',
  PLAYED_GAMES: '/played-games',
  GAME_DETAILS: '/games/:gameId',
  PLAYER_DETAILS: '/players/:playerId',
  gameDetails: (id) => `/games/${id}`,
  playerDetails: (id) => `/players/${id}`,
  editPlayer: (id) => `/players/${id}/edit`,
};
