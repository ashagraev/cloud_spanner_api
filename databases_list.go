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

func ListProjects() []string {
	ctx := context.Background()
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
		result = append(result, "projects/" + p.ProjectId)
	}
	return result
}

func ListInstances(projects []string) []string {
	ctx := context.Background()
	client, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	defer client.Close()

	var instancesList []string
	for _, project := range projects {
		req := &instancepb.ListInstancesRequest{
			Parent: project,
		}

		it := client.ListInstances(ctx, req)

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

func ListDatabases(instances []string) []string {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	defer adminClient.Close()

	var listDatabases []string
	for _, instance := range instances {
		req := &databasepb.ListDatabasesRequest{
			Parent: instance,
		}

		it := adminClient.ListDatabases(ctx, req)
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
