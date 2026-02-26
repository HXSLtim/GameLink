import type { PaymentStatus } from "@/services/payment/types"

const paymentIdKeys = ["paymentId", "payment_id", "id", "pid", "paymentNo"]
const statusKeys = ["status", "trade_status", "resultStatus", "payStatus"]
const nestedPayloadKeys = ["attach", "extra", "metadata", "passback_params"]

export interface PaymentCallbackData {
  paymentId: number | null
  status: PaymentStatus | null
}

function normalizeStatus(raw: string | null): PaymentStatus | null {
  if (!raw) {
    return null
  }

  const value = raw.toLowerCase()

  if (value === "paid" || value === "success" || value === "succeeded" || value === "trade_success") {
    return "paid"
  }

  if (value === "failed" || value === "fail" || value === "error" || value === "closed") {
    return "failed"
  }

  if (value === "refunded" || value === "refund_success") {
    return "refunded"
  }

  if (value === "pending" || value === "processing" || value === "wait_buyer_pay") {
    return "pending"
  }

  return null
}

function parsePositiveInt(input: string | null): number | null {
  if (!input) {
    return null
  }

  const parsed = Number(input)
  return Number.isInteger(parsed) && parsed > 0 ? parsed : null
}

function getFirst(searchParams: URLSearchParams, keys: string[]): string | null {
  for (const key of keys) {
    const value = searchParams.get(key)
    if (value) {
      return value
    }
  }
  return null
}

function parseNestedPayload(searchParams: URLSearchParams): URLSearchParams | null {
  for (const key of nestedPayloadKeys) {
    const raw = searchParams.get(key)
    if (!raw) {
      continue
    }

    try {
      const json = JSON.parse(raw) as Record<string, unknown>
      const nested = new URLSearchParams()
      for (const [k, v] of Object.entries(json)) {
        if (typeof v === "string" || typeof v === "number" || typeof v === "boolean") {
          nested.set(k, String(v))
        }
      }
      return nested
    } catch {
      const nested = new URLSearchParams(raw)
      if (Array.from(nested.keys()).length > 0) {
        return nested
      }
    }
  }

  return null
}

export function parsePaymentCallback(searchParams: URLSearchParams): PaymentCallbackData {
  const directPaymentId = parsePositiveInt(getFirst(searchParams, paymentIdKeys))
  const directStatus = normalizeStatus(getFirst(searchParams, statusKeys))

  if (directPaymentId) {
    return { paymentId: directPaymentId, status: directStatus }
  }

  const nested = parseNestedPayload(searchParams)
  if (!nested) {
    return { paymentId: null, status: directStatus }
  }

  const nestedPaymentId = parsePositiveInt(getFirst(nested, paymentIdKeys))
  const nestedStatus = normalizeStatus(getFirst(nested, statusKeys))

  return {
    paymentId: nestedPaymentId,
    status: directStatus ?? nestedStatus,
  }
}
