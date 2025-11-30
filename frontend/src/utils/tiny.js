export function numToChinese(num) {
  if (num === 0) return '零';

  const digits = ['零', '一', '二', '三', '四', '五', '六', '七', '八', '九'];
  const units = ['', '十', '百', '千'];
  const bigUnits = ['', '万', '亿'];

  let result = '';
  let bigUnitIndex = 0;

  while (num > 0) {
    const section = num % 10000;
    if (section > 0) {
      let sectionStr = '';
      let sectionUnitIndex = 0;
      let temp = section;

      while (temp > 0) {
        const digit = temp % 10;
        if (digit > 0) {
          sectionStr = digits[digit] + units[sectionUnitIndex] + sectionStr;
        } else if (sectionStr && !sectionStr.startsWith('零')) {
          sectionStr = '零' + sectionStr;
        }
        temp = Math.floor(temp / 10);
        sectionUnitIndex++;
      }

      result = sectionStr + bigUnits[bigUnitIndex] + result;
    }
    bigUnitIndex++;
    num = Math.floor(num / 10000);
  }

  // 处理 "一十..." → "十..."
  if (result.startsWith('一十')) {
    result = result.substring(1);
  }

  return result;
}