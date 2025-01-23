package rbac

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist/file-adapter"
)

var enforcer *casbin.Enforcer

// InitializeRBAC initializes the Casbin enforcer
func InitializeRBAC() {
	m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
	`)
	if err != nil {
		log.Fatalf("Failed to load RBAC model: %v", err)
	}

	adapter := fileadapter.NewAdapter("configs/rbac_policy.csv")
	enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin enforcer: %v", err)
	}
}

// Enforce checks if a user has permission to perform an action
func Enforce(subject, object, action string) bool {
	ok, err := enforcer.Enforce(subject, object, action)
	if err != nil {
		log.Printf("Error enforcing policy: %v", err)
		return false
	}
	return ok
}
