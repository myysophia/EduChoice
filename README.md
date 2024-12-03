# EduChoice


## 问题
### 1. json.Unmarshal 无法正确解析。
major_score_his映射整个json过程中json.Unmarshal 无法正确解析。这种情况通常发生在 JSON 中的数字可能被表示为字符串，或者反之，以及可能存在的空值或特殊格式。
为了彻底解决这个问题，我们可以采用更灵活的解析策略，使用 interface{} 来接收所有字段，然后根据实际类型进行转换。这种方法可以避免类型不匹配的问题。