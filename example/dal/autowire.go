package dal

import (
	"github.com/non1996/autowire/anno"
)

var _ = anno.Annotations{
	anno.Component[ADaoImpl]{
		anno.Implement[ADao]{},
		anno.Primary{},
		anno.Autowired{Field: "db"},
	},
	anno.Component[ADaoMock]{
		anno.Implement[ADao]{},
		anno.ConditionalOnConfig{Key: "Mock.Dao", Value: "1"},
	},
	anno.Component[BDaoImpl]{
		anno.Implement[BDao]{},
		anno.Autowired{Field: "w"},
		anno.Autowired{Field: "r"},
	},
}
