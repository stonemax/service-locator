package ServiceLocator

// 作用域
const (
	// Singleton 单例
	Singleton = 1 << iota

	// PerLookup 多实例：每次查询均创建服务
	PerLookup

	// Immediate 立即创建：仅作用域单例服务
	Immediate
)
