package main

type RouteResponse struct {
	// annotate it to json to allow marshalling and unmarshalling
	Message string `json:"message"`
}
