import { useState, useEffect, useCallback } from 'react'
import { createPortal } from 'react-dom'
import { useCorp } from '../../hooks/useCorp'
import { useGameContext } from '../../contexts/GameContext'
import LoadingButton from '../common/LoadingButton'
import '../../styles/common.css'
import '../../styles/chat.css'

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
    if (!isInCorp) {
      return (
        <>
          <button
            className={`chat-tab ${activeTab === 'create' ? 'active' : ''}`}
            onClick={() => setActiveTab('create')}
          >
            Create
          </button>
          <button
            className={`chat-tab ${activeTab === 'join' ? 'active' : ''}`}
            onClick={() => setActiveTab('join')}
          >
            Join
          </button>
        </>
      )
    }

    return (
      <>
        <button
          className={`chat-tab ${activeTab === 'overview' ? 'active' : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          Overview
        </button>
        <button
          className={`chat-tab ${activeTab === 'members' ? 'active' : ''}`}
          onClick={() => setActiveTab('members')}
        >
          Members
        </button>
        <button
          className={`chat-tab ${activeTab === 'donate' ? 'active' : ''}`}
          onClick={() => setActiveTab('donate')}
        >
          Donate
        </button>
        {isLeader && (
          <button
            className={`chat-tab ${activeTab === 'manage' ? 'active' : ''}`}
            onClick={() => setActiveTab('manage')}
          >
            Manage
          </button>
        )}
      </>
    )
  }

  const renderContent = () => {
    if (!isInCorp) {
      if (activeTab === 'create') {
        return (
          <div style={{ padding: '20px' }}>
            <h3>Create a New Corp</h3>
            <div style={{ marginTop: '16px' }}>
              <label style={{ display: 'block', marginBottom: '8px' }}>
                Corp Name:
                <input
                  type="text"
                  value={createForm.name}
                  onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                  style={{
                    width: '100%',
                    padding: '8px',
                    marginTop: '4px',
                    background: 'rgba(0,0,0,0.3)',
                    border: '1px solid #333',
                    borderRadius: '4px',
                    color: '#fff',
                  }}
                  placeholder="Enter corp name"
                  maxLength={50}
                />
              </label>
              <label style={{ display: 'block', marginTop: '12px', marginBottom: '8px' }}>
                Corp Tag (3-5 chars):
                <input
                  type="text"
                  value={createForm.tag}
                  onChange={(e) => setCreateForm({ ...createForm, tag: e.target.value.toUpperCase() })}
                  style={{
                    width: '100%',
                    padding: '8px',
                    marginTop: '4px',
                    background: 'rgba(0,0,0,0.3)',
                    border: '1px solid #333',
                    borderRadius: '4px',
                    color: '#fff',
                  }}
                  placeholder="ABC"
                  maxLength={5}
                />
              </label>
              <label style={{ display: 'block', marginTop: '12px', marginBottom: '8px' }}>
                Description:
                <textarea
                  value={createForm.description}
                  onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                  style={{
                    width: '100%',
                    padding: '8px',
                    marginTop: '4px',
                    background: 'rgba(0,0,0,0.3)',
                    border: '1px solid #333',
                    borderRadius: '4px',
                    color: '#fff',
                    minHeight: '80px',
                  }}
                  placeholder="Describe your corp..."
                  maxLength={200}
                />
              </label>
              <LoadingButton
                className="btn btn-primary"
                onClick={handleCreate}
                loading={loading}
                disabled={
                  !createForm.name.trim() ||
                  createForm.tag.length < 3 ||
                  createForm.tag.length > 5
                }
                style={{ marginTop: '16px' }}
              >
                Create Corp
              </LoadingButton>
            </div>
          </div>
        )
      }

      if (activeTab === 'join') {
        return (
          <div style={{ padding: '20px' }}>
            <h3>Join a Corp</h3>
            <div style={{ marginTop: '16px' }}>
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => handleSearch(e.target.value)}
                style={{
                  width: '100%',
                  padding: '8px',
                  background: 'rgba(0,0,0,0.3)',
                  border: '1px solid #333',
                  borderRadius: '4px',
                  color: '#fff',
                }}
                placeholder="Search corps by name or tag..."
              />
              <div style={{ marginTop: '16px', maxHeight: '300px', overflowY: 'auto' }}>
                {searchLoading && <div className="inline-spinner" />}
                {!searchLoading && searchResults.length === 0 && searchQuery && (
                  <p style={{ color: '#888' }}>No corps found.</p>
                )}
                {!searchLoading && searchResults.length === 0 && !searchQuery && (
                  <p style={{ color: '#888' }}>Type to search for corps.</p>
                )}
                {!searchLoading &&
                  searchResults.map((corp) => (
                    <div
                      key={corp.id}
                      style={{
                        padding: '12px',
                        marginBottom: '8px',
                        background: 'rgba(0,0,0,0.3)',
                        border: '1px solid #333',
                        borderRadius: '4px',
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                      }}
                    >
                      <div>
                        <div style={{ fontWeight: 'bold' }}>
                          [{corp.tag}] {corp.name}
                        </div>
                        <div style={{ fontSize: '0.85rem', color: '#aaa', marginTop: '4px' }}>
                          Level {corp.level} | {corp.member_count}/{corp.max_members} members
                        </div>
                        {corp.description && (
                          <div style={{ fontSize: '0.8rem', color: '#999', marginTop: '4px' }}>
                            {corp.description}
                          </div>
                        )}
                      </div>
                      <LoadingButton
                        className="btn btn-primary"
                        onClick={() => handleJoin(corp.id)}
                        loading={loading}
                      >
                        Join
                      </LoadingButton>
                    </div>
                  ))}
              </div>
            </div>
          </div>
        )
      }
    }

    // In corp tabs
    if (activeTab === 'overview') {
      return (
        <div style={{ padding: '20px' }}>
          <h3>Corp Overview</h3>
          {corpData?.corp && (
            <div style={{ marginTop: '16px' }}>
              <div style={{ marginBottom: '12px' }}>
                <strong>Name:</strong> [{corpData.corp.tag}] {corpData.corp.name}
              </div>
              <div style={{ marginBottom: '12px' }}>
                <strong>Level:</strong> {corpData.corp.level}
              </div>
              <div style={{ marginBottom: '12px' }}>
                <strong>Wealth:</strong> {corpData.corp.wealth.toLocaleString()}
              </div>
              <div style={{ marginBottom: '12px' }}>
                <strong>Members:</strong> {members.length}/{corpData.corp.max_members}
              </div>
              {corpData.bonuses && (
                <>
                  <div style={{ marginBottom: '12px' }}>
                    <strong>Controlled RBPs:</strong> {corpData.bonuses.rbp_count}
                  </div>
                  <div style={{ marginBottom: '12px' }}>
                    <strong>Total RBP Bonus:</strong> +{corpData.bonuses.total_rbp_bonus}%
                  </div>
                </>
              )}
              <div style={{ marginBottom: '12px' }}>
                <strong>Description:</strong> {corpData.corp.description || 'No description'}
              </div>
              {!isLeader && (
                <LoadingButton
                  className="btn btn-secondary"
                  onClick={handleLeave}
                  loading={loading}
                  style={{ marginTop: '16px' }}
                >
                  Leave Corp
                </LoadingButton>
              )}
            </div>
          )}
        </div>
      )
    }

    if (activeTab === 'members') {
      return (
        <div style={{ padding: '20px' }}>
          <h3>Corp Members</h3>
          <div style={{ marginTop: '16px', maxHeight: '300px', overflowY: 'auto' }}>
            {membersLoading && <div className="inline-spinner" />}
            {!membersLoading && members.length === 0 && (
              <p style={{ color: '#888' }}>No members found.</p>
            )}
            {!membersLoading &&
              members.map((member) => (
                <div
                  key={member.player_id}
                  style={{
                    padding: '12px',
                    marginBottom: '8px',
                    background: 'rgba(0,0,0,0.3)',
                    border: '1px solid #333',
                    borderRadius: '4px',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <div>
                      <div style={{ fontWeight: 'bold' }}>{member.player_name}</div>
                      <div style={{ fontSize: '0.85rem', color: '#aaa', marginTop: '4px' }}>
                        {member.role.charAt(0).toUpperCase() + member.role.slice(1)} | Contribution: {member.contribution_points.toLocaleString()}
                      </div>
                      <div style={{ fontSize: '0.8rem', color: '#999', marginTop: '2px' }}>
                        Daily: {member.daily_contribution.toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
          </div>
        </div>
      )
    }

    if (activeTab === 'donate') {
      const myMember = members.length > 0 ? members.find(m => {
        // Find current player by checking state.player
        return state.player && m.player_id === state.player.id
      }) : null
      const dailyUsed = myMember?.daily_contribution ?? 0
      const dailyMax = 200
      const dailyRemaining = dailyMax - dailyUsed
      const dailyPct = Math.min(100, (dailyUsed / dailyMax) * 100)

      return (
        <div style={{ padding: '20px' }}>
          <h3>Donate Resources</h3>
          <div style={{ marginTop: '16px' }}>
            <div style={{
              padding: '12px',
              background: 'rgba(0,0,0,0.3)',
              border: '1px solid #333',
              borderRadius: '4px',
              marginBottom: '16px',
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px', fontSize: '0.85rem' }}>
                <span>Daily Contribution: {dailyUsed}/{dailyMax} pts</span>
                <span style={{ color: dailyRemaining <= 0 ? '#f87171' : dailyRemaining < 50 ? '#fbbf24' : '#4ade80' }}>
                  {dailyRemaining > 0 ? `${dailyRemaining} pts remaining` : 'Limit reached'}
                </span>
              </div>
              <div style={{ height: '6px', background: '#333', borderRadius: '3px', overflow: 'hidden' }}>
                <div style={{
                  height: '100%',
                  width: `${dailyPct}%`,
                  background: dailyPct >= 95 ? '#f87171' : dailyPct >= 80 ? '#fbbf24' : '#4ade80',
                  borderRadius: '3px',
                  transition: 'width 0.3s',
                }} />
              </div>
            </div>
            <p style={{ color: '#aaa', fontSize: '0.9rem', marginBottom: '16px' }}>
              Donate resources to your corp to earn contribution points and increase corp wealth. (1 pt per 10,000 resources)
            </p>
            <label style={{ display: 'block', marginBottom: '8px' }}>
              Metal:
              <input
                type="number"
                value={donateForm.metal}
                onChange={(e) => setDonateForm({ ...donateForm, metal: parseInt(e.target.value) || 0 })}
                style={{
                  width: '100%',
                  padding: '8px',
                  marginTop: '4px',
                  background: 'rgba(0,0,0,0.3)',
                  border: '1px solid #333',
                  borderRadius: '4px',
                  color: '#fff',
                }}
                min={0}
              />
            </label>
            <label style={{ display: 'block', marginTop: '12px', marginBottom: '8px' }}>
              He3:
              <input
                type="number"
                value={donateForm.he3}
                onChange={(e) => setDonateForm({ ...donateForm, he3: parseInt(e.target.value) || 0 })}
                style={{
                  width: '100%',
                  padding: '8px',
                  marginTop: '4px',
                  background: 'rgba(0,0,0,0.3)',
                  border: '1px solid #333',
                  borderRadius: '4px',
                  color: '#fff',
                }}
                min={0}
              />
            </label>
            <label style={{ display: 'block', marginTop: '12px', marginBottom: '8px' }}>
              Gold:
              <input
                type="number"
                value={donateForm.gold}
                onChange={(e) => setDonateForm({ ...donateForm, gold: parseInt(e.target.value) || 0 })}
                style={{
                  width: '100%',
                  padding: '8px',
                  marginTop: '4px',
                  background: 'rgba(0,0,0,0.3)',
                  border: '1px solid #333',
                  borderRadius: '4px',
                  color: '#fff',
                }}
                min={0}
              />
            </label>
            <LoadingButton
              className="btn btn-primary"
              onClick={handleDonate}
              loading={loading}
              disabled={donateForm.metal === 0 && donateForm.he3 === 0 && donateForm.gold === 0}
              style={{ marginTop: '16px' }}
            >
              Donate
            </LoadingButton>
          </div>
        </div>
      )
    }

    if (activeTab === 'manage' && isLeader) {
      return (
        <div style={{ padding: '20px' }}>
          <h3>Manage Members</h3>
          <div style={{ marginTop: '16px', maxHeight: '300px', overflowY: 'auto' }}>
            {membersLoading && <div className="inline-spinner" />}
            {!membersLoading && members.length === 0 && (
              <p style={{ color: '#888' }}>No members found.</p>
            )}
            {!membersLoading &&
              members
                .filter((m) => m.role !== 'leader')
                .map((member) => (
                  <div
                    key={member.player_id}
                    style={{
                      padding: '12px',
                      marginBottom: '8px',
                      background: 'rgba(0,0,0,0.3)',
                      border: '1px solid #333',
                      borderRadius: '4px',
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                    }}
                  >
                    <div>
                      <div style={{ fontWeight: 'bold' }}>{member.player_name}</div>
                      <div style={{ fontSize: '0.85rem', color: '#aaa', marginTop: '4px' }}>
                        {member.role.charAt(0).toUpperCase() + member.role.slice(1)} | Contribution: {member.contribution_points.toLocaleString()}
                      </div>
                    </div>
                    <select
                      value={member.role}
                      onChange={(e) => handleRoleChange(member.player_id, e.target.value)}
                      style={{
                        padding: '6px 12px',
                        background: 'rgba(0,0,0,0.3)',
                        border: '1px solid #333',
                        borderRadius: '4px',
                        color: '#fff',
                        cursor: 'pointer',
                      }}
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
      className="chat-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div className="panel chat-panel">
        <div className="panel-header">
          <h2>Corps</h2>
          <button className="p2-modal-close" onClick={onClose}>
            X
          </button>
        </div>
        <div className="chat-channel-tabs">{renderTabs()}</div>
        <div className="panel-content" style={{ overflowY: 'auto' }}>
          {renderContent()}
        </div>
      </div>
    </div>,
    document.body,
  )
}
