import { Component, type ErrorInfo, type ReactNode } from "react";
import { Button } from "@/components/ui/button";
import { AlertCircle, RefreshCw } from "lucide-react";

interface Props {
    children?: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div className="min-h-screen flex flex-col items-center justify-center p-4 bg-background text-foreground animate-in fade-in duration-500">
                    <div className="max-w-md w-full text-center space-y-6">
                        <div className="w-20 h-20 bg-destructive/10 text-destructive rounded-full flex items-center justify-center mx-auto mb-6">
                            <AlertCircle className="w-10 h-10" />
                        </div>

                        <h1 className="text-3xl font-bold tracking-tight">Something went wrong</h1>

                        <div className="bg-muted/40 p-4 rounded-lg text-left overflow-auto max-h-[200px] border border-border/50">
                            <p className="text-sm font-mono text-muted-foreground break-all">
                                {this.state.error?.message || "An unexpected error occurred."}
                            </p>
                        </div>

                        <div className="flex justify-center gap-4">
                            <Button
                                onClick={() => window.location.reload()}
                                className="gap-2 min-w-[140px]"
                                size="lg"
                            >
                                <RefreshCw className="w-4 h-4" />
                                Reload Page
                            </Button>
                            <Button
                                variant="outline"
                                onClick={() => window.location.href = '/'}
                                size="lg"
                            >
                                Back to Home
                            </Button>
                        </div>

                        <p className="text-xs text-muted-foreground pt-4">
                            If the problem persists, please contact support.
                        </p>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
