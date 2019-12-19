package corral

import (
	"reflect"
	"strings"
)

// ConditionFunc is a function that returns a bool representing whether or not the subject meets some pre-defined criteria on the object.
type ConditionFunc func(subject interface{}, object interface{}) bool

// Permission is the core permission struct.
type Permission struct {
	// SubjectID is the ID of the entity that may perform an action on an object.
	SubjectID int
	// ObjectType is the name of the entity that may have an action performed on it.
	ObjectType string
	// Action is an individual CRUDM action.
	Action Action
	// Condition is a more granular control of whether or not the subject meets some pre-defined criteria on the object.
	Condition ConditionFunc
}

type subjectWithIDMethod interface {
	ID() int
}

type subjectWithSubjectIDMethod interface {
	SubjectID() int
}

// Permissions is a simple helper for a slice of Permission structs.
type Permissions []*Permission

// permissionSet is the core global permission set.
var permissionSet Permissions

// Action is a CRUDM action.
type Action int

const (
	CreateAction = iota
	ReadAction
	UpdateAction
	DeleteAction
	ManageAction
)

func noConditionFunc(subject interface{}, object interface{}) bool {
	return true
}

func getTypeString(v interface{}) string {
	return strings.TrimPrefix(reflect.TypeOf(v).String(), "*")
}

func getSubjectID(v interface{}) int {
	withIDMethod, ok := v.(subjectWithIDMethod)
	if ok {
		return withIDMethod.ID()
	}

	withSubjectIDMethod, ok := v.(subjectWithSubjectIDMethod)
	if ok {
		return withSubjectIDMethod.SubjectID()
	}

	subjV := reflect.ValueOf(v)
	id := reflect.Indirect(subjV).FieldByName("ID")

	return int(id.Int())
}

// Reset purges the permission set.
func Reset() {
	permissionSet = Permissions{}
}

// Authorize adds a permission allowing the `subject` Subject to perform the `action` Action on the `object` Object.
// The `conditionFunc` condition allows for more granular checks, such as ensuring the object is owned by the subject.
func ConditionalAuthorize(subject interface{}, object interface{}, action Action, conditionFunc ConditionFunc) {
	permission := Permission{
		SubjectID: getSubjectID(subject),
		ObjectType: getTypeString(object),
		Action: action,
		Condition: conditionFunc,
	}

	permissionSet = append(permissionSet, &permission)
}

// Authorize adds a permission allowing the `subject` Subject to perform the `action` Action on the `object` Object.
func Authorize(subject interface{}, object interface{}, action Action) {
	ConditionalAuthorize(subject, object, action, noConditionFunc)
}

// Can checks to see whether the `subject` Subject can perform the `action` Action on the `object` Object.
func Can(subject interface{}, object interface{}, action Action) bool {
	if len(permissionSet) == 0 {
		return false
	}

	for _, permission := range permissionSet {
		if permission.SubjectID == getSubjectID(subject) && permission.ObjectType == getTypeString(object) {
			if permission.Action == ManageAction {
				return true
			}

			if permission.Action == action {
		   		return permission.Condition(subject, object)
		   	}
		}
	}

	return false
}

// Cannot is the inverse of Can.
func Cannot(subject interface{}, object interface{}, action Action) bool {
	return !Can(subject, object, action)
}