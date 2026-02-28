import { httpClient } from "@/services/http/client"
import { isRecord, unwrapApiData, type UnknownRecord } from "@/services/http/envelope"

export async function getMyProfile(): Promise<UnknownRecord | null> {
  const response = await httpClient.get<unknown>("/user/profile")
  const data = unwrapApiData(response.data)
  return isRecord(data) ? data : null
}
