package mockgcp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/api/cloudresourcemanager/v3"
	googleapi "google.golang.org/api/googleapi"
	option "google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

type SetPolicyCallItf interface {
	Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error)
}

type GCPClient struct {
	service *MockService
}

func NewClient() *GCPClient {
	service, _ := NewService(context.TODO())
	return &GCPClient{service: service}
}

func (client *GCPClient) ProjectSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) SetPolicyCallItf {
	return client.service.Projects.SetIamPolicy(resource, setiampolicyrequest)
}

func (client *GCPClient) FolderSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) SetPolicyCallItf {
	return client.service.Folders.SetIamPolicy(resource, setiampolicyrequest)
}

type MockService struct {
	Projects *ProjectsService
	Folders  *FoldersService
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
	//s.Organizations = NewOrganizationsService(s)
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

type Project struct {
	projectID string
	policy    *cloudresourcemanager.Policy
}

type Folder struct {
	folderID string
	policy   *cloudresourcemanager.Policy
}

type ProjectsService struct {
	s *MockService
	p []*Project
}

func NewProjectsService(s *MockService) *ProjectsService {
	rs := &ProjectsService{s: s}
	return rs
}

func (r *ProjectsService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *ProjectsGetIamPolicyCall {
	c := &ProjectsGetIamPolicyCall{s: r.s}
	c.resource = resource
	c.getiampolicyrequest = getiampolicyrequest
	return c
}

func (r *ProjectsService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *ProjectsSetIamPolicyCall {
	c := &ProjectsSetIamPolicyCall{s: r.s}
	c.resource = resource
	c.setiampolicyrequest = setiampolicyrequest
	return c
}

type ProjectsGetIamPolicyCall struct {
	s                   *MockService
	resource            string
	getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *ProjectsGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var policy *cloudresourcemanager.Policy
	for _, project := range c.s.Projects.p {
		if project.projectID == c.resource {
			policy = project.policy
		}
	}
	if policy == nil {
        return nil, fmt.Errorf("resource %v does not exist", c.resource)
	}
	return policy, nil
}

type ProjectsSetIamPolicyCall struct {
	s                   *MockService
	resource            string
	setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *ProjectsSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var found bool

	for _, project := range c.s.Projects.p {
		if project.projectID == c.resource {
			found = true
			project.policy = c.setiampolicyrequest.Policy
		}
	}

	if !found {
        return nil, fmt.Errorf("resource %v does not exist", c.resource)
	}

	return c.setiampolicyrequest.Policy, nil
}

type FoldersService struct {
	s *MockService
	f []*Folder
}

func NewFoldersService(s *MockService) *FoldersService {
	rs := &FoldersService{s: s}
	return rs
}

func (r *FoldersService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *FoldersGetIamPolicyCall {
	c := &FoldersGetIamPolicyCall{s: r.s}
	c.resource = resource
	c.getiampolicyrequest = getiampolicyrequest
	return c
}

func (r *FoldersService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *FoldersSetIamPolicyCall {
	c := &FoldersSetIamPolicyCall{s: r.s}
	c.resource = resource
	c.setiampolicyrequest = setiampolicyrequest
	return c
}

type FoldersGetIamPolicyCall struct {
	s                   *MockService
	resource            string
	getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *FoldersGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var policy *cloudresourcemanager.Policy
	for _, folder := range c.s.Folders.f {
		if folder.folderID == c.resource {
			policy = folder.policy
		}
	}

	if policy == nil {
        return nil, fmt.Errorf("resource %v does not exist", c.resource)
	}

	return policy, nil
}

type FoldersSetIamPolicyCall struct {
	s                   *MockService
	resource            string
	setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *FoldersSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	var found bool

	for _, folder := range c.s.Folders.f {
		if folder.folderID == c.resource {
			found = true
			folder.policy = c.setiampolicyrequest.Policy
		}
	}

	if !found {
		return nil, fmt.Errorf("resource does not exist")
	}

	return c.setiampolicyrequest.Policy, nil
}

// These functions below don't emulate anything in the GCP API, they're just for making test data easily

func StringGenerator(seed int) string {
	rand.Seed(time.Now().UnixNano() + int64(seed))
	return fmt.Sprintf("randomString-%d%d", rand.Intn(99999), rand.Intn(99999))
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
		projectID: projectID,
		policy:    policy,
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
		folderID: folderID,
		policy:   policy,
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
