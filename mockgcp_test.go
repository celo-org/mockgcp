package mockgcp

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v3"
)

func TestGetIamPolicy_Do(t *testing.T) {
	t.Run("should return blank policy if project doesn't exist", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		want := &cloudresourcemanager.Policy{}
		got, _ := service.Projects.GetIamPolicy(projectID, request).Do()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})
	t.Run("should get policy if it exists for project", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())

		request := new(cloudresourcemanager.GetIamPolicyRequest)
        policy := GeneratePolicy(nil)
        project := NewProject(projectID, policy)
        projects := append(GenerateProjects(5), project)

        service.Projects.p = append(service.Projects.p, projects...)
        
        want := policy
		got, _ := service.Projects.GetIamPolicy(projectID, request).Do()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
func TestSetIamPolicy_Do(t *testing.T) {
	t.Run("should create policy", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())
        request := new(cloudresourcemanager.SetIamPolicyRequest)
        policy := GeneratePolicy(nil)
        request.Policy = policy
        
        service.Projects.p = append(service.Projects.p, NewProject(projectID, nil))

        service.Projects.SetIamPolicy(projectID, request).Do()

		want := policy
		got := service.Projects.p[0].policy

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})
	t.Run("should return err if project doesn't exist", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())
        request := new(cloudresourcemanager.SetIamPolicyRequest)
        policy := GeneratePolicy(nil)
        request.Policy = policy
        

        _, err := service.Projects.SetIamPolicy(projectID, request).Do()

        if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}

	})


}
