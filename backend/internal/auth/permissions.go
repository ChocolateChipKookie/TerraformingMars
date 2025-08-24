package auth

import (
	"fmt"
	"terraforming-mars-backend/internal/models"
)

// CanCreateGames checks if a player can create games
func CanCreateGames(player models.Player) bool {
	return player.Role == models.RoleAdmin || player.Role == models.RoleUser
}

// CanCreatePlayers checks if a player can create other players
func CanCreatePlayers(player models.Player, targetRole models.PlayerRole) error {
	switch player.Role {
	case models.RoleAdmin:
		// Admin can create any type of player
		return nil
	case models.RoleUser:
		// User can only create "player" role
		if targetRole != models.RolePlayer {
			return fmt.Errorf("users can only create players with 'player' role")
		}
		return nil
	case models.RolePlayer:
		// Player cannot create other players
		return fmt.Errorf("players cannot create other players")
	default:
		return fmt.Errorf("unknown role: %s", player.Role)
	}
}

// CanUpdatePlayer checks if a player can update another player
func CanUpdatePlayer(actor models.Player, target models.Player) error {
	switch actor.Role {
	case models.RoleAdmin:
		// Admin can update anyone
		return nil
	case models.RoleUser:
		// User can update themselves or players they created
		if actor.ID == target.ID {
			return nil // Can update self
		}
		if target.CreatedBy != nil && *target.CreatedBy == actor.ID && target.Role == models.RolePlayer {
			return nil // Can update players they created
		}
		return fmt.Errorf("users can only update themselves or players they created")
	case models.RolePlayer:
		// Player can only update themselves
		if actor.ID == target.ID {
			return nil
		}
		return fmt.Errorf("players can only update themselves")
	default:
		return fmt.Errorf("unknown role: %s", actor.Role)
	}
}

// ValidateRoleTransition checks if a role change is allowed
func ValidateRoleTransition(actor models.Player, currentRole, newRole models.PlayerRole) error {
	// Only admins can change roles
	if actor.Role != models.RoleAdmin {
		if currentRole != newRole {
			return fmt.Errorf("only admins can change player roles")
		}
		return nil
	}

	// Admin can change any role
	return nil
}

// RequiresPassword checks if a role requires a password
func RequiresPassword(role models.PlayerRole) bool {
	return role == models.RoleAdmin || role == models.RoleUser
}

// CanModifyGame checks if a player can modify a game
func CanModifyGame(actor models.Player, gameCreatedBy int) error {
	switch actor.Role {
	case models.RoleAdmin:
		// Admin can modify any game
		return nil
	case models.RoleUser, models.RolePlayer:
		// User/Player can only modify games they created
		if actor.ID == gameCreatedBy {
			return nil
		}
		return fmt.Errorf("only the game creator or admin can modify this game")
	default:
		return fmt.Errorf("unknown role: %s", actor.Role)
	}
}

