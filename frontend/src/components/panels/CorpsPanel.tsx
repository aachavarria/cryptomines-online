import { useState, useEffect, useCallback, type CSSProperties } from 'react'
import { createPortal } from 'react-dom'
import { Users, Shield, X } from 'lucide-react'
import { useCorp } from '../../hooks/useCorp'
import { useGameContext } from '../../contexts/GameContext'
import LoadingButton from '../common/LoadingButton'

const styles: Record<string, CSSProperties> = {
  modalBody: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-4)',
  },
  formGroup: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-1)',
    marginBottom: 'var(--sp-3)',
  },
  list: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-2)',
    maxHeight: 320,
    overflowY: 'auto',
  },
  rowBetween: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 'var(--sp-3)',
  },
  overviewRow: {
    display: 'flex',
    justifyContent: 'space-between',
    paddingBottom: 'var(--sp-2)',
    borderBottom: '1px solid var(--ds-border)',
  },
  donateBar: {
    height: 6,
    background: 'var(--ds-border)',
    borderRadius: 'var(--r-pill)',
    overflow: 'hidden',
    marginTop: 'var(--sp-2)',
  },
  donateBarFill: {
    height: '100%',
    transition: 'width var(--motion-slow)',
    borderRadius: 'inherit',
  },
}

interface CorpsPanelProps {
  onClose: () => void
}

export default function CorpsPanel({ onClose }: CorpsPanelProps) {
  const { state } = useGameContext()
  const {
    corpData,
    members,
    searchResults,
    loading,
    membersLoading,
    searchLoading,
    create,
    join,
    leave,
    donate,
    updateRole,
    search,
    fetchMembers,
  } = useCorp()

  const [activeTab, setActiveTab] = useState<'create' | 'join' | 'overview' | 'members' | 'donate' | 'manage'>('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [createForm, setCreateForm] = useState({ name: '', tag: '', description: '' })
  const [donateForm, setDonateForm] = useState({ metal: 0, he3: 0, gold: 0 })

  const isInCorp = !!corpData?.corp
  const isLeader = corpData?.role === 'leader'

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  // Fetch members when in corp
  useEffect(() => {
    if (isInCorp && !membersLoading && members.length === 0) {
      fetchMembers()
    }
  }, [isInCorp, membersLoading, members.length, fetchMembers])

  // Set default tab based on corp status
  useEffect(() => {
    if (isInCorp) {
      setActiveTab('overview')
    } else {
      setActiveTab('create')
    }
  }, [isInCorp])

  const handleCreate = useCallback(async () => {
    if (!createForm.name || !createForm.tag || createForm.tag.length < 3 || createForm.tag.length > 5) {
      return
    }
    const success = await create(createForm)
    if (success) {
      setCreateForm({ name: '', tag: '', description: '' })
    }
  }, [create, createForm])

  const handleJoin = useCallback(
    async (corpId: string) => {
      await join(corpId)
    },
    [join],
  )

  const handleLeave = useCallback(async () => {
    if (confirm('Are you sure you want to leave this corp?')) {
      await leave()
    }
  }, [leave])

  const handleDonate = useCallback(async () => {
    if (donateForm.metal === 0 && donateForm.he3 === 0 && donateForm.gold === 0) {
      return
    }
    const success = await donate(donateForm)
    if (success) {
      setDonateForm({ metal: 0, he3: 0, gold: 0 })
    }
  }, [donate, donateForm])

  const handleSearch = useCallback(
    (query: string) => {
      setSearchQuery(query)
      search(query)
    },
    [search],
  )

  const handleRoleChange = useCallback(
    async (memberId: string, newRole: string) => {
      await updateRole(memberId, newRole)
    },
    [updateRole],
  )

  const renderTabs = () => {
    const tabs: { id: typeof activeTab; label: string; show: boolean }[] = isInCorp
      ? [
          { id: 'overview', label: 'Overview', show: true },
          { id: 'members', label: 'Members', show: true },
          { id: 'donate', label: 'Donate', show: true },
          { id: 'manage', label: 'Manage', show: !!isLeader },
        ]
      : [
          { id: 'create', label: 'Create', show: true },
          { id: 'join', label: 'Join', show: true },
        ]

    return (
      <div className="ds-tabs" role="tablist">
        {tabs
          .filter((t) => t.show)
          .map((t) => (
            <button
              key={t.id}
              className="ds-tab"
              role="tab"
              aria-selected={activeTab === t.id}
              onClick={() => setActiveTab(t.id)}
            >
              {t.label}
            </button>
          ))}
      </div>
    )
  }

  const renderContent = () => {
    if (!isInCorp) {
      if (activeTab === 'create') {
        return (
          <div style={styles.modalBody}>
            <h3 className="ds-h3" style={{ margin: 0 }}>Create a New Corp</h3>
            <div style={styles.formGroup}>
              <label className="ds-caption">Corp Name</label>
              <input
                type="text"
                className="ds-input"
                value={createForm.name}
                onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                placeholder="Enter corp name"
                maxLength={50}
              />
            </div>
            <div style={styles.formGroup}>
              <label className="ds-caption">Corp Tag (3-5 chars)</label>
              <input
                type="text"
                className="ds-input"
                value={createForm.tag}
                onChange={(e) => setCreateForm({ ...createForm, tag: e.target.value.toUpperCase() })}
                placeholder="ABC"
                maxLength={5}
              />
            </div>
            <div style={styles.formGroup}>
              <label className="ds-caption">Description</label>
              <textarea
                className="ds-textarea"
                value={createForm.description}
                onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                placeholder="Describe your corp..."
                maxLength={200}
                style={{ minHeight: 80 }}
              />
            </div>
            <LoadingButton
              className="ds-btn ds-btn-primary"
              onClick={handleCreate}
              loading={loading}
              disabled={
                !createForm.name.trim() ||
                createForm.tag.length < 3 ||
                createForm.tag.length > 5
              }
            >
              Create Corp
            </LoadingButton>
          </div>
        )
      }

      if (activeTab === 'join') {
        return (
          <div style={styles.modalBody}>
            <h3 className="ds-h3" style={{ margin: 0 }}>Join a Corp</h3>
            <input
              type="text"
              className="ds-input"
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              placeholder="Search corps by name or tag..."
            />
            <div style={styles.list}>
              {searchLoading && <div className="loading-spinner" />}
              {!searchLoading && searchResults.length === 0 && searchQuery && (
                <p className="ds-text-muted">No corps found.</p>
              )}
              {!searchLoading && searchResults.length === 0 && !searchQuery && (
                <p className="ds-text-muted">Type to search for corps.</p>
              )}
              {!searchLoading &&
                searchResults.map((corp) => (
                  <div key={corp.id} className="ds-card" style={styles.rowBetween}>
                    <div>
                      <div style={{ fontWeight: 600 }}>
                        <span className="ds-badge ds-badge--neutral" style={{ marginRight: 'var(--sp-2)' }}>
                          {corp.tag}
                        </span>
                        {corp.name}
                      </div>
                      <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', marginTop: 4 }}>
                        Level {corp.level} | {corp.member_count}/{corp.max_members} members
                      </div>
                      {corp.description && (
                        <div className="ds-text-soft" style={{ fontSize: 'var(--fs-sm)', marginTop: 4 }}>
                          {corp.description}
                        </div>
                      )}
                    </div>
                    <LoadingButton
                      className="ds-btn ds-btn-primary ds-btn--sm"
                      onClick={() => handleJoin(corp.id)}
                      loading={loading}
                    >
                      Join
                    </LoadingButton>
                  </div>
                ))}
            </div>
          </div>
        )
      }
    }

    // In corp tabs
    if (activeTab === 'overview') {
      return (
        <div style={styles.modalBody}>
          <h3 className="ds-h3" style={{ margin: 0 }}>Corp Overview</h3>
          {corpData?.corp && (
            <div className="ds-card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-2)' }}>
              <div style={styles.overviewRow}>
                <strong>Name</strong>
                <span>
                  <span className="ds-badge ds-badge--teal" style={{ marginRight: 'var(--sp-2)' }}>
                    {corpData.corp.tag}
                  </span>
                  {corpData.corp.name}
                </span>
              </div>
              <div style={styles.overviewRow}>
                <strong>Level</strong>
                <span className="ds-mono">{corpData.corp.level}</span>
              </div>
              <div style={styles.overviewRow}>
                <strong>Wealth</strong>
                <span className="ds-mono">{corpData.corp.wealth.toLocaleString()}</span>
              </div>
              <div style={styles.overviewRow}>
                <strong>Members</strong>
                <span className="ds-mono">{members.length}/{corpData.corp.max_members}</span>
              </div>
              {corpData.bonuses && (
                <>
                  <div style={styles.overviewRow}>
                    <strong>Controlled RBPs</strong>
                    <span className="ds-mono">{corpData.bonuses.rbp_count}</span>
                  </div>
                  <div style={styles.overviewRow}>
                    <strong>Total RBP Bonus</strong>
                    <span className="ds-mono" style={{ color: 'var(--ds-teal-dark)' }}>
                      +{corpData.bonuses.total_rbp_bonus}%
                    </span>
                  </div>
                </>
              )}
              <div style={{ ...styles.overviewRow, borderBottom: 'none', alignItems: 'flex-start' }}>
                <strong>Description</strong>
                <span className="ds-text-muted" style={{ textAlign: 'right', maxWidth: '60%' }}>
                  {corpData.corp.description || 'No description'}
                </span>
              </div>
              {!isLeader && (
                <div style={{ marginTop: 'var(--sp-3)' }}>
                  <LoadingButton
                    className="ds-btn ds-btn-ghost"
                    onClick={handleLeave}
                    loading={loading}
                  >
                    Leave Corp
                  </LoadingButton>
                </div>
              )}
            </div>
          )}
        </div>
      )
    }

    if (activeTab === 'members') {
      return (
        <div style={styles.modalBody}>
          <h3 className="ds-h3" style={{ margin: 0 }}>Corp Members</h3>
          <div style={styles.list}>
            {membersLoading && <div className="loading-spinner" />}
            {!membersLoading && members.length === 0 && (
              <p className="ds-text-muted">No members found.</p>
            )}
            {!membersLoading &&
              members.map((member) => {
                const roleBadge =
                  member.role === 'leader' ? 'ds-badge--orange' :
                  member.role === 'officer' ? 'ds-badge--teal' :
                  'ds-badge--neutral'
                return (
                  <div key={member.player_id} className="ds-card">
                    <div style={styles.rowBetween}>
                      <div>
                        <div style={{ fontWeight: 600, display: 'inline-flex', alignItems: 'center', gap: 'var(--sp-2)' }}>
                          {member.player_name}
                          <span className={`ds-badge ${roleBadge}`}>{member.role}</span>
                        </div>
                        <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', marginTop: 4 }}>
                          Contribution: <span className="ds-mono">{member.contribution_points.toLocaleString()}</span>
                          {' | '}
                          Daily: <span className="ds-mono">{member.daily_contribution.toLocaleString()}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                )
              })}
          </div>
        </div>
      )
    }

    if (activeTab === 'donate') {
      const myMember = members.length > 0 ? members.find(m => {
        return state.player && m.player_id === state.player.id
      }) : null
      const dailyUsed = myMember?.daily_contribution ?? 0
      const dailyMax = 200
      const dailyRemaining = dailyMax - dailyUsed
      const dailyPct = Math.min(100, (dailyUsed / dailyMax) * 100)
      const dailyBarColor =
        dailyPct >= 95 ? 'var(--ds-danger)' :
        dailyPct >= 80 ? 'var(--ds-orange)' :
        'var(--ds-success)'

      return (
        <div style={styles.modalBody}>
          <h3 className="ds-h3" style={{ margin: 0 }}>Donate Resources</h3>
          <div className="ds-card">
            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 'var(--fs-sm)' }}>
              <span>Daily Contribution: <span className="ds-mono">{dailyUsed}/{dailyMax}</span> pts</span>
              <span style={{ color: dailyBarColor, fontWeight: 600 }}>
                {dailyRemaining > 0 ? `${dailyRemaining} pts remaining` : 'Limit reached'}
              </span>
            </div>
            <div className="ds-bar" style={styles.donateBar}>
              <div
                className="ds-bar-fill"
                style={{ ...styles.donateBarFill, width: `${dailyPct}%`, background: dailyBarColor }}
              />
            </div>
          </div>
          <p className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', margin: 0 }}>
            Donate resources to your corp to earn contribution points and increase corp wealth. (1 pt per 10,000 resources)
          </p>
          <div style={styles.formGroup}>
            <label className="ds-caption">Metal</label>
            <input
              type="number"
              className="ds-input"
              value={donateForm.metal}
              onChange={(e) => setDonateForm({ ...donateForm, metal: parseInt(e.target.value) || 0 })}
              min={0}
            />
          </div>
          <div style={styles.formGroup}>
            <label className="ds-caption">He3</label>
            <input
              type="number"
              className="ds-input"
              value={donateForm.he3}
              onChange={(e) => setDonateForm({ ...donateForm, he3: parseInt(e.target.value) || 0 })}
              min={0}
            />
          </div>
          <div style={styles.formGroup}>
            <label className="ds-caption">Gold</label>
            <input
              type="number"
              className="ds-input"
              value={donateForm.gold}
              onChange={(e) => setDonateForm({ ...donateForm, gold: parseInt(e.target.value) || 0 })}
              min={0}
            />
          </div>
          <LoadingButton
            className="ds-btn ds-btn-primary"
            onClick={handleDonate}
            loading={loading}
            disabled={donateForm.metal === 0 && donateForm.he3 === 0 && donateForm.gold === 0}
          >
            Donate
          </LoadingButton>
        </div>
      )
    }

    if (activeTab === 'manage' && isLeader) {
      return (
        <div style={styles.modalBody}>
          <h3 className="ds-h3" style={{ margin: 0 }}>Manage Members</h3>
          <div style={styles.list}>
            {membersLoading && <div className="loading-spinner" />}
            {!membersLoading && members.length === 0 && (
              <p className="ds-text-muted">No members found.</p>
            )}
            {!membersLoading &&
              members
                .filter((m) => m.role !== 'leader')
                .map((member) => (
                  <div key={member.player_id} className="ds-card" style={styles.rowBetween}>
                    <div>
                      <div style={{ fontWeight: 600 }}>{member.player_name}</div>
                      <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', marginTop: 4 }}>
                        {member.role.charAt(0).toUpperCase() + member.role.slice(1)} | Contribution:{' '}
                        <span className="ds-mono">{member.contribution_points.toLocaleString()}</span>
                      </div>
                    </div>
                    <select
                      className="ds-select"
                      style={{ width: 'auto' }}
                      value={member.role}
                      onChange={(e) => handleRoleChange(member.player_id, e.target.value)}
                    >
                      <option value="member">Member</option>
                      <option value="officer">Officer</option>
                    </select>
                  </div>
                ))}
          </div>
        </div>
      )
    }

    return null
  }

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div className="ds-modal ds-modal--lg" role="dialog" aria-modal="true">
        <div className="ds-modal-header">
          <h2 className="ds-modal-title" style={{ display: 'inline-flex', alignItems: 'center', gap: 'var(--sp-2)' }}>
            {isInCorp ? (
              <Shield size={20} strokeWidth={1.75} aria-hidden="true" />
            ) : (
              <Users size={20} strokeWidth={1.75} aria-hidden="true" />
            )}
            Corps
          </h2>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={16} strokeWidth={2} aria-hidden="true" />
          </button>
        </div>
        <div style={{ padding: '0 var(--sp-6)' }}>
          {renderTabs()}
        </div>
        <div className="ds-modal-body">
          {renderContent()}
        </div>
      </div>
    </div>,
    document.body,
  )
}
