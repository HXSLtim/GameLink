import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { X, AlertTriangle, CheckCircle, Info, AlertCircle } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

export interface ModalProps {
  /** 是否打开 */
  open: boolean;
  /** 关闭回调 */
  onOpenChange: (open: boolean) => void;
  /** 标题 */
  title?: string;
  /** 描述 */
  description?: string;
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
  /** 是否显示关闭按钮 */
  showClose?: boolean;
  /** 是否可通过点击遮罩关闭 */
  closeOnOverlay?: boolean;
  /** 页脚 */
  footer?: React.ReactNode;
  /** 内容 */
  children?: React.ReactNode;
  /** 自定义类名 */
  className?: string;
}

const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  full: 'max-w-[90vw] h-[90vh]',
};

export const Modal = forwardRef<HTMLDivElement, ModalProps>(({
  open,
  onOpenChange,
  title,
  description,
  size = 'md',
  showClose = true,
  closeOnOverlay = true,
  footer,
  children,
  className,
}, ref) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        ref={ref}
        className={cn(sizeClasses[size], className)}
        onPointerDownOutside={(e) => {
          if (!closeOnOverlay) {
            e.preventDefault();
          }
        }}
      >
        {(title || description) && (
          <DialogHeader>
            {title && <DialogTitle>{title}</DialogTitle>}
            {description && <DialogDescription>{description}</DialogDescription>}
          </DialogHeader>
        )}
        
        {showClose && (
          <button
            onClick={() => onOpenChange(false)}
            className="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
          >
            <X className="h-4 w-4" />
            <span className="sr-only">关闭</span>
          </button>
        )}
        
        <div className="py-2">
          {children}
        </div>
        
        {footer && (
          <DialogFooter>
            {footer}
          </DialogFooter>
        )}
      </DialogContent>
    </Dialog>
  );
});

Modal.displayName = 'Modal';

// ConfirmModal 组件 - 确认对话框
export interface ConfirmModalProps extends Omit<ModalProps, 'footer' | 'children'> {
  /** 类型 */
  type?: 'info' | 'success' | 'warning' | 'danger';
  /** 消息内容 */
  message?: string;
  /** 确认按钮文本 */
  confirmText?: string;
  /** 取消按钮文本 */
  cancelText?: string;
  /** 确认回调 */
  onConfirm?: () => void | Promise<void>;
  /** 取消回调 */
  onCancel?: () => void;
  /** 是否加载中 */
  loading?: boolean;
  /** 自定义图标 */
  icon?: LucideIcon;
}

const typeConfig: Record<string, { icon: LucideIcon; color: string; buttonVariant: 'default' | 'destructive' }> = {
  info: { icon: Info, color: 'text-blue-500', buttonVariant: 'default' },
  success: { icon: CheckCircle, color: 'text-green-500', buttonVariant: 'default' },
  warning: { icon: AlertTriangle, color: 'text-amber-500', buttonVariant: 'default' },
  danger: { icon: AlertCircle, color: 'text-red-500', buttonVariant: 'destructive' },
};

export function ConfirmModal({
  open,
  onOpenChange,
  type = 'info',
  title,
  description,
  message,
  confirmText = '确认',
  cancelText = '取消',
  onConfirm,
  onCancel,
  loading = false,
  icon,
  size = 'sm',
  ...props
}: ConfirmModalProps) {
  const config = typeConfig[type];
  const Icon = icon || config.icon;

  const handleConfirm = async () => {
    if (onConfirm) {
      await onConfirm();
    }
    onOpenChange(false);
  };

  const handleCancel = () => {
    onCancel?.();
    onOpenChange(false);
  };

  return (
    <Modal
      open={open}
      onOpenChange={onOpenChange}
      title={title}
      description={description}
      size={size}
      showClose={false}
      footer={
        <div className="flex justify-end gap-2 w-full">
          <Button variant="outline" onClick={handleCancel} disabled={loading}>
            {cancelText}
          </Button>
          <Button
            variant={config.buttonVariant}
            onClick={handleConfirm}
            disabled={loading}
          >
            {loading && <span className="animate-spin mr-2">⟳</span>}
            {confirmText}
          </Button>
        </div>
      }
      {...props}
    >
      <div className="flex items-start gap-4">
        <div className={cn('p-2 rounded-full bg-muted', config.color)}>
          <Icon className="w-6 h-6" />
        </div>
        <div className="flex-1">
          {message && (
            <p className="text-sm text-muted-foreground">{message}</p>
          )}
        </div>
      </div>
    </Modal>
  );
}
