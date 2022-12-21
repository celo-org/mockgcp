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

		_, err:= service.Projects.GetIamPolicy(projectID, request).Do()

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

		service.Projects.ProjectList = append(service.Projects.ProjectList, projects...)

		want := policy
		got, _ := service.Projects.GetIamPolicy(projectID, request).Do()

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

		service.Projects.ProjectList = append(service.Projects.ProjectList, NewProject(projectID, nil))

		service.Projects.SetIamPolicy(projectID, request).Do()

		want := policy
		got := service.Projects.ProjectList[0].Policy

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
			t.Errorf("expected an error but got none")
		}

	})
}

func TestFolder_GetIamPolicy_Do(t *testing.T) {
	t.Run("should err if folder  doesn't exist", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		_, err:= service.Folders.GetIamPolicy(folderID, request).Do()

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

		service.Folders.FolderList = append(service.Folders.FolderList, folders...)

		want := policy
		got, _ := service.Folders.GetIamPolicy(folderID, request).Do()

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

		_, err := service.Folders.SetIamPolicy(folderID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}

	})

}
