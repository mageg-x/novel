// 生成一次性请求token
export async function genToken() {
  const secret = "secret-key-rao-bang-lin!"; // ← 必须和后端一致！

  const timestamp = Math.floor(Date.now() / 1000); // Unix 秒

  // 构造 8 字节大端时间戳
  const tsBuffer = new ArrayBuffer(8);
  new DataView(tsBuffer).setBigUint64(0, BigInt(timestamp), false); // big-endian

  // 计算 HMAC-SHA256
  const key = await crypto.subtle.importKey(
    "raw",
    new TextEncoder().encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"]
  );

  const signature = await crypto.subtle.sign(
    "HMAC",
    key,
    new Uint8Array(tsBuffer)
  );

  // 拼接：时间戳 + 签名
  const combined = new Uint8Array(40);
  combined.set(new Uint8Array(tsBuffer), 0);
  combined.set(new Uint8Array(signature), 8);

  // Base64URL 编码（无 padding）
  return btoa(String.fromCharCode(...combined))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

// 数字转中文
export function numToChinese(num) {
  if (num === 0) return "零";

  const digits = ["零", "一", "二", "三", "四", "五", "六", "七", "八", "九"];
  const units = ["", "十", "百", "千"];
  const bigUnits = ["", "万", "亿"];

  let result = "";
  let bigUnitIndex = 0;

  while (num > 0) {
    const section = num % 10000;
    if (section > 0) {
      let sectionStr = "";
      let sectionUnitIndex = 0;
      let temp = section;

      while (temp > 0) {
        const digit = temp % 10;
        if (digit > 0) {
          sectionStr = digits[digit] + units[sectionUnitIndex] + sectionStr;
        } else if (sectionStr && !sectionStr.startsWith("零")) {
          sectionStr = "零" + sectionStr;
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
  if (result.startsWith("一十")) {
    result = result.substring(1);
  }

  return result;
}
