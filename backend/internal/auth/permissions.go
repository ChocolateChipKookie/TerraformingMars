package auth

import (
	"fmt"
	"terraforming-mars-backend/internal/models"
)

// CanCreateGames checks if a player can create games
func CanCreateGames(player models.Player) error {
	if player.Role == models.RoleAdmin || player.Role == models.RoleUser {
		return nil
	}
	return fmt.Errorf("only admins and users can create games")
}

// CanCreatePlayers checks if a player can create other players
func CanCreatePlayers(player models.Player, targetRole models.PlayerRole) error {
	// Validate target role
	if err := IsValidRole(targetRole); err != nil {
		return err
	}

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
		return fmt.Errorf("players cannot update anyting")
	default:
		return fmt.Errorf("unknown role: %s", actor.Role)
	}
}

// CanUpdatePlayerName checks if a player can update another player's name
func CanUpdatePlayerName(actor models.Player, target models.Player) error {
	switch actor.Role {
	case models.RoleAdmin:
		// Admin can update anyone's name
		return nil
	case models.RoleUser:
		// Users cannot update any names (to prevent name history issues)
		return fmt.Errorf("users cannot update player names")
	case models.RolePlayer:
		return fmt.Errorf("players cannot update anything")
	default:
		return fmt.Errorf("unknown role: %s", actor.Role)
	}
}

// ValidateRoleTransition checks if a role change is allowed
func ValidateRoleTransition(actor models.Player, currentRole models.PlayerRole, newRole models.PlayerRole) error {
	// Validate new role
	if err := IsValidRole(newRole); err != nil {
		return err
	}

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

// IsValidRole checks if a role is valid
func IsValidRole(role models.PlayerRole) error {
	if role == models.RoleAdmin || role == models.RoleUser || role == models.RolePlayer {
		return nil
	}
	return fmt.Errorf("invalid role '%s': must be 'admin', 'user', or 'player'", role)
}

// CanModifyGame checks if a player can modify a game
func CanModifyGame(actor models.Player, gameCreatedBy int) error {
	switch actor.Role {
	case models.RoleAdmin:
		// Admin can modify any game
		return nil
	case models.RoleUser:
		// User/Player can only modify games they created
		if actor.ID == gameCreatedBy {
			return nil
		}
		return fmt.Errorf("only the game creator or admin can modify this game")
	case models.RolePlayer:
		return fmt.Errorf("players cannot modify games")
	default:
		return fmt.Errorf("unknown role: %s", actor.Role)
	}
}
