
import { Button } from "@/components/ui/button"
import { PageContainer, PageHeader, PageHeaderSkeleton } from "@/components/page-container"
import { Skeleton } from "@/components/ui/skeleton"
import { Separator } from "@/components/ui/separator"
import { Card, CardContent } from "@/components/ui/card"
import { useState } from "react"
import { Loader2 } from "lucide-react"

export default function PageShowcase() {
    const [loading, setLoading] = useState(false);

    const toggleLoading = () => {
        setLoading(true);
        setTimeout(() => setLoading(false), 2000);
    }

    if (loading) {
        return (
            <PageContainer>
                <PageHeaderSkeleton />
                <Separator />
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {Array.from({ length: 9 }).map((_, i) => (
                        <div key={i} className="flex flex-col space-y-3">
                            <Skeleton className="h-[125px] w-full rounded-xl" />
                            <div className="space-y-2">
                                <Skeleton className="h-4 w-[250px]" />
                                <Skeleton className="h-4 w-[200px]" />
                            </div>
                        </div>
                    ))}
                </div>
            </PageContainer>
        )
    }

    return (
        <PageContainer>
            <PageHeader
                title="Page Title"
                description="This describes the standard page structure (Header + Content)."
                action={
                    <Button onClick={toggleLoading}>
                        {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                        Simulate Loading
                    </Button>
                }
            />
            <Separator />
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {Array.from({ length: 6 }).map((_, i) => (
                    <Card key={i}>
                        <CardContent className="h-[125px] flex items-center justify-center p-6 text-muted-foreground">
                            Content Card {i + 1}
                        </CardContent>
                    </Card>
                ))}
            </div>
        </PageContainer>
    )
}
