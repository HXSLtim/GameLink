import * as Icons from '@ant-design/icons';

type IconsType = typeof Icons;

export const getIcon = (iconName?: string) => {
    if (!iconName) return null;

    // Handle specific mappings or generic lookup
    const IconComponent = (Icons as IconsType)[iconName as keyof IconsType] || 
                          (Icons as IconsType)[(iconName + 'Outlined') as keyof IconsType];

    if (IconComponent && typeof IconComponent === 'function') {
        const Icon = IconComponent as React.ComponentType;
        return <Icon />;
    }

    return null;
};
