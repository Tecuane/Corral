# Corral

[![Documentation](https://godoc.org/github.com/tecuane/corral?status.svg)](https://godoc.org/github.com/tecuane/corral)

A dead-simple, quite-opinionated RBAC library for web applications.

Roughly similar to the basics of [CanCanCan](https://github.com/cancancommunity/cancancan) for Rails.

The library code is messy and brittle. Check the test for example usage.

## Overview

Corral operates around the concept of subjects and objects. A subject is an entity that may perform an action, an an object is an entity that may have an action performed upon it.

In order to use Corral, only two additions need to be made:

1. Subjects must implement the following method:
  ```go
  SubjectKey() string
  ```
  SubjectKey is named as such because it may not necessarily be an identifer. It can be `role-administrator`, `1`, or a UUIDv4 string. The only mandate (for stable usage) is that it is consistent throughout the parent application using this library. This key is used to lookup the subject's permissions.

2. Objects must implement the following method:
  ```go
  ObjectType() string
  ```
  ObjectType is named as such because a subject's permission is applied to an entire class of objects. Like `SubjectKey()`, it can be any arbitrary string, so long as it is possible to differentiate between two different object types.

For conditions that may have complex unique cases, a permission can be granted with a matcher function:
```go
type ConditionFunc func(subject interface{}, object interface{}) bool
```

## Example

```go
type Role struct {
	ID int64
	Name string
}

func (r *Role) SubjectKey() string {
	return string(r.ID)
}

type Profile struct {
	ID int64
	Role *Role
	Name string
}

type Post struct {
	ID int64
	ProfileID int64
	Title string
	Hidden bool
}

func (p *Post) ObjectType() string {
	return "post"
}

// Returns false if the post is hidden.
func notHidden(profile interface{}, post interface{}) bool {
	return !post.(*Post).Hidden
}

// Returns false if the post is not owned by the profile.
func ownedByUser(profile interface{}, post interface{}) bool {
	return post.(*Post).ProfileID == profile.(*Profile).ID
}

func main() {
	var adminRole = &Role{ID: 1, Name: "Administrator"}
	var userRole = &Role{ID: 2, Name: "User"}

	// Allow an administrator to do anything with no special conditions.
	corral.Authorize(adminRole.SubjectKey(), "post", ManageAction)

	// Allow a user to read any post, so long as it is not hidden.
	corral.ConditionalAuthorize(userRole.SubjectKey(), "post", ReadAction, notHidden)

	// Allow a user to update any post, so long as it is theirs.
	corral.ConditionalAuthorize(userRole.SubjectKey(), "post", UpdateAction, ownedByUser)
}
```

Note: The use of `.SubjectKey()` above can be replaced with a literal string (e.g., `"user"`, `"administrator"`, or in the case of the above code, `"1"`), so long as it is known at development time.

See [the test file](./corral_test.go) for a more thorough example using the Go testing suite. This test includes benchmarks.

## License

Corral is licensed under the MIT license. A copy of the license text can be found in the [LICENSE file](./LICENSE).
