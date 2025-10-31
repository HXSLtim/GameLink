-- Migrate legacy order status spelling
UPDATE orders SET status='canceled' WHERE status='cancelled';
