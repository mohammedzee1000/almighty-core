package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/application"
	"github.com/almighty/almighty-core/errors"
	"github.com/almighty/almighty-core/jsonapi"
	query "github.com/almighty/almighty-core/query/simple"
	"github.com/almighty/almighty-core/workitem"
	"github.com/goadesign/goa"
)

const (
	pageSizeDefault = 20
	pageSizeMax     = 100
)

// Workitem2Controller implements the workitem.2 resource.
type Workitem2Controller struct {
	*goa.Controller
	db application.DB
}

// NewWorkitem2Controller creates a workitem.2 controller.
func NewWorkitem2Controller(service *goa.Service, db application.DB) *Workitem2Controller {
	if db == nil {
		panic("db must not be nil")
	}
	return &Workitem2Controller{Controller: service.NewController("WorkitemController"), db: db}
}

func setPagingLinks(links *app.PagingLinks, path string, resultLen, offset, limit, count int) {

	// prev link
	if offset > 0 && count > 0 {
		var prevStart int
		// we do have a prev link
		if offset <= count {
			prevStart = offset - limit
		} else {
			// the first range that intersects the end of the useful range
			prevStart = offset - (((offset-count)/limit)+1)*limit
		}
		realLimit := limit
		if prevStart < 0 {
			// need to cut the range to start at 0
			realLimit = limit + prevStart
			prevStart = 0
		}
		prev := fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d", path, prevStart, realLimit)
		links.Prev = &prev
	}

	// next link
	nextStart := offset + resultLen
	if nextStart < count {
		// we have a next link
		next := fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d", path, nextStart, limit)
		links.Next = &next
	}

	// first link
	var firstEnd int
	if offset > 0 {
		firstEnd = offset % limit // this is where the second page starts
	} else {
		// offset == 0, first == current
		firstEnd = limit
	}
	first := fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d", path, 0, firstEnd)
	links.First = &first

	// last link
	var lastStart int
	if offset < count {
		// advance some pages until touching the end of the range
		lastStart = offset + (((count - offset - 1) / limit) * limit)
	} else {
		// retreat at least one page until covering the range
		lastStart = offset - ((((offset - count) / limit) + 1) * limit)
	}
	realLimit := limit
	if lastStart < 0 {
		// need to cut the range to start at 0
		realLimit = limit + lastStart
		lastStart = 0
	}
	last := fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d", path, lastStart, realLimit)
	links.Last = &last
}

// List runs the list action.
// Prev and Next links will be present only when there actually IS a next or previous page.
// Last will always be present. Total Item count needs to be computed from the "Last" link.
func (c *Workitem2Controller) List(ctx *app.ListWorkitem2Context) error {
	// Workitem2Controller_List: start_implement

	exp, err := query.Parse(ctx.Filter)
	if err != nil {
		jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrBadRequest(fmt.Sprintf("could not parse filter: %s", err.Error())))
		return ctx.BadRequest(jerrors)
	}
	var offset int
	var limit int

	if ctx.PageOffset == nil {
		offset = 0
	} else {
		offsetValue, err := strconv.Atoi(*ctx.PageOffset)
		if err != nil {
			offset = 0
		} else {
			offset = offsetValue
		}
	}
	if offset < 0 {
		offset = 0
	}

	if ctx.PageLimit == nil {
		limit = pageSizeDefault
	} else {
		limit = *ctx.PageLimit
	}

	if limit <= 0 {
		limit = pageSizeDefault
	} else if limit > pageSizeMax {
		limit = pageSizeMax
	}

	return application.Transactional(c.db, func(tx application.Application) error {
		result, c, err := tx.WorkItems().List(ctx.Context, exp, &offset, &limit)
		count := int(c)
		if err != nil {
			switch err := err.(type) {
			case errors.BadParameterError:
				jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrBadRequest(fmt.Sprintf("Error listing work items: %s", err.Error())))
				return ctx.BadRequest(jerrors)
			default:
				log.Printf("Error listing work items: %s", err.Error())
				jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrInternal(fmt.Sprintf("Error listing work items: %s", err.Error())))
				return ctx.InternalServerError(jerrors)
			}
		}

		response := app.WorkItemListResponse{
			Links: &app.PagingLinks{},
			Meta:  &app.WorkItemListResponseMeta{TotalCount: count},
			Data:  result,
		}

		setPagingLinks(response.Links, buildAbsoluteURL(ctx.RequestData), len(result), offset, limit, count)

		return ctx.OK(&response)
	})

	// Workitem2Controller_List: end_implement
}

func buildAbsoluteURL(req *goa.RequestData) string {
	scheme := "http"
	if req.TLS != nil { // isHTTPS
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, req.Host, req.URL.Path)
}

// ConvertWorkItemToJSONAPI is responsible for converting given WorkItem model object into a
// response resource object by jsonapi.org specifications
func (c *Workitem2Controller) ConvertWorkItemToJSONAPI(ctx *app.UpdateWorkitem2Context, wi app.WorkItem) *app.WorkItem2 {
	// construct default values from input WI

	absoluteURL := buildAbsoluteURL(ctx.RequestData) // it includes path hence no modifications needed
	op := &app.WorkItem2{
		Links: &app.WorkItemResourceLinksForJSONAPI{
			Self: &absoluteURL,
		},
		Data: &app.WorkItemDataForUpdate{
			ID:   wi.ID,
			Type: workitem.APIStinrgTypeWorkItem,
			Attributes: map[string]interface{}{
				"version": wi.Version,
			},
			Relationships: &app.WorkItemRelationships{
				BaseType: &app.RelationshipBaseType{
					Data: &app.BaseTypeData{
						ID:   wi.Type,
						Type: workitem.APIStinrgTypeWorkItemType,
					},
				},
			},
		},
	}
	// Move fields into Relationships or Attributes as needed
	for name, val := range wi.Fields {
		switch name {
		case workitem.SystemAssignee:
			if val != nil {
				valStr := val.(string)
				op.Data.Relationships.Assignee = &app.RelationAssignee{
					Data: &app.AssigneeData{
						ID:   &valStr,
						Type: workitem.APIStinrgTypeAssignee,
					},
				}
			}
		default:
			op.Data.Attributes[name] = val
		}
	}
	return op
}

// Update does PATCH workitem
func (c *Workitem2Controller) Update(ctx *app.UpdateWorkitem2Context) error {
	return application.Transactional(c.db, func(appl application.Application) error {

		toSave := app.WorkItemDataForUpdate{
			ID:            ctx.Payload.Data.ID,
			Type:          ctx.Payload.Data.Type,
			Relationships: ctx.Payload.Data.Relationships,
			Attributes:    ctx.Payload.Data.Attributes,
		}
		wi, err := appl.WorkItems2().Save(ctx, toSave)
		if err != nil {
			switch err := err.(type) {
			case errors.BadParameterError:
				jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrBadRequest(fmt.Sprintf("Error updating work item: %s", err.Error())))
				return ctx.BadRequest(jerrors)
			case errors.NotFoundError:
				jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrNotFound(err.Error()))
				return ctx.NotFound(jerrors)
			case errors.VersionConflictError:
				jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrBadRequest(fmt.Sprintf("Error updating work item: %s", err.Error())))
				return ctx.BadRequest(jerrors)
			default:
				log.Printf("Error updating work items: %s", err.Error())
				jerrors, _ := jsonapi.ErrorToJSONAPIErrors(goa.ErrInternal(err.Error()))
				return ctx.InternalServerError(jerrors)
			}
		}
		return ctx.OK(c.ConvertWorkItemToJSONAPI(ctx, *wi))
	})
}
