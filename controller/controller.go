package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/controller"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/brokerapi"
	"github.com/rancher/rds-broker/client"
	"github.com/sirupsen/logrus"
)

type errNoSuchInstance struct {
	instanceID string
}

func (e errNoSuchInstance) Error() string {
	return fmt.Sprintf("no such instance with ID %s", e.instanceID)
}

type rdsServiceInstance struct {
	Name       string
	Credential *brokerapi.Credential
}

type rdsController struct {
	rwMutex     sync.RWMutex
	instanceMap map[string]*rdsServiceInstance
}

// CreateController creates an instance of a User Provided service broker controller.
func CreateController() controller.Controller {
	var instanceMap = make(map[string]*rdsServiceInstance)
	return &rdsController{
		instanceMap: instanceMap,
	}
}

func (c *rdsController) Catalog() (*brokerapi.Catalog, error) {
	logrus.Info("Catalog()")
	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				Name:        "mysql",
				ID:          "3533e2f0-6001-xxxx-9d15-d7c0b90b75b5",
				Description: "MySQL engine",
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          "b9600ecb-c1c1-4621-b450-a0fa1738e632",
						Description: "MySQL database default plan",
						Free:        true,
					},
				},
				Bindable:       true,
				PlanUpdateable: true,
			},
			{
				Name:        "mariadb",
				ID:          "3533e2f0-6002-xxxx-9d15-d7c0b90b75b5",
				Description: "MariaDB engine",
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          "b9600ecb-c1c2-4621-b450-a0fa1738e632",
						Description: "MariaDB database default plan",
						Free:        true,
					},
				},
				Bindable:       true,
				PlanUpdateable: true,
			},
			{
				Name:        "postgresql",
				ID:          "3533e2f0-6335-xxxx-9d15-d7c0b90b75b5",
				Description: "PostgreSQL engine",
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          "b9600ecb-c1c3-4621-b450-a0fa1738e632",
						Description: "PostgreSQL database default Plan",
						Free:        true,
					},
				},
				Bindable:       true,
				PlanUpdateable: true,
			},
		},
	}, nil
}

// func (c *rdsController) CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (*brokerapi.CreateServiceInstanceResponse, error) {
// if err := client.Install(releaseName(id), id); err != nil {
// 	return nil, err
// }
// logrus.Infof("Created RDS Service Instance:\n%v\n", id)
// return &brokerapi.CreateServiceInstanceResponse{}, nil
// }

func (c *rdsController) CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (*brokerapi.CreateServiceInstanceResponse, error) {
	logrus.Info("CreateServiceInstance()")
	credString, ok := req.Parameters["credentials"]

	str, err := json.Marshal(req.Parameters)
	if err != nil {
		logrus.Error("Error encoding JSON")
	}

	logrus.Info(string(str))
	logrus.Info(req.OrgID, req.PlanID, req.ServiceID, req.SpaceID)

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	if ok {
		jsonCred, err := json.Marshal(credString)
		if err != nil {
			logrus.Errorf("Failed to marshal credentials: %v", err)
			return nil, err
		}
		var cred brokerapi.Credential
		err = json.Unmarshal(jsonCred, &cred)
		if err != nil {
			logrus.Errorf("Failed to unmarshal credentials: %v", err)
			return nil, err
		}

		c.instanceMap[id] = &rdsServiceInstance{
			Name:       id,
			Credential: &cred,
		}
	} else {
		c.instanceMap[id] = &rdsServiceInstance{
			Name: id,
			Credential: &brokerapi.Credential{
				"special-key-1": "special-value-1",
				"special-key-2": "special-value-2",
			},
		}
	}

	if err := client.Install(releaseName(id), req.ServiceID, "default", req.Parameters); err != nil {
		return nil, err
	}
	logrus.Infof("Created User Provided Service Instance:\n%v\n", c.instanceMap[id])
	return &brokerapi.CreateServiceInstanceResponse{}, nil
}

func (c *rdsController) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *rdsController) GetServiceInstanceLastOperation(instanceID, serviceID, planID, operation string) (*brokerapi.LastOperationResponse, error) {
	logrus.Info("GetServiceInstanceLastOperation()")
	return nil, errors.New("Unimplemented")
}

func (c *rdsController) RemoveServiceInstance(instanceID, serviceID, planID string, acceptsIncomplete bool) (*brokerapi.DeleteServiceInstanceResponse, error) {
	logrus.Info("RemoveServiceInstance()")
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	_, ok := c.instanceMap[instanceID]
	if ok {
		delete(c.instanceMap, instanceID)
		if err := client.Delete(releaseName(instanceID)); err != nil {
			return nil, err
		}
		return &brokerapi.DeleteServiceInstanceResponse{}, nil
	}
	return &brokerapi.DeleteServiceInstanceResponse{}, nil
}

func (c *rdsController) Bind(instanceID, bindingID string, req *brokerapi.BindingRequest) (*brokerapi.CreateServiceBindingResponse, error) {
	logrus.Info("Bind()")
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	instance, ok := c.instanceMap[instanceID]
	if !ok {
		return nil, errNoSuchInstance{instanceID: instanceID}
	}
	cred := instance.Credential
	return &brokerapi.CreateServiceBindingResponse{Credentials: *cred}, nil
}

func (c *rdsController) UnBind(instanceID, bindingID, serviceID, planID string) error {
	logrus.Info("UnBind()")
	// Since we don't persist the binding, there's nothing to do here.
	return nil
}

func releaseName(id string) string {
	return "i-" + id
}
