/**
 * 主题感知的消息/弹窗 Hook
 * 
 * Ant Design 5.x 中，直接使用 message.xxx() 和 Modal.xxx() 静态方法
 * 不会继承 ConfigProvider 的主题配置。
 * 
 * 使用此 hook 获取主题感知的 message、modal、notification 实例。
 * 
 * @example
 * const { message, modal, notification } = useAppMessage();
 * message.success('操作成功');
 * modal.confirm({ title: '确认', content: '确定要删除吗？' });
 */
import { App } from 'antd';

export const useAppMessage = () => {
    const { message, modal, notification } = App.useApp();
    return { message, modal, notification };
};

export default useAppMessage;
