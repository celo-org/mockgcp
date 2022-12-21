package mockgcp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
    "regexp"

	"google.golang.org/api/cloudresourcemanager/v3"
	googleapi "google.golang.org/api/googleapi"
	option "google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

const (
    ResourceNotFoundError = "resource not found"
)

type SetPolicyCallItf interface {
	Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error)
}

type GCPClient struct {
	Service *MockService
}

func NewClient() *GCPClient {
	service, _ := NewService(context.TODO())
	return &GCPClient{Service: service}
}

func (client *GCPClient) ProjectSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) SetPolicyCallItf {
	return client.Service.Projects.SetIamPolicy(resource, setiampolicyrequest)
}

func (client *GCPClient) FolderSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) SetPolicyCallItf {
	return client.Service.Folders.SetIamPolicy(resource, setiampolicyrequest)
}

func (client *GCPClient) OrganizationSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) SetPolicyCallItf {
	return client.Service.Organizations.SetIamPolicy(resource, setiampolicyrequest)
}

type MockService struct {
	Projects *ProjectsService
	Folders  *FoldersService
    Organizations *OrganizationsService
}

func NewService(ctx context.Context, opts ...option.ClientOption) (*MockService, error) {
	client, _, err := htransport.NewClient(ctx, opts...)
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
	s.Folders = NewFoldersService(s)
	s.Organizations = NewOrganizationsService(s)
	/*
		s.Liens = NewLiensService(s)
		s.Operations = NewOperationsService(s)
	*/
	s.Projects = NewProjectsService(s)
	/*
		s.TagBindings = NewTagBindingsService(s)
		s.TagKeys = NewTagKeysService(s)
		s.TagValues = NewTagValuesService(s)
	*/
	return s, nil
}

type Organization struct {
	OrganizationID string
	Policy    *cloudresourcemanager.Policy
}

type Project struct {
	ProjectID string
	Policy    *cloudresourcemanager.Policy
}

type Folder struct {
	FolderID string
	Policy   *cloudresourcemanager.Policy
}

type OrganizationsService struct {
	Service *MockService
	OrganizationList []*Organization
}

func NewOrganizationsService(s *MockService) *OrganizationsService {
	rs := &OrganizationsService{Service: s}
	return rs
}

func (r *OrganizationsService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *OrganizationsGetIamPolicyCall {
	c := &OrganizationsGetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Getiampolicyrequest = getiampolicyrequest
	return c
}

func (r *OrganizationsService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *OrganizationsSetIamPolicyCall {
	c := &OrganizationsSetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Setiampolicyrequest = setiampolicyrequest
	return c
}



type OrganizationsGetIamPolicyCall struct {
	Service                   *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *OrganizationsGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var policy *cloudresourcemanager.Policy
	for _, organization := range c.Service.Organizations.OrganizationList {
		if organization.OrganizationID == c.Resource {
			policy = organization.Policy
		}
	}
	if policy == nil {
        return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
	}
	return policy, nil
}

type OrganizationsSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *OrganizationsSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var found bool

    match, _ := regexp.MatchString("organizations/.*", c.Resource) 
    if !match {
        return nil, fmt.Errorf("resource format invalid")
    }

	for _, organization := range c.Service.Organizations.OrganizationList {
		if organization.OrganizationID == c.Resource {
			found = true
			organization.Policy = c.Setiampolicyrequest.Policy
		}
	}

	if !found {
        return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
	}

	return c.Setiampolicyrequest.Policy, nil
}

type ProjectsService struct {
	Service *MockService
	ProjectList []*Project
}

func NewProjectsService(s *MockService) *ProjectsService {
	rs := &ProjectsService{Service: s}
	return rs
}

func (r *ProjectsService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *ProjectsGetIamPolicyCall {
	c := &ProjectsGetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Getiampolicyrequest = getiampolicyrequest
	return c
}

func (r *ProjectsService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *ProjectsSetIamPolicyCall {
	c := &ProjectsSetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Setiampolicyrequest = setiampolicyrequest
	return c
}

type ProjectsGetIamPolicyCall struct {
	Service                   *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *ProjectsGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var policy *cloudresourcemanager.Policy
	for _, project := range c.Service.Projects.ProjectList {
		if project.ProjectID == c.Resource {
			policy = project.Policy
		}
	}
	if policy == nil {
        return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
	}
	return policy, nil
}

type ProjectsSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *ProjectsSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var found bool

    match, _ := regexp.MatchString("projects/.*", c.Resource) 
    if !match {
        return nil, fmt.Errorf("resource format invalid")
    }

	for _, project := range c.Service.Projects.ProjectList {
		if project.ProjectID == c.Resource {
			found = true
			project.Policy = c.Setiampolicyrequest.Policy
		}
	}

	if !found {
        return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
	}

	return c.Setiampolicyrequest.Policy, nil
}

type FoldersService struct {
	Service *MockService
	FolderList []*Folder
}

func NewFoldersService(s *MockService) *FoldersService {
	rs := &FoldersService{Service: s}
	return rs
}

func (r *FoldersService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *FoldersGetIamPolicyCall {
	c := &FoldersGetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Getiampolicyrequest = getiampolicyrequest
	return c
}

func (r *FoldersService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *FoldersSetIamPolicyCall {
	c := &FoldersSetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Setiampolicyrequest = setiampolicyrequest
	return c
}

type FoldersGetIamPolicyCall struct {
	Service                   *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *FoldersGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var policy *cloudresourcemanager.Policy
	for _, folder := range c.Service.Folders.FolderList {
		if folder.FolderID == c.Resource {
			policy = folder.Policy
		}
	}

	if policy == nil {
        return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)

	}

	return policy, nil
}

type FoldersSetIamPolicyCall struct {
	Service                   *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *FoldersSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var found bool

    match, _ := regexp.MatchString("folders/.*", c.Resource)
    if !match {
        return nil, fmt.Errorf("resource format invalid")
    }
	for _, folder := range c.Service.Folders.FolderList {
		if folder.FolderID == c.Resource {
			found = true
			folder.Policy = c.Setiampolicyrequest.Policy
		}
	}

	if !found {
        return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
	}

	return c.Setiampolicyrequest.Policy, nil
}

// These functions below don't emulate anything in the GCP API, they're just for making test data easily

func StringGenerator(seed int) string {
	rand.Seed(time.Now().UnixNano() + int64(seed))
	return fmt.Sprintf("randomString-%d%d", rand.Intn(99999), rand.Intn(99999))
}


func GenerateOrganizations(count int) (organizations []*Organization) {
	for i := 0; i < count; i++ {
		organizations = append(organizations, GenerateOrganization())
	}
	return organizations
}

func GenerateOrganization() *Organization {
	rand.Seed(time.Now().UnixNano())

	return NewOrganization(StringGenerator(0), GeneratePolicy(nil))
}

func NewOrganization(organizationID string, policy *cloudresourcemanager.Policy) *Organization {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	return &Organization{
		OrganizationID: organizationID,
		Policy:   policy,
	}
}





func GenerateProjects(count int) (projects []*Project) {
	for i := 0; i < count; i++ {
		projects = append(projects, GenerateProject())
	}
	return projects
}

func GenerateProject() *Project {
	rand.Seed(time.Now().UnixNano())

	return NewProject(StringGenerator(0), GeneratePolicy(nil))
}

func NewProject(projectID string, policy *cloudresourcemanager.Policy) *Project {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	return &Project{
		ProjectID: projectID,
		Policy:    policy,
	}
}

func GenerateFolders(count int) (folders []*Folder) {
	for i := 0; i < count; i++ {
		folders = append(folders, GenerateFolder())
	}
	return folders
}

func GenerateFolder() *Folder {
	rand.Seed(time.Now().UnixNano())

	return NewFolder(StringGenerator(0), GeneratePolicy(nil))
}

func NewFolder(folderID string, policy *cloudresourcemanager.Policy) *Folder {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	return &Folder{
		FolderID: folderID,
		Policy:   policy,
	}
}

func GenerateBindings(number int) (bindings []*cloudresourcemanager.Binding) {
	for i := 0; i < number; i++ {
		bindings = append(bindings, GenerateBinding(""))
	}
	return bindings
}

func GenerateBinding(role string, members ...string) *cloudresourcemanager.Binding {
	rand.Seed(time.Now().UnixNano())

	if role == "" {
		role = StringGenerator(0)
	}

	if members == nil {
		for i := 0; i < rand.Intn(10); i++ {
			members = append(members, StringGenerator(len(members)+1))
		}
	}

	return &cloudresourcemanager.Binding{
		Role:    role,
		Members: members,
	}
}

func GeneratePolicy(bindings []*cloudresourcemanager.Binding) *cloudresourcemanager.Policy {
	rand.Seed(time.Now().UnixNano())
	if bindings == nil {
		for i := 0; i < rand.Intn(10); i++ {
			bindings = append(bindings, GenerateBinding("", ""))
		}
	}

	return &cloudresourcemanager.Policy{
		Bindings: bindings,
	}
}

func NewBinding(role string, members ...string) *cloudresourcemanager.Binding {
	return &cloudresourcemanager.Binding{
		Role:    role,
		Members: members,
	}
}

func NewPolicy(bindings []*cloudresourcemanager.Binding) *cloudresourcemanager.Policy {
	return &cloudresourcemanager.Policy{
		Bindings: bindings,
	}
}
