package team

import (
	"fmt"
	"strings"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	repo        Repository
	teamService interfaces.TeamService
}

// NewService returns an object that implements the Service interface
func NewService(repo Repository, r interfaces.Registry) interfaces.TeamService {
	s := &service{repo: repo}

	r.RegisterService(s, "team")
	return s
}

// GetRouteNameFromName returns a route name from a name
func (s *service) GetRouteNameFromName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "", -1))
}

// FindFiltered returns a list of teams filtered by limit offset and search params
func (s *service) FindFiltered(limit, offset int, search string) ([]*models.Team, error) {
	return s.repo.FindFiltered(limit, offset, search)
}

// CreateTeam creates a team :)
func (s *service) CreateTeam(creatorID bson.ObjectId, name, image, description string) (*models.Team, error) {
	team := &models.Team{Name: name, Image: image, Description: description}
	team.RouteName = s.GetRouteNameFromName(team.Name)
	team.Members = []bson.ObjectId{creatorID}
	team.Admins = []bson.ObjectId{creatorID}

	if _, err := s.repo.Store(team); err != nil {
		return nil, err
	}

	return team, nil
}

// FindByID returns a team by its ID or error if not found
func (s *service) FindByID(id bson.ObjectId) (*models.Team, error) {
	return s.repo.FindByID(id)
}

// GetTeamMembers returns the list of members that belong to a team or an error if the user is not in the team
func (s *service) GetTeamMembers(userID bson.ObjectId, team *models.Team) ([]*models.User, error) {
	if team.IsUserMember(userID) {
		return team.Members.([]*models.User), nil
	}
	return nil, fmt.Errorf("you must be in the team to see members")
}

// GetTeamMembersCount returns number members that belong to a team or an error if the user is not in the team
func (s *service) GetTeamMembersCount(team *models.Team) (int, error) {
	return len(team.Members.([]*models.User)), nil
}

// GetTeamAdmins returns the list of admins that belong to a team or an error if the user is not in the team
func (s *service) GetTeamAdmins(userID bson.ObjectId, team *models.Team) ([]*models.User, error) {
	if team.IsUserMember(userID) {
		return team.Members.([]*models.User), nil
	}
	return nil, fmt.Errorf("you must be in the team to see admins")
}

// OnAllServicesRegistered is the method called when all services are registered, used to get dependencies in execution time
func (s *service) OnAllServicesRegistered(r interfaces.Registry) {
	s.teamService = r.GetService("team").(interfaces.TeamService)
}
