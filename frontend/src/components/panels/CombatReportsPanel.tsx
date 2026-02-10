import React, { useState } from 'react';
import { useCombatReports, useCombatReport } from '../../hooks/useCombatReports';
import { combatReportsAPI, type RoundData } from '../../services/api';
import './CombatReportsPanel.css';

export const CombatReportsPanel: React.FC = () => {
  const { reports, loading, error } = useCombatReports();
  const [selectedReportId, setSelectedReportId] = useState<string | null>(null);
  const { report: selectedReport, loading: reportLoading } = useCombatReport(selectedReportId);

  if (loading) {
    return <div className="combat-reports-panel">Loading combat reports...</div>;
  }

  if (error) {
    return <div className="combat-reports-panel error">Error: {error}</div>;
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const getResultClass = (result: string) => {
    if (result === 'attacker_win') return 'result-victory';
    if (result === 'defender_win') return 'result-defeat';
    return 'result-draw';
  };

  const getResultLabel = (result: string) => {
    if (result === 'attacker_win') return 'Victory';
    if (result === 'defender_win') return 'Defeat';
    return 'Draw';
  };

  const renderReportList = () => (
    <div className="reports-list">
      <h3>Combat Reports</h3>
      {reports.length === 0 ? (
        <p className="no-reports">No combat reports yet</p>
      ) : (
        <div className="reports-grid">
          {reports.map((report) => {
            const loot = combatReportsAPI.parseLoot(report);
            return (
              <div
                key={report.id}
                className={`report-card ${getResultClass(report.result)}`}
                onClick={() => setSelectedReportId(report.id)}
              >
                <div className="report-header">
                  <span className="report-type">{report.combat_type}</span>
                  <span className={`report-result ${getResultClass(report.result)}`}>
                    {getResultLabel(report.result)}
                  </span>
                </div>
                <div className="report-stats">
                  <div className="stat">
                    <span className="label">Rounds:</span>
                    <span className="value">{report.total_rounds}</span>
                  </div>
                  <div className="stat">
                    <span className="label">He3 Used:</span>
                    <span className="value">{report.he3_consumed.toLocaleString()}</span>
                  </div>
                </div>
                {loot && (
                  <div className="report-loot">
                    <div className="loot-item">
                      <span className="resource-icon">⚙️</span>
                      {loot.metal.toLocaleString()}
                    </div>
                    <div className="loot-item">
                      <span className="resource-icon">⚗️</span>
                      {loot.he3.toLocaleString()}
                    </div>
                    <div className="loot-item">
                      <span className="resource-icon">💰</span>
                      {loot.gold.toLocaleString()}
                    </div>
                  </div>
                )}
                <div className="report-date">{formatDate(report.created_at)}</div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );

  const renderReportDetails = () => {
    if (!selectedReport) return null;
    if (reportLoading) return <div className="report-details">Loading...</div>;

    const loot = combatReportsAPI.parseLoot(selectedReport);
    const rounds = combatReportsAPI.parseRounds(selectedReport);

    return (
      <div className="report-details">
        <button className="back-button" onClick={() => setSelectedReportId(null)}>
          ← Back to List
        </button>

        <h3>Battle Report</h3>

        <div className="detail-section">
          <h4>Summary</h4>
          <div className="detail-grid">
            <div className="detail-item">
              <span className="label">Type:</span>
              <span>{selectedReport.combat_type}</span>
            </div>
            <div className="detail-item">
              <span className="label">Result:</span>
              <span className={getResultClass(selectedReport.result)}>
                {getResultLabel(selectedReport.result)}
              </span>
            </div>
            <div className="detail-item">
              <span className="label">Total Rounds:</span>
              <span>{selectedReport.total_rounds}</span>
            </div>
            <div className="detail-item">
              <span className="label">He3 Consumed:</span>
              <span>{selectedReport.he3_consumed.toLocaleString()}</span>
            </div>
            <div className="detail-item">
              <span className="label">Date:</span>
              <span>{formatDate(selectedReport.created_at)}</span>
            </div>
          </div>
        </div>

        {loot && (
          <div className="detail-section">
            <h4>Resources Gained</h4>
            <div className="loot-summary">
              <div className="loot-item-large">
                <span className="resource-icon">⚙️</span>
                <span className="resource-name">Metal</span>
                <span className="resource-amount">{loot.metal.toLocaleString()}</span>
              </div>
              <div className="loot-item-large">
                <span className="resource-icon">⚗️</span>
                <span className="resource-name">He3</span>
                <span className="resource-amount">{loot.he3.toLocaleString()}</span>
              </div>
              <div className="loot-item-large">
                <span className="resource-icon">💰</span>
                <span className="resource-name">Gold</span>
                <span className="resource-amount">{loot.gold.toLocaleString()}</span>
              </div>
            </div>
          </div>
        )}

        {rounds && rounds.length > 0 && (
          <div className="detail-section">
            <h4>Round-by-Round Details</h4>
            <div className="rounds-list">
              {rounds.map((round: RoundData) => (
                <div key={round.RoundNumber} className="round-card">
                  <h5>Round {round.RoundNumber}</h5>
                  <div className="attacks-list">
                    {round.Attacks.map((attack, idx) => (
                      <div key={idx} className={`attack-log ${attack.Hit ? 'hit' : 'miss'}`}>
                        <span className="attacker-side">[{attack.AttackerSide}]</span>
                        <span className="attack-action">
                          {attack.Hit ? '⚔️ HIT' : '❌ MISS'}
                        </span>
                        <span className="defender-side">[{attack.DefenderSide}]</span>
                        {attack.Hit && (
                          <>
                            <span className="damage">
                              Dmg: {attack.Damage} ({attack.ShieldDamage} shield, {attack.StructureDamage} structure)
                            </span>
                            {attack.ShipsDestroyed > 0 && (
                              <span className="casualties">💥 {attack.ShipsDestroyed} destroyed</span>
                            )}
                          </>
                        )}
                      </div>
                    ))}
                  </div>
                  {Object.keys(round.Casualties).length > 0 && (
                    <div className="round-casualties">
                      <strong>Round Casualties:</strong>
                      {Object.entries(round.Casualties).map(([stackId, count]) => (
                        <span key={stackId} className="casualty-item">
                          {stackId}: {count} ships
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="combat-reports-panel">
      {selectedReportId ? renderReportDetails() : renderReportList()}
    </div>
  );
};
