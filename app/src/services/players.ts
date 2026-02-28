import { httpClient } from "@/services/http/client"
import { extractRecordList, unwrapApiData, type UnknownRecord } from "@/services/http/envelope"

export async function listPlayers(): Promise<UnknownRecord[]> {
  const response = await httpClient.get<unknown>("/user/players")
  return extractRecordList(unwrapApiData(response.data), ["items", "players"])
}
