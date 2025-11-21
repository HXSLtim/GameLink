/**
 * URL参数处理工具
 */

/**
 * 从URL查询参数中提取值并合并到初始参数
 *
 * @param searchParams - URLSearchParams对象
 * @param initialParams - 初始参数对象
 * @param paramKeys - 需要从URL读取的参数键数组
 * @returns 合并后的参数对象
 */
export function mergeUrlParams<T extends Record<string, any>>(
  searchParams: URLSearchParams,
  initialParams: T,
  paramKeys: string[],
): T {
  const urlParams: Partial<T> = {};

  paramKeys.forEach((key) => {
    const value = searchParams.get(key);
    if (value !== null && value !== '') {
      // 尝试转换为数字（如果是数字字符串）
      if (!isNaN(Number(value))) {
        urlParams[key as keyof T] = Number(value) as any;
      } else {
        urlParams[key as keyof T] = value as any;
      }
    }
  });

  return {
    ...initialParams,
    ...urlParams,
  };
}

