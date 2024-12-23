# 专业推荐相关接口

## 获取专业推荐
GET /api/v1/major/recommend

### 请求参数
Query参数:
- year: 年份，如"2023" (可选，默认当年)
- type: 考生类型，如"理科"/"文科" (可选，默认从用户信息获取)

Header参数:
- Authorization: Bearer {token}

### 响应数据
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "chong_schools": [
            {
                "name": "学校名称",
                "history_infos": {
                    "2023": {
                        "lowest_score": 600,
                        "lowest_rank": 1000,
                        "enrollment_num": 100
                    }
                },
                "parts": {
                    "提前批": [
                        {
                            "name": "专业名称",
                            "rate": 0.95,  // 专业匹配度
                            "weight": 0.9, // 专业权重
                            "match_result": {
                                "score": 0.85,
                                "factors": [
                                    {
                                        "name": "兴趣匹配度",
                                        "score": 0.8,
                                        "weight": 0.3,
                                        "analysis": "与您的兴趣较为匹配"
                                    }
                                ],
                                "suggestions": [
                                    "该专业与您的特点高度匹配，建议考虑报考"
                                ]
                            }
                        }
                    ]
                }
            }
        ],
        "wen_schools": [],
        "bao_schools": []
    }
}
``` 