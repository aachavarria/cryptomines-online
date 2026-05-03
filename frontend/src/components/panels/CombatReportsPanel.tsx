import React, { useState } from 'react';
import { useCombatReports, useCombatReport } from '../../hooks/useCombatReports';
import { combatReportsAPI, type RoundData } from '../../services/api';
import BattlePlayback from './BattlePlayback';

type ReportFilter = 'all' | 'pvp' | 'instance' | 'rbp'

export const CombatReportsPanel: React.FC = () => {
  const { reports, loading, error } = useCombatReports();
  const [selectedReportId, setSelectedReportId] = useState<string | null>(null);
  const [filter, setFilter] = useState<ReportFilter>('all');
  const [showStaticLog, setShowStaticLog] = useState(false);
  const { report: selectedReport, loading: reportLoading } = useCombatReport(selectedReportId);

  const filteredReports = filter === 'all'
    ? reports
    : reports.filter(r => r.combat_type === filter);

  if (loading) {
    return (
      <div className="ds-panel" style={{ maxWidth: 1200, margin: '0 auto' }}>
        Loading combat reports...
      </div>
    );
  }

  if (error) {
    return (
      <div className="ds-panel" style={{ maxWidth: 1200, margin: '0 auto', color: 'var(--ds-danger)' }}>
        Error: {error}
      </div>
    );
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const isVictory = (result: string) => result === 'attacker_win';
  const isDefeat = (result: string) => result === 'defender_win';

  const getResultBadge = (result: string) => {
    if (isVictory(result)) return <span className="ds-badge ds-badge--success">Victory</span>;
    if (isDefeat(result)) return <span className="ds-badge ds-badge--danger">Defeat</span>;
    return <span className="ds-badge ds-badge--neutral">Draw</span>;
  };

  const getResultLabel = (result: string) => {
    if (isVictory(result)) return 'Victory';
    if (isDefeat(result)) return 'Defeat';
    return 'Draw';
  };

  const getResultColor = (result: string) => {
    if (isVictory(result)) return 'var(--ds-success)';
    if (isDefeat(result)) return 'var(--ds-danger)';
    return 'var(--ds-text-muted)';
  };

  const renderReportList = () => (
    <div>
      <h2 className="ds-h2" style={{ marginBottom: 'var(--sp-4)' }}>Combat Reports</h2>
      <div className="ds-tabs" style={{ marginBottom: 'var(--sp-4)' }}>
        {(['all', 'pvp', 'instance', 'rbp'] as ReportFilter[]).map(f => {
          const count = f === 'all' ? reports.length : reports.filter(r => r.combat_type === f).length;
          return (
            <button
              key={f}
              className="ds-tab"
              aria-selected={filter === f}
              onClick={() => setFilter(f)}
            >
              {f === 'all' ? 'All' : f.toUpperCase()} ({count})
            </button>
          );
        })}
      </div>
      {filteredReports.length === 0 ? (
        <p className="ds-text-muted" style={{ textAlign: 'center', padding: 'var(--sp-10)' }}>
          No combat reports found
        </p>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-2)' }}>
          {filteredReports.map((report) => {
            const loot = combatReportsAPI.parseLoot(report);
            return (
              <div
                key={report.id}
                className="ds-list-item"
                onClick={() => setSelectedReportId(report.id)}
                role="button"
                tabIndex={0}
              >
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: 'var(--sp-2)',
                  }}
                >
                  <span className="ds-caption">{report.combat_type}</span>
                  {getResultBadge(report.result)}
                </div>
                <div
                  style={{
                    display: 'grid',
                    gridTemplateColumns: '1fr 1fr',
                    gap: 'var(--sp-2)',
                    fontSize: 'var(--fs-sm)',
                    marginBottom: 'var(--sp-2)',
                  }}
                >
                  <div className="ds-row--between">
                    <span className="ds-text-muted">Rounds:</span>
                    <span className="ds-mono">{report.total_rounds}</span>
                  </div>
                  <div className="ds-row--between">
                    <span className="ds-text-muted">He3 Used:</span>
                    <span className="ds-mono">{report.he3_consumed.toLocaleString()}</span>
                  </div>
                </div>
                {loot && (
                  <div
                    style={{
                      display: 'flex',
                      gap: 'var(--sp-3)',
                      padding: 'var(--sp-2)',
                      background: 'var(--ds-surface-2)',
                      borderRadius: 'var(--r-md)',
                      marginBottom: 'var(--sp-1)',
                      fontSize: 'var(--fs-sm)',
                    }}
                  >
                    <span className="ds-row" style={{ gap: 4 }}>
                      <span className="ds-resource-dot ds-resource-dot--metal" />
                      <span className="ds-mono">{loot.metal.toLocaleString()}</span>
                    </span>
                    <span className="ds-row" style={{ gap: 4 }}>
                      <span className="ds-resource-dot ds-resource-dot--he3" />
                      <span className="ds-mono">{loot.he3.toLocaleString()}</span>
                    </span>
                    <span className="ds-row" style={{ gap: 4 }}>
                      <span className="ds-resource-dot ds-resource-dot--gold" />
                      <span className="ds-mono">{loot.gold.toLocaleString()}</span>
                    </span>
                  </div>
                )}
                <div
                  className="ds-text-soft"
                  style={{ fontSize: 'var(--fs-caption)', textAlign: 'right' }}
                >
                  {formatDate(report.created_at)}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );

  const renderReportDetails = () => {
    if (!selectedReport) return null;
    if (reportLoading) return <div className="ds-panel">Loading...</div>;

    const loot = combatReportsAPI.parseLoot(selectedReport);
    const rounds = combatReportsAPI.parseRounds(selectedReport);

    return (
      <div className="ds-panel">
        <button
          className="ds-btn-ghost"
          onClick={() => setSelectedReportId(null)}
          style={{ marginBottom: 'var(--sp-4)' }}
        >
          ← Back to List
        </button>

        <h2 className="ds-h2" style={{ marginBottom: 'var(--sp-4)' }}>Battle Report</h2>

        <section style={{ marginBottom: 'var(--sp-5)' }}>
          <h3 className="ds-h3" style={{ marginBottom: 'var(--sp-3)' }}>Summary</h3>
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
              gap: 'var(--sp-3)',
            }}
          >
            <div className="ds-card ds-row--between">
              <span className="ds-text-muted">Type:</span>
              <span>{selectedReport.combat_type}</span>
            </div>
            <div className="ds-card ds-row--between">
              <span className="ds-text-muted">Result:</span>
              <span style={{ color: getResultColor(selectedReport.result), fontWeight: 'var(--fw-semibold)' }}>
                {getResultLabel(selectedReport.result)}
              </span>
            </div>
            <div className="ds-card ds-row--between">
              <span className="ds-text-muted">Total Rounds:</span>
              <span className="ds-mono">{selectedReport.total_rounds}</span>
            </div>
            <div className="ds-card ds-row--between">
              <span className="ds-text-muted">He3 Consumed:</span>
              <span className="ds-mono">{selectedReport.he3_consumed.toLocaleString()}</span>
            </div>
            <div className="ds-card ds-row--between">
              <span className="ds-text-muted">Date:</span>
              <span>{formatDate(selectedReport.created_at)}</span>
            </div>
          </div>
        </section>

        {loot && (
          <section style={{ marginBottom: 'var(--sp-5)' }}>
            <h3 className="ds-h3" style={{ marginBottom: 'var(--sp-3)' }}>Resources Gained</h3>
            <div
              style={{
                display: 'flex',
                gap: 'var(--sp-4)',
                justifyContent: 'space-around',
                flexWrap: 'wrap',
              }}
            >
              {[
                { name: 'Metal', amount: loot.metal, dot: 'ds-resource-dot--metal' },
                { name: 'He3', amount: loot.he3, dot: 'ds-resource-dot--he3' },
                { name: 'Gold', amount: loot.gold, dot: 'ds-resource-dot--gold' },
              ].map(({ name, amount, dot }) => (
                <div
                  key={name}
                  className="ds-card"
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    gap: 'var(--sp-2)',
                    minWidth: 120,
                  }}
                >
                  <span className={`ds-resource-dot ${dot}`} style={{ width: 20, height: 20 }} />
                  <span className="ds-caption">{name}</span>
                  <span
                    className="ds-mono"
                    style={{ fontSize: 'var(--fs-h3)', fontWeight: 'var(--fw-semibold)', color: 'var(--ds-success)' }}
                  >
                    {amount.toLocaleString()}
                  </span>
                </div>
              ))}
            </div>
          </section>
        )}

        {rounds && rounds.length > 0 && (
          <section style={{ marginBottom: 'var(--sp-5)' }}>
            <div className="ds-row--between" style={{ marginBottom: 'var(--sp-3)' }}>
              <h3 className="ds-h3" style={{ margin: 0 }}>Battle Playback</h3>
              <button
                className="ds-btn-ghost ds-btn--sm"
                onClick={() => setShowStaticLog((s) => !s)}
              >
                {showStaticLog ? 'Show playback' : 'Show full log'}
              </button>
            </div>
            {showStaticLog ? (
              <div style={{ maxHeight: 600, overflowY: 'auto' }}>
                {rounds.map((round: RoundData) => (
                  <div
                    key={round.RoundNumber}
                    className="ds-card"
                    style={{
                      borderLeft: '3px solid var(--ds-teal)',
                      marginBottom: 'var(--sp-3)',
                    }}
                  >
                    <h4 className="ds-h3" style={{ marginBottom: 'var(--sp-2)' }}>
                      Round {round.RoundNumber}
                    </h4>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                      {round.Attacks.map((attack, idx) => (
                        <div
                          key={idx}
                          style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 'var(--sp-2)',
                            padding: '6px 8px',
                            background: 'var(--ds-surface)',
                            borderRadius: 'var(--r-sm)',
                            borderLeft: `2px solid ${attack.Hit ? 'var(--ds-success)' : 'var(--ds-text-soft)'}`,
                            fontSize: 'var(--fs-sm)',
                            opacity: attack.Hit ? 1 : 0.7,
                          }}
                        >
                          <span className="ds-caption">[{attack.AttackerSide}]</span>
                          <span style={{ fontWeight: 'var(--fw-semibold)', minWidth: 60 }}>
                            {attack.Hit ? 'HIT' : 'MISS'}
                          </span>
                          <span className="ds-caption">[{attack.DefenderSide}]</span>
                          {attack.Hit && (
                            <>
                              <span style={{ color: 'var(--ds-orange-strong)', fontSize: 'var(--fs-caption)' }}>
                                Dmg: {attack.Damage} ({attack.ShieldDamage} shield, {attack.StructureDamage} structure)
                              </span>
                              {attack.ShipsDestroyed > 0 && (
                                <span style={{ color: 'var(--ds-danger)', fontWeight: 'var(--fw-semibold)' }}>
                                  {attack.ShipsDestroyed} destroyed
                                </span>
                              )}
                            </>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <BattlePlayback rounds={rounds} totalRounds={selectedReport.total_rounds} />
            )}
          </section>
        )}
      </div>
    );
  };

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto', padding: 'var(--sp-5)' }}>
      {selectedReportId ? renderReportDetails() : renderReportList()}
    </div>
  );
};
