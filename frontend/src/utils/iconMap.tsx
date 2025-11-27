import * as Icons from '@ant-design/icons';

export const getIcon = (iconName?: string) => {
    if (!iconName) return null;

    // Handle specific mappings or generic lookup
    const IconComponent = (Icons as any)[iconName] || (Icons as any)[iconName + 'Outlined'];

    if (IconComponent) {
        return <IconComponent />;
    }

    return null;
};
