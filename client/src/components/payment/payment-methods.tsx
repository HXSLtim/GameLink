import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { Wallet, CreditCard, QrCode } from "lucide-react";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

export type PaymentMethodType = 'balance' | 'wechat' | 'alipay';

interface PaymentMethodsProps {
    value: PaymentMethodType;
    onChange: (value: PaymentMethodType) => void;
    balance?: number;
    amount?: number;
    className?: string;
}

export function PaymentMethods({ value, onChange, balance = 0, amount = 0, className }: PaymentMethodsProps) {
    const { t } = useTranslation();
    const isBalanceSufficient = balance >= amount;

    return (
        <RadioGroup value={value} onValueChange={(v: string) => onChange(v as PaymentMethodType)} className={cn("grid gap-4", className)}>
            {/* Wallet Balance */}
            <div className={cn(
                "relative flex items-center space-x-4 rounded-xl border p-4 transition-all hover:bg-muted/50",
                value === 'balance' ? "border-primary bg-primary/5 ring-1 ring-primary" : "border-border",
                !isBalanceSufficient && "opacity-60 cursor-not-allowed"
            )}>
                <RadioGroupItem value="balance" id="balance" className="sr-only" disabled={!isBalanceSufficient} />
                <Label htmlFor="balance" className="flex flex-1 items-center space-x-4 cursor-pointer">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary">
                        <Wallet className="h-5 w-5" />
                    </div>
                    <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium leading-none">{t('payment.wallet_balance', { defaultValue: 'Wallet Balance' })}</p>
                        <p className="text-xs text-muted-foreground">
                            {t('payment.balance_available', { defaultValue: 'Available' })}: ¥{balance.toFixed(2)}
                            {!isBalanceSufficient && <span className="text-destructive ml-2">({t('payment.insufficient', { defaultValue: 'Insufficient' })})</span>}
                        </p>
                    </div>
                    <div className={cn(
                        "h-4 w-4 rounded-full border border-primary/50",
                        value === 'balance' && "border-primary bg-primary"
                    )}>
                        {value === 'balance' && <div className="h-full w-full rounded-full bg-primary" />}
                    </div>
                </Label>
            </div>

            {/* WeChat Pay */}
            <div className={cn(
                "relative flex items-center space-x-4 rounded-xl border p-4 transition-all hover:bg-muted/50",
                value === 'wechat' ? "border-[#07C160] bg-[#07C160]/5 ring-1 ring-[#07C160]" : "border-border"
            )}>
                <RadioGroupItem value="wechat" id="wechat" className="sr-only" />
                <Label htmlFor="wechat" className="flex flex-1 items-center space-x-4 cursor-pointer">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#07C160]/10 text-[#07C160]">
                        <QrCode className="h-5 w-5" />
                    </div>
                    <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium leading-none">{t('payment.wechat_pay', { defaultValue: 'WeChat Pay' })}</p>
                        <p className="text-xs text-muted-foreground">{t('payment.wechat_desc', { defaultValue: 'Scan QR code to pay' })}</p>
                    </div>
                    <div className={cn(
                        "h-4 w-4 rounded-full border",
                        value === 'wechat' ? "border-[#07C160] bg-[#07C160]" : "border-primary/50"
                    )} />
                </Label>
            </div>

            {/* Alipay */}
            <div className={cn(
                "relative flex items-center space-x-4 rounded-xl border p-4 transition-all hover:bg-muted/50",
                value === 'alipay' ? "border-[#1677FF] bg-[#1677FF]/5 ring-1 ring-[#1677FF]" : "border-border"
            )}>
                <RadioGroupItem value="alipay" id="alipay" className="sr-only" />
                <Label htmlFor="alipay" className="flex flex-1 items-center space-x-4 cursor-pointer">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#1677FF]/10 text-[#1677FF]">
                        <CreditCard className="h-5 w-5" />
                    </div>
                    <div className="flex-1 space-y-1">
                        <p className="text-sm font-medium leading-none">{t('payment.alipay', { defaultValue: 'Alipay' })}</p>
                        <p className="text-xs text-muted-foreground">{t('payment.alipay_desc', { defaultValue: 'Secure connection' })}</p>
                    </div>
                    <div className={cn(
                        "h-4 w-4 rounded-full border",
                        value === 'alipay' ? "border-[#1677FF] bg-[#1677FF]" : "border-primary/50"
                    )} />
                </Label>
            </div>
        </RadioGroup>
    );
}
