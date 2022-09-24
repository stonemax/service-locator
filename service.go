package ServiceLocator

import (
	"sync"
)

// scopeMap 作用域MAP
var scopeMap = map[int]struct{}{
	Singleton: {},
	PerLookup: {},
	Immediate: {},
}

// ServiceCreator 服务创建者
type ServiceCreator func(*ServiceLocator) (interface{}, error)

// Service 服务
type Service struct {
	// name 名称
	name string

	// scopes 作用域
	scopes int

	// creator 创建者
	creator ServiceCreator

	// 创建错误
	createError error

	// instance 已生成实例，仅用于单例模式
	instance interface{}

	// onceLocker 单次锁
	onceLocker *sync.Once
}

//
// InScope
// @desc 设置作用域
// @receiver s *Service
// @param scope int
//
func (s *Service) InScope(scope int) {
	if _, ok := scopeMap[scope]; !ok {
		return
	}

	if s.HasScope(scope) {
		return
	}

	s.scopes += scope
}

//
// HasScope
// @desc 检查是有拥有作用域
// @receiver s *Service
// @param scope int
// @return bool
//
func (s *Service) HasScope(scope int) bool {
	return (s.scopes & scope) == scope
}

//
// NewService
// @desc 创建服务
// @param name string
// @param creator ServiceCreator
// @param scopeList []int
// @return *Service
//
func NewService(name string, creator ServiceCreator, scopeList ...int) *Service {
	service := &Service{
		name:    name,
		creator: creator,
	}

	if len(scopeList) > 0 {
		for _, scope := range scopeList {
			if _, ok := scopeMap[scope]; !ok {
				continue
			}

			service.scopes += scope
		}
	}

	if service.scopes == 0 {
		service.scopes = Singleton
	}

	if (service.scopes & Singleton) == Singleton {
		service.onceLocker = new(sync.Once)
	}

	return service
}
