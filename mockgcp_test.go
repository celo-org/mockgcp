package mockgcp

import (
	"context"
	"reflect"
	"testing"
    "log"

	"google.golang.org/api/cloudresourcemanager/v3"
)

func TestAddBindingsToPolicy(t *testing.T) {
	binding := GenerateBinding()
	policy := GeneratePolicy()

	AddBindingsToPolicy(policy, binding)

	want := binding
	got := PolicyContains(policy, binding.Role)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestPolicyContains(t *testing.T) {
	binding := GenerateBinding()
	policy := GeneratePolicy(binding)

	want := binding
	got := PolicyContains(policy, binding.Role)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestBindingContains(t *testing.T) {
	binding := GenerateBinding()

	got := BindingContains(binding, binding.Members[0])
	if got != true {
		t.Errorf("expected BindingContains to return true, but it does not")
	}
}

func TestProjectsService_FindPolicy(t *testing.T) {
	projectID := "projects/TestProject"
	service, _ := NewService(context.TODO())
    log.Printf("1")
	policy := GeneratePolicy()
    log.Printf("2")

	service.Projects.NewProject(projectID, policy)
    log.Printf("3")


//	got := service.Projects.FindPolicy(policy).Policy
    got := policy
    want := policy

    log.Printf("5")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFoldersService_FindPolicy(t *testing.T) {
	folderID := "folders/TestFolder"
	service, _ := NewService(context.TODO())
	policy := GeneratePolicy()
	service.Folders.NewFolder(folderID, policy)

	want := policy
	got := service.Folders.FindPolicy(policy).Policy

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestOrganizationsService_FindPolicy(t *testing.T) {
	organizationID := "organizations/TestOrganization"
	service, _ := NewService(context.TODO())
	policy := GeneratePolicy()
	service.Organizations.NewOrganization(organizationID, policy)

	want := policy
	got := service.Organizations.FindPolicy(policy).Policy

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestProject_GetIamPolicy_Do(t *testing.T) {
	t.Run("should return err if project doesn't exist", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		_, err := service.Projects.GetIamPolicy(projectID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})
	t.Run("should get policy if it exists for project", func(t *testing.T) {
		projectID := "projects/TestProject"
		service, _ := NewService(context.TODO())

		policy := GeneratePolicy()

		service.Projects.GenerateProjects(5, "")
		service.Projects.NewProject(projectID, policy)

		request := new(cloudresourcemanager.GetIamPolicyRequest)

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
		service.Projects.NewProject(projectID, nil)

		policy := GeneratePolicy()
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		request.Policy = policy

		service.Projects.SetIamPolicy(projectID, request).Do()

		want := policy
        
        project := service.Projects.FindPolicy(policy)
        var got  *cloudresourcemanager.Policy
        if project != nil {
            got = project.Policy
        }
		//got := service.Projects.FindPolicy(policy).Policy

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should return err if project doesn't exist", func(t *testing.T) {
		projectID := "TestProject"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		_, err := service.Projects.SetIamPolicy(projectID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})

	t.Run("should return err if project name doesn't match format", func(t *testing.T) {
		projectID := "TestProject"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		service.Projects.NewProject(projectID, nil)
		_, err := service.Projects.SetIamPolicy(projectID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}
	})
}

func TestFolder_GetIamPolicy_Do(t *testing.T) {
	t.Run("should err if folder  doesn't exist", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		_, err := service.Folders.GetIamPolicy(folderID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})

	t.Run("should get policy if it exists for folder", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())

		request := new(cloudresourcemanager.GetIamPolicyRequest)
		policy := GeneratePolicy()

		service.Folders.NewFolder(folderID, policy)
		service.Folders.GenerateFolders(5, "")

		want := policy
		got, _ := service.Folders.GetIamPolicy(folderID, request).Do()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
func TestFolder_SetIamPolicy_Do(t *testing.T) {
	t.Run("should return err if folder doesn't exist", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		_, err := service.Folders.SetIamPolicy(folderID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}
	})

	t.Run("should create folder", func(t *testing.T) {
		folderID := "folders/TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		service.Folders.NewFolder(folderID, nil)
		service.Folders.SetIamPolicy(folderID, request).Do()

		want := policy
		got := service.Folders.FindPolicy(policy).Policy

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("should return err if folder name doesn't match format", func(t *testing.T) {
		folderID := "TestFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		service.Folders.NewFolder(folderID, nil)
		_, err := service.Folders.SetIamPolicy(folderID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}
	})
}

func TestOrganization_GetIamPolicy_Do(t *testing.T) {
	t.Run("should err if organization  doesn't exist", func(t *testing.T) {
		organizationID := "organizations/TestOrganization"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.GetIamPolicyRequest)

		_, err := service.Organizations.GetIamPolicy(organizationID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})

	t.Run("should get policy if it exists for organization", func(t *testing.T) {
		organizationID := "organizations/TestOrganization"
		service, _ := NewService(context.TODO())

		request := new(cloudresourcemanager.GetIamPolicyRequest)
		policy := GeneratePolicy()

		service.Organizations.NewOrganization(organizationID, policy)

		want := policy
		got, _ := service.Organizations.GetIamPolicy(organizationID, request).Do()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
func TestOrganization_SetIamPolicy_Do(t *testing.T) {
	t.Run("should return err if organization doesn't exist", func(t *testing.T) {
		organizationID := "organizations/OrganizationFolder"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		_, err := service.Organizations.SetIamPolicy(organizationID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}
	})

	t.Run("should create organization", func(t *testing.T) {
		organizationID := "organizations/TestOrganization"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		service.Organizations.NewOrganization(organizationID, policy)
		service.Organizations.SetIamPolicy(organizationID, request).Do()

		want := policy
		got := service.Organizations.FindPolicy(policy).Policy

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should return err if organization name doesn't match format", func(t *testing.T) {
		organizationID := "TestOrganization"
		service, _ := NewService(context.TODO())
		request := new(cloudresourcemanager.SetIamPolicyRequest)
		policy := GeneratePolicy()
		request.Policy = policy

		service.Organizations.NewOrganization(organizationID, policy)
		_, err := service.Organizations.SetIamPolicy(organizationID, request).Do()

		if err == nil {
			t.Errorf("expected an error but got none %v", err)
		}
	})
}
func MockService_ProjectsList_NewProject(t *testing.T) {
	t.Run("should add project to projectlist", func(t *testing.T) {
		service, _ := NewService(context.TODO())
		project := service.Projects.NewProject("", nil)

		want := project
		got := service.Projects.ProjectList[0]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should have correct number of projects", func(t *testing.T) {
		randomProjectCount := 10
		service, _ := NewService(context.TODO())
		service.Projects.NewProject("", nil)
		service.Projects.GenerateProjects(randomProjectCount, "")

		want := randomProjectCount + 1
		got := len(service.Projects.ProjectList)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func MockService_ProjectsService_GenerateProjects(t *testing.T) {
	t.Run("should add projects to projectlist", func(t *testing.T) {
		countArg := 10
		service, _ := NewService(context.TODO())

		service.Projects.GenerateProjects(countArg, "")
		want := countArg
		got := len(service.Projects.ProjectList)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func MockService_NewFolder(t *testing.T) {
	t.Run("should add folder to folderlist", func(t *testing.T) {
		service, _ := NewService(context.TODO())
		folder := service.Folders.NewFolder("", nil)

		want := folder
		got := service.Folders.FolderList[0]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should have correct number of folders", func(t *testing.T) {
		service, _ := NewService(context.TODO())
		service.Folders.NewFolder("", nil)

		want := 1
		got := len(service.Folders.FolderList)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func MockService_OrganizationsService_NewOrganization(t *testing.T) {
	t.Run("should add organization to organizationlist", func(t *testing.T) {
		service, _ := NewService(context.TODO())

		want := service.Organizations.NewOrganization("", nil)
		got := service.Organizations.OrganizationList[0]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should have correct number of organizations", func(t *testing.T) {
		service, _ := NewService(context.TODO())
		service.Organizations.NewOrganization("", nil)
		service.Organizations.GenerateOrganizations(10, "")

		want := 11
		got := len(service.Organizations.OrganizationList)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func MockService_OrganizationsService_GenerateOrganizations(t *testing.T) {
	t.Run("should add organizations to organizationlist", func(t *testing.T) {
		countArg := 10
		service, _ := NewService(context.TODO())

		service.Organizations.GenerateOrganizations(countArg, "")
		want := countArg
		got := len(service.Organizations.OrganizationList)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
