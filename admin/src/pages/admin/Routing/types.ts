/**
 * Routing Rule Management - Type Definitions
 */

import type {
    RoutingRule,
    RoutingCondition,
    CollectionEntity,
    ConditionField,
    ConditionOperator,
    RuleStatus,
} from '@/api/routing';

// Re-export types from API
export type {
    RoutingRule,
    RoutingCondition,
    CollectionEntity,
    ConditionField,
    ConditionOperator,
    RuleStatus,
};

/**
 * Form value types for condition builder
 */
export interface ConditionFormItem {
    id: string;
    field: ConditionField;
    operator: ConditionOperator;
    value: string | number | string[] | number[];
}

/**
 * Field option for condition builder dropdown
 */
export interface FieldOption {
    value: ConditionField;
    label: string;
    type: 'string' | 'number' | 'select' | 'multiselect';
    options?: Array<{ label: string; value: string }>;
    operators: ConditionOperator[];
}

/**
 * Operator option for condition builder dropdown
 */
export interface OperatorOption {
    value: ConditionOperator;
    label: string;
    requiresMultiple: boolean;
}

/**
 * Priority reorder item
 */
export interface PriorityReorderItem {
    id: number;
    name: string;
    priority: number;
    status: RuleStatus;
}

/**
 * Test request form values
 */
export interface TestFormValues {
    gameType?: string;
    serviceType?: string;
    amountCents?: number;
    region?: string;
}

/**
 * Form values for creating/editing routing rule
 */
export interface RoutingRuleFormValues {
    name: string;
    description?: string;
    priority: number;
    targetEntityId: number;
    conditions: ConditionFormItem[];
    status?: RuleStatus;
}
