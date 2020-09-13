package main

import (
	"context"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func ListProjects(ctx context.Context) ([]string, error) {
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return []string{}, err
	}

	request := cloudresourcemanagerService.Projects.List()
	response, err := request.Do()
	if err != nil {
		return []string{}, err
	}

	var result []string
	for _, p := range response.Projects {
		result = append(result, "projects/"+p.ProjectId)
		LogProjectLoad(ctx, p.ProjectId)
	}
	return result, nil
}

func ListInstances(ctx context.Context, instanceClient *instance.InstanceAdminClient, projects []string) ([]string, error) {
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
				return []string{}, err
			}
			instancesList = append(instancesList, resp.Name)
			LogInstanceLoad(ctx, resp.Name)
		}
	}

	return instancesList, nil
}

func ListDatabases(ctx context.Context) ([]string, error) {
	instanceClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return []string{}, err
	}
	databaseClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return []string{}, err
	}

	projects, err := ListProjects(ctx)
	if err != nil {
		return []string{}, err
	}
	instances, err := ListInstances(ctx, instanceClient, projects)
	if err != nil {
		return []string{}, err
	}

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
				return []string{}, err
			}
			listDatabases = append(listDatabases, resp.Name)
			LogDataBaseLoad(ctx, resp.Name)
		}
	}
	return listDatabases, nil
}
