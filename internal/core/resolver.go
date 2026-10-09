/*
|--------------------------------------------------------------------------
| Resolver map
|--------------------------------------------------------------------------
|
| ResolverMap for Query/Mutation fields, dynamic root composition, and
| merging module resolvers with subscription roots (compositeRoot).
|
*/

package core

import (
	"fmt"
	"reflect"
	"unicode"
)

// ResolverMap collects field resolver functions in a graphql-modules style (Query.user, User.name, ...).
type ResolverMap struct {
	query        map[string]any
	mutation     map[string]any
	subscription map[string]any
	objects      map[string]map[string]any
}

// NewResolverMap creates an empty resolver map.
func NewResolverMap() *ResolverMap {
	return &ResolverMap{
		query:        make(map[string]any),
		mutation:     make(map[string]any),
		subscription: make(map[string]any),
		objects:      make(map[string]map[string]any),
	}
}

// Query registers a Query field resolver function.
func (m *ResolverMap) Query(field string, fn any) *ResolverMap {
	m.query[field] = fn
	return m
}

// Mutation registers a Mutation field resolver function.
func (m *ResolverMap) Mutation(field string, fn any) *ResolverMap {
	m.mutation[field] = fn
	return m
}

// Subscription registers a Subscription field resolver function.
func (m *ResolverMap) Subscription(field string, fn any) *ResolverMap {
	m.subscription[field] = fn
	return m
}

// Object registers a field resolver on a GraphQL object type.
func (m *ResolverMap) Object(typeName, field string, fn any) *ResolverMap {
	if m.objects[typeName] == nil {
		m.objects[typeName] = make(map[string]any)
	}
	m.objects[typeName][field] = fn
	return m
}

func (m *ResolverMap) merge(other *ResolverMap) error {
	if other == nil {
		return nil
	}
	for k, v := range other.query {
		if _, exists := m.query[k]; exists {
			return fmt.Errorf("gogql: duplicate Query resolver %q", k)
		}
		m.query[k] = v
	}
	for k, v := range other.mutation {
		if _, exists := m.mutation[k]; exists {
			return fmt.Errorf("gogql: duplicate Mutation resolver %q", k)
		}
		m.mutation[k] = v
	}
	for k, v := range other.subscription {
		if _, exists := m.subscription[k]; exists {
			return fmt.Errorf("gogql: duplicate Subscription resolver %q", k)
		}
		m.subscription[k] = v
	}
	for typeName, fields := range other.objects {
		if m.objects[typeName] == nil {
			m.objects[typeName] = make(map[string]any)
		}
		for field, fn := range fields {
			if _, exists := m.objects[typeName][field]; exists {
				return fmt.Errorf("gogql: duplicate %s resolver %q", typeName, field)
			}
			m.objects[typeName][field] = fn
		}
	}
	return nil
}

func (m *ResolverMap) buildRoot() (any, error) {
	if m == nil {
		return struct{}{}, nil
	}

	ops := make(map[string]reflect.Value, 3)

	buildOp := func(name string, fields map[string]any) error {
		if len(fields) == 0 {
			return nil
		}
		opStruct, values, err := buildFieldStruct(fields)
		if err != nil {
			return err
		}
		v := reflect.New(opStruct).Elem()
		for k, val := range values {
			v.FieldByName(k).Set(val)
		}
		ops[name] = v
		return nil
	}

	if err := buildOp("Query", m.query); err != nil {
		return nil, err
	}
	if err := buildOp("Mutation", m.mutation); err != nil {
		return nil, err
	}
	if len(m.subscription) > 0 {
		return nil, fmt.Errorf("gogql: ResolverMap.Subscription is not supported (graph-gophers requires methods); set ModuleConfig.SubscriptionResolvers to a struct with subscription field methods")
	}

	if len(ops) == 0 && len(m.objects) == 0 {
		return struct{}{}, nil
	}
	if len(ops) == 1 {
		for _, v := range ops {
			return v.Addr().Interface(), nil
		}
	}

	return buildMultiOpRoot(ops)
}

func buildMultiOpRoot(ops map[string]reflect.Value) (any, error) {
	q, hasQ := ops["Query"]
	m, hasM := ops["Mutation"]
	s, hasS := ops["Subscription"]

	switch {
	case hasQ && hasM && hasS:
		return &rootQMS{query: q, mutation: m, subscription: s}, nil
	case hasQ && hasM:
		return &rootQM{query: q, mutation: m}, nil
	case hasQ && hasS:
		return &rootQS{query: q, subscription: s}, nil
	case hasM && hasS:
		return &rootMS{mutation: m, subscription: s}, nil
	case hasQ:
		return q.Addr().Interface(), nil
	case hasM:
		return m.Addr().Interface(), nil
	case hasS:
		return s.Addr().Interface(), nil
	default:
		return struct{}{}, nil
	}
}

type rootQMS struct{ query, mutation, subscription reflect.Value }

func (r *rootQMS) Query() any        { return r.query.Addr().Interface() }
func (r *rootQMS) Mutation() any     { return r.mutation.Addr().Interface() }
func (r *rootQMS) Subscription() any { return r.subscription.Addr().Interface() }

type rootQM struct{ query, mutation reflect.Value }

func (r *rootQM) Query() any    { return r.query.Addr().Interface() }
func (r *rootQM) Mutation() any { return r.mutation.Addr().Interface() }

type rootQS struct{ query, subscription reflect.Value }

func (r *rootQS) Query() any        { return r.query.Addr().Interface() }
func (r *rootQS) Subscription() any { return r.subscription.Addr().Interface() }

type rootMS struct{ mutation, subscription reflect.Value }

func (r *rootMS) Mutation() any     { return r.mutation.Addr().Interface() }
func (r *rootMS) Subscription() any { return r.subscription.Addr().Interface() }

func buildFieldStruct(fields map[string]any) (reflect.Type, map[string]reflect.Value, error) {
	structFields := make([]reflect.StructField, 0, len(fields))
	values := make(map[string]reflect.Value, len(fields))

	for graphQLName, fn := range fields {
		fnVal := reflect.ValueOf(fn)
		if fnVal.Kind() != reflect.Func {
			return nil, nil, fmt.Errorf("gogql: resolver %q must be a function", graphQLName)
		}
		goName := fieldNameForGraphQL(graphQLName)
		structFields = append(structFields, reflect.StructField{
			Name: goName,
			Type: fnVal.Type(),
			Tag:  reflect.StructTag(`graphql:"` + graphQLName + `"`),
		})
		values[goName] = fnVal
	}

	return reflect.StructOf(structFields), values, nil
}

func fieldNameForGraphQL(name string) string {
	if name == "" {
		return "X"
	}
	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func asResolverMap(v any) (*ResolverMap, bool) {
	m, ok := v.(*ResolverMap)
	return m, ok
}

func usesResolverMap(modules []*Module) bool {
	for _, mod := range modules {
		if _, ok := asResolverMap(mod.resolvers); ok {
			return true
		}
	}
	return false
}

func mergeResolvers(modules []*Module) (any, error) {
	combined := NewResolverMap()
	hasMap := false
	var structResolver any
	var subscriptionResolver any

	for _, mod := range modules {
		if mod.subscriptionResolvers != nil {
			if subscriptionResolver != nil {
				return nil, fmt.Errorf("gogql: module %q: only one SubscriptionResolvers root is supported across the application", mod.id)
			}
			subscriptionResolver = mod.subscriptionResolvers
		}
		if mod.resolvers == nil {
			continue
		}
		if rm, ok := asResolverMap(mod.resolvers); ok {
			hasMap = true
			if err := combined.merge(rm); err != nil {
				return nil, err
			}
			continue
		}
		if structResolver != nil {
			return nil, fmt.Errorf("gogql: module %q: only one struct-based root resolver is supported; use ResolverMap in other modules", mod.id)
		}
		structResolver = mod.resolvers
	}

	var queryRoot any
	var err error

	if hasMap {
		queryRoot, err = combined.buildRoot()
		if err != nil {
			return nil, err
		}
		if structResolver != nil {
			return nil, fmt.Errorf("gogql: cannot mix struct resolvers with ResolverMap across modules")
		}
	} else if structResolver != nil {
		queryRoot = structResolver
	}

	if subscriptionResolver != nil {
		return attachSubscriptionRoot(queryRoot, subscriptionResolver)
	}

	if queryRoot != nil {
		return queryRoot, nil
	}

	return struct{}{}, nil
}

func attachSubscriptionRoot(queryRoot, subscriptionRoot any) (any, error) {
	if queryRoot == nil {
		return subscriptionRoot, nil
	}
	return &compositeRoot{query: queryRoot, subscription: subscriptionRoot}, nil
}

type compositeRoot struct {
	query        any
	subscription any
}

func (c *compositeRoot) Query() any {
	if c == nil || c.query == nil {
		return struct{}{}
	}
	if q, ok := c.query.(interface{ Query() any }); ok {
		return q.Query()
	}
	return c.query
}

func (c *compositeRoot) Mutation() any {
	if c == nil || c.query == nil {
		return struct{}{}
	}
	if q, ok := c.query.(interface{ Mutation() any }); ok {
		return q.Mutation()
	}
	return struct{}{}
}

func (c *compositeRoot) Subscription() any {
	if c == nil {
		return struct{}{}
	}
	if s, ok := c.subscription.(interface{ Subscription() any }); ok {
		return s.Subscription()
	}
	return c.subscription
}
