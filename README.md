# DreamBridge - 高考志愿智能推荐系统

## 最近更新
- 新增学校历年招生计划数据采集功能
  - 支持按年份、批次、科类采集
  - 包含专业名称、招生人数、学费等信息
  - 支持多线程并发采集
  - 自动重试和代理切换机制
- 优化数据采集错误处理
  - 改进反爬限制检测
  - 优化错误提示信息
  - 完善数据格式验证

## 系统架构

### 整体架构

## 功能特性

### 1. 智能代理管理
- 支持快代理和公开代理源的自动切换
- 代理自动检测与切换功能
- 支持配置化的代理开关控制
- 支持本地代理配置文件

### 2. 数据采集系统
- 学校基本信息采集
- 专业分数线历史数据采集
- 招生计划数据采集
- 支持多源数据整合

### 3. 智能推荐系统
- 基于用户分数的院校匹配
- 考虑地理位置偏好
- 专业兴趣匹配
- 历史数据分析支持

### 4. 数据库设计
- 学校信息表(school)
- 专业信息表(major)
- 历史分数线表(major_score)
- 招生计划表(school_num)
- 用户信息表(user)

## 技术实现要点

### 1. 代理系统
- 多代理源支持
  - 快代理(付费)
  - ProxyScrape API(免费)
- 代理自动切换机制
- 代理可用性检测

### 2. 数据采集
- 分布式爬虫系统
- 数据清洗和标准化
- 错误重试机制
- 并发控制

### 3. 数据存储
- MySQL数据库
- 数据模型设计
- 数据完整性保证

### 4. 配置管理
- TOML配置文件
- 运行时配置切换
- 环境变量支持

## 部署要求
- Go 1.23+
- MySQL 5.7+
- 配置文件设置
  - proxy.switchon: 代理开关
  - proxy.kuaidaili: 快代理优先开关


## 问题
### 1. json.Unmarshal 无法正确解析。
major_score_his映射整个json过程中json.Unmarshal 无法正确解析。这种情况通常发生在 JSON 中的数字可能被表示为字符串，或者反之，以及可能存在的空值或特殊格式。
为了彻底解决这个问题，我们可以采用更灵活的解析策略，使用 interface{} 来接收所有字段，然后根据实际类型进行转换。这种方法可以避免类型不匹配的问题。

### 2. 反爬限制检测
- 检测访问频率限制 (code 1069)
- 自动识别 IP 封禁
- 准确的错误提示

### 3. 数据验证
- JSON 格式验证
- 字段类型检查
- 数据完整性验证

### 4. 错误恢复
- 自动重试机制
- 代理切换策略
- 错误日志记录

### 5. 最佳实践
- 错误提示应准确反映问题本质
- 区分反爬限制和数据格式错误
- 提供有效的错误诊断信息

## Go语言知识点

### 1. 类型断言与错误处理
在 `must/special_scores_his.go` 中的错误处理提供了很好的学习案例：
```go
// 先检查是否被反爬
if code, ok := rawResponse["code"].(string); ok && code == "1069" {
return nil, fmt.Errorf("访问频率限制: %s", rawResponse["message"])
}
// 检查data字段
dataMap, ok := rawResponse["data"].(map[string]interface{})
if !ok {
return nil, fmt.Errorf("data字段类型不匹配: 期望 map[string]interface{}, 实际类型 %T, 值: %+v",
rawResponse["data"], rawResponse["data"])
}
```

知识要点：
1. 类型断言基础
   - 语法：`value, ok := x.(T)`
   - ok 表示断言是否成功
   - 失败时 value 为类型零值

2. 错误处理最佳实践
   - 错误信息应该包含上下文
   - 使用 fmt.Errorf 格式化错误
   - 按优先级检查错误（如先检查反爬）

3. 接口类型处理
   - interface{} 可以接收任何类型
   - 使用类型断言获取具体类型
   - 处理 map[string]interface{} 常用模式

4. 调试技巧
   - 使用 %T 打印类型信息
   - 使用 %+v 打印详细值
   - 分层次处理复杂数据结构

这种错误处理方式特别适合处理：
- API 响应解析
- JSON 数据处理
- 类型安全转换
- 分级错误处理

### 2. 并发安全与超时控制
在 `safe/school.go` 中的 GetSchoolInfoSafe 函数提供了并发编程的典型案例。这个函数的主要作用是安全地获取学校信息，它通过以下机制来保证可靠性：
- 并发请求：同时发起两个相同的请求，提高成功率
- 超时控制：设置5秒超时，避免请求阻塞
- 错误恢复：自动处理panic，确保程序稳定性
- 自动重试：失败后自动重新尝试

代码实现：
```go
func GetSchoolInfoSafe(schoolId int) *response.SchoolInfoResponse {
    var info *response.SchoolInfoResponse
this:
    for {
        infoChan := make(chan *response.SchoolInfoResponse, 1)
        for i := 0; i < 2; i++ {
            go func() {
                defer func() {
                    if r := recover(); r != nil {
                        common.LOG.Error(fmt.Sprintf("%v", r))
                    }
                }()
                info, err := scraper.SchoolInfo(schoolId)
                if err != nil || info == nil {
                    return
                }
                infoChan <- info
            }()
        }
        ticker := time.NewTicker(5 * time.Second)
        select {
        case info = <-infoChan:
            break this
        case <-ticker.C:
            common.LOG.Error("get school info: time out 5 s")
            proxy.ChangeHttpProxyIP()
        }
    }
    return info
}
```

知识要点分析：

1. goroutine 并发设计
   - 项目使用场景：同时发起两个相同请求，谁先成功用谁的结果
   - 实现方式：使用 for 循环启动两个 goroutine
   - 优势：提高请求成功率，减少等待时间
   - 注意事项：需要考虑 goroutine 泄漏问题

2. channel 通信机制
   - 项目使用场景：goroutine 之间传递学校信息数据
   - 实现方式：使用带缓冲的 channel（容量为1）
   - 优势：避免 goroutine 阻塞，确保数据安全传递
   - 实际应用：infoChan := make(chan *response.SchoolInfoResponse, 1)

3. select 多路复用
   - 项目使用场景：同时处理数据接收和超时检测
   - 实现方式：使用 select 监听多个 channel
   - 优势：优雅地处理超时情况，避免请求卡死
   - 实际应用：在接收数据和定时器之间进行选择

4. defer 错误恢复
   - 项目使用场景：防止爬虫请求panic导致程序崩溃
   - 实现方式：在每个 goroutine 中使用 defer recover
   - 优势：保证程序稳定性，不会因单个请求失败而中断
   - 实际应用：捕获并记录错误，允许程序继续运行

5. 代理自动切换
   - 项目使用场景：请求超时时自动更换代理IP
   - 实现方式：在超时分支中调用 proxy.ChangeHttpProxyIP()
   - 优势：提高请求成功率，绕过反爬限制
   - 实际应用：配合定时器实现自动切换

这种并发设计模式在项目中的其他场景也有应用：
- 专业分数线数据采集
- 招生计划信息获取
- 学校列表更新
- 任何需要可靠性和容错性的网络请求

通过这种设计，我们实现了：
1. 高可用性：多次尝试和自动重试
2. 容错性：错误恢复和超时处理
3. 性能优化：并发请求和超时控制
4. 反爬处理：代理自动切换

### 3. defer 使用注意事项

1. defer 的基本作用
   - 将函数调用推迟到当前函数返回之前执行
   - 常用于资源清理、解锁等操作
   - defer 语句会按照后进先出（LIFO）的顺序执行

2. 在循环中使用 defer 的问题
```go
go
// 错误示例
func badExample() {
for {
ticker := time.NewTicker(1 time.Second)
defer ticker.Stop() // 这个 defer 要等函数返回才会执行，而不是在每次循环结束时执行
// ...
}
}
// 正确示例
func goodExample() {
for {
ticker := time.NewTicker(1 time.Second)
select {
case <-ticker.C:
ticker.Stop() // 立即停止不再需要的 ticker
// 处理逻辑
}
}
}

3. 主要注意点：
   - defer 语句只在函数返回时执行，而不是在代码块或循环结束时
   - 在循环中使用 defer 可能导致资源泄漏
   - 对于需要及时清理的资源，应该在不再需要时立即进行清理，而不是依赖 defer

4. 最佳实践：
   - 在函数级别使用 defer 处理资源清理
   - 在循环中需要立即手动清理资源
   - 特别注意定时器、文件句柄等系统资源的及时释放

### 4. 数据库事务处理最佳实践

在处理数据库事务时，需要确保事务的完整性和一致性。以下是一个典型的事务处理示例：

```go
func CreateRecords(records []*Record) error {
    tx := db.Begin()
    // 使用 defer 确保事务一定会被处理
    committed := false
    defer func() {
        if !committed {
            tx.Rollback()
        }
    }()

    for _, record := range records {
        if err := tx.Create(record).Error; err != nil {
            return err
        }
    }

    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("提交事务失败: %v", err)
    }
    committed = true
    return nil
}
```

主要注意点：
1. 事务状态跟踪
   - 使用 committed 标志记录事务状态
   - defer 配合标志确保事务正确处理

2. 错误处理
   - 创建失败时及时返回错误
   - 提交失败时提供详细错误信息

3. 资源清理
   - 使用 defer 确保事务一定会被处理
   - 避免事务悬而未决

4. 最佳实践
   - 在函数开始时初始化事务
   - 使用 defer 处理回滚
   - 仅在所有操作都成功后提交
   - 检查提交操作的返回值
