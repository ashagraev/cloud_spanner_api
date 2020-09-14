package main

import (
	"context"
	"fmt"

	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"

	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

func ListProjects(ctx context.Context) ([]string, error) {
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudresourcemanager.NewService error: %v", err)
	}

	request := cloudresourcemanagerService.Projects.List()
	response, err := request.Do()
	if err != nil {
		return nil, fmt.Errorf("cloudresourcemanagerService.Projects.List() error: %v", err)
	}

	var result []string
	for _, p := range response.Projects {
		result = append(result, "projects/"+p.ProjectId)
		LogProjectLoad(ctx, p.ProjectId)
	}
	return result, nil
}

func (db *DatabaseClient) ListInstances(ctx context.Context, projects []string) ([]string, error) {
	var instancesList []string
	for _, project := range projects {
		req := &instancepb.ListInstancesRequest{
			Parent: project,
		}

		it := db.instanceAdminClient.ListInstances(ctx, req)
		for {
			resp, err := it.Next()

			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("instance.InstanceIterator.Next() error: %v", err)
			}
			instancesList = append(instancesList, resp.Name)
			LogInstanceLoad(ctx, resp.Name)
		}
	}

	return instancesList, nil
}

func (db *DatabaseClient) ListDatabases(ctx context.Context) ([]string, error) {
	projects, err := ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	instances, err := db.ListInstances(ctx, projects)
	if err != nil {
		return nil, err
	}

	var listDatabases []string
	for _, instance := range instances {
		req := &databasepb.ListDatabasesRequest{
			Parent: instance,
		}

		it := db.databaseAdminClient.ListDatabases(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("database.DatabaseIterator.Next() error for instance %v: %v", instance, err)
			}
			listDatabases = append(listDatabases, resp.Name)
			LogDataBaseLoad(ctx, resp.Name)
		}
	}
	return listDatabases, nil
}
