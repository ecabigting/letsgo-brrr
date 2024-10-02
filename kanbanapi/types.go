package main

type RouteResponse struct {
	// annotate it to json to allow marshalling and unmarshalling
	Message string `json:"message"`
	//for endpoints that doesnt have the id as part of the url request
	ID string `json:"id,omitempty"`
}
