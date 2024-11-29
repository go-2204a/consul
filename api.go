package consul

import (
	"github.com/hashicorp/consul/api"
	"log"
	"time"
)

// InitConsulService 初始化Consul服务注册
// address: Consul服务的地址
// serviceID: 服务的唯一ID
// serviceName: 服务的名称
// servicePort: 服务监听的端口
// ttl: 服务健康检查的TTL（生存时间）
// checkID: 健康检查ID
func InitConsulService(address string, serviceID string, serviceName string, servicePort int, ttl string, checkID string) func() error {
	// 创建并初始化 Consul 客户端
	config := api.DefaultConfig()
	config.Address = address //地址
	consul, err := api.NewClient(config)
	if err != nil {
		log.Printf("初始化 Consul 客户端失败: %v", err)
		return nil
	}
	// 创建健康检查配置
	check := &api.AgentServiceCheck{
		TTL:                            ttl,
		CheckID:                        checkID,
		DeregisterCriticalServiceAfter: ttl,
	}
	// 创建服务注册信息
	a := new(api.AgentServiceRegistration)
	a.ID = serviceID
	a.Name = serviceName
	a.Tags = []string{}
	a.Port = servicePort
	a.Address = "127.0.0.1"
	a.Check = check

	// 服务注册函数
	return func() error {
		// 注册服务到 Consul
		err := consul.Agent().ServiceRegister(a)
		if err != nil {
			return err
		}

		// 定期更新服务的健康状态
		go func() {
			for {
				err := consul.Agent().UpdateTTL(a.ID, "服务正常", api.HealthPassing)
				if err != nil {
					log.Printf("更新健康检查失败: %v", err)
				}
				time.Sleep(10 * time.Second)
			}
		}()
		log.Println("服务注册成功并开始健康检查定期更新")
		return nil
	}
}
