import { httpClient } from "@/services/http/client"
import { extractRecordList, unwrapApiData, type UnknownRecord } from "@/services/http/envelope"

export async function listOrders(): Promise<UnknownRecord[]> {
  const response = await httpClient.get<unknown>("/user/orders")
  return extractRecordList(unwrapApiData(response.data), ["items", "orders"])
}
