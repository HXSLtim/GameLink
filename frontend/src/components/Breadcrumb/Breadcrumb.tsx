import { Link } from 'react-router-dom';
import styles from './Breadcrumb.module.less';

export interface BreadcrumbItem {
  /** 标题 */
  label: string;
  /** 路径 */
  path?: string;
}

export interface BreadcrumbProps {
  /** 面包屑项 */
  items: BreadcrumbItem[];
  /** 分隔符 */
  separator?: string;
}

// 箭头图标
const ArrowIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <polyline
      points="9 18 15 12 9 6"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const Breadcrumb: React.FC<BreadcrumbProps> = ({ items, separator }) => {
  return (
    <nav className={styles.breadcrumb} aria-label="面包屑导航">
      <ol className={styles.list}>
        {items.map((item, index) => {
          const isLast = index === items.length - 1;

          return (
            <li key={index} className={styles.item}>
              {item.path && !isLast ? (
                <Link to={item.path} className={styles.link}>
                  {item.label}
                </Link>
              ) : (
                <span className={`${styles.label} ${isLast ? styles.current : ''}`}>
                  {item.label}
                </span>
              )}

              {!isLast && <span className={styles.separator}>{separator || <ArrowIcon />}</span>}
            </li>
          );
        })}
      </ol>
    </nav>
  );
};
