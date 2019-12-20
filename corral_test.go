package corral

import (
	"testing")

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

func BenchmarkNoPermissions(b *testing.B) {
	for _, post := range testPosts {
		if Can(userRole, post, ReadAction) {
			b.Fatalf("Was able to perform an action without permissions.")
		}
	}
}

func TestFullCRUD(t *testing.T) {
	defer Reset()
	Authorize(adminRole.SubjectKey(), "post", ManageAction)

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

func BenchmarkFullCRUD(b *testing.B) {
	defer Reset()
	Authorize(adminRole.SubjectKey(), "post", ManageAction)

	for _, post := range testPosts {
		if Cannot(adminProfile.Role, post, CreateAction) {
			b.Fatalf("Admin was marked as manage, but cannot create.")
		}

		if Cannot(adminProfile.Role, post, ReadAction) {
			b.Fatalf("Admin was marked as manage, but cannot read.")
		}

		if Cannot(adminProfile.Role, post, UpdateAction) {
			b.Fatalf("Admin was marked as manage, but cannot update.")
		}

		if Cannot(adminProfile.Role, post, DeleteAction) {
			b.Fatalf("Admin was marked as manage, but cannot delete.")
		}

		if Can(userProfile.Role, post, CreateAction) {
			b.Fatalf("User was not authorized, but can create.")
		}

		if Can(userProfile.Role, post, ReadAction) {
			b.Fatalf("User was not authorized, but can read.")
		}

		if Can(userProfile.Role, post, UpdateAction) {
			b.Fatalf("User was not authorized, but can update.")
		}

		if Can(userProfile.Role, post, DeleteAction) {
			b.Fatalf("User was not authorized, but can delete.")
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

func TestUserComplex(t *testing.T) {
	defer Reset()
	ConditionalAuthorize(userRole.SubjectKey(), "post", ReadAction, notHidden)

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

func BenchmarkUserComplex(b *testing.B) {
	defer Reset()
	ConditionalAuthorize(userRole.SubjectKey(), "post", ReadAction, notHidden)
	ConditionalAuthorize(userRole.SubjectKey(), "post", UpdateAction, owned)

	if Cannot(userProfile.Role, testPosts[0], ReadAction) {
		b.Fatalf("User was allowed to read all posts but cannot read.")
	}

	if Cannot(userProfile.Role, testPosts[1], ReadAction) {
		b.Fatalf("User was allowed to read all posts but cannot read.")
	}

	if Can(userProfile.Role, testPosts[2], ReadAction) {
		b.Fatalf("User was allowed to read a post they shouldn't be able to see.")
	}

	if Can(userProfile.Role, testPosts[3], ReadAction) {
		b.Fatalf("User was allowed to read a post they shouldn't be able to see.")
	}
}
