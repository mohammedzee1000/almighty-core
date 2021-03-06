package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.Resource("workitem", func() {
	a.BasePath("/workitems")

	a.Action("show", func() {
		a.Routing(
			a.GET("/:id"),
		)
		a.Description("Retrieve work item with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK, func() {
			a.Media(workItem)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List work items.")
		a.Params(func() {
			a.Param("filter", d.String, "a query language expression restricting the set of found work items")
			a.Param("page", d.String, "Paging in the format <start>,<limit>")
		})
		a.Response(d.OK, func() {
			a.Media(a.CollectionOf(workItem))
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Description("create work item with type and id.")
		a.Payload(CreateWorkItemPayload)
		a.Response(d.Created, "/workitems/.*", func() {
			a.Media(workItem)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("delete", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE("/:id"),
		)
		a.Description("Delete work item with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("update", func() {
		a.Security("jwt")
		a.Routing(
			a.PUT("/:id"),
		)
		a.Description("update the given work item with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Payload(UpdateWorkItemPayload)
		a.Response(d.OK, func() {
			a.Media(workItem)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

})

// new version of "list" for migration
var _ = a.Resource("workitem.2", func() {
	a.BasePath("/workitems.2")
	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List work items.")
		a.Params(func() {
			a.Param("filter", d.String, "a query language expression restricting the set of found work items")
			a.Param("page[offset]", d.String, "Paging start position")
			a.Param("page[limit]", d.Integer, "Paging size")
		})
		a.Response(d.OK, func() {
			a.Media(workItemListResponse)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("update", func() {
		a.Security("jwt")
		a.Routing(
			a.PATCH("/:id"),
		)
		a.Description("update the work item with the given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Payload(updateWorkItemJSONAPIPayload)
		a.Response(d.OK, func() {
			// Still using workitem in MediaTypes.
			// ToDo update to struct which complies to jsonapi
			a.Media(workItem2)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
})

var _ = a.Resource("workitemtype", func() {

	a.BasePath("/workitemtypes")

	a.Action("show", func() {

		a.Routing(
			a.GET("/:name"),
		)
		a.Description("Retrieve work item type with given name.")
		a.Params(func() {
			a.Param("name", d.String, "name")
		})
		a.Response(d.OK, func() {
			a.Media(workItemType)
		})
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Description("Create work item type.")
		a.Payload(CreateWorkItemTypePayload)
		a.Response(d.Created, "/workitemtypes/.*", func() {
			a.Media(workItemType)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List work item types.")
		a.Params(func() {
			a.Param("page", d.String, "Paging in the format <start>,<limit>")
		})
		a.Response(d.OK, func() {
			a.Media(a.CollectionOf(workItemType))
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})

var _ = a.Resource("user", func() {
	a.BasePath("/user")

	a.Action("show", func() {
		a.Security("jwt")
		a.Routing(
			a.GET(""),
		)
		a.Description("Get the authenticated user")
		a.Response(d.OK, func() {
			a.Media(User)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

})

var _ = a.Resource("identity", func() {
	a.BasePath("/identities")

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List all identities.")
		a.Response(d.OK, func() {
			a.Media(identityArray)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})

var _ = a.Resource("status", func() {

	a.DefaultMedia(ALMStatus)
	a.BasePath("/status")

	a.Action("show", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("Show the status of the current running instance")
		a.Response(d.OK)
		a.Response(d.ServiceUnavailable, ALMStatus)
	})
})

var _ = a.Resource("login", func() {

	a.BasePath("/login")

	a.Action("authorize", func() {
		a.Routing(
			a.GET("authorize"),
		)
		a.Description("Authorize with the ALM")
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.TemporaryRedirect)
	})

	a.Action("generate", func() {
		a.Routing(
			a.GET("generate"),
		)
		a.Description("Generates a set of Tokens for different Auth levels. NOT FOR PRODUCTION. Only available if server is running in dev mode")
		a.Response(d.OK, func() {
			a.Media(a.CollectionOf(AuthToken))
		})
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
})

var _ = a.Resource("tracker", func() {
	a.BasePath("/trackers")

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List all tracker configurations.")
		a.Params(func() {
			a.Param("filter", d.String, "a query language expression restricting the set of found items")
			a.Param("page", d.String, "Paging in the format <start>,<limit>")
		})
		a.Response(d.OK, func() {
			a.Media(a.CollectionOf(Tracker))
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("show", func() {
		a.Routing(
			a.GET("/:id"),
		)
		a.Description("Retrieve tracker configuration for the given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK, func() {
			a.Media(Tracker)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Description("Add new tracker configuration.")
		a.Payload(CreateTrackerAlternatePayload)
		a.Response(d.Created, "/trackers/.*", func() {
			a.Media(Tracker)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("delete", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE("/:id"),
		)
		a.Description("Delete tracker configuration.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("update", func() {
		a.Security("jwt")
		a.Routing(
			a.PUT("/:id"),
		)
		a.Description("Update tracker configuration.")
		a.Payload(UpdateTrackerAlternatePayload)
		a.Response(d.OK, func() {
			a.Media(Tracker)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

})

var _ = a.Resource("trackerquery", func() {
	a.BasePath("/trackerqueries")
	a.Action("show", func() {
		a.Routing(
			a.GET("/:id"),
		)
		a.Description("Retrieve tracker configuration for the given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK, func() {
			a.Media(TrackerQuery)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})
	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Description("Add new tracker query.")
		a.Payload(CreateTrackerQueryAlternatePayload)
		a.Response(d.Created, "/trackerqueries/.*", func() {
			a.Media(TrackerQuery)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("update", func() {
		a.Security("jwt")
		a.Routing(
			a.PUT("/:id"),
		)
		a.Description("Update tracker query.")
		a.Payload(UpdateTrackerQueryAlternatePayload)
		a.Response(d.OK, func() {
			a.Media(TrackerQuery)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("delete", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE("/:id"),
		)
		a.Description("Delete tracker query")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List all tracker queries.")
		a.Response(d.OK, func() {
			a.Media(a.CollectionOf(TrackerQuery))
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})
})

var _ = a.Resource("search", func() {
	a.BasePath("/search")

	a.Action("show", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("Search by ID, URL, full text capability")
		a.Params(func() {
			a.Param("q", d.String,
				`Following are valid input for seach query
				1) "id:100" :- Look for work item hainvg id 100
				2) "url:http://demo.almighty.io/details/500" :- Search on WI having id 500 and check 
					if this URL is mentioned in searchable columns of work item
				3) "simple keywords seperated by space" :- Search in Work Items based on these keywords.`)
			a.Param("page[offset]", d.String, "Paging start position") // #428
			a.Param("page[limit]", d.Integer, "Paging size")
			a.Required("q")
		})
		a.Response(d.OK, func() {
			a.Media(searchResponse)
		})

		a.Response(d.BadRequest, JSONAPIErrors)

		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})

var _ = a.Resource("work-item-link-category", func() {
	a.BasePath("/workitemlinkcategories")

	a.Action("show", func() {
		a.Routing(
			a.GET("/:id"),
		)
		a.Description("Retrieve work item link category (as JSONAPI) for the given ID.")
		a.Params(func() {
			a.Param("id", d.String, "ID of the work item link category")
		})
		a.Response(d.OK, func() {
			a.Media(WorkItemLinkCategory)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List work item link categories.")
		a.Response(d.OK, func() {
			a.Media(WorkItemLinkCategoryArray)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Description("Create a work item link category")
		a.Payload(CreateWorkItemLinkCategoryPayload)
		a.Response(d.Created, "/workitemlinkcategories/.*", func() {
			a.Media(WorkItemLinkCategory)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("delete", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE("/:id"),
		)
		a.Description("Delete work item link category with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("update", func() {
		a.Security("jwt")
		a.Routing(
			a.PATCH("/:id"),
		)
		a.Description("Update the given work item link category with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Payload(UpdateWorkItemLinkCategoryPayload)
		a.Response(d.OK, func() {
			a.Media(WorkItemLinkCategory)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
})

var _ = a.Resource("work-item-link-type", func() {
	a.BasePath("/workitemlinktypes")

	a.Action("show", func() {
		a.Routing(
			a.GET("/:id"),
		)
		a.Description("Retrieve work item link type (as JSONAPI) for the given link ID.")
		a.Params(func() {
			a.Param("id", d.String, "ID of the work item link type")
		})
		a.Response(d.OK, func() {
			a.Media(WorkItemLinkType)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List work item link types.")
		a.Response(d.OK, func() {
			a.Media(WorkItemLinkTypeArray)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Description("Create a work item link type")
		a.Payload(CreateWorkItemLinkTypePayload)
		a.Response(d.Created, "/workitemlinktypes/.*", func() {
			a.Media(WorkItemLinkType)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("delete", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE("/:id"),
		)
		a.Description("Delete work item link type with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Response(d.OK)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("update", func() {
		a.Security("jwt")
		a.Routing(
			a.PATCH("/:id"),
		)
		a.Description("Update the given work item link type with given id.")
		a.Params(func() {
			a.Param("id", d.String, "id")
		})
		a.Payload(UpdateWorkItemLinkTypePayload)
		a.Response(d.OK, func() {
			a.Media(WorkItemLinkType)
		})
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
})

var _ = a.Resource("work-item-link", func() {
	a.BasePath("/workitemlinks")
	a.Action("show", showWorkItemLink)
	a.Action("list", listWorkItemLinks)
	a.Action("create", createWorkItemLink)
	a.Action("delete", deleteWorkItemLink)
	a.Action("update", updateWorkItemLink)
})

var _ = a.Resource("work-item-relationships-links", func() {
	a.BasePath("/relationships/links")
	a.Parent("workitem")
	a.Action("show", showWorkItemLink)
	a.Action("delete", deleteWorkItemLink)
	a.Action("update", updateWorkItemLink)
	a.Action("list", func() {
		listWorkItemLinks()
		a.Description("List work item links associated with the given work item (either as source or as target work item).")
		a.Response(d.NotFound, JSONAPIErrors, func() {
			a.Description("This error arises when the given work item does not exist.")
		})
	})
	a.Action("create", func() {
		createWorkItemLink()
		a.Response(d.NotFound, JSONAPIErrors, func() {
			a.Description("This error arises when the given work item does not exist.")
		})
	})
})

// listWorkItemLinks defines the list action for endpoints that return an array
// of work item links.
func listWorkItemLinks() {
	a.Description("Retrieve work item link (as JSONAPI) for the given link ID.")
	a.Routing(
		a.GET(""),
	)
	a.Response(d.OK, func() {
		a.Media(WorkItemLinkArray)
	})
	a.Response(d.BadRequest, JSONAPIErrors)
	a.Response(d.InternalServerError, JSONAPIErrors)
}

func showWorkItemLink() {
	a.Description("Retrieve work item link (as JSONAPI) for the given link ID.")
	a.Routing(
		a.GET("/:linkId"),
	)
	a.Params(func() {
		a.Param("linkId", d.String, "ID of the work item link to show")
	})
	a.Response(d.OK, func() {
		a.Media(WorkItemLink)
	})
	a.Response(d.BadRequest, JSONAPIErrors)
	a.Response(d.InternalServerError, JSONAPIErrors)
	a.Response(d.NotFound, JSONAPIErrors)
}

func createWorkItemLink() {
	a.Description("Create a work item link")
	a.Security("jwt")
	a.Routing(
		a.POST(""),
	)
	a.Payload(CreateWorkItemLinkPayload)
	a.Response(d.Created, "/workitemlinks/.*", func() {
		a.Media(WorkItemLink)
	})
	a.Response(d.BadRequest, JSONAPIErrors)
	a.Response(d.InternalServerError, JSONAPIErrors)
	a.Response(d.Unauthorized, JSONAPIErrors)
}

func deleteWorkItemLink() {
	a.Description("Delete work item link with given id.")
	a.Security("jwt")
	a.Routing(
		a.DELETE("/:linkId"),
	)
	a.Params(func() {
		a.Param("linkId", d.String, "ID of the work item link to be deleted")
	})
	a.Response(d.OK)
	a.Response(d.BadRequest, JSONAPIErrors)
	a.Response(d.InternalServerError, JSONAPIErrors)
	a.Response(d.NotFound, JSONAPIErrors)
	a.Response(d.Unauthorized, JSONAPIErrors)
}

func updateWorkItemLink() {
	a.Description("Update the given work item link with given id.")
	a.Security("jwt")
	a.Routing(
		a.PATCH("/:linkId"),
	)
	a.Params(func() {
		a.Param("linkId", d.String, "ID of the work item link to be updated")
	})
	a.Payload(UpdateWorkItemLinkPayload)
	a.Response(d.OK, func() {
		a.Media(WorkItemLink)
	})
	a.Response(d.BadRequest, JSONAPIErrors)
	a.Response(d.InternalServerError, JSONAPIErrors)
	a.Response(d.NotFound, JSONAPIErrors)
	a.Response(d.Unauthorized, JSONAPIErrors)
}
