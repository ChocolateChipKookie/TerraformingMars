package auth

import (
	"fmt"
	"terraforming-mars-backend/internal/models"
)

func CanCreateGames(player models.Player) error {
	if player.Role == models.RoleAdmin || player.Role == models.RoleUser {
		return nil
	}
	return fmt.Errorf("only admins and users can create games")
}

func CanCreatePlayers(player models.Player, targetRole models.PlayerRole) error {
	if !targetRole.IsValid() {
		return fmt.Errorf("Invalid role '%s'", targetRole)
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

func CanUpdatePlayer(actor models.Player, target models.Player) error {
	switch actor.Role {
	case models.RoleAdmin:
		// Admin can update anyone
		return nil
	case models.RoleUser:
		// User can update themselves or players they created
		if actor.ID == target.ID {
			return nil
		}
		if target.CreatedBy != nil && *target.CreatedBy == actor.ID && target.Role == models.RolePlayer {
			return nil
		}
		return fmt.Errorf("users can only update themselves or players they created")
	case models.RolePlayer:
		return fmt.Errorf("players cannot update anyting")
	default:
		return fmt.Errorf("unknown role: %s", actor.Role)
	}
}

func ValidateRoleTransition(actor models.Player, currentRole models.PlayerRole, newRole models.PlayerRole) error {
	if !newRole.IsValid() {
		return fmt.Errorf("Invalid role '%s'", newRole)
	}

	if actor.Role != models.RoleAdmin {
		if currentRole != newRole {
			return fmt.Errorf("only admins can change player roles")
		}
		return nil
	}

	// Admin can change any role
	return nil
}

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

func RequiresPassword(role models.PlayerRole) bool {
	return role == models.RoleAdmin || role == models.RoleUser
}
