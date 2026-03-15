import http from 'k6/http';
import { check } from 'k6';

// 1. 压测配置：100个并发虚拟用户，总共发 1000 次请求
export const options = {
    vus: 500,         // Virtual Users (并发数)
    iterations: 1000, // 总请求次数
};

// 2. 压测主逻辑
export default function () {
    const url = 'http://127.0.0.1:8888/douyin/user/refresh';
    
    // 准备 Payload
    const payload = JSON.stringify({
        // 👇 替换为你真实有效的 refresh_token
        refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM3NTQzOTksImlhdCI6MTc3MzE0OTU5OSwianRpIjoiZTg3OTZmMzUtNDg1YS00MDY0LTkwMDMtMDZlNTliMGMyMjhhIiwidXNlcklkIjoxMDMwNjQ1NzIwMzEyNDUxMDcyfQ.RMBO48TLkuunelTKysFhkdR1Oeqxkwt-wPV6g2ymR2o" 
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // 发起 POST 请求
    const res = http.post(url, payload, params);

    // 3. 核心：业务断言 (Checks)
    // 假设你的接口成功时返回 "status_code": "0"，重放拦截时返回 "status_code": "10005"
    check(res, {
        // 第一层：网络与框架层有没有崩？(不崩就是 200)
        '网络连通正常 (HTTP 200)': (r) => r.status === 200,
        
        // 第二层：业务逻辑层剖析
        '👑 成功拿到新Token (仅应有1次)': (r) => r.json("status_code") === "0",
        '🛡️ 惨遭重放拦截 (应占绝大多数)': (r) => r.json("status_code") === "10005",
        '❓ 其他业务异常 (过期/Token错误)': (r) => r.json("status_code") !== "0" && r.json("status_code") !== "10005",
    });
}