import { Card } from '../../components/Card/Card';
import styles from './Reports.module.less';

export const ReportDashboard = () => {
  return (
    <div className={styles.container}>
      <Card title="报表中心" description="关键业务指标与趋势概览">
        <p>报表功能即将上线，敬请期待。</p>
      </Card>
    </div>
  );
};

export default ReportDashboard;
