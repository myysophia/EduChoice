package major

// ZYMockResp 保持与原有模拟数据格式兼容
type ZYMockResp struct {
    ChongSchools []School `json:"chong_schools"`
    WenSchools   []School `json:"wen_schools"`
    BaoSchools   []School `json:"bao_schools"`
}

// ToMockResp 将推荐结果转换为模拟数据格式
func (r *RecommendResp) ToMockResp() *ZYMockResp {
    return &ZYMockResp{
        ChongSchools: r.ChongSchools,
        WenSchools:   r.WenSchools,
        BaoSchools:   r.BaoSchools,
    }
} 