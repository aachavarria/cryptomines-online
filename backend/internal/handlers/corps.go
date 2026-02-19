package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cryptomines-online/backend/internal/combat"
	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/services"
)

// Corp represents a corporation with its members and bonuses
type Corp struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Tag         string                `json:"tag"`
	Description string                `json:"description"`
	LeaderID    string                `json:"leader_id"`
	Level       int                   `json:"level"`
	Wealth      int64                 `json:"wealth"`
	MaxMembers  int                   `json:"max_members"`
	CreatedAt   time.Time             `json:"created_at"`
	Bonuses     *services.CorpBonuses `json:"bonuses,omitempty"`
}

// CorpMember represents a member of a corporation
type CorpMember struct {
	PlayerID            string    `json:"player_id"`
	PlayerName          string    `json:"player_name"`
	Role                string    `json:"role"`
	ContributionPoints  int64     `json:"contribution_points"`
	DailyContribution   int       `json:"daily_contribution"`
	JoinedAt            time.Time `json:"joined_at"`
}

// GetCorp handles GET /api/corp
// Returns player's corp info or null if not in a corp
func GetCorp(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Check if player is in a corp
	var corpID string
	err := database.DB.QueryRow(`
		SELECT corp_id FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&corpID)

	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"corp": nil})
		return
	}

	if err != nil {
		log.Printf("GetCorp: query failed: %v", err)
		errs.InternalError("Failed to get corp membership").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Get corp details
	var corp Corp
	err = database.DB.QueryRow(`
		SELECT id, name, tag, description, leader_id, level, wealth, max_members, created_at
		FROM corps WHERE id = $1
	`, corpID).Scan(&corp.ID, &corp.Name, &corp.Tag, &corp.Description,
		&corp.LeaderID, &corp.Level, &corp.Wealth, &corp.MaxMembers, &corp.CreatedAt)

	if err != nil {
		log.Printf("GetCorp: get corp details failed: %v", err)
		errs.InternalError("Failed to get corp details").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Get corp bonuses
	bonuses, err := services.GetCorpBonuses(corpID)
	if err == nil {
		corp.Bonuses = bonuses
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"corp": corp})
}

// CreateCorp handles POST /api/corp
// Creates a new corporation
func CreateCorp(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var body struct {
		Name        string `json:"name"`
		Tag         string `json:"tag"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errs.InvalidRequest("Invalid request body").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Validate name (3-30 chars)
	if len(body.Name) < 3 || len(body.Name) > 30 {
		errs.InvalidRequest("Corp name must be 3-30 characters").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Validate tag (3-5 chars uppercase)
	body.Tag = strings.ToUpper(strings.TrimSpace(body.Tag))
	if len(body.Tag) < 3 || len(body.Tag) > 5 {
		errs.InvalidRequest("Corp tag must be 3-5 characters").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Check player has Alliance Center Lv1+
	var allianceCenterLevel int
	err := database.DB.QueryRow(`
		SELECT b.level FROM buildings b
		JOIN building_types bt ON b.building_type = bt.id
		WHERE bt.name = 'alliance_center' AND b.planet_id IN (
			SELECT id FROM planets WHERE player_id = $1
		)
		LIMIT 1
	`, playerID).Scan(&allianceCenterLevel)

	if err == sql.ErrNoRows || allianceCenterLevel < 1 {
		errs.PrerequisiteNotMet("Alliance Center Level 1 required", "alliance_center", 0, 1).WriteJSON(w, http.StatusForbidden)
		return
	}

	if err != nil && err != sql.ErrNoRows {
		log.Printf("CreateCorp: check alliance center failed: %v", err)
		errs.InternalError("Failed to check prerequisites").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Check player not already in a corp
	var existingCorpID string
	err = database.DB.QueryRow(`
		SELECT corp_id FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&existingCorpID)

	if err != sql.ErrNoRows {
		errs.AlreadyExists("You are already in a corp").WriteJSON(w, http.StatusConflict)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("CreateCorp: tx begin failed: %v", err)
		errs.InternalError("Failed to create corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert into corps
	var corpID string
	err = tx.QueryRow(`
		INSERT INTO corps (name, tag, description, leader_id, level, wealth, max_members)
		VALUES ($1, $2, $3, $4, 1, 0, 50)
		RETURNING id
	`, body.Name, body.Tag, body.Description, playerID).Scan(&corpID)

	if err != nil {
		log.Printf("CreateCorp: insert corps failed: %v", err)
		errs.InternalError("Failed to create corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Insert into corp_members with leader role
	_, err = tx.Exec(`
		INSERT INTO corp_members (corp_id, player_id, role, contribution_points, daily_contribution)
		VALUES ($1, $2, 'leader', 0, 0)
	`, corpID, playerID)

	if err != nil {
		log.Printf("CreateCorp: insert corp_members failed: %v", err)
		errs.InternalError("Failed to add leader to corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("CreateCorp: commit failed: %v", err)
		errs.InternalError("Failed to create corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Return created corp
	var corp Corp
	corp.ID = corpID
	corp.Name = body.Name
	corp.Tag = body.Tag
	corp.Description = body.Description
	corp.LeaderID = playerID
	corp.Level = 1
	corp.Wealth = 0
	corp.MaxMembers = 50
	corp.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(corp)
}

// JoinCorp handles POST /api/corp/join
// Join an existing corp
func JoinCorp(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var body struct {
		CorpID string `json:"corp_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errs.InvalidRequest("Invalid request body").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Check player has Alliance Center Lv1+
	var allianceCenterLevel int
	err := database.DB.QueryRow(`
		SELECT b.level FROM buildings b
		JOIN building_types bt ON b.building_type = bt.id
		WHERE bt.name = 'alliance_center' AND b.planet_id IN (
			SELECT id FROM planets WHERE player_id = $1
		)
		LIMIT 1
	`, playerID).Scan(&allianceCenterLevel)

	if err == sql.ErrNoRows || allianceCenterLevel < 1 {
		errs.PrerequisiteNotMet("Alliance Center Level 1 required", "alliance_center", 0, 1).WriteJSON(w, http.StatusForbidden)
		return
	}

	if err != nil && err != sql.ErrNoRows {
		log.Printf("JoinCorp: check alliance center failed: %v", err)
		errs.InternalError("Failed to check prerequisites").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Check player not already in a corp
	var existingCorpID string
	err = database.DB.QueryRow(`
		SELECT corp_id FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&existingCorpID)

	if err != sql.ErrNoRows {
		errs.AlreadyExists("You are already in a corp").WriteJSON(w, http.StatusConflict)
		return
	}

	// Check corp exists and has room
	var maxMembers int
	var memberCount int
	err = database.DB.QueryRow(`
		SELECT c.max_members, COUNT(cm.player_id)
		FROM corps c
		LEFT JOIN corp_members cm ON cm.corp_id = c.id
		WHERE c.id = $1
		GROUP BY c.id, c.max_members
	`, body.CorpID).Scan(&maxMembers, &memberCount)

	if err == sql.ErrNoRows {
		errs.NotFound("Corp not found").WriteJSON(w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("JoinCorp: check corp failed: %v", err)
		errs.InternalError("Failed to check corp capacity").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	if memberCount >= maxMembers {
		errs.MaxCountReached("Corp is full", maxMembers, memberCount).WriteJSON(w, http.StatusConflict)
		return
	}

	// Insert into corp_members
	_, err = database.DB.Exec(`
		INSERT INTO corp_members (corp_id, player_id, role, contribution_points, daily_contribution)
		VALUES ($1, $2, 'member', 0, 0)
	`, body.CorpID, playerID)

	if err != nil {
		log.Printf("JoinCorp: insert failed: %v", err)
		errs.InternalError("Failed to join corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// LeaveCorp handles POST /api/corp/leave
// Leave current corp
func LeaveCorp(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Check if player is in a corp
	var corpID, role string
	err := database.DB.QueryRow(`
		SELECT corp_id, role FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&corpID, &role)

	if err == sql.ErrNoRows {
		errs.NotFound("You are not in a corp").WriteJSON(w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("LeaveCorp: query failed: %v", err)
		errs.InternalError("Failed to check corp membership").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Leader cannot leave
	if role == "leader" {
		errs.InvalidRequest("Leader cannot leave, transfer leadership or disband corp").WriteJSON(w, http.StatusForbidden)
		return
	}

	// Delete from corp_members
	_, err = database.DB.Exec(`
		DELETE FROM corp_members WHERE player_id = $1
	`, playerID)

	if err != nil {
		log.Printf("LeaveCorp: delete failed: %v", err)
		errs.InternalError("Failed to leave corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ListCorpMembers handles GET /api/corp/members
// Returns all members of player's corp
func ListCorpMembers(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Get player's corp
	var corpID string
	err := database.DB.QueryRow(`
		SELECT corp_id FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&corpID)

	if err == sql.ErrNoRows {
		errs.NotFound("You are not in a corp").WriteJSON(w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("ListCorpMembers: get corp failed: %v", err)
		errs.InternalError("Failed to get corp").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Get all members
	rows, err := database.DB.Query(`
		SELECT cm.player_id, p.username, cm.role, cm.contribution_points, cm.daily_contribution, cm.joined_at
		FROM corp_members cm
		JOIN players p ON p.id = cm.player_id
		WHERE cm.corp_id = $1
		ORDER BY cm.role DESC, cm.contribution_points DESC
	`, corpID)

	if err != nil {
		log.Printf("ListCorpMembers: query failed: %v", err)
		errs.InternalError("Failed to get members").WriteJSON(w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	members := []CorpMember{}
	for rows.Next() {
		var m CorpMember
		if err := rows.Scan(&m.PlayerID, &m.PlayerName, &m.Role, &m.ContributionPoints, &m.DailyContribution, &m.JoinedAt); err != nil {
			log.Printf("ListCorpMembers: scan failed: %v", err)
			continue
		}
		members = append(members, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// DonateResources handles POST /api/corp/donate
// Donate resources to corp
func DonateResources(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var body struct {
		Metal int64 `json:"metal"`
		He3   int64 `json:"he3"`
		Gold  int64 `json:"gold"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errs.InvalidRequest("Invalid request body").WriteJSON(w, http.StatusBadRequest)
		return
	}

	if body.Metal < 0 || body.He3 < 0 || body.Gold < 0 {
		errs.InvalidRequest("Negative donations not allowed").WriteJSON(w, http.StatusBadRequest)
		return
	}

	if body.Metal == 0 && body.He3 == 0 && body.Gold == 0 {
		errs.InvalidRequest("Must donate at least one resource").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Get player's corp and daily contribution
	var corpID string
	var dailyContribution int
	err := database.DB.QueryRow(`
		SELECT corp_id, daily_contribution FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&corpID, &dailyContribution)

	if err == sql.ErrNoRows {
		errs.NotFound("You are not in a corp").WriteJSON(w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("DonateResources: get corp failed: %v", err)
		errs.InternalError("Failed to check corp membership").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Calculate contribution points (1 pt per 10,000 resources)
	totalResources := body.Metal + body.He3 + body.Gold
	contributionPts := int(totalResources / 10000)

	// Check daily contribution limit (200 pts max)
	if dailyContribution+contributionPts > 200 {
		remaining := 200 - dailyContribution
		errs.MaxCountReached("Daily contribution limit reached", 200, dailyContribution).WriteJSON(w, http.StatusConflict)
		log.Printf("DonateResources: daily limit. Current: %d, Attempting: %d, Max: 200, Remaining: %d", dailyContribution, contributionPts, remaining)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DonateResources: tx begin failed: %v", err)
		errs.InternalError("Failed to process donation").WriteJSON(w, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Deduct resources from player's homeworld
	result, err := tx.Exec(`
		UPDATE resources
		SET metal = metal - $1, he3 = he3 - $2, gold = gold - $3, updated_at = now()
		WHERE planet_id = (
			SELECT id FROM planets WHERE player_id = $4 AND is_homeworld = true LIMIT 1
		) AND metal >= $1 AND he3 >= $2 AND gold >= $3
	`, body.Metal, body.He3, body.Gold, playerID)

	if err != nil {
		log.Printf("DonateResources: deduct resources failed: %v", err)
		errs.InternalError("Failed to deduct resources").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		errs.InsufficientResources("Not enough resources", map[string]int64{
			"metal": body.Metal,
			"he3":   body.He3,
			"gold":  body.Gold,
		}, nil).WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Add to corp wealth
	_, err = tx.Exec(`
		UPDATE corps
		SET wealth = wealth + $1, updated_at = now()
		WHERE id = $2
	`, contributionPts, corpID)

	if err != nil {
		log.Printf("DonateResources: add wealth failed: %v", err)
		errs.InternalError("Failed to add corp wealth").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Update corp_members contribution
	_, err = tx.Exec(`
		UPDATE corp_members
		SET contribution_points = contribution_points + $1,
		    daily_contribution = daily_contribution + $1,
		    updated_at = now()
		WHERE player_id = $2
	`, contributionPts, playerID)

	if err != nil {
		log.Printf("DonateResources: update member failed: %v", err)
		errs.InternalError("Failed to update contribution").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Insert donation record
	_, err = tx.Exec(`
		INSERT INTO corp_donations (corp_id, player_id, metal, he3, gold, contribution_points)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, corpID, playerID, body.Metal, body.He3, body.Gold, contributionPts)

	if err != nil {
		log.Printf("DonateResources: insert donation failed: %v", err)
		// Non-critical, continue
	}

	// Recalculate corp level
	var wealth int64
	err = tx.QueryRow(`SELECT wealth FROM corps WHERE id = $1`, corpID).Scan(&wealth)
	if err == nil {
		newLevel := services.CalculateCorpLevel(wealth)
		_, err = tx.Exec(`UPDATE corps SET level = $1 WHERE id = $2`, newLevel, corpID)
		if err != nil {
			log.Printf("DonateResources: update corp level failed: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DonateResources: commit failed: %v", err)
		errs.InternalError("Failed to process donation").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	newDaily := dailyContribution + contributionPts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "ok",
		"contribution_pts":   contributionPts,
		"daily_contribution": newDaily,
		"daily_max":          200,
		"daily_remaining":    200 - newDaily,
	})
}

// UpdateMemberRole handles PUT /api/corp/members/{id}/role
// Change member role (leader only)
func UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	targetPlayerID := r.PathValue("id")

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errs.InvalidRequest("Invalid request body").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Validate role
	if body.Role != "officer" && body.Role != "member" {
		errs.InvalidRequest("Role must be 'officer' or 'member'").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Check if player is leader
	var corpID, role string
	err := database.DB.QueryRow(`
		SELECT corp_id, role FROM corp_members WHERE player_id = $1
	`, playerID).Scan(&corpID, &role)

	if err == sql.ErrNoRows {
		errs.NotFound("You are not in a corp").WriteJSON(w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("UpdateMemberRole: get player role failed: %v", err)
		errs.InternalError("Failed to check permissions").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	if role != "leader" {
		errs.InvalidRequest("Only leader can change roles").WriteJSON(w, http.StatusForbidden)
		return
	}

	// Cannot change own role
	if targetPlayerID == playerID {
		errs.InvalidRequest("Cannot change own role").WriteJSON(w, http.StatusForbidden)
		return
	}

	// Update target player role
	result, err := database.DB.Exec(`
		UPDATE corp_members
		SET role = $1, updated_at = now()
		WHERE player_id = $2 AND corp_id = $3 AND role != 'leader'
	`, body.Role, targetPlayerID, corpID)

	if err != nil {
		log.Printf("UpdateMemberRole: update failed: %v", err)
		errs.InternalError("Failed to update role").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		errs.NotFound("Member not found or is leader").WriteJSON(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// SearchCorps handles GET /api/corp/search?q=searchterm
// Search corps by name or tag
func SearchCorps(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		errs.InvalidRequest("Query parameter 'q' required").WriteJSON(w, http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.name, c.tag, c.description, c.level, c.wealth, c.max_members, c.created_at,
		       COUNT(cm.player_id) as member_count
		FROM corps c
		LEFT JOIN corp_members cm ON cm.corp_id = c.id
		WHERE c.name ILIKE $1 OR c.tag ILIKE $1
		GROUP BY c.id
		ORDER BY c.level DESC, c.wealth DESC
		LIMIT 50
	`, "%"+query+"%")

	if err != nil {
		log.Printf("SearchCorps: query failed: %v", err)
		errs.InternalError("Failed to search corps").WriteJSON(w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type searchResult struct {
		Corp        Corp `json:"corp"`
		MemberCount int  `json:"member_count"`
	}

	results := []searchResult{}
	for rows.Next() {
		var r searchResult
		if err := rows.Scan(&r.Corp.ID, &r.Corp.Name, &r.Corp.Tag, &r.Corp.Description,
			&r.Corp.Level, &r.Corp.Wealth, &r.Corp.MaxMembers, &r.Corp.CreatedAt,
			&r.MemberCount); err != nil {
			log.Printf("SearchCorps: scan failed: %v", err)
			continue
		}
		results = append(results, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// AttackRBP handles POST /api/corp/rbp/{id}/attack
// Attack an RBP planet
func AttackRBP(w http.ResponseWriter, r *http.Request) {
	attackerID := middleware.GetPlayerID(r)
	rbpPlanetID := r.PathValue("id")

	var body struct {
		FleetIDs []string `json:"fleet_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errs.InvalidRequest("Invalid request body").WriteJSON(w, http.StatusBadRequest)
		return
	}

	if len(body.FleetIDs) == 0 {
		errs.InvalidRequest("At least one fleet required").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Get player's corp
	var attackerCorpID string
	err := database.DB.QueryRow(`
		SELECT corp_id FROM corp_members WHERE player_id = $1
	`, attackerID).Scan(&attackerCorpID)

	if err == sql.ErrNoRows {
		errs.NotFound("You must be in a corp to attack RBPs").WriteJSON(w, http.StatusForbidden)
		return
	}

	if err != nil {
		log.Printf("AttackRBP: get attacker corp failed: %v", err)
		errs.InternalError("Failed to check corp membership").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Get RBP planet info
	var isRBP bool
	var rbpLevel int
	var controllingCorpID sql.NullString
	var protectionUntil sql.NullTime
	var rbpName string

	err = database.DB.QueryRow(`
		SELECT name, is_rbp, rbp_level, controlling_corp_id, protection_until
		FROM planets WHERE id = $1
	`, rbpPlanetID).Scan(&rbpName, &isRBP, &rbpLevel, &controllingCorpID, &protectionUntil)

	if err == sql.ErrNoRows {
		errs.NotFound("RBP planet not found").WriteJSON(w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("AttackRBP: get rbp planet failed: %v", err)
		errs.InternalError("Failed to get RBP planet").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	if !isRBP {
		errs.InvalidRequest("Target is not an RBP").WriteJSON(w, http.StatusBadRequest)
		return
	}

	// Check protection
	if protectionUntil.Valid && protectionUntil.Time.After(time.Now()) {
		remaining := int(time.Until(protectionUntil.Time).Seconds())
		errs.OnCooldown("RBP is under protection", "rbp_attack", protectionUntil.Time.Format(time.RFC3339), remaining).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check if corp can control more RBPs (if RBP is unowned or owned by different corp)
	if !controllingCorpID.Valid || controllingCorpID.String != attackerCorpID {
		canControl, err := services.CanCorpControlMoreRBPs(attackerCorpID)
		if err != nil {
			log.Printf("AttackRBP: check rbp capacity failed: %v", err)
			errs.InternalError("Failed to check corp capacity").WriteJSON(w, http.StatusInternalServerError)
			return
		}

		if !canControl {
			count, _ := services.GetCorpRBPCount(attackerCorpID)
			var level int
			database.DB.QueryRow(`SELECT level FROM corps WHERE id = $1`, attackerCorpID).Scan(&level)
			errs.MaxCountReached("Corp cannot control more RBPs (max = corp level)", level, count).WriteJSON(w, http.StatusConflict)
			return
		}
	}

	// Validate attacker fleets
	for _, fid := range body.FleetIDs {
		fleet, err := getOwnedFleet(fid, attackerID)
		if err != nil {
			errs.NotFound("Fleet not found: " + fid).WriteJSON(w, http.StatusNotFound)
			return
		}
		if fleet.Status != "stationed" {
			errs.InvalidRequest("Fleet not stationed: " + fid).WriteJSON(w, http.StatusConflict)
			return
		}
		stacks := getFleetStacks(fid)
		hasShips := false
		for _, s := range stacks {
			if s.ShipCount > 0 {
				hasShips = true
				break
			}
		}
		if !hasShips {
			errs.InvalidRequest("Fleet has no ships: " + fid).WriteJSON(w, http.StatusConflict)
			return
		}
	}

	// Load attacker fleets
	var attackerStacks []*combat.FleetStack
	for _, fid := range body.FleetIDs {
		attackerFleet, err := combat.LoadPlayerFleet(fid, attackerID)
		if err != nil {
			log.Printf("AttackRBP: load attacker fleet failed: %v", err)
			errs.InternalError("Failed to load attacker fleet").WriteJSON(w, http.StatusInternalServerError)
			return
		}
		attackerStacks = append(attackerStacks, attackerFleet.Stacks...)
	}

	// Create combined attacker fleet
	attackerFleet := &combat.Fleet{
		PlayerID:       attackerID,
		FleetID:        "rbp_attacker",
		CommanderBonus: nil,
		TechBonuses:    &services.TechBonuses{},
		Stacks:         attackerStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "attacker",
	}

	// Load attacker tech bonuses
	if attackerTechBonuses, err := services.GetPlayerTechBonuses(attackerID); err == nil {
		attackerFleet.TechBonuses = attackerTechBonuses
	}

	// Load RBP defenses (defense buildings only, no player fleets)
	defenderFleet, err := loadRBPDefenses(rbpPlanetID, rbpLevel)
	if err != nil {
		log.Printf("AttackRBP: load rbp defenses failed: %v", err)
		errs.InternalError("Failed to load RBP defenses").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Execute combat
	engine := combat.NewCombatEngine(time.Now().UnixNano())
	combatResult, err := engine.ExecuteCombat(attackerFleet, defenderFleet)
	if err != nil {
		log.Printf("AttackRBP: combat execution failed: %v", err)
		errs.InternalError("Combat execution failed").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Determine result
	var result string
	if combatResult.Winner == "attacker" {
		result = "attacker_win"
	} else if combatResult.Winner == "defender" {
		result = "defender_win"
	} else {
		result = "draw"
	}

	// Apply casualties to attacker fleets
	attackerShipsDestroyed := make(map[string]int)
	for _, stack := range attackerFleet.Stacks {
		destroyed := stack.ShipCount - stack.CurrentShips
		if destroyed > 0 {
			attackerShipsDestroyed[stack.ID] = destroyed
			database.DB.Exec(`UPDATE fleet_stacks SET ship_count = $1 WHERE id = $2`, stack.CurrentShips, stack.ID)
		}
	}

	// On attacker victory, transfer RBP control
	if result == "attacker_win" {
		_, err = database.DB.Exec(`
			UPDATE planets
			SET controlling_corp_id = $1, protection_until = now() + interval '24 hours', updated_at = now()
			WHERE id = $2
		`, attackerCorpID, rbpPlanetID)

		if err != nil {
			log.Printf("AttackRBP: update planet control failed: %v", err)
		}

		// Insert rbp_attacks record
		_, err = database.DB.Exec(`
			INSERT INTO rbp_attacks (rbp_planet_id, attacker_corp_id, defender_corp_id, result, created_at)
			VALUES ($1, $2, $3, $4, now())
		`, rbpPlanetID, attackerCorpID, controllingCorpID, result)

		if err != nil {
			log.Printf("AttackRBP: insert rbp_attacks failed: %v", err)
		}
	}

	// Create combat report
	reportID := createRBPReport(attackerID, attackerCorpID, rbpPlanetID, result, combatResult.TotalRounds, int64(combatResult.AttackerCasualties))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"report_id":        reportID,
		"result":           result,
		"total_rounds":     combatResult.TotalRounds,
		"attacker_losses":  battleLosses{ShipsDestroyed: attackerShipsDestroyed, He3Consumed: int64(combatResult.AttackerCasualties)},
	})
}

// loadRBPDefenses creates a defensive fleet from RBP defense buildings
func loadRBPDefenses(planetID string, rbpLevel int) (*combat.Fleet, error) {
	// Load defense buildings
	defenseStacks, err := loadDefenseBuildings(planetID, nil)
	if err != nil {
		return nil, err
	}

	// RBP defenses scale with RBP level
	// Each building's stats are multiplied by (1 + rbpLevel * 0.1)
	multiplier := 1.0 + float64(rbpLevel)*0.1
	for _, stack := range defenseStacks {
		stack.BaseAttack = int(float64(stack.BaseAttack) * multiplier)
		stack.BaseDefense = int(float64(stack.BaseDefense) * multiplier)
		stack.BaseShield = int(float64(stack.BaseShield) * multiplier)
		stack.BaseStructure = int(float64(stack.BaseStructure) * multiplier)
		stack.CurrentShield = stack.BaseShield
		stack.CurrentStructure = stack.BaseStructure
	}

	fleet := &combat.Fleet{
		PlayerID:       "rbp_defense",
		FleetID:        "rbp_defender",
		CommanderBonus: nil,
		TechBonuses:    &services.TechBonuses{},
		Stacks:         defenseStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "defender",
	}

	return fleet, nil
}

// createRBPReport creates a combat report for RBP attack
func createRBPReport(attackerID, attackerCorpID, rbpPlanetID, result string, totalRounds int, he3Consumed int64) string {
	var reportID string
	database.DB.QueryRow(`
		INSERT INTO combat_reports (attacker_id, defender_id, combat_type, result, total_rounds, he3_consumed, created_at)
		VALUES ($1, $2, 'rbp_attack', $3, $4, $5, now())
		RETURNING id
	`, attackerID, rbpPlanetID, result, totalRounds, he3Consumed).Scan(&reportID)

	return reportID
}

// corpBrief represents minimal corp info for galaxy map
type corpBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

// galaxyZone represents a zone in the galaxy map
type galaxyZone struct {
	X               int        `json:"x"`
	Y               int        `json:"y"`
	RBPPlanetID     string     `json:"rbp_planet_id"`
	RBPName         string     `json:"rbp_name"`
	RBPLevel        int        `json:"rbp_level"`
	ControllingCorp *corpBrief `json:"controlling_corp"`
	ProtectionUntil *time.Time `json:"protection_until"`
}

// GetGalaxyMap handles GET /api/galaxy/map
// Returns 7x7 galaxy map with RBPs
func GetGalaxyMap(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT gz.zone_x, gz.zone_y, gz.rbp_planet_id,
		       p.name as rbp_name, p.rbp_level, p.controlling_corp_id, p.protection_until,
		       c.name as corp_name, c.tag as corp_tag
		FROM galaxy_zones gz
		JOIN planets p ON p.id = gz.rbp_planet_id
		LEFT JOIN corps c ON c.id = p.controlling_corp_id
		ORDER BY gz.zone_y, gz.zone_x
	`)

	if err != nil {
		log.Printf("GetGalaxyMap: query failed: %v", err)
		errs.InternalError("Failed to get galaxy map").WriteJSON(w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	zones := []galaxyZone{}
	for rows.Next() {
		var z galaxyZone
		var corpID, corpName, corpTag sql.NullString
		var protectionUntil sql.NullTime

		if err := rows.Scan(&z.X, &z.Y, &z.RBPPlanetID,
			&z.RBPName, &z.RBPLevel, &corpID, &protectionUntil,
			&corpName, &corpTag); err != nil {
			log.Printf("GetGalaxyMap: scan failed: %v", err)
			continue
		}

		if protectionUntil.Valid {
			z.ProtectionUntil = &protectionUntil.Time
		}

		if corpID.Valid {
			z.ControllingCorp = &corpBrief{
				ID:   corpID.String,
				Name: corpName.String,
				Tag:  corpTag.String,
			}
		}

		zones = append(zones, z)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}
