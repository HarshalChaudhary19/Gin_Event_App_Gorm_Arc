package graph

import (
	"Gin_Event_App_Arc/repositories"
	"Gin_Event_App_Arc/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CategoryRepo      repositories.GraphRepoI
	CourseRepo        repositories.GraphRepoCourseI
	FileUploadService services.FileServiceI
}
