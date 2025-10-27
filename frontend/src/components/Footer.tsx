import { Layout, Typography } from '@arco-design/web-react';
import styles from './Footer.module.less';

/**
 * Footer component
 *
 * @component
 * @description Displays the application footer with copyright information
 */
export const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <Layout.Footer className={styles.footer}>
      <Typography.Text type="secondary">© {currentYear} GameLink Admin</Typography.Text>
    </Layout.Footer>
  );
};
