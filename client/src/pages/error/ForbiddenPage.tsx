import { Button } from '@/components/ui/button';
import { useNavigate } from 'react-router-dom';

export default function ForbiddenPage() {
    const navigate = useNavigate();

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-background text-center p-4">
            <h1 className="text-4xl font-bold mb-4">403</h1>
            <p className="text-xl mb-8">Access Forbidden</p>
            <p className="text-muted-foreground mb-8">
                You do not have permission to access this page.
            </p>
            <Button onClick={() => navigate('/')}>
                Go Home
            </Button>
        </div>
    );
}
