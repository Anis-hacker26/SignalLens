type ErrorBannerProps = {
  message: string
}

export default function ErrorBanner({
  message,
}: ErrorBannerProps) {
  return (
    <div className="common-error" role="alert">
      {message}
    </div>
  )
}
