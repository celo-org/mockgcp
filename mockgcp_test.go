package mockgcp

import(
    "testing"
)

func TestAddPolicyToProject(t *testing.T) {
    request := new(cloudresourcemanager.SetIamPolicyRequest)
    request.Policy = policy
    projectID = "testProject"
    crmService.Projects.SetIamPolicy("testProject", request).Do()
