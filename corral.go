package corral

// ConditionFunc is a function that returns a bool representing whether or not the subject meets some pre-defined criteria on the object.
type ConditionFunc func(subject interface{}, object interface{}) bool

// Permission is the core permission struct.
type Permission struct {
	// SubjectKey is the ID of the entity that may perform an action on an object.
	SubjectKey string
	// ObjectType is the name of the entity that may have an action performed on it.
	ObjectType string
	// Action is an individual CRUDM action.
	Action Action
	// Condition is a more granular control of whether or not the subject meets some pre-defined criteria on the object.
	Condition ConditionFunc
}

type subjectWithKeyMethod interface {
	SubjectKey() string
}

type objectWithObjectType interface {
	ObjectType() string
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
	withObjectTypeMethod, ok := v.(objectWithObjectType)
	if ok {
		return withObjectTypeMethod.ObjectType()
	}

	return ""
}

func getSubjectKey(v interface{}) string {
	subjectWithKeyMethod, ok := v.(subjectWithKeyMethod)
	if ok {
		return subjectWithKeyMethod.SubjectKey()
	}

	return ""
}

// Reset purges the permission set.
func Reset() {
	permissionSet = Permissions{}
}

// Authorize adds a permission allowing the `subject` Subject to perform the `action` Action on the `object` Object.
// The `conditionFunc` condition allows for more granular checks, such as ensuring the object is owned by the subject.
func ConditionalAuthorize(subjectKey string, objectType string, action Action, conditionFunc ConditionFunc) {
	permission := Permission{
		SubjectKey: subjectKey,
		ObjectType: objectType,
		Action: action,
		Condition: conditionFunc,
	}

	permissionSet = append(permissionSet, &permission)
}

// Authorize adds a permission allowing the `subject` Subject to perform the `action` Action on the `object` Object.
func Authorize(subjectKey string, objectType string, action Action) {
	ConditionalAuthorize(subjectKey, objectType, action, noConditionFunc)
}

// Can checks to see whether the `subject` Subject can perform the `action` Action on the `object` Object.
func Can(subject interface{}, object interface{}, action Action) bool {
	if len(permissionSet) == 0 {
		return false
	}

	for _, permission := range permissionSet {
		if permission.SubjectKey == getSubjectKey(subject) && permission.ObjectType == getTypeString(object) {
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
