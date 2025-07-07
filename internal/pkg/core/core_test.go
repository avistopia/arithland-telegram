package core_test

import (
	"github.com/avistopia/arithland-telegram/internal/models"
	"github.com/avistopia/arithland-telegram/internal/pkg/core"
	"github.com/avistopia/arithland-telegram/internal/pkg/test"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestFlow() {
	db, err := test.DB()
	s.Require().NoError(err)

	userRepo, err := models.NewUserRepo(db)
	s.Require().NoError(err)

	// TODO mock and pass the tg bot object

	_, err = core.NewService(nil, userRepo).Flow()
	s.Require().NoError(err)
}
