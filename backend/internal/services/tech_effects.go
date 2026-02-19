package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cryptomines-online/backend/internal/database"
)

// TechBonuses represents all active technology bonuses for a player.
// These bonuses are calculated from completed research levels.
type TechBonuses struct {
	// Production bonuses (percent)
	MetalOutput  float64 `json:"metal_output"`
	He3Output    float64 `json:"he3_output"`
	GoldOutput   float64 `json:"gold_output"`

	// Construction bonuses
	BuildSpeed           float64 `json:"build_speed"`             // percent bonus to building construction speed
	BuildCostReduction   float64 `json:"build_cost_reduction"`    // percent reduction in building costs
	ConstructionSlots    int     `json:"construction_slots"`      // additional construction slots
	ShipBuildSpeed       float64 `json:"ship_build_speed"`        // percent bonus to ship construction speed
	ShipBuildCostReduction float64 `json:"ship_build_cost_reduction"` // percent reduction in ship costs
	ShipProductionSlots  int     `json:"ship_production_slots"`   // additional shipbuilding slots

	// Capacity bonuses
	WarehouseCapacity int64 `json:"warehouse_capacity"` // flat bonus to warehouse storage

	// Ship repair
	ShipRepairPercent float64 `json:"ship_repair_percent"` // percent of max HP repaired

	// Combat bonuses - Ballistics
	BallisticDamage      float64 `json:"ballistic_damage"`
	BallisticCritRate    float64 `json:"ballistic_crit_rate"`
	BallisticCritDamage  float64 `json:"ballistic_crit_damage"`
	BallisticHitRate     float64 `json:"ballistic_hit_rate"`
	WeaponSpaceReduction float64 `json:"weapon_space_reduction"`
	ShieldBypass         float64 `json:"shield_bypass"`

	// Combat bonuses - Directional
	DirectionalDamage   float64 `json:"directional_damage"`
	DirectionalCritRate float64 `json:"directional_crit_rate"`
	DirectionalAccuracy float64 `json:"directional_accuracy"`
	SteeringPower       float64 `json:"steering_power"`

	// Combat bonuses - Missile
	MissileDamage     float64 `json:"missile_damage"`
	MissileHitRate    float64 `json:"missile_hit_rate"`
	InterceptReduction float64 `json:"intercept_reduction"`

	// Combat bonuses - Fighter
	FighterDamage   float64 `json:"fighter_damage"`
	FighterHitRate  float64 `json:"fighter_hit_rate"`
	FuelOptimization float64 `json:"fuel_optimization"` // He3 cost reduction

	// Defense bonuses
	BaseShield     float64 `json:"base_shield"`     // percent bonus
	BaseStructure  float64 `json:"base_structure"`  // percent bonus
	BaseAgility    float64 `json:"base_agility"`    // percent bonus
	BaseDefense    float64 `json:"base_defense"`    // percent bonus
	BaseStability  float64 `json:"base_stability"`  // percent bonus

	// Planetary defense bonuses
	DefenseCostReduction   float64 `json:"defense_cost_reduction"`
	DefenseBuildSpeed      float64 `json:"defense_build_speed"`
	DefenseValue           float64 `json:"defense_value"`
	EmplacementAttack      float64 `json:"emplacement_attack"`
	MaxDefenseStructures   float64 `json:"max_defense_structures"`
	MaxThorCannon          int     `json:"max_thor_cannon"`

	// Combat bonuses - Scatter
	ScatterDamage       float64 `json:"scatter_damage"`
	ScatterRate         float64 `json:"scatter_rate"`
	ScatterWeaponChance float64 `json:"scatter_weapon_chance"`
	ScatterAll          float64 `json:"scatter_all"`
	ScatterBonus        float64 `json:"scatter_bonus"`

	// Combat bonuses - Shield
	ShieldDamageReduction float64 `json:"shield_damage_reduction"`
	ShieldPenResist       float64 `json:"shield_pen_resist"`
	ShieldFlatReduction   float64 `json:"shield_flat_reduction"`

	// Combat bonuses - Absorb & Reflect
	AbsorbNoHe3          float64 `json:"absorb_no_he3"`
	ReflectDamage         float64 `json:"reflect_damage"`
	AbsorbNoHe3Structure  float64 `json:"absorb_no_he3_structure"`
	ReflectStructureDamage float64 `json:"reflect_structure_damage"`

	// Combat bonuses - Structure defense
	StructureDamageReduction float64 `json:"structure_damage_reduction"`
	StructurePenResist       float64 `json:"structure_pen_resist"`
	StructureFlatReduction   float64 `json:"structure_flat_reduction"`

	// Combat bonuses - Piercing
	PiercingDamage      float64 `json:"piercing_damage"`
	PiercingDamageBonus float64 `json:"piercing_damage_bonus"`

	// Combat bonuses - Evasion & Movement
	EnemyHitReduction    float64 `json:"enemy_hit_reduction"`
	MovementCritBonus    float64 `json:"movement_crit_bonus"`
	EnemyAttackReduction float64 `json:"enemy_attack_reduction"`

	// Planetary defense - range and movement
	DefenseRange    int `json:"defense_range"`    // additional range for Particle Cannon and Anti-Aircraft Gun
	DefenseMovement int `json:"defense_movement"` // flat ship movement bonus when defending own planet

	// Combat bonuses - Armor
	ArmorBonusNeutral float64 `json:"armor_bonus_neutral"`
	ArmorBonusLight   float64 `json:"armor_bonus_light"`
	ArmorBonusRegen   float64 `json:"armor_bonus_regen"`

	// Combat bonuses - Shield penetration
	ShieldPenChance     float64 `json:"shield_pen_chance"`
	ShieldPenLightBonus float64 `json:"shield_pen_light_bonus"`

	// Combat bonuses - Range damage
	RangeDamageRanges    []float64 `json:"range_damage_ranges,omitempty"`
	RangeDamageCritRate  float64   `json:"range_damage_crit_rate"`
	RangeDamageCritDmg   float64   `json:"range_damage_crit_damage"`

	// Combat bonuses - Composite directional
	HitAndShieldPenHit        float64 `json:"hit_and_shield_pen_hit"`
	HitAndShieldPenShieldPen  float64 `json:"hit_and_shield_pen_shield_pen"`
	SteeringReductionChance   float64 `json:"steering_reduction_chance"`
	SteeringReductionPerLvl   float64 `json:"steering_reduction_per_level"`
	PiercingCriticalEnabled   bool    `json:"piercing_critical_enabled"`
	BaseDamageBonus           float64 `json:"base_damage_bonus"`
	BonusAccuracy             float64 `json:"bonus_accuracy"`
	IgnoreAgilityPerLevel     float64 `json:"ignore_agility_per_level"`
	IgnoreAgilityChance       float64 `json:"ignore_agility_chance"`
	IgnoreDefensePerLevel     float64 `json:"ignore_defense_per_level"`
	IgnoreDefenseMovReduction float64 `json:"ignore_defense_mov_reduction"`
	IgnoreDefenseChance       float64 `json:"ignore_defense_chance"`

	// Combat bonuses - Composite missile
	MissileShieldPenChance      float64 `json:"missile_shield_pen_chance"`
	ScatterVsHighStructure      float64 `json:"scatter_vs_high_structure"`
	ScatterVsHighStructureCrit  float64 `json:"scatter_vs_high_structure_crit"`
	KnockbackChance             float64 `json:"knockback_chance"`
	KnockbackDistance           int     `json:"knockback_distance"`

	// Combat bonuses - Fighter advanced
	BaseAttack             float64 `json:"base_attack"`
	He3AndShieldDamageHe3  float64 `json:"he3_and_shield_damage_he3"`
	ShieldDamageBonus      float64 `json:"shield_damage_bonus"`
	ShipShieldDamage       float64 `json:"ship_shield_damage"`
	ReloadChance           float64 `json:"reload_chance"`
	ReloadReduction        int     `json:"reload_reduction"`
	LongRangeAttack        float64 `json:"long_range_attack"`
	LongRangeBonusDblDmg   float64 `json:"long_range_bonus_double_damage"`
	LongRangeBonusCritRate float64 `json:"long_range_bonus_crit_rate"`
	SwarmChancePerLevel    float64 `json:"swarm_chance_per_level"`
	SwarmAttackBonusBase   float64 `json:"swarm_attack_bonus_base"`
	SwarmAttackBonusMax    float64 `json:"swarm_attack_bonus_max"`
	SwarmHe3CostBase       float64 `json:"swarm_he3_cost_base"`
	SwarmHe3CostMax        float64 `json:"swarm_he3_cost_max"`

	// Combat bonuses - Multi-bonus aggregated
	MultiBonusSwarm        float64 `json:"multi_bonus_swarm"`
	MultiBonusAttack       float64 `json:"multi_bonus_attack"`
	MultiBonusHe3          float64 `json:"multi_bonus_he3"`
	MultiBonusIntercept    float64 `json:"multi_bonus_intercept"`
	MultiBonusUnshielded   float64 `json:"multi_bonus_unshielded"`
	MultiBonusCritDamage   float64 `json:"multi_bonus_crit_damage"`
	MultiBonusWeaponSpace  float64 `json:"multi_bonus_weapon_space"`
	MultiBonusReload       float64 `json:"multi_bonus_reload"`
	MultiBonusShieldDmg    float64 `json:"multi_bonus_shield_damage"`
	MultiBonusAttackRange  float64 `json:"multi_bonus_attack_range"`
	MultiBonusInterception float64 `json:"multi_bonus_interception"`

	// Defense bonuses - Shield/Structure restore
	ShieldRestore             float64 `json:"shield_restore"`
	ShieldRestoreIntercept    float64 `json:"shield_restore_intercept"`
	StructureRestore          float64 `json:"structure_restore"`
	AbsorbDouble              float64 `json:"absorb_double"`
	AbsorbDoubleCollateral    float64 `json:"absorb_double_collateral"`
	AbsorbDoubleStructure     float64 `json:"absorb_double_structure"`
	AbsorbDoubleStructureColl float64 `json:"absorb_double_structure_collateral"`
}

// techEffectData represents the parsed effects_json structure from tech_types.
type techEffectData struct {
	Type     string  `json:"type"`
	PerLevel float64 `json:"per_level"`
	Unit     string  `json:"unit"`
	Flat     float64 `json:"flat"`
}

// GetPlayerTechBonuses calculates and returns all active technology bonuses for a player.
// It queries all completed research and aggregates the effects.
func GetPlayerTechBonuses(playerID string) (*TechBonuses, error) {
	bonuses := &TechBonuses{
		ConstructionSlots: 0,
		ShipProductionSlots: 0,
		MaxThorCannon: 0,
	}

	// Query all completed techs with their effects
	rows, err := database.DB.Query(`
		SELECT t.level, tt.effects_json, tt.name
		FROM technologies t
		JOIN tech_types tt ON t.tech_type = tt.id
		WHERE t.player_id = $1 AND t.level > 0
		ORDER BY tt.name
	`, playerID)
	if err != nil {
		log.Printf("Failed to query player research: %v", err)
		return nil, fmt.Errorf("failed to query player research: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var level int
		var effectsJSON []byte
		var techName string

		if err := rows.Scan(&level, &effectsJSON, &techName); err != nil {
			log.Printf("Failed to scan tech row: %v", err)
			continue
		}

		// Parse effects JSON
		var effects techEffectData
		if err := json.Unmarshal(effectsJSON, &effects); err != nil {
			log.Printf("Failed to parse effects for tech %s: %v", techName, err)
			continue
		}

		// Apply effect based on type
		applyTechEffect(bonuses, &effects, effectsJSON, level, techName)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating tech rows: %v", err)
		return nil, fmt.Errorf("error iterating tech rows: %w", err)
	}

	return bonuses, nil
}

// applyTechEffect applies a single tech effect to the bonuses struct.
func applyTechEffect(bonuses *TechBonuses, effect *techEffectData, rawJSON []byte, level int, techName string) {
	switch effect.Type {
	// Production
	case "metal_output":
		bonuses.MetalOutput += effect.PerLevel * float64(level)
	case "he3_output":
		bonuses.He3Output += effect.PerLevel * float64(level)
	case "gold_output":
		bonuses.GoldOutput += effect.PerLevel * float64(level)

	// Construction
	case "build_speed":
		bonuses.BuildSpeed += effect.PerLevel * float64(level)
	case "build_cost_reduction":
		bonuses.BuildCostReduction += effect.PerLevel * float64(level)
	case "construction_slots":
		bonuses.ConstructionSlots += int(effect.PerLevel * float64(level))
	case "ship_build_speed":
		bonuses.ShipBuildSpeed += effect.PerLevel * float64(level)
	case "ship_build_cost_reduction":
		bonuses.ShipBuildCostReduction += effect.PerLevel * float64(level)
	case "ship_production_slots":
		bonuses.ShipProductionSlots += int(effect.PerLevel * float64(level))

	// Capacity
	case "warehouse_capacity":
		bonuses.WarehouseCapacity += int64(effect.PerLevel * float64(level))

	// Ship repair
	case "ship_repair_percent":
		bonuses.ShipRepairPercent += effect.PerLevel * float64(level)

	// Ballistics
	case "ballistic_damage":
		bonuses.BallisticDamage += effect.PerLevel * float64(level)
	case "ballistic_crit_rate":
		bonuses.BallisticCritRate += effect.PerLevel * float64(level)
	case "ballistic_crit_damage":
		bonuses.BallisticCritDamage += effect.PerLevel * float64(level)
	case "ballistic_hit_rate":
		bonuses.BallisticHitRate += effect.PerLevel * float64(level)
	case "weapon_space_reduction":
		bonuses.WeaponSpaceReduction += effect.PerLevel * float64(level)
	case "shield_bypass":
		bonuses.ShieldBypass += effect.Flat

	// Directional
	case "directional_damage":
		bonuses.DirectionalDamage += effect.PerLevel * float64(level)
	case "directional_crit_rate":
		bonuses.DirectionalCritRate += effect.PerLevel * float64(level)
	case "directional_accuracy":
		bonuses.DirectionalAccuracy += effect.PerLevel * float64(level)
	case "steering_power":
		bonuses.SteeringPower += effect.PerLevel * float64(level)

	// Missile
	case "missile_damage":
		bonuses.MissileDamage += effect.PerLevel * float64(level)
	case "missile_hit_rate":
		bonuses.MissileHitRate += effect.PerLevel * float64(level)
	case "intercept_reduction":
		bonuses.InterceptReduction += effect.PerLevel * float64(level)

	// Fighter
	case "fighter_damage":
		bonuses.FighterDamage += effect.PerLevel * float64(level)
	case "fighter_hit_rate":
		bonuses.FighterHitRate += effect.PerLevel * float64(level)
	case "he3_cost_reduction":
		bonuses.FuelOptimization += effect.PerLevel * float64(level)

	// Defense stats
	case "base_defense_stats":
		var baseDefenseEffect struct {
			Shield    float64 `json:"shield"`
			Structure float64 `json:"structure"`
			Agility   float64 `json:"agility"`
			Defense   float64 `json:"defense"`
			Stability float64 `json:"stability"`
		}
		if err := json.Unmarshal(rawJSON, &baseDefenseEffect); err != nil {
			log.Printf("Failed to parse base_defense_stats for tech %s: %v", techName, err)
			break
		}
		bonuses.BaseShield += baseDefenseEffect.Shield * float64(level)
		bonuses.BaseStructure += baseDefenseEffect.Structure * float64(level)
		bonuses.BaseAgility += baseDefenseEffect.Agility * float64(level)
		bonuses.BaseDefense += baseDefenseEffect.Defense * float64(level)
		bonuses.BaseStability += baseDefenseEffect.Stability * float64(level)
	case "base_shield":
		bonuses.BaseShield += effect.PerLevel * float64(level)
	case "base_structure":
		bonuses.BaseStructure += effect.PerLevel * float64(level)
	case "base_agility":
		bonuses.BaseAgility += effect.PerLevel * float64(level)
	case "base_defense":
		bonuses.BaseDefense += effect.PerLevel * float64(level)
	case "base_stability":
		bonuses.BaseStability += effect.PerLevel * float64(level)

	// Planetary defense
	case "defense_cost_reduction":
		bonuses.DefenseCostReduction += effect.PerLevel * float64(level)
	case "defense_build_speed":
		bonuses.DefenseBuildSpeed += effect.PerLevel * float64(level)
	case "defense_value":
		bonuses.DefenseValue += effect.PerLevel * float64(level)
	case "emplacement_attack":
		bonuses.EmplacementAttack += effect.PerLevel * float64(level)
	case "max_defense_structures":
		bonuses.MaxDefenseStructures += effect.PerLevel * float64(level)
	case "max_thor_cannon":
		bonuses.MaxThorCannon += int(effect.Flat)
	case "defense_range":
		bonuses.DefenseRange += int(effect.PerLevel * float64(level))
	case "defense_movement":
		bonuses.DefenseMovement += int(effect.PerLevel * float64(level))

	// Scatter
	case "scatter_damage":
		bonuses.ScatterDamage += effect.PerLevel * float64(level)
	case "scatter_rate":
		bonuses.ScatterRate += effect.PerLevel * float64(level)
	case "scatter_weapon_chance":
		bonuses.ScatterWeaponChance += effect.PerLevel * float64(level)
	case "scatter_all":
		bonuses.ScatterAll += effect.PerLevel * float64(level)
	case "scatter_bonus":
		bonuses.ScatterBonus += effect.PerLevel * float64(level)

	// Shield
	case "shield_damage_reduction":
		bonuses.ShieldDamageReduction += effect.PerLevel * float64(level)
	case "shield_pen_resist":
		bonuses.ShieldPenResist += effect.PerLevel * float64(level)
	case "shield_flat_reduction":
		bonuses.ShieldFlatReduction += effect.PerLevel * float64(level)

	// Absorb & Reflect
	case "absorb_no_he3":
		bonuses.AbsorbNoHe3 += effect.PerLevel * float64(level)
	case "reflect_damage":
		bonuses.ReflectDamage += effect.PerLevel * float64(level)
	case "absorb_no_he3_structure":
		bonuses.AbsorbNoHe3Structure += effect.PerLevel * float64(level)
	case "reflect_structure_damage":
		bonuses.ReflectStructureDamage += effect.PerLevel * float64(level)

	// Structure defense
	case "structure_damage_reduction":
		bonuses.StructureDamageReduction += effect.PerLevel * float64(level)
	case "structure_pen_resist":
		bonuses.StructurePenResist += effect.PerLevel * float64(level)
	case "structure_flat_reduction":
		bonuses.StructureFlatReduction += effect.PerLevel * float64(level)

	// Piercing
	case "piercing_damage":
		bonuses.PiercingDamage += effect.PerLevel * float64(level)
	case "piercing_damage_bonus":
		bonuses.PiercingDamageBonus += effect.PerLevel * float64(level)

	// Evasion & Movement
	case "enemy_hit_reduction":
		bonuses.EnemyHitReduction += effect.PerLevel * float64(level)
	case "movement_crit_bonus":
		bonuses.MovementCritBonus += effect.PerLevel * float64(level)
	case "enemy_attack_reduction":
		bonuses.EnemyAttackReduction += effect.PerLevel * float64(level)

	// --- Complex / composite effect handlers ---

	// Armor bonus: {"type":"armor_bonus","neutral":10,"light":1,"unit":"percent_per_level"}
	case "armor_bonus":
		var ab struct {
			Neutral float64 `json:"neutral"`
			Light   float64 `json:"light"`
			Regen   float64 `json:"regen"`
		}
		if err := json.Unmarshal(rawJSON, &ab); err != nil {
			log.Printf("Failed to parse armor_bonus for tech %s: %v", techName, err)
			break
		}
		bonuses.ArmorBonusNeutral += ab.Neutral * float64(level)
		bonuses.ArmorBonusLight += ab.Light * float64(level)
		bonuses.ArmorBonusRegen += ab.Regen * float64(level)

	// Shield pen chance: {"type":"shield_pen_chance","per_level":3,"light_bonus":1,"unit":"percent"}
	case "shield_pen_chance":
		var sp struct {
			PerLevel   float64 `json:"per_level"`
			LightBonus float64 `json:"light_bonus"`
		}
		if err := json.Unmarshal(rawJSON, &sp); err != nil {
			log.Printf("Failed to parse shield_pen_chance for tech %s: %v", techName, err)
			break
		}
		bonuses.ShieldPenChance += sp.PerLevel * float64(level)
		bonuses.ShieldPenLightBonus += sp.LightBonus * float64(level)

	// Range damage: {"type":"range_damage","ranges":[220,180,150,120],"crit_rate":5,"crit_damage":5}
	case "range_damage":
		var rd struct {
			Ranges    []float64 `json:"ranges"`
			CritRate  float64   `json:"crit_rate"`
			CritDamage float64  `json:"crit_damage"`
		}
		if err := json.Unmarshal(rawJSON, &rd); err != nil {
			log.Printf("Failed to parse range_damage for tech %s: %v", techName, err)
			break
		}
		bonuses.RangeDamageRanges = rd.Ranges
		bonuses.RangeDamageCritRate += rd.CritRate
		bonuses.RangeDamageCritDmg += rd.CritDamage

	// Base shield bonus (level-indexed): {"type":"base_shield_bonus","values":[6,12,20],"unit":"percent"}
	case "base_shield_bonus":
		var bs struct {
			Values []float64 `json:"values"`
		}
		if err := json.Unmarshal(rawJSON, &bs); err != nil {
			log.Printf("Failed to parse base_shield_bonus for tech %s: %v", techName, err)
			break
		}
		if level > 0 && level <= len(bs.Values) {
			bonuses.BaseShield += bs.Values[level-1]
		}

	// Base structure bonus (level-indexed): {"type":"base_structure_bonus","values":[6,12,20],"unit":"percent"}
	case "base_structure_bonus":
		var bst struct {
			Values []float64 `json:"values"`
		}
		if err := json.Unmarshal(rawJSON, &bst); err != nil {
			log.Printf("Failed to parse base_structure_bonus for tech %s: %v", techName, err)
			break
		}
		if level > 0 && level <= len(bst.Values) {
			bonuses.BaseStructure += bst.Values[level-1]
		}

	// Hit and shield pen: {"type":"hit_and_shield_pen","hit":8,"shield_pen":8,"unit":"percent"}
	case "hit_and_shield_pen":
		var hp struct {
			Hit      float64 `json:"hit"`
			ShieldPen float64 `json:"shield_pen"`
		}
		if err := json.Unmarshal(rawJSON, &hp); err != nil {
			log.Printf("Failed to parse hit_and_shield_pen for tech %s: %v", techName, err)
			break
		}
		bonuses.HitAndShieldPenHit += hp.Hit
		bonuses.HitAndShieldPenShieldPen += hp.ShieldPen

	// Steering reduction: {"type":"steering_reduction","chance":5,"per_level":3,"unit":"percent"}
	case "steering_reduction":
		var sr struct {
			Chance   float64 `json:"chance"`
			PerLevel float64 `json:"per_level"`
		}
		if err := json.Unmarshal(rawJSON, &sr); err != nil {
			log.Printf("Failed to parse steering_reduction for tech %s: %v", techName, err)
			break
		}
		bonuses.SteeringReductionChance += sr.Chance * float64(level)
		bonuses.SteeringReductionPerLvl += sr.PerLevel * float64(level)

	// Piercing critical: {"type":"piercing_critical","enabled":true}
	case "piercing_critical":
		bonuses.PiercingCriticalEnabled = true

	// Base damage and accuracy: {"type":"base_damage_and_accuracy","damage":5,"accuracy":3,"unit":"percent_per_level"}
	case "base_damage_and_accuracy":
		var da struct {
			Damage   float64 `json:"damage"`
			Accuracy float64 `json:"accuracy"`
		}
		if err := json.Unmarshal(rawJSON, &da); err != nil {
			log.Printf("Failed to parse base_damage_and_accuracy for tech %s: %v", techName, err)
			break
		}
		bonuses.BaseDamageBonus += da.Damage * float64(level)
		bonuses.BonusAccuracy += da.Accuracy * float64(level)

	// Ignore agility: {"type":"ignore_agility","per_level":8,"chance":10,"unit":"percent"}
	case "ignore_agility":
		var ia struct {
			PerLevel float64 `json:"per_level"`
			Chance   float64 `json:"chance"`
		}
		if err := json.Unmarshal(rawJSON, &ia); err != nil {
			log.Printf("Failed to parse ignore_agility for tech %s: %v", techName, err)
			break
		}
		bonuses.IgnoreAgilityPerLevel += ia.PerLevel * float64(level)
		bonuses.IgnoreAgilityChance += ia.Chance * float64(level)

	// Ignore defense: {"type":"ignore_defense","per_level":8,"mov_reduction":1,"chance":6,"unit":"percent"}
	case "ignore_defense":
		var id struct {
			PerLevel     float64 `json:"per_level"`
			MovReduction float64 `json:"mov_reduction"`
			Chance       float64 `json:"chance"`
		}
		if err := json.Unmarshal(rawJSON, &id); err != nil {
			log.Printf("Failed to parse ignore_defense for tech %s: %v", techName, err)
			break
		}
		bonuses.IgnoreDefensePerLevel += id.PerLevel * float64(level)
		bonuses.IgnoreDefenseMovReduction += id.MovReduction * float64(level)
		bonuses.IgnoreDefenseChance += id.Chance * float64(level)

	// Missile damage and pen: {"type":"missile_damage_and_pen","damage":3,"pen":1,"unit":"percent_per_level"}
	case "missile_damage_and_pen":
		var mp struct {
			Damage float64 `json:"damage"`
			Pen    float64 `json:"pen"`
		}
		if err := json.Unmarshal(rawJSON, &mp); err != nil {
			log.Printf("Failed to parse missile_damage_and_pen for tech %s: %v", techName, err)
			break
		}
		bonuses.MissileDamage += mp.Damage * float64(level)
		bonuses.MissileShieldPenChance += mp.Pen * float64(level)

	// Scatter vs high structure: {"type":"scatter_vs_high_structure","per_level":3,"crit":3,"unit":"percent"}
	case "scatter_vs_high_structure":
		var sv struct {
			PerLevel float64 `json:"per_level"`
			Crit     float64 `json:"crit"`
		}
		if err := json.Unmarshal(rawJSON, &sv); err != nil {
			log.Printf("Failed to parse scatter_vs_high_structure for tech %s: %v", techName, err)
			break
		}
		bonuses.ScatterVsHighStructure += sv.PerLevel * float64(level)
		bonuses.ScatterVsHighStructureCrit += sv.Crit * float64(level)

	// Knockback: {"type":"knockback","chance":25,"distance":2,"per_round":true}
	case "knockback":
		var kb struct {
			Chance   float64 `json:"chance"`
			Distance int     `json:"distance"`
		}
		if err := json.Unmarshal(rawJSON, &kb); err != nil {
			log.Printf("Failed to parse knockback for tech %s: %v", techName, err)
			break
		}
		bonuses.KnockbackChance += kb.Chance
		bonuses.KnockbackDistance = kb.Distance

	// Base attack: {"type":"base_attack","flat":5,"unit":"percent"}
	case "base_attack":
		bonuses.BaseAttack += effect.Flat

	// He3 and shield damage: {"type":"he3_and_shield_damage","he3":2,"shield_damage":3,"unit":"percent_per_level"}
	case "he3_and_shield_damage":
		var hs struct {
			He3         float64 `json:"he3"`
			ShieldDamage float64 `json:"shield_damage"`
		}
		if err := json.Unmarshal(rawJSON, &hs); err != nil {
			log.Printf("Failed to parse he3_and_shield_damage for tech %s: %v", techName, err)
			break
		}
		bonuses.He3AndShieldDamageHe3 += hs.He3 * float64(level)
		bonuses.ShieldDamageBonus += hs.ShieldDamage * float64(level)

	// Ship shield damage: {"type":"ship_shield_damage","flat":10,"unit":"percent"}
	case "ship_shield_damage":
		bonuses.ShipShieldDamage += effect.Flat

	// Reload chance: {"type":"reload_chance","per_level":10,"unit":"percent"}
	case "reload_chance":
		bonuses.ReloadChance += effect.PerLevel * float64(level)

	// Reload reduction: {"type":"reload_reduction","flat":1,"unit":"rounds"}
	case "reload_reduction":
		bonuses.ReloadReduction += int(effect.Flat)

	// Long range attack: {"type":"long_range_attack","flat":15,"unit":"percent","range":"6-10 slots"}
	case "long_range_attack":
		bonuses.LongRangeAttack += effect.Flat

	// Long range bonus: {"type":"long_range_bonus","double_damage_per_level":10,"crit_rate_per_level":10,"unit":"percent"}
	case "long_range_bonus":
		var lr struct {
			DoubleDamagePerLevel float64 `json:"double_damage_per_level"`
			CritRatePerLevel     float64 `json:"crit_rate_per_level"`
		}
		if err := json.Unmarshal(rawJSON, &lr); err != nil {
			log.Printf("Failed to parse long_range_bonus for tech %s: %v", techName, err)
			break
		}
		bonuses.LongRangeBonusDblDmg += lr.DoubleDamagePerLevel * float64(level)
		bonuses.LongRangeBonusCritRate += lr.CritRatePerLevel * float64(level)

	// Swarm chance: {"type":"swarm_chance","chance_per_level":10,"attack_bonus_base":8,"attack_bonus_max":25,...}
	case "swarm_chance":
		var sc struct {
			ChancePerLevel    float64 `json:"chance_per_level"`
			AttackBonusBase   float64 `json:"attack_bonus_base"`
			AttackBonusMax    float64 `json:"attack_bonus_max"`
			He3CostBase       float64 `json:"he3_cost_increase_base"`
			He3CostMax        float64 `json:"he3_cost_increase_max"`
		}
		if err := json.Unmarshal(rawJSON, &sc); err != nil {
			log.Printf("Failed to parse swarm_chance for tech %s: %v", techName, err)
			break
		}
		bonuses.SwarmChancePerLevel += sc.ChancePerLevel * float64(level)
		bonuses.SwarmAttackBonusBase = sc.AttackBonusBase
		bonuses.SwarmAttackBonusMax = sc.AttackBonusMax
		bonuses.SwarmHe3CostBase = sc.He3CostBase
		bonuses.SwarmHe3CostMax = sc.He3CostMax

	// Multi bonus: variable sub-fields depending on tech
	case "multi_bonus":
		var mb struct {
			Swarm                     float64 `json:"swarm"`
			Attack                    float64 `json:"attack"`
			He3                       float64 `json:"he3"`
			Intercept                 float64 `json:"intercept"`
			Unshielded                float64 `json:"unshielded"`
			CritDamagePerLevel        float64 `json:"crit_damage_per_level"`
			WeaponSpaceReductionPerLvl float64 `json:"weapon_space_reduction_per_level"`
			Reload                    float64 `json:"reload"`
			ShieldDamage              float64 `json:"shield_damage"`
			AttackRange               float64 `json:"attack_range"`
			Interception              float64 `json:"interception"`
			MissileDamage             float64 `json:"missile_damage"`
			ScatterAll                float64 `json:"scatter_all"`
			MissileHitRate            float64 `json:"missile_hit_rate"`
			InterceptReduction        float64 `json:"intercept_reduction"`
		}
		if err := json.Unmarshal(rawJSON, &mb); err != nil {
			log.Printf("Failed to parse multi_bonus for tech %s: %v", techName, err)
			break
		}
		bonuses.MultiBonusSwarm += mb.Swarm
		bonuses.MultiBonusAttack += mb.Attack
		bonuses.MultiBonusHe3 += mb.He3
		bonuses.MultiBonusIntercept += mb.Intercept
		bonuses.MultiBonusUnshielded += mb.Unshielded
		bonuses.MultiBonusCritDamage += mb.CritDamagePerLevel * float64(level)
		bonuses.MultiBonusWeaponSpace += mb.WeaponSpaceReductionPerLvl * float64(level)
		bonuses.MultiBonusReload += mb.Reload * float64(level)
		bonuses.MultiBonusShieldDmg += mb.ShieldDamage * float64(level)
		bonuses.MultiBonusAttackRange += mb.AttackRange * float64(level)
		bonuses.MultiBonusInterception += mb.Interception * float64(level)
		bonuses.MissileDamage += mb.MissileDamage
		bonuses.ScatterAll += mb.ScatterAll
		bonuses.MissileHitRate += mb.MissileHitRate
		bonuses.InterceptReduction += mb.InterceptReduction

	// Shield restore: {"type":"shield_restore","per_level":30,"intercept":1,"unit":"percent"}
	case "shield_restore":
		var sr struct {
			PerLevel  float64 `json:"per_level"`
			Intercept float64 `json:"intercept"`
		}
		if err := json.Unmarshal(rawJSON, &sr); err != nil {
			log.Printf("Failed to parse shield_restore for tech %s: %v", techName, err)
			break
		}
		bonuses.ShieldRestore += sr.PerLevel * float64(level)
		bonuses.ShieldRestoreIntercept += sr.Intercept * float64(level)

	// Structure restore: {"type":"structure_restore","per_level":30,"unit":"percent"}
	case "structure_restore":
		bonuses.StructureRestore += effect.PerLevel * float64(level)

	// Absorb double: {"type":"absorb_double","per_level":10,"collateral_reduction":15,"unit":"percent"}
	case "absorb_double":
		var ad struct {
			PerLevel            float64 `json:"per_level"`
			CollateralReduction float64 `json:"collateral_reduction"`
		}
		if err := json.Unmarshal(rawJSON, &ad); err != nil {
			log.Printf("Failed to parse absorb_double for tech %s: %v", techName, err)
			break
		}
		bonuses.AbsorbDouble += ad.PerLevel * float64(level)
		bonuses.AbsorbDoubleCollateral += ad.CollateralReduction * float64(level)

	// Absorb double structure: {"type":"absorb_double_structure","per_level":10,"collateral_reduction":15,"unit":"percent"}
	case "absorb_double_structure":
		var ads struct {
			PerLevel            float64 `json:"per_level"`
			CollateralReduction float64 `json:"collateral_reduction"`
		}
		if err := json.Unmarshal(rawJSON, &ads); err != nil {
			log.Printf("Failed to parse absorb_double_structure for tech %s: %v", techName, err)
			break
		}
		bonuses.AbsorbDoubleStructure += ads.PerLevel * float64(level)
		bonuses.AbsorbDoubleStructureColl += ads.CollateralReduction * float64(level)

	default:
		// Unknown effect type - log for future implementation
		log.Printf("Unknown tech effect type: %s (tech: %s)", effect.Type, techName)
	}
}

// GetTechLevel returns the current level of a specific technology for a player.
// Returns 0 if the technology hasn't been researched yet.
// This is used for prerequisite checks and condition evaluation.
func GetTechLevel(playerID string, techName string) (int, error) {
	var level int
	err := database.DB.QueryRow(`
		SELECT COALESCE(t.level, 0)
		FROM tech_types tt
		LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
		WHERE tt.name = $2
	`, playerID, techName).Scan(&level)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Printf("Failed to get tech level for %s: %v", techName, err)
		return 0, fmt.Errorf("failed to get tech level: %w", err)
	}

	return level, nil
}
