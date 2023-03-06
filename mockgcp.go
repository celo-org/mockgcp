package mockgcp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v3"
	googleapi "google.golang.org/api/googleapi"
	option "google.golang.org/api/option"
)

const (
	resourceNotFoundError = "resource not found"
)

// GCPClient is the final wrapper we make with locally assigned methods, so we can match it on
// an interface easier.  Since the google cloud methods aren't directly on the client, but on
// the services (such as client.Projects.ProjectSetIamPolicy), and I need to create wrappers for
// them, this lets me call client.ProjectSetIamPolicy instead, and make clients with those to
// match the interface.  This is the wrappr you should probably be using if you want to use
// the library.
type GCPClient struct {
	Service *MockService
}

// NewClient returns the mock GCPClient client above
func NewClient() *GCPClient {
	service, _ := NewService(context.TODO())
	return &GCPClient{Service: service}
}

// PolicyCallItf interface will match for Do() so we can call it on SetPolicyCall and GetPolicyCall
type PolicyCallItf interface {
	Do(opts ...googleapi.CallOption) (*cloudresourcemanager.Policy, error)
}

// SearchCallItf interface will match for Do() so we can match
//type SearchCallItf interface {
//	Do(opts ...googleapi.CallOption) (ResponseItf, error)
//}

//type ResponseItf interface {
//   
//}

// Wrapper methods for Google Clouds API

// ProjectSetIamPolicy is a wrapper for the Projects.SetIamPolicy method so we can create and interface to match
// our mock client to the GCP client
func (client *GCPClient) ProjectSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) PolicyCallItf {
	return client.Service.Projects.SetIamPolicy(resource, setiampolicyrequest)
}

// ProjectsSearch Searches for folders by Name to get the ID
func (client *GCPClient) ProjectsSearch() *ProjectsSearchCall {
	return client.Service.Projects.Search()
}

// ProjectGetIamPolicy is a wrapper for the Projects.GetIamPolicy method so we can create and interface to match
// our mock client to the GCP client
func (client *GCPClient) ProjectGetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) PolicyCallItf {
	return client.Service.Projects.GetIamPolicy(resource, getiampolicyrequest)
}

// FoldersSearch Searches for folders by Name to get the ID
/*
func (client *GCPClient) FoldersSearch() *FoldersSearchCall {
	return client.Service.Folders.Search()
}
*/

// FolderSetIamPolicy is a wrapper for the Folders.SetIamPolicy method so we can create and interface to match
// our mock client to the GCP client
func (client *GCPClient) FolderSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) PolicyCallItf {
	return client.Service.Folders.SetIamPolicy(resource, setiampolicyrequest)
}

// FolderGetIamPolicy is a wrapper for the Folders.SetIamPolicy method so we can create and interface to match
// our mock client to the GCP client
func (client *GCPClient) FolderGetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) PolicyCallItf {
	return client.Service.Folders.GetIamPolicy(resource, getiampolicyrequest)
}
/*
// OrganizationsSearch Searches for folders by Name to get the ID
func (client *GCPClient) OrganizationsSearch() *OrganizationsSearchCall {
	return client.Service.Organizations.Search()
}
*/
// OrganizationSetIamPolicy is a wrapper for the Organizations.SetIamPolicy method so we can create and interface to match
// our mock client to the GCP client
func (client *GCPClient) OrganizationSetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) PolicyCallItf {
	return client.Service.Organizations.SetIamPolicy(resource, setiampolicyrequest)
}

// OrganizationGetIamPolicy is a wrapper for the Organizations.GetIamPolicy method so we can create and interface to match
// our mock client to the GCP client
func (client *GCPClient) OrganizationGetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) PolicyCallItf {
	return client.Service.Organizations.GetIamPolicy(resource, getiampolicyrequest)
}

// MockService is a mockup of cloud resource manager's service (wrapper for it's client)
type MockService struct {
	Projects      *ProjectsService
	Folders       *FoldersService
	Organizations *OrganizationsService
}

// NewService creates a MockService and returns it with an http client Wrapper
func NewService(ctx context.Context, opts ...option.ClientOption) (*MockService, error) {
	client := &http.Client{}
	s, err := New(client)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// New is the client which NewService will call to create a new service.
// This wil be wrapped with an http wrapper with NewService
func New(client *http.Client) (*MockService, error) {
	s := &MockService{}
	s.Folders = NewFoldersService(s)
	s.Organizations = NewOrganizationsService(s)
	s.Projects = NewProjectsService(s)
	return s, nil
}

// Organization is a mock of a google cloud Organization
type Organization struct {
	OrganizationID string
	Domain         string
	Policy         *cloudresourcemanager.Policy
}

// Project is a mock of a google cloud Project
type Project struct {
	ProjectID   string
	DisplayName string
	Policy      *cloudresourcemanager.Policy
}

// Folder is a mock of a google cloud Folder
type Folder struct {
	FolderID    string
	DisplayName string
	Policy      *cloudresourcemanager.Policy
}

// OrganizationsService is a mock of google Cloud's Organization Service
type OrganizationsService struct {
	Service          *MockService
	OrganizationList []*Organization
}

// NewOrganizationsService will return a new Organization Service
func NewOrganizationsService(s *MockService) *OrganizationsService {
	rs := &OrganizationsService{Service: s}
	return rs
}

// NewOrganization creates a new organization with the specified ID and policy on the Organizations Service
// and returns a pointer to the created organization.  If policy isn't specified it will generate a blank one
func (r *OrganizationsService) NewOrganization(orgID, domain string, policy *cloudresourcemanager.Policy) *Organization {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	organization := &Organization{
		OrganizationID: orgID,
		Domain:         domain,
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
		organizations = append(organizations, r.NewOrganization(orgID, "", policy))
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
/*
// OrganizationsSearchCall contains the query information for a organization search
type OrganizationsSearchCall struct {
	query   string
	service OrganizationsService
}

// SearchOrganizationsResponse contains the organizations from a search
type SearchOrganizationsResponse struct {
	Organizations []*Organization
}

// Search Creates a Organizations Search Call with the search parameters
func (r *OrganizationsService) Search() *OrganizationsSearchCall {
	c := &OrganizationsSearchCall{service: *r}
	return c
}

// Query Adds the query parameter to the search call
func (call *OrganizationsSearchCall) Query(query string) *OrganizationsSearchCall {
	c := call
	c.query = query
	return c
}

// Do executes the Organizations Search Call and returns the response
func (call *OrganizationsSearchCall) Do(opts ...googleapi.CallOption) (*SearchOrganizationsResponse, error) {
	query := strings.Split(call.query, "=")
	response := &SearchOrganizationsResponse{}
	if len(query) != 2 || query[0] != "domain" {
		return response, fmt.Errorf("invalid organization query")
	}
	for _, org := range call.service.OrganizationList {
		if org.Domain == query[1] {
			response.Organizations = append(response.Organizations, org)
		}
	}
	return response, nil
}
*/

// OrganizationsGetIamPolicyCall is a structure that is returned by Organizations.GetIamPolicy which contains the Request
// to get a policy.  Then we call Do() on it to actually return the policy
type OrganizationsGetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

// Do will be called on OrganizationsGetIamPolicyCall and return the policy found
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
	return nil, fmt.Errorf("%v: %v", resourceNotFoundError, c.Resource)
}

// OrganizationsSetIamPolicyCall is a structure that is returned by Organizations.SetIamPolicy which contains the Request
// to get a policy.  Then we call Do() on in it to Set the Organization Policy
type OrganizationsSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

// Do will be called on OrganizationsGetIamPolicyCall to process the policy change and returns the policy it sets
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
	return nil, fmt.Errorf("%v: %v", resourceNotFoundError, c.Resource)
}

// ProjectsService is a mock of google Cloud's Project Service
type ProjectsService struct {
	Service     *MockService
	ProjectList []*Project
}

// NewProjectsService will return a new Project Service
func NewProjectsService(s *MockService) *ProjectsService {
	rs := &ProjectsService{Service: s}
	return rs
}

// ProjectsSearchCall contains the query information for a project search
type ProjectsSearchCall struct {
	query   string
	service ProjectsService
}

// SearchProjectsResponse contains the projects from a search
type SearchProjectsResponse struct {
	Projects []*Project
}

// Search Creates a Projects Search Call with the search parameters
func (r *ProjectsService) Search() *ProjectsSearchCall {
	c := &ProjectsSearchCall{service: *r}
	return c
}

// Query Adds the query parameter to the search call
func (call *ProjectsSearchCall) Query(query string) *ProjectsSearchCall {
	c := call
	c.query = query
	return c
}

// Do executes the Projects Search Call and returns the response
func (call *ProjectsSearchCall) Do(opts ...googleapi.CallOption) (*SearchProjectsResponse, error) {
	query := strings.Split(call.query, "=")
	response := &SearchProjectsResponse{}
	if len(query) != 2 || query[0] != "displayName" {
		return response, fmt.Errorf("invalid project query")
	}

	for _, project := range call.service.ProjectList {
		if project.DisplayName == query[1] {
			response.Projects = append(response.Projects, project)
		}
	}
	return response, nil
}

// NewProject creates a new project with the specified ID and policy on the Projects Service
// and returns a pointer to the created project.  If policy isn't specified it will generate a blank one
func (r *ProjectsService) NewProject(projectID, projectName string, policy *cloudresourcemanager.Policy) *Project {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	project := &Project{
		ProjectID:   projectID,
		DisplayName: projectName,
		Policy:      policy,
	}
	r.ProjectList = append(r.ProjectList, project)
	return project
}

// GenerateProjects takes a count of Projects to create, and a basename, and will generate random
// data for the Projects and add them to the Projects Service
func (r *ProjectsService) GenerateProjects(count int, baseName string) (projects []*Project) {
	rand.Seed(time.Now().UnixNano())
	startNumber := rand.Intn(9999)

	for i := 0; i < count; i++ {
		projectID := fmt.Sprintf("%v%v-%d", "projects/", baseName, startNumber+i)
		policy := GeneratePolicy(nil)
		projects = append(projects, r.NewProject(projectID, "", policy))
	}
	return projects
}

// FindPolicy will Search a Project Service and return the project with that policy.
// It will only return the first one found, so this should only be used for testing
// where you need to return the project added, and not a reliable way of determining
// which projects have a policy
func (r *ProjectsService) FindPolicy(policy *cloudresourcemanager.Policy) *Project {
	for _, project := range r.ProjectList {
		if reflect.DeepEqual(policy, project.Policy) {
			return project
		}
	}
	return nil
}

// GetIamPolicy will take a resource name (project ID), and a getiampolicyrequest
// and returns a GetIamPolicy Call, so we can run a Do() method on it.
func (r *ProjectsService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *ProjectsGetIamPolicyCall {
	c := &ProjectsGetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Getiampolicyrequest = getiampolicyrequest
	return c
}

// SetIamPolicy will take a resource name (project ID), and a setiampolicyrequest
// and returns a SetIamPolicy Call, so we can run a Do() method on it.
func (r *ProjectsService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *ProjectsSetIamPolicyCall {
	c := &ProjectsSetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Setiampolicyrequest = setiampolicyrequest
	return c
}

// ProjectsGetIamPolicyCall is a structure that is returned by Projects.GetIamPolicy which contains the Request
// to get a policy.  Then we call Do() on it to actually return the policy
type ProjectsGetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

// Do will be called on OrganizationsGetIamPolicyCall and return the policy found
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
	return nil, fmt.Errorf("%v: %v", resourceNotFoundError, c.Resource)
}

// ProjectsSetIamPolicyCall is a structure that is returned by Projects.SetIamPolicy which contains the Request
// to get a policy.  Then we call Do() on in it to Set the Project Policy
type ProjectsSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

// Do will be called on ProjectsGetIamPolicyCall to process the policy change and returns the policy it sets
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
	return nil, fmt.Errorf("%v: %v", resourceNotFoundError, c.Resource)
}

// FoldersService is a mock of google Cloud's Folder Service
type FoldersService struct {
	Service    *MockService
	FolderList []*Folder
}

// NewFoldersService will return a new Folder Service
func NewFoldersService(s *MockService) *FoldersService {
	rs := &FoldersService{Service: s}
	return rs
}

// NewFolder creates a new folder with the specified ID and policy on the Folders Service
// and returns a pointer to the created folder.  If policy isn't specified it will generate a blank one
func (r *FoldersService) NewFolder(folderID string, folderName string, policy *cloudresourcemanager.Policy) *Folder {
	if policy == nil {
		policy = &cloudresourcemanager.Policy{}
	}
	folder := &Folder{
		FolderID:    folderID,
		DisplayName: folderName,
		Policy:      policy,
	}
	r.FolderList = append(r.FolderList, folder)

	return folder
}

// GenerateFolders takes a count of Folders to create, and a basename, and will generate random
// data for the Folders and add them to the Folders Service
func (r *FoldersService) GenerateFolders(count int, baseName string) (folders []*Folder) {
	rand.Seed(time.Now().UnixNano())
	startNumber := rand.Intn(9999)

	for i := 0; i < count; i++ {
		folderID := fmt.Sprintf("%v%v-%d", "folders/", baseName, startNumber+i)
		policy := GeneratePolicy(nil)
		folders = append(folders, r.NewFolder(folderID, folderID, policy))
	}
	return folders
}

// FindPolicy will Search a Folder Service and return the folder with that policy.
// It will only return the first one found, so this should only be used for testing
// where you need to return the folder added, and not a reliable way of determining
// which folders have a policy
func (r *FoldersService) FindPolicy(policy *cloudresourcemanager.Policy) *Folder {
	for _, folder := range r.FolderList {
		if reflect.DeepEqual(policy, folder.Policy) {
			return folder
		}
	}
	return nil
}

// GetIamPolicy will take a resource name (folder ID), and a getiampolicyrequest
// and returns a GetIamPolicy Call, so we can run a Do() method on it.
func (r *FoldersService) GetIamPolicy(resource string, getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest) *FoldersGetIamPolicyCall {
	c := &FoldersGetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Getiampolicyrequest = getiampolicyrequest
	return c
}

// SetIamPolicy will take a resource name (folder ID), and a setiampolicyrequest
// and returns a SetIamPolicy Call, so we can run a Do() method on it.
func (r *FoldersService) SetIamPolicy(resource string, setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest) *FoldersSetIamPolicyCall {
	c := &FoldersSetIamPolicyCall{Service: r.Service}
	c.Resource = resource
	c.Setiampolicyrequest = setiampolicyrequest
	return c
}
/*
// FoldersSearchCall contains the query information for a folder search
type FoldersSearchCall struct {
	query   string
	service FoldersService
}

// SearchFoldersResponse contains the folders from a search
type SearchFoldersResponse struct {
	Folders []*Folder
}

// Search Creates a Folders Search Call with the search parameters
func (r *FoldersService) Search() *FoldersSearchCall {
	c := &FoldersSearchCall{service: *r}
	return c
}

// Query Adds the query parameter to the search call
func (call *FoldersSearchCall) Query(query string) *FoldersSearchCall {
	c := call
	c.query = query
	return c
}

// Do executes the Folders Search Call and returns the response
func (call *FoldersSearchCall) Do(opts ...googleapi.CallOption) (*SearchFoldersResponse, error) {
	query := strings.Split(call.query, "=")
	response := &SearchFoldersResponse{}
	if len(query) != 2 || query[0] != "displayName" {
		return response, fmt.Errorf("invalid folder query")
	}
	for _, folder := range call.service.FolderList {
		if folder.DisplayName == query[1] {
			response.Folders = append(response.Folders, folder)
		}
	}
	return response, nil
}
*/
// FoldersGetIamPolicyCall is a structure that is returned by Folders.GetIamPolicy which contains the Request
// to get a policy.  Then we call Do() on it to actually return the policy
type FoldersGetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Getiampolicyrequest *cloudresourcemanager.GetIamPolicyRequest
}

// Do will be called on OrganizationsGetIamPolicyCall and return the policy found
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
	return nil, fmt.Errorf("%v: %v", resourceNotFoundError, c.Resource)
}

// FoldersSetIamPolicyCall is a structure that is returned by Folders.SetIamPolicy which contains the Request
// to get a policy.  Then we call Do() on in it to Set the Folder Policy
type FoldersSetIamPolicyCall struct {
	Service             *MockService
	Resource            string
	Setiampolicyrequest *cloudresourcemanager.SetIamPolicyRequest
}

// Do will be called on FoldersGetIamPolicyCall to process the policy change and returns the policy it sets
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
	return nil, fmt.Errorf("%v: %v", resourceNotFoundError, c.Resource)
}

// NewPolicy creates a policy with the specified bindings
func NewPolicy(bindings []*cloudresourcemanager.Binding) *cloudresourcemanager.Policy {
	return &cloudresourcemanager.Policy{
		Bindings: bindings,
	}
}

// GeneratePolicy takes a number of bindings and generates a policy.  If no bindings are
// supplied, 10 of them will be generated.  For use with testing
func GeneratePolicy(bindings ...*cloudresourcemanager.Binding) *cloudresourcemanager.Policy {
	rand.Seed(time.Now().UnixNano())
	if bindings == nil {
		for i := 0; i < rand.Intn(10)+10; i++ {
			bindings = append(bindings, GenerateBinding())
		}
	}
	return NewPolicy(bindings)
}

// AddBindingsToPolicy will add bindings to given policy and return the list of pointers to the bindings
// that were added
func AddBindingsToPolicy(policy *cloudresourcemanager.Policy, bindings ...*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	policy.Bindings = append(policy.Bindings, bindings...)
	return policy.Bindings
}

// NewBinding takes a role and members and returns a binding
func NewBinding(role string, members ...string) *cloudresourcemanager.Binding {
	m := make([]string, len(members))
	copy(m, members)
	return &cloudresourcemanager.Binding{
		Role:    role,
		Members: m,
	}
}

// GenerateBinding Generates a binding with 10 members for use in testing
func GenerateBinding() *cloudresourcemanager.Binding {
	rand.Seed(time.Now().UnixNano())

	role := GenerateRole(StringGenerator())
	var members []string

	for i := 0; i < rand.Intn(10)+1; i++ {
		members = append(members, GenerateMember(StringGenerator()))
	}

	return NewBinding(role, members...)
}

// GenerateBindings will create the specified number of bindings
func GenerateBindings(number int) (bindings []*cloudresourcemanager.Binding) {
	for i := 0; i < number; i++ {
		bindings = append(bindings, GenerateBinding())
	}
	return bindings
}

// GenerateMember creates a binding member with a random name off a base principal in the format of an email address
func GenerateMember(principal string) string {
	rand.Seed(time.Now().UnixNano() + int64(len(principal)))
	return fmt.Sprintf("%v-%d-%d-%v", principal, rand.Intn(99999), rand.Intn(99999), "@testdomain.co")
}

// GenerateRole creates a random string and returns it to be used as a role
func GenerateRole(role string) string {
	rand.Seed(time.Now().UnixNano() + int64(len(role)))
	return fmt.Sprintf("%v-%d-%d", role, rand.Intn(99999), rand.Intn(99999))
}

// StringGenerator returns a random string
func StringGenerator() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("randomString-%d%d", rand.Intn(99999), rand.Intn(99999))
}

// PolicyContains searches a policy for a role and returns its binding
func PolicyContains(policy *cloudresourcemanager.Policy, role string) *cloudresourcemanager.Binding {
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			return binding
		}
	}
	return nil
}

// BindingContains searches a binding for a member string and returns a boolean to indicate if found
func BindingContains(binding *cloudresourcemanager.Binding, member string) bool {
	for _, m := range binding.Members {
		if m == member {
			return true
		}
	}
	return false
}

// PolicyRoleMembers searches a Policy for a role and returns the members if they exist
func PolicyRoleMembers(policy *cloudresourcemanager.Policy, role string) ([]string, error) {
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			return binding.Members, nil
		}
	}
	return nil, fmt.Errorf("binding not found")
}
