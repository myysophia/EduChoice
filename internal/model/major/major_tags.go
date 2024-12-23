package major

// MajorTag 专业特征标签
type MajorTag struct {
    Name        string   // 标签名称
    Keywords    []string // 关键词
    Weight      float64  // 权重
    Description string   // 描述
}

// 专业特征标签库
var MajorTags = map[string]MajorTag{
    "理论型": {
        Name:     "理论型",
        Keywords: []string{"研究", "理论", "分析", "科学"},
        Weight:   0.8,
        Description: "偏重理论研究和学术探索",
    },
    "应用型": {
        Name:     "应用型",
        Keywords: []string{"应用", "实践", "操作", "技术"},
        Weight:   0.8,
        Description: "注重实践应用和技能培养",
    },
    "创新型": {
        Name:     "创新型",
        Keywords: []string{"创新", "创造", "设计", "研发"},
        Weight:   0.7,
        Description: "强调创新能力和创造思维",
    },
    "管理型": {
        Name:     "管理型",
        Keywords: []string{"管理", "决策", "规划", "组织"},
        Weight:   0.7,
        Description: "培养管理能力和��策能力",
    },
    "服务型": {
        Name:     "服务型",
        Keywords: []string{"服务", "咨询", "指导", "帮助"},
        Weight:   0.6,
        Description: "注重服务意识和沟通能力",
    },
}

// 专业-标签映射
var MajorTagMap = map[string][]string{
    "计算机科学与技术": {"理论型", "应用型", "创新型"},
    "软件工程":      {"应用型", "创新型"},
    "人工智能":      {"理论型", "创新型"},
    "数据科学":      {"理论型", "应用型"},
    // ... 可以继续添加更多专业
}

// GetMajorTags 获取专业标签
func GetMajorTags(majorName string) []MajorTag {
    tagNames, exists := MajorTagMap[majorName]
    if !exists {
        return nil
    }
    
    tags := make([]MajorTag, 0, len(tagNames))
    for _, name := range tagNames {
        if tag, ok := MajorTags[name]; ok {
            tags = append(tags, tag)
        }
    }
    return tags
} 