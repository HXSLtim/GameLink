import { httpClient } from "@/services/http/client"
import { extractRecordList, isRecord, unwrapApiData, type UnknownRecord } from "@/services/http/envelope"

export async function getWalletBalance(): Promise<UnknownRecord | null> {
  const response = await httpClient.get<unknown>("/user/wallet/balance")
  const data = unwrapApiData(response.data)
  return isRecord(data) ? data : null
}

export async function listWalletTransactions(): Promise<UnknownRecord[]> {
  const response = await httpClient.get<unknown>("/user/wallet/transactions")
  return extractRecordList(unwrapApiData(response.data), ["items", "transactions"])
}
