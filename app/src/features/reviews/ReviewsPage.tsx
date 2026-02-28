import { useEffect, useState, type ReactElement } from "react"
import { Button } from "@/components/ui/button"
import { readNumber, readString, type UnknownRecord } from "@/services/http/envelope"
import { listMyReviews } from "@/services/reviews"

export function ReviewsPage(): ReactElement {
  const [reviews, setReviews] = useState<UnknownRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadReviews = async () => {
      setLoading(true)
      setError(null)
      try {
        setReviews(await listMyReviews())
      } catch (fetchError) {
        const message = fetchError instanceof Error ? fetchError.message : "获取评价失败"
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadReviews()
  }, [])

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <div className="flex items-center justify-between">
        <h1 className="m-0 text-xl font-semibold text-slate-900">我的评价</h1>
        <Button variant="outline" onClick={() => void listMyReviews().then(setReviews).catch(() => setError("刷新失败"))}>
          刷新
        </Button>
      </div>
      {loading ? <p className="m-0 text-sm text-slate-600">加载中...</p> : null}
      {error ? <p className="m-0 text-sm text-red-700">{error}</p> : null}
      {!loading && reviews.length === 0 ? <p className="m-0 text-sm text-slate-600">暂无评价</p> : null}
      <ul className="m-0 list-none space-y-2 p-0">
        {reviews.map((review, index) => {
          const reviewId = readNumber(review, "id")
          const rating = readNumber(review, "rating") ?? readNumber(review, "score")
          const content = readString(review, "content") ?? readString(review, "comment") ?? "(无评价内容)"
          return (
            <li key={reviewId ?? index} className="rounded-lg border border-slate-200 p-3 text-sm">
              <p className="m-0 font-medium text-slate-900">评价#{reviewId ?? "-"} · 评分: {rating ?? "-"}</p>
              <p className="m-0 text-slate-600">{content}</p>
            </li>
          )
        })}
      </ul>
    </section>
  )
}
