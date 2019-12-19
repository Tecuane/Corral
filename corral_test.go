package corral

import (
	"testing"
)

type Role struct {
	ID int64
	Name string
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

// Returns the ID of the post's owner.
func (p *Post) OwnerID() int64 {
	return p.ProfileID
}

var adminRole = &Role{ID: 1, Name: "Administrator"}
var userRole = &Role{ID: 2, Name: "User"}

var adminProfile = &Profile{ID: 1000, Role: adminRole, Name: "Administrator Profile"}
var userProfile = &Profile{ID: 2000, Role: userRole, Name: "User Profile"}

var testPosts = []*Post {
	&Post{ID: 1, ProfileID: 1, Hidden: false, Title: "Post by Administrator"},
	&Post{ID: 2, ProfileID: 2, Hidden: false,  Title: "Post by User"},
	&Post{ID: 3, ProfileID: 1, Hidden: true,  Title: "Hidden Post by Administrator"},
	&Post{ID: 4, ProfileID: 2, Hidden: true,  Title: "Hidden Post by User"},
}

func TestNoPermissions(t *testing.T) {
	for _, post := range testPosts {
		if Can(userRole, post, ReadAction) {
			t.Fatalf("Was able to perform an action without permissions.")
		}
	}
}

func TestAdminCanDoAnything(t *testing.T) {
	defer Reset()
	Authorize(adminRole, &Post{}, ManageAction)

	for _, post := range testPosts {
		if Cannot(adminProfile.Role, post, CreateAction) {
			t.Fatalf("Admin was marked as manage, but cannot create.")
		}

		if Cannot(adminProfile.Role, post, ReadAction) {
			t.Fatalf("Admin was marked as manage, but cannot read.")
		}

		if Cannot(adminProfile.Role, post, UpdateAction) {
			t.Fatalf("Admin was marked as manage, but cannot update.")
		}

		if Cannot(adminProfile.Role, post, DeleteAction) {
			t.Fatalf("Admin was marked as manage, but cannot delete.")
		}
	}
}

// Returns false if the post is hidden.
func notHidden(profile interface{}, post interface{}) bool {
	return !post.(*Post).Hidden
}

// Returns false if the post is not owned by the profile.
func owned(profile interface{}, post interface{}) bool {
	return post.(*Post).ProfileID == profile.(*Profile).ID
}

func TestUserCannotReadHidden(t *testing.T) {
	defer Reset()
	ConditionalAuthorize(userRole, &Post{}, ReadAction, notHidden)

	if Cannot(userProfile.Role, testPosts[0], ReadAction) {
		t.Fatalf("User was allowed to read all posts but cannot read.")
	}

	if Cannot(userProfile.Role, testPosts[1], ReadAction) {
		t.Fatalf("User was allowed to read all posts but cannot read.")
	}

	if Can(userProfile.Role, testPosts[2], ReadAction) {
		t.Fatalf("User was allowed to read a post they shouldn't be able to see.")
	}

	if Can(userProfile.Role, testPosts[3], ReadAction) {
		t.Fatalf("User was allowed to read a post they shouldn't be able to see.")
	}
}

func TestUserCannotCreate(t *testing.T) {
	defer Reset()
	ConditionalAuthorize(userRole, &Post{}, ReadAction, notHidden)

	if Can(userProfile.Role, &Post{}, CreateAction) {
		t.Fatalf("User could create a post.")
	}
}

func TestAdminCanCreate(t *testing.T) {
	Reset()
	Authorize(adminRole, &Post{}, CreateAction)

	if Can(userProfile.Role, &Post{}, CreateAction) {
		t.Fatalf("User could create a post.")
	}

	if Cannot(adminProfile.Role, &Post{}, CreateAction) {
		t.Fatalf("Admin could not create a post.")
	}
}