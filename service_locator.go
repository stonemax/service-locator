package ServiceLocator

import (
	"fmt"
)

// ServiceLocator 服务定位器
type ServiceLocator struct {
	// services 已注册服务
	services map[string]*Service
}

//
// Get
// @desc 获取服务
// @receiver sl *ServiceLocator
// @param name string
// @return interface{}
// @return error
//
func (sl *ServiceLocator) Get(name string) (interface{}, error) {
	service, ok := sl.services[name]

	if !ok {
		return nil, fmt.Errorf("the service(%s) hasn't been registered", name)
	}

	if service.HasScope(Singleton) {
		return sl.getSingleton(service)
	}

	return sl.getPerLookup(service)
}

//
// getSingleton
// @desc 获取单例服务
// @receiver sl *ServiceLocator
// @param service *Service
// @return interface{}
// @return error
//
func (sl *ServiceLocator) getSingleton(service *Service) (interface{}, error) {
	if service.instance != nil {
		return service.instance, nil
	}

	// 执行创建
	service.onceLocker.Do(func() {
		service.instance, service.createError = service.creator(sl)
	})

	// 创建失败
	if service.createError != nil {
		return nil, service.createError
	}

	return service.instance, nil
}

//
// getPerLookup
// @desc 获取多实例服务
// @receiver sl *ServiceLocator
// @param service *Service
// @return interface{}
// @return error
//
func (sl *ServiceLocator) getPerLookup(service *Service) (interface{}, error) {
	return service.creator(sl)
}

//
// Register
// @desc 注册服务
// @receiver sl *ServiceLocator
// @param service *Service
// @return error
//
func (sl *ServiceLocator) Register(service *Service) error {
	// 多实例
	if service.HasScope(PerLookup) {
		sl.services[service.name] = service

		return nil
	}

	// 单例
	if service.HasScope(Singleton) {
		// 立即创建
		if service.HasScope(Immediate) {
			if service.instance, service.createError = service.creator(sl); service.createError != nil {
				return service.createError
			}
		}

		sl.services[service.name] = service

		return nil
	}

	return fmt.Errorf("the service(%s)'s scopes(%d) is invalid", service.name, service.scopes)
}

//
// NewServiceLocator
// @desc 创建服务定位器
// @return *ServiceLocator
//
func NewServiceLocator() *ServiceLocator {
	return &ServiceLocator{
		services: make(map[string]*Service),
	}
}
