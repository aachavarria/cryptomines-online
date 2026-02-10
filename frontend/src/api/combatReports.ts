import { listCombatReports, getCombatReport, type CombatReport } from '../services/api';

export type { CombatReport };

export interface LootData {
  metal: number;
  he3: number;
  gold: number;
}

export interface AttackData {
  AttackerStackID: string;
  DefenderStackID: string;
  AttackerSide: string;
  DefenderSide: string;
  Hit: boolean;
  Damage: number;
  ShieldDamage: number;
  StructureDamage: number;
  ShipsDestroyed: number;
}

export interface RoundData {
  RoundNumber: number;
  Attacks: AttackData[];
  Casualties: Record<string, number>;
}

export const combatReportsAPI = {
  list: listCombatReports,
  get: getCombatReport,

  parseLoot: (report: CombatReport): LootData | null => {
    if (!report.loot_json) return null;
    try {
      return JSON.parse(report.loot_json) as LootData;
    } catch {
      return null;
    }
  },

  parseRounds: (report: CombatReport): RoundData[] | null => {
    if (!report.rounds_json) return null;
    try {
      return JSON.parse(report.rounds_json) as RoundData[];
    } catch {
      return null;
    }
  },
};
