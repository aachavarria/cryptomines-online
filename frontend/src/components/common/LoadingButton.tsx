import '../../styles/common.css'

interface LoadingButtonProps {
  onClick: () => void
  loading?: boolean
  disabled?: boolean
  className?: string
  children: React.ReactNode
  type?: 'button' | 'submit'
}

export default function LoadingButton({
  onClick,
  loading = false,
  disabled = false,
  className = 'btn btn-primary',
  children,
  type = 'button',
}: LoadingButtonProps) {
  return (
    <button
      type={type}
      className={`${className} ${loading ? 'loading' : ''}`}
      onClick={onClick}
      disabled={loading || disabled}
    >
      {loading && <span className="spinner" />}
      <span className={loading ? 'loading-text' : ''}>{children}</span>
    </button>
  )
}
