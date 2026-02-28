import { httpClient } from "@/services/http/client"
import { extractRecordList, unwrapApiData, type UnknownRecord } from "@/services/http/envelope"

export async function listMyReviews(): Promise<UnknownRecord[]> {
  const response = await httpClient.get<unknown>("/user/reviews/my")
  return extractRecordList(unwrapApiData(response.data), ["items", "reviews"])
}
