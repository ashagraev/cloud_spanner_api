package main

import (
	"context"
	"fmt"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func ListProjects(ctx context.Context) []string {
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	request := cloudresourcemanagerService.Projects.List()
	response, err := request.Do()
	if err != nil {
		fmt.Print(err)
		return []string{}
	}

	var result []string
	for _, p := range response.Projects {
		result = append(result, "projects/"+p.ProjectId)
	}
	return result
}

func ListInstances(ctx context.Context, instanceClient *instance.InstanceAdminClient, projects []string) []string {
	var instancesList []string
	for _, project := range projects {
		req := &instancepb.ListInstancesRequest{
			Parent: project,
		}

		it := instanceClient.ListInstances(ctx, req)

		for {
			resp, err := it.Next()

			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			instancesList = append(instancesList, resp.Name)
		}
	}

	return instancesList
}

func ListDatabases(ctx context.Context) []string {
	instanceClient, _ := instance.NewInstanceAdminClient(ctx)
	databaseClient, _ := database.NewDatabaseAdminClient(ctx)

	projects := ListProjects(ctx)
	instances := ListInstances(ctx, instanceClient, projects)

	var listDatabases []string
	for _, instance := range instances {
		req := &databasepb.ListDatabasesRequest{
			Parent: instance,
		}

		it := databaseClient.ListDatabases(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			listDatabases = append(listDatabases, resp.Name)
		}
	}
	return listDatabases
}
