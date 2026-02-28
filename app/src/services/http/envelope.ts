export type UnknownRecord = Record<string, unknown>

export interface ApiEnvelope<T = unknown> {
  success?: boolean
  code?: number
  message?: string
  data?: T
}

const DEFAULT_LIST_KEYS = ["items", "list", "records", "rows", "players", "orders", "groups", "messages", "reviews", "transactions"]

export function isRecord(value: unknown): value is UnknownRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value)
}

export function readString(record: UnknownRecord, key: string): string | undefined {
  const value = record[key]
  return typeof value === "string" ? value : undefined
}

export function readNumber(record: UnknownRecord, key: string): number | undefined {
  const value = record[key]
  if (typeof value === "number" && Number.isFinite(value)) {
    return value
  }
  if (typeof value === "string" && value.trim() !== "") {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : undefined
  }
  return undefined
}

export function unwrapApiData(payload: unknown): unknown {
  if (!isRecord(payload)) {
    throw new Error("响应格式无效")
  }

  if (payload.success === false) {
    throw new Error(readString(payload, "message") ?? "请求失败")
  }

  return "data" in payload ? payload.data : payload
}

export function extractRecordList(value: unknown, preferredKeys: string[] = []): UnknownRecord[] {
  if (Array.isArray(value)) {
    return value.filter(isRecord)
  }

  if (!isRecord(value)) {
    return []
  }

  const candidateKeys = [...preferredKeys, ...DEFAULT_LIST_KEYS]
  for (const key of candidateKeys) {
    const candidate = value[key]
    if (Array.isArray(candidate)) {
      return candidate.filter(isRecord)
    }
  }

  if ("data" in value && value.data !== value) {
    return extractRecordList(value.data, preferredKeys)
  }

  return []
}
