package models

import (
	"encoding/json"
	"time"
)

type QuestType struct {
	ID                 int             `json:"id"`
	QuestKey           string          `json:"quest_key"`
	Category           string          `json:"category"`
	DisplayName        string          `json:"display_name"`
	Description        string          `json:"description"`
	RequirementType    string          `json:"requirement_type"`
	RequirementTarget  string          `json:"requirement_target"`
	RequirementValue   int             `json:"requirement_value"`
	ChainOrder         *int            `json:"chain_order"`
	PrerequisiteQuestID *int           `json:"prerequisite_quest_id"`
	RewardMetal        int64           `json:"reward_metal"`
	RewardHe3          int64           `json:"reward_he3"`
	RewardGold         int64           `json:"reward_gold"`
	RewardItemJSON     json.RawMessage `json:"reward_item_json"`
	Phase              int             `json:"phase"`
	IsActive           bool            `json:"is_active"`
}

type PlayerQuest struct {
	ID            string     `json:"id"`
	PlayerID      string     `json:"player_id"`
	QuestTypeID   int        `json:"quest_type_id"`
	Status        string     `json:"status"`
	ProgressValue int        `json:"progress_value"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	ClaimedAt     *time.Time `json:"claimed_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type DailyQuestProgress struct {
	ID                    string          `json:"id"`
	PlayerID              string          `json:"player_id"`
	QuestDate             string          `json:"quest_date"`
	DailyPoints           int             `json:"daily_points"`
	QuestsCompletedJSON   json.RawMessage `json:"quests_completed_json"`
	TierRewardsClaimedJSON json.RawMessage `json:"tier_rewards_claimed_json"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

// PlayerQuestWithType combines player quest progress with quest type details for API responses.
type PlayerQuestWithType struct {
	ID               string          `json:"id"`
	QuestKey         string          `json:"quest_key"`
	DisplayName      string          `json:"display_name"`
	Description      string          `json:"description"`
	Category         string          `json:"category"`
	Status           string          `json:"status"`
	ProgressValue    int             `json:"progress_value"`
	RequirementValue int             `json:"requirement_value"`
	ChainOrder       *int            `json:"chain_order"`
	RewardMetal      int64           `json:"reward_metal"`
	RewardHe3        int64           `json:"reward_he3"`
	RewardGold       int64           `json:"reward_gold"`
	RewardItemJSON   json.RawMessage `json:"reward_item_json"`
}
