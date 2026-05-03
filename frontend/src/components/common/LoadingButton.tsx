interface LoadingButtonProps {
  onClick: () => void
  loading?: boolean
  disabled?: boolean
  className?: string
  children: React.ReactNode
  type?: 'button' | 'submit'
}

/**
 * Button with an inline spinner. Pure inline styling so it carries no
 * dependency on any global stylesheet — works on top of any DS button class.
 */
export default function LoadingButton({
  onClick,
  loading = false,
  disabled = false,
  className = 'ds-btn-primary',
  children,
  type = 'button',
}: LoadingButtonProps) {
  return (
    <button
      type={type}
      className={className}
      onClick={onClick}
      disabled={loading || disabled}
      aria-busy={loading || undefined}
      style={loading ? { cursor: 'not-allowed', opacity: 0.7 } : undefined}
    >
      {loading && (
        <span
          aria-hidden="true"
          style={{
            display: 'inline-block',
            width: 14,
            height: 14,
            borderRadius: '50%',
            border: '2px solid currentColor',
            borderRightColor: 'transparent',
            animation: 'lb-spin 0.6s linear infinite',
            marginRight: 8,
            verticalAlign: 'middle',
          }}
        />
      )}
      <span style={{ verticalAlign: 'middle' }}>{children}</span>
      <style>{`@keyframes lb-spin { to { transform: rotate(360deg); } }`}</style>
    </button>
  )
}
