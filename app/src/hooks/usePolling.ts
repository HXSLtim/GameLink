import { useCallback, useEffect, useRef, useState } from "react"

export interface PollingOptions<T> {
  /** Polling function that fetches data */
  fetcher: () => Promise<T>
  /** Condition to stop polling (returns true to stop) */
  shouldStop: (data: T | null, error: Error | null) => boolean
  /** Initial delay in milliseconds */
  initialDelay?: number
  /** Maximum delay in milliseconds */
  maxDelay?: number
  /** Backoff multiplier (e.g., 1.5 means 50% increase each retry) */
  backoffFactor?: number
  /** Maximum total polling duration in milliseconds */
  timeout?: number
  /** Whether to start polling immediately */
  immediate?: boolean
  /** Callback when data is fetched successfully */
  onSuccess?: (data: T) => void
  /** Callback when an error occurs */
  onError?: (error: Error) => void
  /** Callback when polling stops */
  onStop?: (reason: "completed" | "timeout" | "error" | "manual") => void
}

export interface PollingState<T> {
  data: T | null
  error: Error | null
  isLoading: boolean
  isPolling: boolean
  attemptCount: number
  currentDelay: number
}

export function usePolling<T>(options: PollingOptions<T>): {
  state: PollingState<T>
  start: () => void
  stop: () => void
  reset: () => void
} {
  const {
    fetcher,
    shouldStop,
    initialDelay = 3000,
    maxDelay = 30000,
    backoffFactor = 1.5,
    timeout = 300000, // 5 minutes default
    immediate = true,
    onSuccess,
    onError,
    onStop,
  } = options

  const [state, setState] = useState<PollingState<T>>({
    data: null,
    error: null,
    isLoading: false,
    isPolling: false,
    attemptCount: 0,
    currentDelay: initialDelay,
  })

  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const startTimeRef = useRef<number | null>(null)
  const currentDelayRef = useRef(initialDelay)
  const isPollingRef = useRef(false)
  const executePollRef = useRef<(() => Promise<void>) | null>(null)
  const mountedRef = useRef(true)

  const clearCurrentTimer = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
      timeoutRef.current = null
    }
  }, [])

  const stop = useCallback(
    (reason: "completed" | "timeout" | "error" | "manual" = "manual") => {
      clearCurrentTimer()
      isPollingRef.current = false
      if (mountedRef.current) {
        setState((prev) => ({ ...prev, isPolling: false, isLoading: false }))
      }
      onStop?.(reason)
    },
    [clearCurrentTimer, onStop],
  )

  const scheduleNextPoll = useCallback(
    (delay: number) => {
      clearCurrentTimer()
      if (!mountedRef.current || !isPollingRef.current) return
      timeoutRef.current = setTimeout(() => {
        void executePollRef.current?.()
      }, delay)
    },
    [clearCurrentTimer],
  )

  const executePoll = useCallback(async () => {
    if (!mountedRef.current || !isPollingRef.current) return

    // Check timeout
    if (startTimeRef.current && Date.now() - startTimeRef.current >= timeout) {
      stop("timeout")
      return
    }

    setState((prev) => ({ ...prev, isLoading: true }))

    try {
      const data = await fetcher()

      if (!mountedRef.current) return

      setState((prev) => ({
        ...prev,
        data,
        error: null,
        isLoading: false,
        attemptCount: prev.attemptCount + 1,
      }))

      onSuccess?.(data)

      if (shouldStop(data, null)) {
        stop("completed")
        return
      }

      // Schedule next poll with exponential backoff
      const nextDelay = Math.min(currentDelayRef.current * backoffFactor, maxDelay)
      currentDelayRef.current = nextDelay
      setState((prev) => ({ ...prev, currentDelay: nextDelay }))
      scheduleNextPoll(nextDelay)
    } catch (err) {
      if (!mountedRef.current) return

      const error = err instanceof Error ? err : new Error(String(err))

      setState((prev) => ({
        ...prev,
        error,
        isLoading: false,
        attemptCount: prev.attemptCount + 1,
      }))

      onError?.(error)

      if (shouldStop(null, error)) {
        stop("error")
        return
      }

      // Continue polling on error with backoff
      const nextDelay = Math.min(currentDelayRef.current * backoffFactor, maxDelay)
      currentDelayRef.current = nextDelay
      setState((prev) => ({ ...prev, currentDelay: nextDelay }))
      scheduleNextPoll(nextDelay)
    }
  }, [fetcher, shouldStop, timeout, backoffFactor, maxDelay, stop, onSuccess, onError, scheduleNextPoll])

  useEffect(() => {
    executePollRef.current = executePoll
  }, [executePoll])

  const start = useCallback(() => {
    if (state.isPolling) return

    clearCurrentTimer()
    startTimeRef.current = Date.now()
    currentDelayRef.current = initialDelay
    isPollingRef.current = true

    setState((prev) => ({
      ...prev,
      isPolling: true,
      currentDelay: initialDelay,
    }))

    void executePoll()
  }, [state.isPolling, clearCurrentTimer, initialDelay, executePoll])

  const reset = useCallback(() => {
    clearCurrentTimer()
    startTimeRef.current = null
    currentDelayRef.current = initialDelay
    isPollingRef.current = false
    setState({
      data: null,
      error: null,
      isLoading: false,
      isPolling: false,
      attemptCount: 0,
      currentDelay: initialDelay,
    })
  }, [clearCurrentTimer, initialDelay])

  useEffect(() => {
    mountedRef.current = true

    if (immediate) {
      start()
    }

    return () => {
      mountedRef.current = false
      clearCurrentTimer()
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return { state, start, stop: () => stop("manual"), reset }
}
