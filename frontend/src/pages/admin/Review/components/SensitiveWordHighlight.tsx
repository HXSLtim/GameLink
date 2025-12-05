/**
 * 敏感词高亮组件
 * 需求: 2.4
 */
import React from 'react';
import { Tooltip } from 'antd';
import type { SensitiveWordCategory, SensitiveWordSeverity } from '@/types/review';
import {
  SENSITIVE_WORD_CATEGORY_TEXT,
  SENSITIVE_WORD_SEVERITY_TEXT,
} from '@/types/review';

interface DetectedWord {
  word: string;
  category: SensitiveWordCategory;
  severity: SensitiveWordSeverity;
  positions: number[];
}

interface SensitiveWordHighlightProps {
  content: string;
  detectedWords: DetectedWord[];
}

// 根据严重程度获取背景颜色
const getSeverityColor = (severity: SensitiveWordSeverity): string => {
  const colorMap: Record<SensitiveWordSeverity, string> = {
    low: '#fff7e6',      // 浅橙色
    medium: '#fff1f0',   // 浅红色
    high: '#ffccc7',     // 红色
  };
  return colorMap[severity] || '#fff7e6';
};

// 根据严重程度获取边框颜色
const getSeverityBorderColor = (severity: SensitiveWordSeverity): string => {
  const colorMap: Record<SensitiveWordSeverity, string> = {
    low: '#ffc069',      // 橙色
    medium: '#ff7875',   // 浅红色
    high: '#ff4d4f',     // 红色
  };
  return colorMap[severity] || '#ffc069';
};

const SensitiveWordHighlight: React.FC<SensitiveWordHighlightProps> = ({
  content,
  detectedWords,
}) => {
  if (!content || !detectedWords || detectedWords.length === 0) {
    return <span>{content || '-'}</span>;
  }

  // 构建需要高亮的位置映射
  const highlightMap = new Map<number, DetectedWord>();
  
  detectedWords.forEach(detected => {
    // 查找所有匹配位置
    let startIndex = 0;
    while (true) {
      const index = content.toLowerCase().indexOf(detected.word.toLowerCase(), startIndex);
      if (index === -1) break;
      
      // 标记每个字符的位置
      for (let i = 0; i < detected.word.length; i++) {
        highlightMap.set(index + i, detected);
      }
      startIndex = index + 1;
    }
  });

  // 渲染内容
  const renderContent = () => {
    const result: React.ReactNode[] = [];
    let currentWord: DetectedWord | null = null;
    let currentText = '';
    let normalText = '';

    for (let i = 0; i <= content.length; i++) {
      const detected = highlightMap.get(i);

      if (i === content.length || detected !== currentWord) {
        // 输出普通文本
        if (normalText) {
          result.push(<span key={`normal-${i}`}>{normalText}</span>);
          normalText = '';
        }

        // 输出高亮文本
        if (currentText && currentWord) {
          result.push(
            <Tooltip
              key={`highlight-${i}`}
              title={
                <div>
                  <div>分类: {SENSITIVE_WORD_CATEGORY_TEXT[currentWord.category]}</div>
                  <div>严重程度: {SENSITIVE_WORD_SEVERITY_TEXT[currentWord.severity]}</div>
                </div>
              }
            >
              <span
                style={{
                  backgroundColor: getSeverityColor(currentWord.severity),
                  border: `1px solid ${getSeverityBorderColor(currentWord.severity)}`,
                  borderRadius: 2,
                  padding: '0 2px',
                  cursor: 'pointer',
                }}
              >
                {currentText}
              </span>
            </Tooltip>
          );
          currentText = '';
        }

        currentWord = detected || null;
      }

      if (i < content.length) {
        if (detected) {
          currentText += content[i];
        } else {
          normalText += content[i];
        }
      }
    }

    return result;
  };

  return <span>{renderContent()}</span>;
};

export default SensitiveWordHighlight;
