package mockgcp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"google.golang.org/api/cloudresourcemanager/v3"
	googleapi "google.golang.org/api/googleapi"
	option "google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

const (
	ResourceNotFoundError = "resource not found"
)

type PolicyCallItf interface {
	Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error)
}

type GCPClient struct {
	Service *MockService
}

func NewClient() *GCPClient {
	service, _ := NewService(context.TODO())
	return &GCPClient{Service: service}
}

// Wrapper methods for Google Clouds API
func (client *GCPClient) ProjectSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) PolicyCallItf {
	return client.Service.Projects.SetIamPolicy(resource, setiampolicyrequest)
}
func (client *GCPClient) ProjectGetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) PolicyCallItf {
	return client.Service.Projects.GetIamPolicy(resource, getiampolicyrequest)
}

func (client *GCPClient) FolderSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) PolicyCallItf {
	return client.Service.Folders.SetIamPolicy(resource, setiampolicyrequest)
}

func (client *GCPClient) FolderGetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) PolicyCallItf {
	return client.Service.Folders.GetIamPolicy(resource, getiampolicyrequest)
}

func (client *GCPClient) OrganizationSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) PolicyCallItf {
	return client.Service.Organizations.SetIamPolicy(resource, setiampolicyrequest)
}

func (client *GCPClient) OrganizationGetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) PolicyCallItf {
	return client.Service.Organizations.GetIamPolicy(resource, getiampolicyrequest)
}

// A mockup of cloud resource manager's service
type MockService struct {
	Projects      *ProjectsService
	Folders       *FoldersService
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

// The client which NewService will call to create a new service.
func New(client *http.Client) (*MockService, error) {
	s := &MockService{}
	s.Folders = NewFoldersService(s)
	s.Organizations = NewOrganizationsService(s)
	s.Projects = NewProjectsService(s)
	return s, nil
}

type Organization struct {
	OrganizationID string
	Policy         *cloudresourcemanager.Policy
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
	Service          *MockService
	OrganizationList []*Organization
}

func NewOrganizationsService(s *MockService) *OrganizationsService {
	rs := &OrganizationsService{Service: s}
	return rs
}

// NewOrganization creates a new organization with the specified ID and policy on the Organizations Service
// and returns a pointer to the created organization.  If policy isn't specified it will generate a blank one
func (r *OrganizationsService) NewOrganization(orgID string, policy *cloudresourcemanager.Policy) *Organization {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	organization := &Organization{
		OrganizationID: orgID,
		Policy:         policy,
	}

	r.OrganizationList = append(r.OrganizationList, organization)

	return organization
}

// GenerateOrganizations takes a count of Organizations to create, and a basename, and will generate random
// data for the Organizations and add them to the Organizations Service
func (r *OrganizationsService) GenerateOrganizations(count int, baseName string) (organizations []*Organization) {
	rand.Seed(time.Now().UnixNano())
	startNumber := rand.Intn(9999)

	for i := 0; i < count; i++ {
		orgID := fmt.Sprintf("%v%v-%d", "organizations/", baseName, startNumber+i)
		policy := GeneratePolicy(nil)
		organizations = append(organizations, r.NewOrganization(orgID, policy))
	}
	return organizations
}

// FindPolicy will search the organizations service for a matching policy, and return
// the organization that contains it
func (r *OrganizationsService) FindPolicy(policy *cloudresourcemanager.Policy) *Organization {
	for _, organization := range r.OrganizationList {
		if reflect.DeepEqual(policy, organization.Policy) {
			return organization
		}
	}
	return nil
}

// GetIamPolicy will take a resource name (organization ID), and a getiampolicyrequest
// and returns a GetIamPolicy Call, so we can run a Do() method
func (r *OrganizationsService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *OrganizationsGetIamPolicyCall {
	c := &OrganizationsGetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Getiampolicyrequest = getiampolicyrequest
	return c
}

// SetIamPolicy will take a resource name (organization ID), and a setiampolicyrequest
// and returns a SetIamPolicy Call, so we can run a Do() method
func (r *OrganizationsService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *OrganizationsSetIamPolicyCall {
	c := &OrganizationsSetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Setiampolicyrequest = setiampolicyrequest
	return c
}

type OrganizationsGetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *OrganizationsGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	for _, organization := range c.Service.Organizations.OrganizationList {
		if organization.OrganizationID == c.Resource {
			policy := *organization.Policy
            bindings := make([]*cloudresourcemanager.Binding, 0, len(organization.Policy.Bindings))
            for _, b := range organization.Policy.Bindings {
                binding := *b
                bindings = append(bindings, &binding)
            }
            policy.Bindings = bindings
            return &policy, nil
		}
	}
	return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
}

type OrganizationsSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *OrganizationsSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	match, _ := regexp.MatchString("organizations/.*", c.Resource)
	if !match {
		return nil, fmt.Errorf("resource format invalid")
	}
	for _, organization := range c.Service.Organizations.OrganizationList {
		if organization.OrganizationID == c.Resource {
			organization.Policy = c.Setiampolicyrequest.Policy
			return organization.Policy, nil
		}
	}
	return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
}

type ProjectsService struct {
	Service     *MockService
	ProjectList []*Project
}

func NewProjectsService(s *MockService) *ProjectsService {
	rs := &ProjectsService{Service: s}
	return rs
}

func (r *ProjectsService) NewProject(projectID string, policy *cloudresourcemanager.Policy) *Project {
/*
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
    */
    /*
	project := &Project{
		ProjectID: projectID,
		Policy:    policy,
	}


*/
    if policy == nil {
    	policy = &cloudresourcemanager.Policy{}
    }

    project := &Project{
            ProjectID: "test",
            Policy: policy,
          }

//	r.ProjectList = append(r.ProjectList, project)

    return nil
    return project

}

func (r *ProjectsService) GenerateProjects(count int, baseName string) (projects []*Project) {
	rand.Seed(time.Now().UnixNano())
	startNumber := rand.Intn(9999)

	for i := 0; i < count; i++ {
		projectID := fmt.Sprintf("%v%v-%d", "projects/", baseName, startNumber+i)
		policy := GeneratePolicy(nil)
		projects = append(projects, r.NewProject(projectID, policy))
	}
	return projects
}

func (r *ProjectsService) FindPolicy(policy *cloudresourcemanager.Policy) *Project {
	for _, project := range r.ProjectList {
		if reflect.DeepEqual(policy, project.Policy) {
			return project
		}
	}
	return nil
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
	Service             *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *ProjectsGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	for _, project := range c.Service.Projects.ProjectList {
		if project.ProjectID == c.Resource {
			policy := *project.Policy
            bindings := make([]*cloudresourcemanager.Binding, 0, len(project.Policy.Bindings))
            for _, b := range project.Policy.Bindings {
                binding := *b
                bindings = append(bindings, &binding)
            }
            policy.Bindings = bindings
            return &policy, nil
		}
	}
	return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
}

type ProjectsSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}


func (c *ProjectsSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	match, _ := regexp.MatchString("projects/.*", c.Resource)
	if !match {
		return nil, fmt.Errorf("resource format invalid")
	}
	for _, project := range c.Service.Projects.ProjectList {
		if project.ProjectID == c.Resource {
			project.Policy = c.Setiampolicyrequest.Policy
            return project.Policy, nil
		}
	}
	return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
}

type FoldersService struct {
	Service    *MockService
	FolderList []*Folder
}

func NewFoldersService(s *MockService) *FoldersService {
	rs := &FoldersService{Service: s}
	return rs
}

func (r *FoldersService) NewFolder(folderID string, policy *cloudresourcemanager.Policy) *Folder {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	folder := &Folder{
		FolderID: folderID,
		Policy:   policy,
	}
	r.FolderList = append(r.FolderList, folder)

	return folder
}

func (r *FoldersService) GenerateFolders(count int, baseName string) (folders []*Folder) {
	rand.Seed(time.Now().UnixNano())
	startNumber := rand.Intn(9999)

	for i := 0; i < count; i++ {
		folderID := fmt.Sprintf("%v%v-%d", "folders/", baseName, startNumber+i)
		policy := GeneratePolicy(nil)
		folders = append(folders, r.NewFolder(folderID, policy))
	}
	return folders
}

func (r *FoldersService) FindPolicy(policy *cloudresourcemanager.Policy) *Folder {
	for _, folder := range r.FolderList {
		if reflect.DeepEqual(policy, folder.Policy) {
			return folder
		}
	}
	return nil
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
	Service             *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

func (c *FoldersGetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {

	for _, folder := range c.Service.Folders.FolderList {
		if folder.FolderID == c.Resource {
			policy := *folder.Policy
            bindings := make([]*cloudresourcemanager.Binding, 0, len(folder.Policy.Bindings))
            for _, b := range folder.Policy.Bindings {
                binding := *b
                bindings = append(bindings, &binding)
            }
            policy.Bindings = bindings
            return &policy, nil
		}
	}
	return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
}

type FoldersSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

func (c *FoldersSetIamPolicyCall) Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error) {
	match, _ := regexp.MatchString("folders/.*", c.Resource)
	if !match {
		return nil, fmt.Errorf("resource format invalid")
	}
	for _, folder := range c.Service.Folders.FolderList {
		if folder.FolderID == c.Resource {
			folder.Policy = c.Setiampolicyrequest.Policy
            return folder.Policy, nil
		}
	}

		return nil, fmt.Errorf("%v: %v", ResourceNotFoundError, c.Resource)
}

func NewPolicy(bindings []*cloudresourcemanager.Binding) *cloudresourcemanager.Policy {
	return &cloudresourcemanager.Policy{
		Bindings: bindings,
	}
}

func GeneratePolicy(bindings ...*cloudresourcemanager.Binding) *cloudresourcemanager.Policy {
	rand.Seed(time.Now().UnixNano())
	if bindings == nil {
		for i := 0; i < rand.Intn(10) + 10; i++ {
			bindings = append(bindings, GenerateBinding())
		}
	}
	return NewPolicy(bindings)
}

func AddBindingsToPolicy(policy *cloudresourcemanager.Policy, bindings ...*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
    policy.Bindings = append(policy.Bindings, bindings...)
	return policy.Bindings
}

func NewBinding(role string, members ...string) *cloudresourcemanager.Binding {
	m := make([]string, len(members))
	copy(m, members)
	return &cloudresourcemanager.Binding{
		Role:    role,
		Members: m,
	}
}

func GenerateBinding() *cloudresourcemanager.Binding {
	rand.Seed(time.Now().UnixNano())

	role := GenerateRole(StringGenerator())
	var members []string

	for i := 0; i < rand.Intn(10) + 1; i++ {
		members = append(members, GenerateMember(StringGenerator()))
	}

	return NewBinding(role, members...)
}

func GenerateBindings(number int) (bindings []*cloudresourcemanager.Binding) {
	for i := 0; i < number; i++ {
		bindings = append(bindings, GenerateBinding())
	}
	return bindings
}

func GenerateMember(principal string) string {
	rand.Seed(time.Now().UnixNano() + int64(len(principal)))
	return fmt.Sprintf("%v-%d-%d-%v", principal, rand.Intn(99999), rand.Intn(99999), "@testdomain.co")
}

func GenerateRole(role string) string {
	rand.Seed(time.Now().UnixNano() + int64(len(role)))
	return fmt.Sprintf("%v-%d-%d", role, rand.Intn(99999), rand.Intn(99999))
}

func StringGenerator() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("randomString-%d%d", rand.Intn(99999), rand.Intn(99999))
}

func PolicyContains(policy *cloudresourcemanager.Policy, role string) *cloudresourcemanager.Binding {
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			return binding
		}
	}
	return nil
}

func BindingContains(binding *cloudresourcemanager.Binding, member string) bool {
	for _, m := range binding.Members {
		if m == member {
			return true
		}
	}
	return false
}
