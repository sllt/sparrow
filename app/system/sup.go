package system

import (
	"github.com/sllt/sparrow/actor"
	"github.com/sllt/sparrow/app/system/inspect"
	"github.com/sllt/sparrow/gen"
)

func factory_sup() gen.ProcessBehavior {
	return &sup{}
}

type sup struct {
	actor.Supervisor
}

func (s *sup) Init(args ...any) (actor.SupervisorSpec, error) {

	spec := actor.SupervisorSpec{
		Type: actor.SupervisorTypeOneForOne,
		Children: []actor.SupervisorChildSpec{
			{
				Factory: factory_metrics,
				Name:    "system_metrics",
			},
			{
				Factory: inspect.Factory,
				Name:    inspect.Name,
			},
		},
	}
	spec.Restart.Strategy = actor.SupervisorStrategyPermanent
	return spec, nil
}
