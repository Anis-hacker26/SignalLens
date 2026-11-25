type EmptyStateProps = {
  message: string
}

export default function EmptyState({
  message,
}: EmptyStateProps) {
  return (
    <div className="common-state">
      {message}
    </div>
  )
}
