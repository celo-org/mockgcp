package mockgcp

import (
	"google.golang.org/api/cloudresourcemanager/v3"
	googleapi "google.golang.org/api/googleapi"
)

type MockService struct {
	Projects *ProjectsService
}

func NewService(ctx context.Context, opts ...option.ClientOption) (*Service, error) {
	client, endpoint, err := htransport.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	s, err := New(client)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func New(client *http.Client) (*MockService, error) {
	s := &MockService{}
	/*
		s.Folders = NewFoldersService(s)
		s.Liens = NewLiensService(s)
		s.Operations = NewOperationsService(s)
		s.Organizations = NewOrganizationsService(s)
	*/
	s.Projects = NewProjectsService(s)
	/*
		s.TagBindings = NewTagBindingsService(s)
		s.TagKeys = NewTagKeysService(s)
		s.TagValues = NewTagValuesService(s)
	*/
	return s, nil
}

type Projects struct {
	projectID string
	policy    *cloudresourcemanager.Policy
}

type ProjectsService struct {
	s *MockService
	p *[]Projects
}

func NewProjectsService(s *MockService) *ProjectsService {
	rs := &ProjectsService{s: s}
	return rs
}

func (r *ProjectsService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *ProjectsSetIamPolicyCall {
	c := &cloudresourcemanager.ProjectsSetIamPolicyCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.resource = resource
	c.setiampolicyrequest = setiampolicyrequest
	return c
}

type ProjectsSetIamPolicyCall struct {
	s                   *MockService
	resource            string
	setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *ProjectsSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*Policy, error) {
	project := &Project{projectID: resource, policy: c.setiampolicyrequest.Policy}
	s.projects.p = append(s.projects.p, project)
    return project.Policy, nil
}




