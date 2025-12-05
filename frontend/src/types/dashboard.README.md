# Dashboard Module Type Definitions

This file contains comprehensive TypeScript type definitions for the Dashboard Module.

## Setup Complete ✓

### Dependencies Installed
- ✅ **Recharts** (v3.5.0) - Already installed
- ✅ **Socket.IO Client** (v4.8.1) - Newly installed for real-time data

### Type Definitions Created
- ✅ **frontend/src/types/dashboard.ts** - Complete type definitions
- ✅ **frontend/src/constants/dashboard.ts** - Constants and configuration

## Type Categories

### Core Types
- `TimeRange` - Time range selection ('7d' | '30d' | '90d')
- `TrendData` - Trend comparison data with percentage changes
- `DashboardStats` - Main dashboard statistics container

### Order Types
- `OrderStatus` - Order status enumeration
- `OrderStatusData` - Order status distribution data
- `Order` - Order entity with full details

### Chart Data Types
- `RevenueData` - Revenue trend data points
- `UserGrowthData` - User growth metrics
- `OperationalTrendData` - Operational metrics over time

### Real-time Monitoring
- `RealtimeMetrics` - System monitoring metrics
- `Alert` - Alert/notification entity
- `AlertType` & `AlertLevel` - Alert classification

### KPI Types
- `KPIMetrics` - Key performance indicators
- `KPIMetric` - Individual KPI metric with targets

### Component Props
All component prop types are defined for:
- StatCard
- RevenueChart
- OrderStatusPie
- UserGrowthChart
- RecentOrders
- TopPlayers
- RealtimeMonitor
- AlertBanner
- KPIPanel
- OperationalOverview
- Dashboard (main)

## Constants Available

### Time Ranges
- `TIME_RANGE_OPTIONS` - Dropdown options for time selection
- `DEFAULT_TIME_RANGE` - Default time range ('7d')

### Status Labels & Colors
- `ORDER_STATUS_LABELS` - Chinese labels for order statuses
- `ORDER_STATUS_COLORS` - Color codes for each status
- `ALERT_TYPE_LABELS` - Alert type labels
- `ALERT_LEVEL_COLORS` - Alert severity colors

### Chart Configuration
- `CHART_COLORS` - Standard chart color palette
- `REVENUE_CHART_COLORS` - Revenue-specific colors
- `USER_GROWTH_CHART_COLORS` - User growth colors

### Monitoring
- `MONITORING_THRESHOLDS` - Warning/critical thresholds for metrics
- `REFRESH_INTERVALS` - Auto-refresh intervals
- `WEBSOCKET_CONFIG` - WebSocket connection settings

### Formatting
- `NUMBER_FORMAT` - Number formatting configuration
- `DATE_FORMATS` - Date/time format strings
- `ANIMATION_DURATION` - Animation timing constants

## Usage Example

```typescript
import type { DashboardStats, TimeRange } from '@/types/dashboard';
import { TIME_RANGE_OPTIONS, ORDER_STATUS_COLORS } from '@/constants/dashboard';

const MyComponent = () => {
  const [timeRange, setTimeRange] = useState<TimeRange>('7d');
  const [stats, setStats] = useState<DashboardStats | null>(null);
  
  // Use constants
  const options = TIME_RANGE_OPTIONS;
  const statusColor = ORDER_STATUS_COLORS['completed'];
  
  return <div>...</div>;
};
```

## Next Steps

With the foundation in place, you can now:
1. Implement individual dashboard components
2. Create API service functions
3. Build custom hooks for data fetching
4. Implement WebSocket connections for real-time data

## Requirements Satisfied

✅ **Requirement 1.1**: TypeScript types configured for all dashboard components
✅ **Recharts**: Chart library ready for data visualization
✅ **Socket.IO**: Real-time data connection capability installed
