package mockgcp

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v3"
)

func TestProject_GetIamPolicy_Do(t *testing.T) {
	t.Run("should err if project doesn't exist", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		_, err:= service.Project.GetIamPolicy(projectID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})
	t.Run("should get policy if it exists for project", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())

		request := new(cloudresourcemanager.GetIamPolicyRequest)
		policy := GeneratePolicy(nil)
		project := NewProject(projectID, policy)
		projects := append(GenerateProjects(5), project)

		service.Project.Projects = append(service.Project.Projects, projects...)

		want := policy
		got, _ := service.Project.GetIamPolicy(projectID, request).Do()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
func TestProject_SetIamPolicy_Do(t *testing.T) {
	t.Run("should create policy", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy(nil)
		request.Policy = policy

		service.Project.Projects = append(service.Project.Projects, NewProject(projectID, nil))

		service.Project.SetIamPolicy(projectID, request).Do()

		want := policy
		got := service.Project.Projects[0].Policy

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

		_, err := service.Project.SetIamPolicy(projectID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}

	})
}

func TestFolder_GetIamPolicy_Do(t *testing.T) {
	t.Run("should err if folder  doesn't exist", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		_, err:= service.Folder.GetIamPolicy(folderID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})

	t.Run("should get policy if it exists for folder", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())

		request := new(cloudresourcemanager.GetIamPolicyRequest)
		policy := GeneratePolicy(nil)
		folder := NewFolder(folderID, policy)
		folders := append(GenerateFolders(5), folder)

		service.Folder.Folders = append(service.Folder.Folders, folders...)

		want := policy
		got, _ := service.Folder.GetIamPolicy(folderID, request).Do()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should return err if folder doesn't exist", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy(nil)
		request.Policy = policy

		_, err := service.Folder.SetIamPolicy(folderID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}

	})

}
