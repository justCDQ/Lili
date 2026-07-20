export function formatCents(value, currency) {
    if (!Number.isSafeInteger(value)) {
        throw new TypeError("value 必须是安全整数");
    }
    return new Intl.NumberFormat(currency === "CNY" ? "zh-CN" : "en-US", {
        style: "currency",
        currency,
    }).format(value / 100);
}
