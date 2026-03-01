import { describe, expect, it } from "vitest"
import { extractRecordList, isRecord, readNumber, unwrapApiData } from "./envelope"

describe("envelope helpers", () => {
  it("unwraps nested data payload", () => {
    const payload = {
      success: true,
      data: {
        items: [{ id: 1 }, { id: 2 }],
      },
    }

    const data = unwrapApiData(payload)
    expect(isRecord(data)).toBe(true)
    expect(extractRecordList(data, ["items"])).toHaveLength(2)
  })

  it("throws when api reports failure", () => {
    expect(() =>
      unwrapApiData({
        success: false,
        message: "boom",
      })
    ).toThrow("boom")
  })

  it("extracts list from nested data containers", () => {
    const rows = extractRecordList(
      {
        data: {
          records: [
            { id: "a" },
            null,
            { id: "b" },
          ],
        },
      },
      ["items"]
    )
    expect(rows).toEqual([{ id: "a" }, { id: "b" }])
  })

  it("throws on invalid payload shape", () => {
    expect(() => unwrapApiData("bad-payload")).toThrow("响应格式无效")
  })

  it("reads numeric values from number or numeric strings", () => {
    expect(readNumber({ v: 7 }, "v")).toBe(7)
    expect(readNumber({ v: "8" }, "v")).toBe(8)
    expect(readNumber({ v: "  " }, "v")).toBeUndefined()
    expect(readNumber({ v: "abc" }, "v")).toBeUndefined()
  })
})
