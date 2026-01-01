/**
 * 陪玩师认证主页面
 * 包含实名认证和段位认证两个子页面
 */
import React, { useState } from 'react';
import { Tabs, Card } from 'antd';
import {
    SafetyOutlined,
    TrophyOutlined,
} from '@ant-design/icons';
import IdentityCertification from './Identity';
import RankCertification from './Rank';

const { TabPane } = Tabs;

/**
 * 陪玩师认证主页面
 */
const CertificationPage: React.FC = () => {
    const [activeTab, setActiveTab] = useState('identity');

    return (
        <div style={{ padding: 24 }}>
            <Card>
                <Tabs
                    activeKey={activeTab}
                    onChange={setActiveTab}
                    items={[
                        {
                            key: 'identity',
                            label: (
                                <span>
                                    <SafetyOutlined />
                                    实名认证
                                </span>
                            ),
                            children: <IdentityCertification />,
                        },
                        {
                            key: 'rank',
                            label: (
                                <span>
                                    <TrophyOutlined />
                                    段位认证
                                </span>
                            ),
                            children: <RankCertification />,
                        },
                    ]}
                />
            </Card>
        </div>
    );
};

export default CertificationPage;
