// Code generated by jrohc. DO NOT EDIT.

package petstoreapi

import (
	context "context"
	apicommon "github.com/go-tk/jroh/go/apicommon"
	http "net/http"
)

type PetClient interface {
	AddPet(ctx context.Context, params *AddPetParams) (err error)
	GetPet(ctx context.Context, params *GetPetParams) (results *GetPetResults, err error)
	GetPets(ctx context.Context, params *GetPetsParams) (results *GetPetsResults, err error)
	UpdatePet(ctx context.Context, params *UpdatePetParams) (err error)
	FindPets(ctx context.Context, params *FindPetsParams) (results *FindPetsResults, err error)
}

type petClient struct {
	apicommon.Client

	rpcFiltersTable [5][]apicommon.RPCHandler
	transportTable  [5]http.RoundTripper
}

func NewPetClient(rpcBaseURL string, options apicommon.ClientOptions) PetClient {
	options.Sanitize()
	var c petClient
	c.Init(rpcBaseURL, options.Timeout)
	apicommon.FillRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
	apicommon.FillTransportTable(c.transportTable[:], options.Transport, options.Middlewares)
	return &c
}

func (c *petClient) AddPet(ctx context.Context, params *AddPetParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      AddPetParams
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Pet_AddPet]
	s.OutgoingRPC.Init("Petstore", "Pet", "AddPet", &s.Params, nil, apicommon.HandleRPC, rpcFilters)
	transport := c.transportTable[Pet_AddPet]
	return c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Pet.AddPet")
}

func (c *petClient) GetPet(ctx context.Context, params *GetPetParams) (*GetPetResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      GetPetParams
		Results     GetPetResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Pet_GetPet]
	s.OutgoingRPC.Init("Petstore", "Pet", "GetPet", &s.Params, &s.Results, apicommon.HandleRPC, rpcFilters)
	transport := c.transportTable[Pet_GetPet]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Pet.GetPet"); err != nil {
		return nil, err
	}
	return &s.Results, nil
}

func (c *petClient) GetPets(ctx context.Context, params *GetPetsParams) (*GetPetsResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      GetPetsParams
		Results     GetPetsResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Pet_GetPets]
	s.OutgoingRPC.Init("Petstore", "Pet", "GetPets", &s.Params, &s.Results, apicommon.HandleRPC, rpcFilters)
	transport := c.transportTable[Pet_GetPets]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Pet.GetPets"); err != nil {
		return nil, err
	}
	return &s.Results, nil
}

func (c *petClient) UpdatePet(ctx context.Context, params *UpdatePetParams) error {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      UpdatePetParams
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Pet_UpdatePet]
	s.OutgoingRPC.Init("Petstore", "Pet", "UpdatePet", &s.Params, nil, apicommon.HandleRPC, rpcFilters)
	transport := c.transportTable[Pet_UpdatePet]
	return c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Pet.UpdatePet")
}

func (c *petClient) FindPets(ctx context.Context, params *FindPetsParams) (*FindPetsResults, error) {
	var s struct {
		OutgoingRPC apicommon.OutgoingRPC
		Params      FindPetsParams
		Results     FindPetsResults
	}
	s.Params = *params
	rpcFilters := c.rpcFiltersTable[Pet_FindPets]
	s.OutgoingRPC.Init("Petstore", "Pet", "FindPets", &s.Params, &s.Results, apicommon.HandleRPC, rpcFilters)
	transport := c.transportTable[Pet_FindPets]
	if err := c.DoRPC(ctx, &s.OutgoingRPC, transport, "/rpc/Petstore.Pet.FindPets"); err != nil {
		return nil, err
	}
	return &s.Results, nil
}

type PetClientFuncs struct {
	AddPetFunc    func(context.Context, *AddPetParams) error
	GetPetFunc    func(context.Context, *GetPetParams) (*GetPetResults, error)
	GetPetsFunc   func(context.Context, *GetPetsParams) (*GetPetsResults, error)
	UpdatePetFunc func(context.Context, *UpdatePetParams) error
	FindPetsFunc  func(context.Context, *FindPetsParams) (*FindPetsResults, error)
}

var _ PetClient = (*PetClientFuncs)(nil)

func (cf *PetClientFuncs) AddPet(ctx context.Context, params *AddPetParams) error {
	if f := cf.AddPetFunc; f != nil {
		return f(ctx, params)
	}
	return apicommon.ErrNotImplemented
}

func (cf *PetClientFuncs) GetPet(ctx context.Context, params *GetPetParams) (*GetPetResults, error) {
	if f := cf.GetPetFunc; f != nil {
		return f(ctx, params)
	}
	return nil, apicommon.ErrNotImplemented
}

func (cf *PetClientFuncs) GetPets(ctx context.Context, params *GetPetsParams) (*GetPetsResults, error) {
	if f := cf.GetPetsFunc; f != nil {
		return f(ctx, params)
	}
	return nil, apicommon.ErrNotImplemented
}

func (cf *PetClientFuncs) UpdatePet(ctx context.Context, params *UpdatePetParams) error {
	if f := cf.UpdatePetFunc; f != nil {
		return f(ctx, params)
	}
	return apicommon.ErrNotImplemented
}

func (cf *PetClientFuncs) FindPets(ctx context.Context, params *FindPetsParams) (*FindPetsResults, error) {
	if f := cf.FindPetsFunc; f != nil {
		return f(ctx, params)
	}
	return nil, apicommon.ErrNotImplemented
}
